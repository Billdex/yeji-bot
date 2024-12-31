CREATE TABLE `combo_recipe`
(
    `recipe_id` int NOT NULL COMMENT '菜谱id',
    `recipe_name` varchar(30) NOT NULL DEFAULT '' COMMENT '菜谱名称',
    `need_recipe_ids` varchar(255) NOT NULL DEFAULT '[]' COMMENT '所需菜谱id列表',
    PRIMARY KEY (`recipe_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;