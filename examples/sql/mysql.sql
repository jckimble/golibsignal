CREATE TABLE `config` (
  `key` VARCHAR(32),
  `value` LONGTEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `config` VALUES ('server','https://textsecure-service.whispersystems.org:443'),('phone','+16016632014');

CREATE TABLE `contacts` (
  `tel` VARCHAR(32),
  `devices` LONGTEXT,
  `name` VARCHAR(128),
  `avatar` VARCHAR(128),
  `identitykey` LONGBLOB,
  `profilekey` LONGBLOB
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `groups` (
  `id` LONGBLOB,
  `hexid` LONGTEXT,
  `flags` int(11) DEFAULT NULL,
  `name` VARCHAR(128),
  `members` LONGTEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `prekeys` (
  `id` int(11) DEFAULT NULL,
  `key` LONGBLOB
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `sessions` (
  `recipient` int(11) DEFAULT NULL,
  `device` int(11) DEFAULT NULL,
  `data` LONGBLOB
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE UNIQUE INDEX `recipient_device` ON `sessions`(`recipient`,`device`);

CREATE TABLE `signedprekeys` (
  `id` int(11) DEFAULT NULL,
  `key` LONGBLOB
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
