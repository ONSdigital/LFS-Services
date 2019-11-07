#ifndef _SAV_READER_H
#define _SAV_READER_H

#include <stdbool.h>
#include "readstat.h"

struct Data* parse_sav(const char *input_file);
void cleanup(struct Data*);

extern void goAddData(char *, char *);

struct Header {
    char *var_name;
    char *var_description;
    readstat_type_t var_type;
    size_t length;
    int precision;
};

struct Rows {
    char **row_data;
    int row_position;
    int row_length;
};

struct Data {
    struct Header **header;
    unsigned long header_count;
    int header_position;

    struct Rows **rows;
    unsigned long row_count;
    unsigned long rows_position;

    char *buffer;
    unsigned long buffer_size;

    int variable_count;
};


#endif
