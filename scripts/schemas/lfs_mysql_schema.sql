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

create table users
(
    username varchar(255) null,
    password varchar(255) null
);

insert into users(username, password)
values ('Paul', '$2a$04$uCR1AINowJXKQxiiPwyLLubTm1k0.PWMhWDHMPE3PNu59ZglB1fLG');
