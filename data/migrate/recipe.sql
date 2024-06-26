CREATE TABLE `recipe`
(
    `recipe_id`           int          NOT NULL COMMENT '菜谱 id',
    `name`                varchar(30)  NOT NULL DEFAULT '' COMMENT '菜谱名称',
    `gallery_id`          varchar(20)  NOT NULL DEFAULT '' COMMENT '图鉴id',
    `rarity`              int          NOT NULL DEFAULT 0 COMMENT '菜谱稀有度',
    `origins`             varchar(255) NOT NULL DEFAULT '[]' COMMENT '获得来源列表',
    `stirfry`             int          NOT NULL DEFAULT 0 COMMENT '炒技法数值',
    `bake`                int          NOT NULL DEFAULT 0 COMMENT '烤技法数值',
    `boil`                int          NOT NULL DEFAULT 0 COMMENT '煮技法数值',
    `steam`               int          NOT NULL DEFAULT 0 COMMENT '蒸技法数值',
    `fry`                 int          NOT NULL DEFAULT 0 COMMENT '炸技法数值',
    `cut`                 int          NOT NULL DEFAULT 0 COMMENT '切技法数值',
    `condiment`           varchar(20)  NOT NULL DEFAULT '' COMMENT '味道',
    `price`               int          NOT NULL DEFAULT 0 COMMENT '价格',
    `ex_price`            int          NOT NULL DEFAULT 0 COMMENT '专精后的额外价格',
    `gold_efficiency`     int          NOT NULL DEFAULT 0 COMMENT '赚钱效率',
    `material_efficiency` int          NOT NULL DEFAULT 0 COMMENT '食材小号效率',
    `gift`                varchar(255) NOT NULL DEFAULT '' COMMENT '第一次做到神级送的礼物',
    `guest_gifts`         varchar(255) NOT NULL DEFAULT '[]' COMMENT '贵客礼物列表',
    `upgrade_guests`      varchar(255) NOT NULL DEFAULT '[]' COMMENT '升阶贵客列表',
    `time`                int          NOT NULL DEFAULT 0 COMMENT '制作时间,单位秒',
    `limit`               int          NOT NULL DEFAULT 0 COMMENT '每组数量限制',
    `total_time`          int          NOT NULL DEFAULT 0 COMMENT '每组总制作时间,单位秒',
    `unlock`              varchar(30)  NOT NULL DEFAULT '' COMMENT '神级可解锁',
    `combos`              varchar(255) NOT NULL DEFAULT '[]' COMMENT '可用于合成套餐列表',
    `img`                 varchar(255) NOT NULL DEFAULT '' COMMENT '图片地址',
    `materials`           varchar(255) NOT NULL DEFAULT '[]' COMMENT '食材列表',
    PRIMARY KEY (`recipe_id`),
    UNIQUE KEY `IDX_recipe_gallery_id` (`gallery_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;