#include "readstat.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "import_sav.h"

int handle_metadata(readstat_metadata_t *metadata, void *ctx) {

  int *my_var_count = (int *)ctx;
  *my_var_count = readstat_get_var_count(metadata);
  return READSTAT_HANDLER_OK;
}

char *data = NULL;
int used = 0;
int have = 0;
int line_number = 0;
const int ACCOC_SIZE = 2 * 1024 * 1024;
const int BUFFER_SIZE = 4 * 1024; // maximum size of a column's value in bytes. Bit over the top but meh

void add_to_buffer(const char *d) {
  int len = strlen(d) + 1;

  if (have < len) {
    data = realloc(data, used + ACCOC_SIZE);
    have += ACCOC_SIZE;
  }

  strcpy(&data[used], d);
  used += len - 1;
  have -= len;
}

int handle_variable(int index, readstat_variable_t *variable,
                    const char *val_labels, void *ctx) {

  return READSTAT_HANDLER_OK;
}

int handle_value(int obs_index, readstat_variable_t *variable,
                 readstat_value_t value, void *ctx) {

  int *my_var_count = (int *) ctx;
  int var_index = readstat_variable_get_index(variable);

  readstat_type_t type = readstat_value_type(value);
  const char *name = readstat_variable_get_name(variable);

   if (var_index == 0) {
     add_to_buffer("{");
   }

    char buf[BUFFER_SIZE];

    switch (type) {

    case READSTAT_TYPE_STRING:
      if (readstat_value_is_system_missing(value)) {
        snprintf(buf, sizeof(buf), "\"%s\":\"\"", name);
      } else {
        snprintf(buf, sizeof(buf), "\"%s\":\"%s\"", name, readstat_string_value(value));
      }
      add_to_buffer(buf);
      break;

    case READSTAT_TYPE_INT8:
      if (readstat_value_is_system_missing(value)) {
        snprintf(buf, sizeof(buf), "\"%s\":0", name);
      } else {
        snprintf(buf, sizeof(buf), "\"%s\":%hhd", name, readstat_int8_value(value));
      }
      add_to_buffer(buf);
      break;

    case READSTAT_TYPE_INT16:
      if (readstat_value_is_system_missing(value)) {
        snprintf(buf, sizeof(buf), "\"%s\":0", name);
      } else {
        snprintf(buf, sizeof(buf), "\"%s\":%d", name, readstat_int16_value(value));
      }
      add_to_buffer(buf);
      break;

    case READSTAT_TYPE_INT32:
      if (readstat_value_is_system_missing(value)) {
         snprintf(buf, sizeof(buf), "\"%s\":0", name);
      } else {
         snprintf(buf, sizeof(buf), "\"%s\":%d", name, readstat_int32_value(value));
      }
      add_to_buffer(buf);
      break;

    case READSTAT_TYPE_FLOAT:
      if (readstat_value_is_system_missing(value)) {
         snprintf(buf, sizeof(buf), "\"%s\":0.0", name);
      } else {
         snprintf(buf, sizeof(buf), "\"%s\":%f", name, readstat_float_value(value));
      }
      add_to_buffer(buf);

      break;

    case READSTAT_TYPE_DOUBLE:
      if (readstat_value_is_system_missing(value)) {
        snprintf(buf, sizeof(buf), "\"%s\":0.0", name);
      } else {
        snprintf(buf, sizeof(buf), "\"%s\":%lf", name, readstat_double_value(value));
      }
      add_to_buffer(buf);

      break;

    default:
      return READSTAT_HANDLER_OK;
    }


  if (var_index == *my_var_count - 1) {
    add_to_buffer("},\n");
  } else {
    add_to_buffer(",");
  }

  return READSTAT_HANDLER_OK;
}

int parse_sav(const char *input_file) {

  if (input_file == 0) {
    return -1;
  }

  int my_var_count = 0;
  readstat_error_t error = READSTAT_OK;
  readstat_parser_t *parser = readstat_parser_init();
  readstat_set_metadata_handler(parser, &handle_metadata);
  readstat_set_variable_handler(parser, &handle_variable);
  readstat_set_value_handler(parser, &handle_value);

  add_to_buffer("[");
  error = readstat_parse_sav(parser, input_file, &my_var_count);
  strcpy(&data[used - 2], "]");
  readstat_parser_free(parser);

  if (error != READSTAT_OK) {
    //printf("Error processing %s: %d\n", input_file, error);
    return -1;
  }

  goAddLine(data);

  if (data != NULL) {
    free(data);
  }

  return 0;
}