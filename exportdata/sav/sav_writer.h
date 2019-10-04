//
// Created by Paul Soule on 2019-07-25.
//

#ifndef READ_SAV_SAV_WRITER_H
#define READ_SAV_SAV_WRITER_H

#include <readstat.h>

typedef struct {
    const int sav_type;
    const char *name;
    const char *label;
} file_header;

typedef struct {
    const int sav_type;
    const int int_value;
    const float float_value;
    const double double_value;
    const char *string_value;
} data_item;

int save_sav(const char *output_file, const char *label,
             file_header **sav_header, const int column_cnt, const int data_rows, const data_item **sav_data);

#endif