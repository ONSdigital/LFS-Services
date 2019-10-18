
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for columns
-- ----------------------------
DROP TABLE IF EXISTS `columns`;
CREATE TABLE `columns` (
                           `table_name` varchar(255) NOT NULL,
                           `column_name` varchar(255) NOT NULL,
                           `column_number` int(11) NOT NULL,
                           `kind` int(255) NOT NULL,
                           `column_rows` longtext NOT NULL,
                           PRIMARY KEY (`table_name`,`column_name`),
                           UNIQUE KEY `columns_pk` (`column_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for upload_audit
-- ----------------------------
DROP TABLE IF EXISTS `upload_audit`;
CREATE TABLE `upload_audit` (
                                `id` int(11) NOT NULL AUTO_INCREMENT,
                                `file_name` varchar(1024) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
                                `reference_date` datetime DEFAULT NULL,
                                `num_var_file` int(11) DEFAULT NULL,
                                `num_var_loaded` int(11) DEFAULT NULL,
                                `num_ob_file` int(11) DEFAULT NULL,
                                `num_ob_loaded` int(11) DEFAULT NULL,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
                         `username` varchar(255) DEFAULT NULL,
                         `password` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` VALUES ('Paul', '$2a$04$uCR1AINowJXKQxiiPwyLLubTm1k0.PWMhWDHMPE3PNu59ZglB1fLG');
COMMIT;

-- ----------------------------
-- Table structure for export_definitions
-- ----------------------------
DROP TABLE IF EXISTS `export_definitions`;
CREATE TABLE `export_definitions` (
                                      `Variables` varchar(10) NOT NULL,
                                      `Research` tinyint(1) NOT NULL,
                                      `Regional_Client` tinyint(1) NOT NULL,
                                      `Government` tinyint(1) NOT NULL,
                                      `Special_License` tinyint(1) NOT NULL,
                                      `End_User` tinyint(1) NOT NULL,
                                      `Adhoc` tinyint(1) NOT NULL,
                                      PRIMARY KEY (`Variables`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
