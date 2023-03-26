

CREATE TABLE IF NOT EXISTS `user` (
	`uid` int unsigned NOT NULL AUTO_INCREMENT,
	`email` varchar(128) NOT NULL DEFAULT '',
	`name` varchar(20) NOT NULL COMMENT '用户名',
	`salt` char(12) NOT NULL DEFAULT '' COMMENT '加密随机数',
	`passwd` char(32) NOT NULL DEFAULT '' COMMENT 'md5密码',
	`last_login_ip` varchar(31) NOT NULL DEFAULT '' COMMENT '最后登录 IP',
	`last_login_unix` int unsigned COMMENT '最后一次登录时间（主动登录）',
	`created_unix` int unsigned COMMENT '账号创建时间',
	`updated_unix` int unsigned COMMENT '账号数据修改时间',
	`is_admin` BOOLEAN DEFAULT FALSE COMMENT '管理员标记',
	`prohibit_login` BOOLEAN DEFAULT FALSE COMMENT '禁止登录标记',
	PRIMARY KEY (`uid`),
	UNIQUE KEY (`name`),
	UNIQUE KEY (`email`),
	KEY `logintime` (`last_login_unix`)
  ) ENGINE=InnoDB AUTO_INCREMENT 1000000 DEFAULT CHARSET=utf8mb4 COMMENT '用户表';