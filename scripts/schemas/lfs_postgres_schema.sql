drop table if exists addresses;
drop table if exists users;
drop table if exists export_definitions;
drop table if exists annual_batch;
drop table if exists quarterly_batch;
drop table if exists survey;
drop table if exists ni_batch_item;
drop table if exists gb_batch_items;
drop table if exists monthly_batch;
drop table if exists survey_audit;
drop table if exists status_values;
drop table if exists definitions;

create table addresses
(
    id                  integer generated always as identity primary key,
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

alter table addresses
    owner to lfs;

create table export_definitions
(
    variables       varchar(10) not null primary key,
    research        bit         not null,
    regional_client bit         not null,
    government      bit         not null,
    special_license bit         not null,
    end_user        bit         not null,
    adhoc           bit         not null
);

alter table export_definitions
    owner to lfs;

create table status_values
(
    id          integer primary key,
    description varchar(255) not null
);

alter table status_values
    owner to lfs;

insert into status_values(id, description)
values (0, 'Not Started');

insert into status_values(id, description)
values (1, 'File Uploaded');

insert into status_values(id, description)
values (2, 'File Reloaded');

insert into status_values(id, description)
values (3, 'Upload Failed');

create table monthly_batch
(
    id          integer generated always as identity primary key,
    month       integer default 0 not null,
    year        integer           not null,
    status      integer default 0 not null,
    description text,
    foreign key (status) references status_values (id)
);

alter table monthly_batch
    owner to lfs;

create table ni_batch_item
(
    id     integer primary key,
    year   integer,
    month  integer,
    status integer,

    foreign key (id) references monthly_batch (id),
    foreign key (status) references status_values (id)
);

alter table ni_batch_item
    owner to lfs;

create table annual_batch
(
    id          integer generated always as identity primary key,
    year        integer null,
    status      integer null,
    description text    null,

    foreign key (status) references status_values (id)
);

create table quarterly_batch
(
    id          integer generated always as identity primary key,
    quarter     integer,
    year        integer,
    status      integer,
    description text,

    foreign key (status) references status_values (id)
);

alter table quarterly_batch
    owner to lfs;

create table gb_batch_items
(
    id     integer not null,
    year   integer not null,
    month  integer,
    week   integer not null,
    status integer,

    primary key (week, id),
    foreign key (id) references monthly_batch (id),
    foreign key (status) references status_values (id)
);

alter table gb_batch_items
    owner to lfs;

create table survey
(
    id          integer      not null,
    file_name   varchar(255) not null,
    file_source char(2),
    week        integer      not null,
    month       integer      not null,
    year        integer      not null,
    columns     jsonb        not null,

    foreign key (week, id) references gb_batch_items (week, id) on delete cascade,
    foreign key (id) references ni_batch_item (id) on delete cascade
);

alter table survey
    owner to lfs;

create index survey_id_name_index
    on survey (id);

create index survey_period_index
    on survey (year, month, week);

create table survey_audit
(
    id             integer       not null,
    file_name      varchar(1024) not null,
    file_source    char(2)       not null,
    week           integer       null,
    month          integer       null,
    year           integer       null,
    reference_date timestamp     not null,
    num_var_file   integer       not null default 0,
    num_var_loaded integer       not null default 0,
    num_ob_file    integer       not null default 0,
    num_ob_loaded  integer       not null default 0,
    status         integer       not null,
    message        text          null,

    foreign key (status) references status_values (id)
);

create index survey_audit_file_name_index
    on survey_audit (file_name);

create table users
(
    username text primary key,
    password text not null
);

alter table users
    owner to lfs;

CREATE TYPE spss_types AS ENUM ('string', 'int8', 'uint8', 'int', 'int32', 'uint32',
    'int64', 'uint64', 'float32', 'float64');

create table definitions
(
    variable    text       not null,
    description text,
    type        spss_types not null default 'string',
    length      integer,
    precision   integer,
    alias       text,
    editable    bool                default false,
    imputation  bool                default false,
    dv          bool                default false
);

create index definitions_name_index
    on definitions (variable);
