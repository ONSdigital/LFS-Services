create schema if not exists LFS collate utf8mb4_0900_ai_ci;

create table if not exists columns
(
    id            int auto_increment primary key,
    table_name    varchar(255) null,
    column_name   varchar(255) null,
    column_number int          null,
    kind          int(255)     null,
    rows          longtext     null
);

create index columns_table_name_index
    on columns (table_name);