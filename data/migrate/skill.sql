create table `skill` (
    `skill_id` int(11) unsigned not null auto_increment comment '技能ID',
    `description` varchar(255) not null default '' comment '技能描述',
    `effects` json comment '技能效果',
    primary key (`skill_id`)
) engine=InnoDB default charset=utf8 comment '技能效果';