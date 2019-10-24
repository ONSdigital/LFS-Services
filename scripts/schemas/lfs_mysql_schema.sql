-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS = @@UNIQUE_CHECKS, UNIQUE_CHECKS = 0;
SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0;
SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE =
        'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema LFS
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `LFS`;

-- -----------------------------------------------------
-- Schema LFS
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `LFS`;
USE `LFS`;

-- -----------------------------------------------------
-- Table `monthly_batch`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `monthly_batch`;

CREATE TABLE IF NOT EXISTS `monthly_batch`
(
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `month`       INT          NOT NULL DEFAULT 0,
    `year`        INT          NOT NULL,
    `status`      INT          NOT NULL DEFAULT 0,
    `description` TEXT         NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idf_UNIQUE` (`id` ASC) VISIBLE
)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `gb_batch_items`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `gb_batch_items`;

CREATE TABLE IF NOT EXISTS `gb_batch_items`
(
    `id`     INT UNSIGNED NOT NULL,
    `year`   INT          NULL,
    `month`  INT          NULL,
    `week`   INT          NOT NULL,
    `status` INT          NULL,
    PRIMARY KEY (`week`, `id`),
    CONSTRAINT `batch`
        FOREIGN KEY (`id`)
            REFERENCES `monthly_batch` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
)
    ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `ni_batch_item`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `ni_batch_item`;

CREATE TABLE IF NOT EXISTS `ni_batch_item`
(
    `id`     INT UNSIGNED NOT NULL,
    `year`   INT          NULL,
    `month`  INT          NULL,
    `status` INT          NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    CONSTRAINT `monthly`
        FOREIGN KEY (`id`)
            REFERENCES `monthly_batch` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
)
    ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `survey`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `survey`;

CREATE TABLE IF NOT EXISTS `survey`
(
    `item`          INT UNSIGNED NOT NULL,
    `table_name`    VARCHAR(255) NOT NULL,
    `column_name`   VARCHAR(255) NOT NULL,
    `column_number` INT(11)      NOT NULL,
    `kind`          INT(255)     NOT NULL,
    `column_rows`   LONGTEXT     NOT NULL,
    PRIMARY KEY (`item`, `table_name`, `column_name`),
    CONSTRAINT `gb_key`
        FOREIGN KEY (`item`)
            REFERENCES `gb_batch_items` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION,
    CONSTRAINT `ni_key`
        FOREIGN KEY (`item`)
            REFERENCES `ni_batch_item` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `upload_audit`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `upload_audit`;

CREATE TABLE IF NOT EXISTS `upload_audit`
(
    `id`             INT UNSIGNED                       NOT NULL AUTO_INCREMENT,
    `file_name`      VARCHAR(1024) CHARACTER SET 'utf8' NULL DEFAULT NULL,
    `reference_date` DATETIME                           NULL DEFAULT NULL,
    `num_var_file`   INT                                NULL DEFAULT NULL,
    `num_var_loaded` INT                                NULL DEFAULT NULL,
    `num_ob_file`    INT                                NULL DEFAULT NULL,
    `num_ob_loaded`  INT                                NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `survey`
        FOREIGN KEY (`id`)
            REFERENCES `survey` (`item`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
)
    ENGINE = InnoDB
    AUTO_INCREMENT = 6
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `users`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `users`;

CREATE TABLE IF NOT EXISTS `users`
(
    `username` VARCHAR(255) NULL DEFAULT NULL,
    `password` VARCHAR(255) NULL DEFAULT NULL
)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `export_definitions`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `export_definitions`;

CREATE TABLE IF NOT EXISTS `export_definitions`
(
    `Variables`       VARCHAR(10) NOT NULL,
    `Research`        TINYINT(1)  NOT NULL,
    `Regional_Client` TINYINT(1)  NOT NULL,
    `Government`      TINYINT(1)  NOT NULL,
    `Special_License` TINYINT(1)  NOT NULL,
    `End_User`        TINYINT(1)  NOT NULL,
    `Adhoc`           TINYINT(1)  NOT NULL,
    PRIMARY KEY (`Variables`)
)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `addresses`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `addresses`;

CREATE TABLE IF NOT EXISTS `addresses`
(
    `pcd7`                VARCHAR(7)     NOT NULL,
    `tlec99`              VARCHAR(3)     NULL DEFAULT NULL,
    `ELWA`                DECIMAL(38, 0) NULL DEFAULT NULL,
    `SCOTER`              VARCHAR(6)     NULL DEFAULT NULL,
    `Walespca`            DECIMAL(38, 0) NULL DEFAULT NULL,
    `ward03`              VARCHAR(6)     NULL DEFAULT NULL,
    `scotpca`             DECIMAL(38, 0) NULL DEFAULT NULL,
    `ukpca`               DECIMAL(38, 0) NULL DEFAULT NULL,
    `TTWA07`              DECIMAL(38, 0) NULL DEFAULT NULL,
    `ttwa08`              DECIMAL(38, 0) NULL DEFAULT NULL,
    `pca2010`             VARCHAR(3)     NULL DEFAULT NULL,
    `nuts2`               VARCHAR(4)     NULL DEFAULT NULL,
    `nuts3`               VARCHAR(5)     NULL DEFAULT NULL,
    `nuts4`               VARCHAR(7)     NULL DEFAULT NULL,
    `nuts10`              VARCHAR(10)    NULL DEFAULT NULL,
    `nuts102`             VARCHAR(4)     NULL DEFAULT NULL,
    `nuts103`             VARCHAR(5)     NULL DEFAULT NULL,
    `nuts104`             VARCHAR(7)     NULL DEFAULT NULL,
    `eregn10`             VARCHAR(2)     NULL DEFAULT NULL,
    `eregn103`            VARCHAR(3)     NULL DEFAULT NULL,
    `NUTS133`             VARCHAR(5)     NULL DEFAULT NULL,
    `NUTS132`             VARCHAR(4)     NULL DEFAULT NULL,
    `eregn133`            VARCHAR(3)     NULL DEFAULT NULL,
    `eregn13`             VARCHAR(2)     NULL DEFAULT NULL,
    `DEGURBA`             DECIMAL(38, 0) NULL DEFAULT NULL,
    `dzone1`              VARCHAR(9)     NULL DEFAULT NULL,
    `dzone2`              VARCHAR(9)     NULL DEFAULT NULL,
    `soa1`                VARCHAR(9)     NULL DEFAULT NULL,
    `soa2`                VARCHAR(9)     NULL DEFAULT NULL,
    `ward05`              VARCHAR(6)     NULL DEFAULT NULL,
    `oacode`              VARCHAR(10)    NULL DEFAULT NULL,
    `urind`               DECIMAL(38, 0) NULL DEFAULT NULL,
    `urindsul`            DECIMAL(38, 0) NULL DEFAULT NULL,
    `lea`                 VARCHAR(3)     NULL DEFAULT NULL,
    `ward98`              VARCHAR(6)     NULL DEFAULT NULL,
    `OSLAUA9d`            VARCHAR(9)     NULL DEFAULT NULL,
    `ctry9d`              VARCHAR(9)     NOT NULL,
    `casward`             VARCHAR(6)     NULL DEFAULT NULL,
    `oa11`                VARCHAR(9)     NULL DEFAULT NULL,
    `CTY`                 VARCHAR(9)     NULL DEFAULT NULL,
    `LAUA`                VARCHAR(9)     NULL DEFAULT NULL,
    `WARD`                VARCHAR(9)     NULL DEFAULT NULL,
    `CED`                 VARCHAR(9)     NULL DEFAULT NULL,
    `GOR9d`               VARCHAR(9)     NULL DEFAULT NULL,
    `PCON9d`              VARCHAR(9)     NULL DEFAULT NULL,
    `TECLEC9d`            VARCHAR(9)     NULL DEFAULT NULL,
    `TTWA9d`              VARCHAR(9)     NULL DEFAULT NULL,
    `lau2`                VARCHAR(9)     NULL DEFAULT NULL,
    `PARK`                VARCHAR(9)     NULL DEFAULT NULL,
    `LSOA11`              VARCHAR(9)     NULL DEFAULT NULL,
    `MSOA11`              VARCHAR(9)     NULL DEFAULT NULL,
    `CCG`                 VARCHAR(9)     NULL DEFAULT NULL,
    `RU11IND`             VARCHAR(2)     NULL DEFAULT NULL,
    `OAC11`               VARCHAR(3)     NULL DEFAULT NULL,
    `LEP1`                VARCHAR(9)     NULL DEFAULT NULL,
    `LEP2`                VARCHAR(9)     NULL DEFAULT NULL,
    `IMD`                 DECIMAL(38, 0) NULL DEFAULT NULL,
    `ru11indsul`          DECIMAL(38, 0) NULL DEFAULT NULL,
    `NUTS163`             VARCHAR(5)     NULL DEFAULT NULL,
    `NUTS162`             VARCHAR(4)     NULL DEFAULT NULL,
    `eregn163`            VARCHAR(3)     NULL DEFAULT NULL,
    `eregn16`             VARCHAR(2)     NOT NULL,
    `METCTY`              VARCHAR(9)     NOT NULL,
    `UTLA`                VARCHAR(9)     NOT NULL,
    `WIMD2014quintile`    DECIMAL(38, 0) NULL DEFAULT NULL,
    `decile2015`          DECIMAL(38, 0) NULL DEFAULT NULL,
    `CombinedAuthorities` VARCHAR(9)     NOT NULL,
    `id`                  INT UNSIGNED   NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    CONSTRAINT `audit`
        FOREIGN KEY (`id`)
            REFERENCES `upload_audit` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `quarterly_batch`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `quarterly_batch`;

CREATE TABLE IF NOT EXISTS `quarterly_batch`
(
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `quarter`     INT          NULL,
    `year`        INT          NULL,
    `status`      INT          NULL,
    `description` MEDIUMTEXT   NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    CONSTRAINT `qui_to_monthly`
        FOREIGN KEY (`id`)
            REFERENCES `monthly_batch` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
)
    ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `annual_batch`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `annual_batch`;

CREATE TABLE IF NOT EXISTS `annual_batch`
(
    `id`          INT UNSIGNED NOT NULL,
    `year`        INT          NULL,
    `status`      INT          NULL,
    `description` MEDIUMTEXT   NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    CONSTRAINT `an_to_monthly`
        FOREIGN KEY (`id`)
            REFERENCES `monthly_batch` (`id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
)
    ENGINE = InnoDB;


SET SQL_MODE = @OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS = @OLD_UNIQUE_CHECKS;
