#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "readstat.h"

#include "sav_generate.h"

int handle_metadata(readstat_metadata_t *metadata, void *ctx) {
  int *var_count = (int *)ctx;
  *var_count = readstat_get_var_count(metadata);
  return READSTAT_HANDLER_OK;
}

int handle_variable(int index, readstat_variable_t *variable, const char *val_labels, void *ctx) {
  int *var_count = (int *)ctx;

  int var_index = readstat_variable_get_index(variable);
  readstat_type_t type = variable->type;

  const char *name = readstat_variable_get_name(variable);

  int user_width =  variable->storage_width;

  if (index == *var_count - 1) {
    goAddHeaderItem(var_index, (char *)name, (int)type, 1, (int)user_width);
  } else {
    goAddHeaderItem(var_index, (char *)name, (int)type, 0, (int)user_width);
  }

  return READSTAT_HANDLER_OK;
}

int read_header(const char *input_file) {

  if (input_file == 0) {
    return -1;
  }

  readstat_error_t error;
  readstat_parser_t *parser = readstat_parser_init();
  readstat_set_metadata_handler(parser, &handle_metadata);
  readstat_set_variable_handler(parser, &handle_variable);

  int cnt = 0;

  error = readstat_parse_sav(parser, input_file, &cnt);

  readstat_parser_free(parser);

  if (error != READSTAT_OK) {
    return -1;
  }

  return 0;
}
