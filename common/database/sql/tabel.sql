-- kira.tbl_file definition

CREATE TABLE `tbl_file` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `owner` varchar(256) DEFAULT NULL,
  `file_id` varchar(256) DEFAULT NULL,
  `file_name` varchar(256) DEFAULT NULL,
  `file_ext` varchar(256) DEFAULT NULL,
  `file_size` bigint(20) DEFAULT NULL,
  `file_hash` varchar(256) DEFAULT NULL,
  `bucket` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tbl_file_deleted_at` (`deleted_at`),
  KEY `idx_owner` (`owner`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_file_hash` (`file_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- kira.tbl_token_user definition

CREATE TABLE `tbl_token_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `user_id` varchar(256) DEFAULT NULL,
  `token` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  KEY `idx_token` (`token`),
  KEY `idx_tbl_token_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- kira.tbl_user definition

CREATE TABLE `tbl_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `user_id` varchar(256) DEFAULT NULL,
  `user_name` varchar(256) DEFAULT NULL,
  `password` varchar(256) DEFAULT NULL,
  `role` varchar(256) DEFAULT NULL,
  `status` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  KEY `idx_tbl_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;