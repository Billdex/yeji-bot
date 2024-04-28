CREATE TABLE `equip`
(
    `equip_id`   int          NOT NULL DEFAULT 0 COMMENT '厨具id',
    `name`       varchar(30)  NOT NULL DEFAULT 0 COMMENT '厨具名称',
    `gallery_id` varchar(20)  NOT NULL DEFAULT 0 COMMENT '图鉴id',
    `origins`    varchar(255) NOT NULL DEFAULT '[]' COMMENT '来源列表',
    `rarity`     int          NOT NULL DEFAULT 0 COMMENT '厨具稀有度',
    `skills`     varchar(255) NOT NULL DEFAULT '[]' COMMENT '技能id列表',
    `img`        varchar(255) NOT NULL DEFAULT '' COMMENT '图片地址',
    PRIMARY KEY (`equip_id`),
    UNIQUE KEY uniq_gallery_id (`gallery_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '厨具数据表';