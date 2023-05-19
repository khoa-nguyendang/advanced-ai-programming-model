CREATE DATABASE IF NOT EXISTS app
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

use app;

CREATE TABLE IF NOT EXISTS logs
  (
     id            BIGINT(64) auto_increment,
     app_user_id   BIGINT(64) NOT NULL DEFAULT 0,
     activity      INT(11) NOT NULL DEFAULT 0,
     device_uuid   NVARCHAR(250) NOT NULL DEFAULT '',
     full_message  NVARCHAR(3000) NOT NULL DEFAULT '',
     date          BIGINT(64) NOT NULL DEFAULT 0,
     PRIMARY KEY (id)
  );

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `company_code` NVARCHAR(250) NULL DEFAULT '',
  `user_id` NVARCHAR(250) NULL DEFAULT '',
  `user_name` NVARCHAR(250) NULL DEFAULT '',
  `user_role_id` BIGINT NULL DEFAULT 0,
  `user_group_ids` NVARCHAR(250) NULL DEFAULT '',
  `user_info` BLOB NULL,
  `user_state` INT NOT NULL DEFAULT 0,
  `thumbnail_image_url`  NVARCHAR(2500) NULL DEFAULT '',
  `last_modified` BIGINT NULL DEFAULT 0,
  `issued_date` BIGINT NULL DEFAULT 0,
  `activation_date` BIGINT NULL DEFAULT 0,
  `expiry_date` BIGINT NULL DEFAULT 0,
  `reference_id` NVARCHAR(250) NULL DEFAULT '',
  PRIMARY KEY (`id`));

CREATE TABLE IF NOT EXISTS `defined_type` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` NVARCHAR(250) NULL DEFAULT '',
  `description`  NVARCHAR(2500) NULL DEFAULT '',
  `is_static` INT NOT NULL DEFAULT 0,
  `target_group` NVARCHAR(250) NULL DEFAULT '',
  `created_at` BIGINT NULL DEFAULT 0,
  `last_modified` BIGINT NULL DEFAULT 0,
  `created_user_id` BIGINT NULL DEFAULT 0,
  `last_modified_user_id` BIGINT NULL DEFAULT 0,
  PRIMARY KEY (`id`));

CREATE TABLE IF NOT EXISTS `user_image` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_column_id` BIGINT NULL DEFAULT 0,
  `path`  NVARCHAR(2500) NULL DEFAULT '',
  `image_type_id` INT NOT NULL DEFAULT 0,
  `created_at` BIGINT NULL DEFAULT 0,
  `last_modified` BIGINT NULL DEFAULT 0,
  `created_user_id` BIGINT NULL DEFAULT 0,
  `last_modified_user_id` BIGINT NULL DEFAULT 0,
  PRIMARY KEY (`id`));


CREATE TABLE IF NOT EXISTS `user_tracking` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `app_user_id` BIGINT NULL DEFAULT 0,
  `device_id` BIGINT NULL DEFAULT 0,
  `activity_id` BIGINT NULL DEFAULT 0,
  `timestamp` BIGINT NULL DEFAULT 0,
  PRIMARY KEY (`id`));

CREATE TABLE IF NOT EXISTS `activity` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` NVARCHAR(250) NULL DEFAULT '',
  `description`  NVARCHAR(2500) NULL DEFAULT '',
  `target_group` NVARCHAR(250) NULL DEFAULT '',
  `created_at` BIGINT NULL DEFAULT 0,
  `last_modified` BIGINT NULL DEFAULT 0,
  `created_user_id` BIGINT NULL DEFAULT 0,
  `last_modified_user_id` BIGINT NULL DEFAULT 0,
  PRIMARY KEY (`id`));

CREATE INDEX  idx_user_id ON users(user_id); 