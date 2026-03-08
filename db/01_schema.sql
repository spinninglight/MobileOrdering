CREATE DATABASE IF NOT EXISTS order_db;
use order_db;

-- 2. 设置编码（推荐，防止中文乱码）
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0; -- 暂时关闭外键检查，方便删表重建

CREATE TABLE IF NOT EXISTS `t_merchant` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商户全局唯一ID',
  `brand_id` BIGINT DEFAULT NULL COMMENT '品牌ID（指向同一个品牌的根商户或品牌聚合ID），若为NULL则表示该商户是独立单店',
  `name` VARCHAR(50) NOT NULL COMMENT '店铺名称 (如：张亮麻辣烫-高新店)',
  `logo_url` VARCHAR(255) DEFAULT NULL COMMENT '店铺Logo图片的云存储链接',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '营业状态：0-打烊，1-营业中',
  `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '商家联系电话',
  `address` VARCHAR(100) DEFAULT NULL COMMENT '店铺物理地址',
  `config_json` JSON DEFAULT NULL COMMENT '扩展配置字典 (小票机配置、主题色等)',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除：0-正常，1-已删除',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户主表 (多租户根基)';

CREATE TABLE IF NOT EXISTS `t_customer` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户全局唯一ID',
  `openid` VARCHAR(64) NOT NULL COMMENT '微信用户的 OpenID (登录核心凭证)',
  `nickname` VARCHAR(50) DEFAULT NULL COMMENT '微信昵称',
  `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '用户头像链接',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除：0-正常，1-已删除',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '首次授权注册时间',
  `last_login` DATETIME DEFAULT NULL COMMENT '最后登录时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_openid` (`openid`) COMMENT '确保一个微信用户只有一条记录'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='C端顾客表 (平台全局共享)';

CREATE TABLE IF NOT EXISTS `t_merchant_admin` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `merchant_id` BIGINT NOT NULL COMMENT '关联的商户ID (多租户隔离字段)',
  `username` VARCHAR(50) NOT NULL COMMENT '登录账号 (手机号或自定义)',
  `password` VARCHAR(100) NOT NULL COMMENT '加密后的密码',
  `role_type` TINYINT NOT NULL DEFAULT 1 COMMENT '角色：1-超级店长，2-普通店员',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除：0-正常，1-已删除',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`) COMMENT '确保登录账号全局唯一',
  KEY `idx_merchant_id` (`merchant_id`) COMMENT '加速按商户查询员工列表'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='B端商户管理员表';

-- 1. 商品分类表
CREATE TABLE IF NOT EXISTS `t_categories` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `merchant_id` BIGINT NOT NULL COMMENT '所属商户ID',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `sort_order` INT DEFAULT 0 COMMENT '排序权重(数值越大越靠前)',
  `is_visible` TINYINT DEFAULT 1 COMMENT '是否显示',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL DEFAULT NULL,
  KEY `idx_category_merchant_visible` (`merchant_id`, `is_visible`),
  KEY `idx_category_sort` (`merchant_id`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 商品基础信息表
CREATE TABLE IF NOT EXISTS `t_products` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `category_id` BIGINT NOT NULL COMMENT '所属分类ID',
  `merchant_id` BIGINT NOT NULL COMMENT '所属商户ID',
  `name` VARCHAR(100) NOT NULL COMMENT '商品名称',
  `description` TEXT COMMENT '商品详情描述',
  `image_url` VARCHAR(255) COMMENT '商品主图',
  `sales_count` INT DEFAULT 0 COMMENT '虚拟销量',
  `lowest_price` DECIMAL(10, 2) NOT NULL COMMENT '最低价',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-下架, 1-上架',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL DEFAULT NULL,
   KEY `idx_product_category_status` (`category_id`, `status`),
  KEY `idx_product_merchant_status` (`merchant_id`, `status`),
  KEY `idx_product_created` (`merchant_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 商品规格表 (Stock Keeping Unit)
-- 处理：大杯 $15, 小杯 $12 这种情况
CREATE TABLE IF NOT EXISTS `t_product_skus` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '所属商品ID',
  `spec_name` VARCHAR(50) NOT NULL COMMENT '规格名称(如:容量)',
  `spec_value` VARCHAR(50) NOT NULL COMMENT '规格值(如:大杯)',
  `price` DECIMAL(10, 2) NOT NULL COMMENT '销售价格',
  `stock` INT DEFAULT -1 COMMENT '库存: -1表示无限',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL DEFAULT NULL,
   KEY `idx_sku_product` (`product_id`),
  KEY `idx_sku_product_stock` (`product_id`, `stock`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. 商品属性表
-- 处理：加糖、冰度等不改变价格的选项
CREATE TABLE IF NOT EXISTS `t_product_attributes` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '所属商品ID',
  `attr_name` VARCHAR(50) NOT NULL COMMENT '属性名(如:甜度)',
  `attr_values` VARCHAR(255) NOT NULL COMMENT '可选属性值(用逗号隔开:正常,半糖,无糖)',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL DEFAULT NULL,
   KEY `idx_attr_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `t_order` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `order_sn` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '订单编号',
    `shop_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属店铺ID',
    `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `total_amount` DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '订单总额',
    `pay_amount` DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '应付金额',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0待支付, 1已支付, 2制作中, 3已完成, 4已取消',
    `pay_type` TINYINT NOT NULL DEFAULT 0 COMMENT '支付方式: 1微信支付, 2余额',
    `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME NULL DEFAULT NULL,
    -- 索引：多商户系统下，按店铺查询订单是最高频的操作
    UNIQUE INDEX `idx_order_sn` (`order_sn`),
    INDEX `idx_shop_id` (`shop_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单主表';

CREATE TABLE IF NOT EXISTS `t_order_item` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `order_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '对应主订单ID',
    `shop_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属店铺ID',
    `product_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品ID',
    `product_name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '商品名称',
    `sku_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '规格ID',
    `sku_info` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '规格描述',
    `price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '单价',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
    `total_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '总价',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME NULL DEFAULT NULL,
    -- 索引
    INDEX `idx_order_id` (`order_id`),
    INDEX `idx_shop_id` (`shop_id`) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单详情表';

SET FOREIGN_KEY_CHECKS = 1; -- 恢复外键检查