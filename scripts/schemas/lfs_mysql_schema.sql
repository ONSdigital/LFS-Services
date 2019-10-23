SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for columns
-- ----------------------------
DROP TABLE IF EXISTS `columns`;
CREATE TABLE `columns`
(
    `table_name`    varchar(255) NOT NULL,
    `column_name`   varchar(255) NOT NULL,
    `column_number` int(11)      NOT NULL,
    `kind`          int(255)     NOT NULL,
    `column_rows`   longtext     NOT NULL,
    PRIMARY KEY (`table_name`, `column_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for upload_audit
-- ----------------------------
DROP TABLE IF EXISTS `upload_audit`;
CREATE TABLE `upload_audit`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT,
    `file_name`      varchar(1024) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
    `reference_date` datetime                                                 DEFAULT NULL,
    `num_var_file`   int(11)                                                  DEFAULT NULL,
    `num_var_loaded` int(11)                                                  DEFAULT NULL,
    `num_ob_file`    int(11)                                                  DEFAULT NULL,
    `num_ob_loaded`  int(11)                                                  DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 6
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `username` varchar(255) DEFAULT NULL,
    `password` varchar(255) DEFAULT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for export_definitions
-- ----------------------------
DROP TABLE IF EXISTS `export_definitions`;
CREATE TABLE `export_definitions`
(
    `Variables`       varchar(10) NOT NULL,
    `Research`        tinyint(1)  NOT NULL,
    `Regional_Client` tinyint(1)  NOT NULL,
    `Government`      tinyint(1)  NOT NULL,
    `Special_License` tinyint(1)  NOT NULL,
    `End_User`        tinyint(1)  NOT NULL,
    `Adhoc`           tinyint(1)  NOT NULL,
    PRIMARY KEY (`Variables`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `addresses`;
CREATE TABLE `addresses`
(
    pcd7                  VARCHAR(7) NOT NULL,
    tlec99                VARCHAR(3),
    `ELWA`                DECIMAL(38, 0),
    `SCOTER`              VARCHAR(6),
    `Walespca`            DECIMAL(38, 0),
    ward03                VARCHAR(6),
    scotpca               DECIMAL(38, 0),
    ukpca                 DECIMAL(38, 0),
    `TTWA07`              DECIMAL(38, 0),
    ttwa08                DECIMAL(38, 0),
    pca2010               VARCHAR(3),
    nuts2                 VARCHAR(4),
    nuts3                 VARCHAR(5),
    nuts4                 VARCHAR(7),
    nuts10                VARCHAR(10),
    nuts102               VARCHAR(4),
    nuts103               VARCHAR(5),
    nuts104               VARCHAR(7),
    eregn10               VARCHAR(2),
    eregn103              VARCHAR(3),
    `NUTS133`             VARCHAR(5),
    `NUTS132`             VARCHAR(4),
    eregn133              VARCHAR(3),
    eregn13               VARCHAR(2),
    `DEGURBA`             DECIMAL(38, 0),
    dzone1                VARCHAR(9),
    dzone2                VARCHAR(9),
    soa1                  VARCHAR(9),
    soa2                  VARCHAR(9),
    ward05                VARCHAR(6),
    oacode                VARCHAR(10),
    urind                 DECIMAL(38, 0),
    urindsul              DECIMAL(38, 0),
    lea                   VARCHAR(3),
    ward98                VARCHAR(6),
    `OSLAUA9d`            VARCHAR(9),
    ctry9d                VARCHAR(9) NOT NULL,
    casward               VARCHAR(6),
    oa11                  VARCHAR(9),
    `CTY`                 VARCHAR(9),
    `LAUA`                VARCHAR(9),
    `WARD`                VARCHAR(9),
    `CED`                 VARCHAR(9),
    `GOR9d`               VARCHAR(9),
    `PCON9d`              VARCHAR(9),
    `TECLEC9d`            VARCHAR(9),
    `TTWA9d`              VARCHAR(9),
    lau2                  VARCHAR(9),
    `PARK`                VARCHAR(9),
    `LSOA11`              VARCHAR(9),
    `MSOA11`              VARCHAR(9),
    `CCG`                 VARCHAR(9),
    `RU11IND`             VARCHAR(2),
    `OAC11`               VARCHAR(3),
    `LEP1`                VARCHAR(9),
    `LEP2`                VARCHAR(9),
    `IMD`                 DECIMAL(38, 0),
    ru11indsul            DECIMAL(38, 0),
    `NUTS163`             VARCHAR(5),
    `NUTS162`             VARCHAR(4),
    eregn163              VARCHAR(3),
    eregn16               VARCHAR(2) NOT NULL,
    `METCTY`              VARCHAR(9) NOT NULL,
    `UTLA`                VARCHAR(9) NOT NULL,
    `WIMD2014quintile`    DECIMAL(38, 0),
    decile2015            DECIMAL(38, 0),
    `CombinedAuthorities` VARCHAR(9) NOT NULL
);
-- ----------------------------
-- Records of users
-- ----------------------------

insert into users(username, password)
values ('Admin', '$2a$04$Su7c9o6E9pLaGut2Nv9FqO2ZUbntDmUweOlO/Vj3hczi86qrnbKK2');
