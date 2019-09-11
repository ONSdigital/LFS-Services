#include "sav_writer.h"
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

const int MAX_STRING = 255;

readstat_variable_t *save_header(file_header *const *sav_header, int column_cnt,
                                 readstat_writer_t *writer);

static ssize_t write_bytes(const void *data, size_t len, void *ctx) {
    int fd = *(int *) ctx;
    return write(fd, data, len);
}

int save_sav(const char *output_file, const char *label, file_header **sav_header, const int column_cnt,
             const int row_count, const data_item **sav_data) {


    readstat_writer_t *writer = readstat_writer_init();
    readstat_set_data_writer(writer, &write_bytes);
    readstat_writer_set_file_label(writer, label);
    readstat_writer_set_compression(writer, READSTAT_COMPRESS_ROWS);

    readstat_variable_t **variables = malloc(sizeof(readstat_variable_t) * column_cnt);

    for (int i = 0; i < column_cnt; i++) {
        unsigned long cnt = 20;
        if (sav_header[i]->sav_type == READSTAT_TYPE_STRING) {
            cnt = MAX_STRING;
        }
        readstat_variable_t *variable =
                readstat_add_variable(writer, sav_header[i]->name, sav_header[i]->sav_type, cnt);
        readstat_variable_set_label(variable, sav_header[i]->label);
        variables[i] = variable;
    }

    int fd = open(output_file, O_WRONLY | O_CREAT | O_TRUNC, 0666);

    if (fd == -1) {
        return -1;
    }

    readstat_begin_writing_sav(writer, &fd, row_count);

    int cnt = 0;

    for (int i = 0; i < row_count; i++) {
        readstat_begin_row(writer);

        for (int j = 0; j < column_cnt; j++) {
            readstat_variable_t *variable = variables[j];
            switch (sav_data[cnt]->sav_type) {
                case READSTAT_TYPE_STRING:
                   {
                        int len = strlen(sav_data[cnt]->string_value);
                        if (len == 0) {
                            readstat_insert_missing_value(writer, variable);
                        } else {
                            int to_allocate = (len > MAX_STRING ? MAX_STRING : len);
                            char *str = malloc(to_allocate + 1);
                            if (!str) { return -1; }

                            memcpy(str, sav_data[cnt]->string_value, to_allocate);

                            str[to_allocate] = 0;
                            readstat_insert_string_value(writer, variable, (const char *) str);
                            free(str);
                        }
                    }
                    break;

                case READSTAT_TYPE_INT8: {
                        int i = sav_data[cnt]->int_value;
                        if (i == 0) {
                           readstat_insert_missing_value(writer, variable);
                        } else {
                            readstat_insert_int8_value(writer, variable, sav_data[cnt]->int_value);
                        }
                    }
                    break;

                case READSTAT_TYPE_INT16:{
                     int i = sav_data[cnt]->int_value;
                     if (i == 0) {
                        readstat_insert_missing_value(writer, variable);
                     } else {
                         readstat_insert_int8_value(writer, variable, sav_data[cnt]->int_value);
                     }
                 }
                 break;

                case READSTAT_TYPE_INT32: {
                    int i = sav_data[cnt]->int_value;
                    if (i == 0) {
                       readstat_insert_missing_value(writer, variable);
                    } else {
                        readstat_insert_int8_value(writer, variable, sav_data[cnt]->int_value);
                    }
                }
                break;

                case READSTAT_TYPE_FLOAT: {
                        float f = sav_data[cnt]->float_value;
                        if (f == 0.0) {
                           readstat_insert_missing_value(writer, variable);
                        } else {
                            readstat_insert_float_value(writer, variable, f);
                        }
                    }
                break;

                case READSTAT_TYPE_DOUBLE: {
                        long double d = sav_data[cnt]->double_value;
                        if (d == 0.0) {
                           readstat_insert_missing_value(writer, variable);
                        } else {
                           readstat_insert_double_value(writer, variable, d);
                        }
                    }
                break;
            }

            cnt++;

        }

        readstat_end_row(writer);
    }

    readstat_end_writing(writer);
    readstat_writer_free(writer);
    close(fd);

    return 0;
}
