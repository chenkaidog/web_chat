CREATE TABLE `account`
(
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `account_id` VARCHAR(128)    NOT NULL COMMENT 'unique account id',
    `username`   VARCHAR(64)     NOT NULL COMMENT '用户帐号',
    `password`   VARCHAR(256)    NOT NULL COMMENT '密码md5',
    `salt`       VARCHAR(256)    NOT NULL COMMENT '盐',
    `status`     VARCHAR(32)     NOT NULL COMMENT '帐号状态',
    `expiration_time`DATETIME        NOT NULL COMMENT '账号有效期',
    `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_account_id` (`account_id`),
    UNIQUE INDEX `uniq_username` (`username`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
    COMMENT '用户帐号表';
