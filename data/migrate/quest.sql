CREATE TABLE `quest` (
    `quest_id` int(11) NOT NULL COMMENT '任务ID',
    `quest_id_disp` decimal(8, 2) NOT NULL DEFAULT 0 COMMENT '系列任务编号顺序',
    `type` varchar(60) NOT NULL DEFAULT '' COMMENT '任务类型',
    `goal` varchar(255) NOT NULL DEFAULT '' COMMENT '任务目标描述',
    `rewards` json COMMENT '任务奖励描述',
    `conditions` json COMMENT '任务条件描述',
    PRIMARY KEY (`quest_id`)
) engine = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT '任务数据表';