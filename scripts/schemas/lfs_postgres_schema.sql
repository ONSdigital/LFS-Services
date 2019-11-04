create table addresses
(
    id                  integer generated always as identity
        constraint addresses_pkey
        primary key
        constraint id_unique
        unique,
    pcd7                varchar(7) not null,
    tlec99              varchar(3),
    elwa                numeric(38),
    scoter              varchar(6),
    walespca            numeric(38),
    ward03              varchar(6),
    scotpca             numeric(38),
    ukpca               numeric(38),
    ttwa07              numeric(38),
    ttwa08              numeric(38),
    pca2010             varchar(3),
    nuts2               varchar(4),
    nuts3               varchar(5),
    nuts4               varchar(7),
    nuts10              varchar(10),
    nuts102             varchar(4),
    nuts103             varchar(5),
    nuts104             varchar(7),
    eregn10             varchar(2),
    eregn103            varchar(3),
    nuts133             varchar(5),
    nuts132             varchar(4),
    eregn133            varchar(3),
    eregn13             varchar(2),
    degurba             numeric(38),
    dzone1              varchar(9),
    dzone2              varchar(9),
    soa1                varchar(9),
    soa2                varchar(9),
    ward05              varchar(6),
    oacode              varchar(10),
    urind               numeric(38),
    urindsul            numeric(38),
    lea                 varchar(3),
    ward98              varchar(6),
    oslaua9d            varchar(9),
    ctry9d              varchar(9) not null,
    casward             varchar(6),
    oa11                varchar(9),
    cty                 varchar(9),
    laua                varchar(9),
    ward                varchar(9),
    ced                 varchar(9),
    gor9d               varchar(9),
    pcon9d              varchar(9),
    teclec9d            varchar(9),
    ttwa9d              varchar(9),
    lau2                varchar(9),
    park                varchar(9),
    lsoa11              varchar(9),
    msoa11              varchar(9),
    ccg                 varchar(9),
    ru11ind             varchar(2),
    oac11               varchar(3),
    lep1                varchar(9),
    lep2                varchar(9),
    imd                 numeric(38),
    ru11indsul          numeric(38),
    nuts163             varchar(5),
    nuts162             varchar(4),
    eregn163            varchar(3),
    eregn16             varchar(2) not null,
    metcty              varchar(9) not null,
    utla                varchar(9) not null,
    wimd2014quintile    numeric(38),
    decile2015          numeric(38),
    combinedauthorities varchar(9) not null
);

alter table addresses owner to lfs;

create table export_definitions
(
    variables       varchar(10) not null
        constraint export_definitions_pkey
        primary key,
    research        bit         not null,
    regional_client bit         not null,
    government      bit         not null,
    special_license bit         not null,
    end_user        bit         not null,
    adhoc           bit         not null
);

alter table export_definitions owner to lfs;

create table status_values
(
    id          integer      not null
        constraint status_values_pkey
        primary key
        constraint status_values_id_uindex
        unique,
    description varchar(255) not null
);

alter table status_values owner to lfs;

create table monthly_batch
(
    id          integer generated always as identity
        constraint monthly_batch_pkey
        primary key
        constraint idf_unique
        unique,
    month       integer default 0 not null,
    year        integer           not null,
    status      integer default 0 not null
        constraint monthly_batch_status_values_id_fk
        references status_values,
    description text
);

alter table monthly_batch owner to lfs;

create table annual_batch
(
    id          integer not null
        constraint annual_batch_pkey
        primary key,
    year        integer,
    status      integer,
    description varchar(255)
);

alter table annual_batch owner to lfs;

create table ni_batch_item
(
    id     integer not null
        constraint ni_batch_item_pkey
        primary key
        constraint ni_id_unique
        unique
        constraint monthly
        references monthly_batch,
    year   integer,
    month  integer,
    status integer
        constraint ni_batch_item_status_values_id_fk
        references status_values
);

alter table ni_batch_item owner to lfs;

create table quarterly_batch
(
    id          integer generated always as identity
        constraint quarterly_batch_pkey
        primary key
        constraint qb_to_mb
        references monthly_batch,
    quarter     integer,
    year        integer,
    status      integer
        constraint quarterly_batch_status_values_id_fk
        references status_values,
    description varchar(255)
);

alter table quarterly_batch owner to lfs;

create table gb_batch_items
(
    id     integer not null
        constraint gb_batch_items_id_key
        unique
        constraint batch
        references monthly_batch,
    year   integer,
    month  integer,
    week   integer not null,
    status integer
        constraint gb_batch_items_status_values_id_fk
        references status_values,
    constraint gb_batch_items_pkey
        primary key (week, id)
);

alter table gb_batch_items owner to lfs;

create table survey
(
    id          integer      not null
        constraint gb_key
        references gb_batch_items (id)
        on delete cascade
        constraint ni_key
        references ni_batch_item
        on delete cascade,
    file_name   varchar(255) not null,
    file_source char(2),
    week        integer      not null,
    month       integer      not null,
    year        integer      not null,
    columns     json         not null
);

alter table survey owner to lfs;

create index survey_id_name_index
    on survey (id);

create index survey_period_index
    on survey (year, month, week);

create table users
(
    username varchar(255) not null
        constraint users_pkey
        primary key
        constraint users_username_uindex
        unique,
    password varchar(255) not null
);

alter table users owner to lfs;

