CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `passport` varchar(45) NOT NULL COMMENT 'User Passport',
  `password` varchar(45) NOT NULL COMMENT 'User Password',
  `nickname` varchar(45) NOT NULL COMMENT 'User Nickname',
  `create_at` datetime DEFAULT NULL COMMENT 'Created Time',
  `update_at` datetime DEFAULT NULL COMMENT 'Updated Time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `file` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'FILE ID',
  `created_at` datetime DEFAULT NULL COMMENT '文件创建时间',
  `name` varchar(50) NOT NULL COMMENT '文件名',
  `md5` varchar(50) NOT NULL COMMENT '文件MD5值',
  `size` bigint NOT NULL COMMENT '文件大小',
  `storage_path` varchar(50) NOT NULL COMMENT '文件存储目录',
  `uncompressed_path` varchar(50) NOT NULL COMMENT '文件解压之后目录',
  `import` tinyint NOT NULL DEFAULT 0 COMMENT '文件是否被解压,默认未处理',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `image` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'image auto increment ID',
    `created_at` datetime DEFAULT NULL COMMENT '镜像记录创建时间',
    `new` varchar(50) NOT NULL COMMENT '更改后的新镜像名称',
    `old` varchar(50) NOT NULL COMMENT '上传的镜像名称',
    `file_id` int NOT NULL COMMENT '对应文件自增ID',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;