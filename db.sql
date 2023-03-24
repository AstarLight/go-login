

CREATE TABLE IF NOT EXISTS `user` (
	`uid` int unsigned NOT NULL,
	`email` varchar(128) NOT NULL DEFAULT '',
	`username` varchar(20) NOT NULL COMMENT '用户名',
	`passcode` char(12) NOT NULL DEFAULT '' COMMENT '加密随机数',
	`passwd` char(32) NOT NULL DEFAULT '' COMMENT 'md5密码',
	`login_ip` varchar(31) NOT NULL DEFAULT '' COMMENT '最后登录 IP',
	`login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后一次登录时间（主动登录或cookie登录）',
	PRIMARY KEY (`uid`),
	UNIQUE KEY (`username`),
	UNIQUE KEY (`email`),
	KEY `logintime` (`login_time`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户登录表';