CREATE TABLE `user` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                        `uid` varchar(128) NOT NULL DEFAULT '' COMMENT 'uid',
                        `email` varchar(128) NOT NULL DEFAULT '' COMMENT 'email',
                        `password` varchar(128) NOT NULL DEFAULT '' COMMENT 'password',
                        `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT 'nickname',
                        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
                        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
                        `deleted_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'deleted time',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `uni_uid` (`uid`)
                            USING BTREE, KEY `idx_email` (`email`),
                        KEY `idx_deleted_at` (`deleted_at`)) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'user table';