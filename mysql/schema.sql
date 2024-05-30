CREATE
DATABASE pcap CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE `record`
(
    `record_id` INT AUTO_INCREMENT,
    `time`      DATETIME(6) NOT NULL,
    `file`      VARCHAR(255) NOT NULL,
    `seq`       INT DEFAULT NULL,
    `rel-id`    INT DEFAULT NULL,
    `src`       VARCHAR(255) NOT NULL,
    `dst`       VARCHAR(255) NOT NULL,
    `data`      VARCHAR(255) NOT NULL,
    PRIMARY KEY (`record_id`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE utf8mb4_unicode_ci
;

CREATE TABLE `record_grouped_rt`
(
    `record_id` INT AUTO_INCREMENT,
    `time`      DATETIME(6) NOT NULL,
    `file`      VARCHAR(255) NOT NULL,
    `seq`       INT DEFAULT NULL,
    `rel-id`    INT DEFAULT NULL,
    `src`       VARCHAR(255) NOT NULL,
    `dst`       VARCHAR(255) NOT NULL,
    `data`      VARCHAR(255) NOT NULL,
    PRIMARY KEY (`record_id`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE utf8mb4_unicode_ci
;

