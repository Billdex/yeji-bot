CREATE TABLE `material`
(
    `material_id` int         NOT NULL DEFAULT 0 COMMENT '食材id',
    `name`        varchar(30) NOT NULL DEFAULT 0 COMMENT '食材名称',
    `origin`      varchar(20) NOT NULL DEFAULT '' COMMENT '食材来源',
    PRIMARY KEY (`material_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '食材数据';