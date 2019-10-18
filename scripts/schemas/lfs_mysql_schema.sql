create schema if not exists LFS collate utf8mb4_0900_ai_ci;

drop table if exists columns;

create table columns
(
    table_name varchar(255) not null,
    column_name varchar(255) not null,
    column_number int not null,
    kind int(255) not null,
    column_rows longtext not null,
    primary key (table_name, column_name),
    constraint columns_pk
        unique (column_name)
);


create table if not exists upload_audit
(
    id             int auto_increment
        primary key,
    file_name      varchar(1024) charset utf8 null,
    reference_date datetime                   null,
    num_var_file   int                        null,
    num_var_loaded int                        null,
    num_ob_file    int                        null,
    num_ob_loaded  int                        null
);

drop table if exists users;

create table users
(
    username varchar(255) null,
    password varchar(255) null
);

insert into users(username, password)
values ('Paul', '$2a$04$uCR1AINowJXKQxiiPwyLLubTm1k0.PWMhWDHMPE3PNu59ZglB1fLG');

create table if not exists export_definitions
(
    Variables varchar(10) not null
        primary key,
    Research tinyint(1) not null,
    Regional_Client tinyint(1) not null,
    Government tinyint(1) not null,
    Special_License tinyint(1) not null,
    End_User tinyint(1) not null,
    Adhoc tinyint(1) not null
);

