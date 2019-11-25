CREATE TABLE `nos_metadata` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `object_name` varchar(150) NOT NULL DEFAULT '' COMMENT '对象名称',
  `sha256_code` varchar(300) NOT NULL DEFAULT '' COMMENT '对象数据的sha256加密码',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '该对象是否已删除，0：未删除，1：已删除',
  `add_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_object_name_is_del` (`object_name`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='分布式对象存储系统NOS的元数据表';
