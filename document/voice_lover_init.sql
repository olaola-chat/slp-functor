CREATE TABLE `voice_lover_audio`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `title`        varchar(128) NOT NULL COMMENT '标题',
    `desc`         text COMMENT '简介',
    `resource`     text COMMENT '资源链接',
    `cover`        text COMMENT '封面链接',
    `from`         tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '作品来源：0-无 1-原创 2-搬运',
    `seconds`      int(10) unsigned NOT NULL DEFAULT '0' COMMENT '作品时长 单位秒',
    `audit_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态：0-待审核 1-审核通过 2-审核不通过',
    `audit_reason` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '审核不通过原因：0-默认',
    `op_uid`       bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '审核人',
    `pub_uid`      bigint(20) unsigned NOT NULL COMMENT '发布人',
    `apply_time`   bigint(20) unsigned NOT NULL COMMENT '审核时间',
    `create_time`  bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time`  bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY            `idx_puid` (`pub_uid`),
    KEY            `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人音频表';

CREATE TABLE `voice_lover_audio_partner`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
    `audio_id`    bigint(20) unsigned NOT NULL COMMENT '音频id',
    `type`        tinyint(3) NOT NULL COMMENT '1:配音 2:文案 3:后期 4:封面设计',
    `uid`         bigint(20) unsigned NOT NULL COMMENT '参与人id',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_adid` (`audio_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '房间恋人音频参与人表';

CREATE TABLE `voice_lover_audio_label`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
    `audio_id`    bigint(20) unsigned NOT NULL COMMENT '音频id',
    `label`       varchar(256) NOT NULL DEFAULT '' COMMENT '标签文案',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_adid` (`audio_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '房间恋人音频标签表';

CREATE TABLE `voice_lover_audio_comment`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `audio_id`    bigint(20) unsigned NOT NULL COMMENT '音频id',
    `uid`         bigint(20) unsigned NOT NULL COMMENT '用户uid',
    `content`     text COMMENT '评论内容',
    `type`        tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态：0-评论 1-弹幕',
    `status`      tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态：0-默认状态 1-举报封禁',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    `address`     varchar(64) NOT NULL DEFAULT '' COMMENT '位置',
    PRIMARY KEY (`id`, `create_time`),
    KEY           `idx_adid` (`audio_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人音频评论表' PARTITION BY RANGE (`create_time`) (
    PARTITION `p202312` VALUES LESS THAN (1704038400),
    PARTITION `p202406` VALUES LESS THAN (1719763200),
    PARTITION `p202412` VALUES LESS THAN (1735660800),
    PARTITION `p202506` VALUES LESS THAN (1751299200),
    PARTITION `p202512` VALUES LESS THAN (1767196800)
);

CREATE TABLE `voice_lover_album`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '专辑编号',
    `name`        varchar(256) NOT NULL COMMENT '专辑名称',
    `intro`       text COMMENT '专辑简介',
    `cover`       text COMMENT '专辑封面',
    `is_deleted`  tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否删除:0-未删除 1-已删除',
    `op_uid`      bigint(20) unsigned NOT NULL COMMENT '操作人',
    `choice`      tinyint(3) NOT NULL DEFAULT '0' COMMENT '类型:0:默认 1:精选',
    `choice_time` bigint(20) unsigned DEFAULT '0' COMMENT '设置类型时间，choice非0时写入',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人专辑表';

CREATE TABLE `voice_lover_album_comment`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `album_id`    bigint(20) unsigned NOT NULL COMMENT '专辑id',
    `uid`         bigint(20) unsigned NOT NULL COMMENT '用户uid',
    `content`     text COMMENT '评论内容',
    `status`      tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态：0-默认状态 1-举报封禁',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    `address`     varchar(64) NOT NULL DEFAULT '' COMMENT '位置',
    PRIMARY KEY (`id`, `create_time`),
    KEY           `idx_abid` (`album_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人专辑评论表' PARTITION BY RANGE (`create_time`) (
    PARTITION `p202312` VALUES LESS THAN (1704038400),
    PARTITION `p202406` VALUES LESS THAN (1719763200),
    PARTITION `p202412` VALUES LESS THAN (1735660800),
    PARTITION `p202506` VALUES LESS THAN (1751299200),
    PARTITION `p202512` VALUES LESS THAN (1767196800)
);

CREATE TABLE `voice_lover_subject`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '专题编号',
    `name`        varchar(256) NOT NULL COMMENT '专题名称',
    `op_id`       bigint(20) unsigned NOT NULL COMMENT '操作人',
    `is_deleted`  tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否删除:0-未删除 1-已删除',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='房间恋人专题表';

CREATE TABLE `voice_lover_audio_album`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `audio_id`    bigint(20) unsigned NOT NULL COMMENT '音频id',
    `album_id`    bigint(20) unsigned NOT NULL COMMENT '专辑id',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_adid_abid` (`audio_id`,`album_id`),
    KEY           `idx_adid` (`audio_id`),
    KEY           `idx_abid` (`album_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人音频-专辑关系表';

CREATE TABLE `voice_lover_album_subject`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `album_id`    bigint(20) unsigned NOT NULL COMMENT '专辑id',
    `subject_id`  bigint(20) unsigned NOT NULL COMMENT '专题id',
    `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_abid_sjid` (`album_id`,`subject_id`),
    KEY           `idx_abid` (`album_id`),
    KEY           `idx_sjid` (`subject_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人专辑-专题关系表';

CREATE TABLE `voice_lover_user_collect`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `uid`          bigint(20) unsigned NOT NULL COMMENT '用户uid',
    `collect_id`   bigint(20) unsigned NOT NULL COMMENT '收藏的资源id',
    `collect_type` tinyint(3) NOT NULL COMMENT '收藏资源类型:0-专辑 1-音频',
    `create_time`  bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`  bigint unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人收藏表';

CREATE TABLE `voice_lover_banner`
(
    `id`       bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `title`    varchar(128) NOT NULL COMMENT '标题',
    `cover`    TEXT COMMENT '专辑封面',
    `schema`   TEXT COMMENT '跳转链接',
    `op_uid`      bigint(20) unsigned NOT NULL COMMENT '操作人',
    `start_time` bigint(20) unsigned NOT NULL COMMENT '开始时间',
    `end_time`  bigint(20) unsigned NOT NULL COMMENT '结束时间',
    `sort`      tinyint(3) unsigned NOT NULL COMMENT '排序',
    `create_time`  bigint(20) unsigned NOT NULL COMMENT '创建时间',
    `update_time`  bigint unsigned NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='声音恋人banner表';