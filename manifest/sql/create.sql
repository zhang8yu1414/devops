CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `passport` varchar(45) NOT NULL COMMENT 'User Passport',
  `password` varchar(45) NOT NULL COMMENT 'User Password',
  `nickname` varchar(45) NOT NULL COMMENT 'User Nickname',
  `create_at` datetime DEFAULT NULL COMMENT 'Created Time',
  `update_at` datetime DEFAULT NULL COMMENT 'Updated Time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `file` (
  `id` int(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'FILE ID',
  `created_at` timestamp NOT NULL COMMENT '文件创建时间',
  `name` varchar(50) NOT NULL COMMENT '文件名',
  `md5` varchar(50) NOT NULL COMMENT '文件MD5值',
  `size` bigint NOT NULL COMMENT '文件大小',
  `path` varchar(50) NOT NULL COMMENT '文件存储目录',
  `import` tinyint NOT NULL DEFAULT 0 COMMENT '文件是否被处理,默认未处理',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;