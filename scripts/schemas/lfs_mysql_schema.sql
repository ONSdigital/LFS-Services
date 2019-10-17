create schema if not exists LFS collate utf8mb4_0900_ai_ci;

drop table if exists columns;

create table columns
(
    id            int auto_increment primary key,
    table_name    varchar(255) null,
    column_name   varchar(255) null,
    column_number int          null,
    kind          int(255)     null,
    column_rows   longtext     null
);

create index columns_table_name_index
    on columns (table_name);

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

create table if not exists addresses
(
    pcd7 varchar(7) not null unique primary key,
    tlec99 varchar(3) null,
    ELWA decimal(38) null,
    SCOTER varchar(6) null,
    Walespca decimal(38) null,
    ward03 varchar(6) null,
    scotpca decimal(38) null,
    ukpca decimal(38) null,
    TTWA07 decimal(38) null,
    ttwa08 decimal(38) null,
    pca2010 varchar(3) null,
    nuts2 varchar(4) null,
    nuts3 varchar(5) null,
    nuts4 varchar(7) null,
    nuts10 varchar(10) null,
    nuts102 varchar(4) null,
    nuts103 varchar(5) null,
    nuts104 varchar(7) null,
    eregn10 varchar(2) null,
    eregn103 varchar(3) null,
    NUTS133 varchar(5) null,
    NUTS132 varchar(4) null,
    eregn133 varchar(3) null,
    eregn13 varchar(2) null,
    DEGURBA decimal(38) null,
    dzone1 varchar(9) null,
    dzone2 varchar(9) null,
    soa1 varchar(9) null,
    soa2 varchar(9) null,
    ward05 varchar(6) null,
    oacode varchar(10) null,
    urind decimal(38) null,
    urindsul decimal(38) null,
    lea varchar(3) null,
    ward98 varchar(6) null,
    OSLAUA9d varchar(9) null,
    ctry9d varchar(9) not null,
    casward varchar(6) null,
    oa11 varchar(9) null,
    CTY varchar(9) null,
    LAUA varchar(9) null,
    WARD varchar(9) null,
    CED varchar(9) null,
    GOR9d varchar(9) null,
    PCON9d varchar(9) null,
    TECLEC9d varchar(9) null,
    TTWA9d varchar(9) null,
    lau2 varchar(9) null,
    PARK varchar(9) null,
    LSOA11 varchar(9) null,
    MSOA11 varchar(9) null,
    CCG varchar(9) null,
    RU11IND varchar(2) null,
    OAC11 varchar(3) null,
    LEP1 varchar(9) null,
    LEP2 varchar(9) null,
    IMD decimal(38) null,
    ru11indsul decimal(38) null,
    NUTS163 varchar(5) null,
    NUTS162 varchar(4) null,
    eregn163 varchar(3) null,
    eregn16 varchar(2) not null,
    METCTY varchar(9) not null,
    UTLA varchar(9) not null,
    WIMD2014quintile decimal(38) null,
    decile2015 decimal(38) null,
    CombinedAuthorities varchar(9) not null
);
