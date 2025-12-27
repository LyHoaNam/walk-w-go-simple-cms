USE simple_golang_db;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Johnson', 'bob@example.com')
ON DUPLICATE KEY UPDATE name=name;

CREATE TABLE IF NOT EXISTS `customer` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `first_name` varchar(225) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `last_name` varchar(225) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `email` varchar(225) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `address` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `phone_number` VARCHAR(32) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT   CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_customer_phone` (`phone_number`),
    KEY `idx_created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `category` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    KEY `idx_created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `product` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` text,
    `dimension` text,
    `category_id` bigint NOT NULL,
    `status` bigint NOT NULL,
    `weight` DECIMAL(10, 2),
    `sku` text NOT NULL,
    `barcode` text,
    `material` text,
    `origin` text,
    `brand` text,
    `img_url` text,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `UQ_name` (`name`),
    KEY `idx_created_at` (`created_at`),
    KEY `category_id` (`category_id`),
    CONSTRAINT `product_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Variant attribute types (dynamic schema: color, size, form, etc.)
CREATE TABLE IF NOT EXISTS `product_variant` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'e.g., color, size, form, material',
    `display_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'e.g., Color, Size, Fit Style',
    `display_order` int DEFAULT 0 COMMENT 'Order for UI display',
    `is_required` TINYINT DEFAULT 0 COMMENT '1=required for all variants, 0=optional',
    `product_id` bigint NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_display_order` (`display_order`),
    CONSTRAINT `product_variant_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Variant attribute values (e.g., Red, Blue for color; S, M, L for size)
CREATE TABLE IF NOT EXISTS `product_variant_value` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `attribute_id` bigint NOT NULL,
    `value` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'e.g., Red, Blue, S, M, L',
    `display_order` int DEFAULT 0,
    `stock_quantity` int DEFAULT 0 NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `UQ_attribute_value` (`attribute_id`, `value`),
    KEY `idx_attribute_id` (`attribute_id`),
    CONSTRAINT `product__value_ibfk_1` FOREIGN KEY (`attribute_id`) REFERENCES `product_variant` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `price` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `variant_id` bigint NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `compare_at_price` DECIMAL(10, 2) DEFAULT NULL COMMENT 'Original price for discounts',
    `cost_price` DECIMAL(10, 2) DEFAULT NULL COMMENT 'Cost from supplier',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 2=inactive',
    `effective_from` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `effective_to` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_variant_id` (`variant_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `chk_price_dates` CHECK (effective_to IS NULL OR effective_to > effective_from),
    CONSTRAINT `price_ibfk_1` FOREIGN KEY (`variant_id`) REFERENCES `product_variant` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `retail_stores` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone_number` VARCHAR(32) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `platform` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `api_endpoint` TEXT DEFAULT NULL,
    `feature_struct` TEXT DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `payment_methods` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL,
    `code` VARCHAR(20) UNIQUE NOT NULL,
    `description` TEXT,
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `orders` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `payment_status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=unpaid, 2=paid, 3=refunded',
    `customer_id` bigint DEFAULT NULL,
    `platform_id` bigint DEFAULT NULL,
    `payment_id` bigint DEFAULT NULL,
    `retail_stores_id` bigint DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`),
    KEY `customer_id` (`customer_id`),
    CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`),
    CONSTRAINT `platform_ibfk_1` FOREIGN KEY (`platform_id`) REFERENCES `platform` (`id`),
    CONSTRAINT `store_ibfk_1` FOREIGN KEY (`retail_stores_id`) REFERENCES `retail_stores` (`id`),
    CONSTRAINT `payment_ibfk_1` FOREIGN KEY (`payment_id`) REFERENCES `payment_methods` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `order_status` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=pending, 2=paid, 3=shipped, 4=completed, 5=canceled',
    `description` text DEFAULT NULL,
    `order_id` bigint DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `orders_status_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `order_items` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `order_id` bigint NOT NULL,
    `price_id` bigint NOT NULL,
    `variant_value_id` bigint NOT NULL,
    `quantity` int NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`),
    KEY `order_id` (`order_id`),
    KEY `price_id` (`price_id`),
    KEY `variant_value_id` (`variant_value_id`),
    CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
    CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`price_id`) REFERENCES `price` (`id`),
    CONSTRAINT `order_item_pvvfk_2` FOREIGN KEY (`variant_value_id`) REFERENCES `product_variant_value` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Insert default fashion categories
INSERT INTO `category` (`name`) VALUES
    ('Jeans'),
    ('T-Shirt'),
    ('Dress'),
    ('Jacket'),
    ('Shoes'),
    ('Sneakers'),
    ('Hoodie'),
    ('Sweater'),
    ('Shorts'),
    ('Skirt')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `retail_stores` (`name`, `phone_number`)
VALUES
  ('Downtown Store', '0123456789'),
  ('Uptown Store', '0987654321'),
  ('Mall Kiosk', NULL)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `platform` (`id`, `name`, `api_endpoint`, `feature_struct`)
VALUES
  (1, 'shopee', 'https://partner.shopeemobile.com/api/v2/', '{"commission": "2.5% per order"}'),
  (2, 'lazada', 'https://open.lazada.com/api/', '{"commission": "2.0% per order"}'),
  (3, 'tiktok shop', 'https://open-api.tiktokglobalshop.com/', '{"commission": "3.0% per order"}'),
  (4, 'offline', NULL, '{"commission": "0%"}')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `api_endpoint` = VALUES(`api_endpoint`),
  `feature_struct` = VALUES(`feature_struct`);

INSERT INTO `payment_methods` (`id`, `name`, `code`, `description`, `is_active`)
VALUES
  (1, 'Cash', 'CASH', 'Cash payment at store', TRUE),
  (2, 'Bank Transfer', 'BANK', 'Bank transfer payment', TRUE),
  (3, 'ShopeePay', 'SHOPEEPAY', 'ShopeePay e-wallet', TRUE),
  (4, 'Lazada Wallet', 'LAZADAWALLET', 'Lazada Wallet payment', TRUE),
  (5, 'TikTokPay', 'TIKTOKPAY', 'TikTok Shop payment', TRUE),
  (6, 'Credit Card', 'CREDITCARD', 'Credit or debit card payment', TRUE)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `code` = VALUES(`code`),
  `description` = VALUES(`description`),
  `is_active` = VALUES(`is_active`);

INSERT INTO `product` (
  `name`,
  `description`,
  `dimension`,
  `category_id`,
  `status`,
  `weight`,
  `sku`,
  `barcode`,
  `material`,
  `origin`,
  `brand`,
  `img_url`
) VALUES (
  'Classic Unisex T-Shirt',
  'Soft 100% cotton crew neck t-shirt with a relaxed unisex fit. Available in Black and White.',
  '100 x 100',
  (SELECT `id` FROM `category` WHERE `name` = 'T-Shirt'),
  1,
  0.20,
  'TSH-CLSSC-UNX-0001',
  '1234567890123',
  '100% Cotton',
  'Vietnam',
  'SimpleWear',
  'https://example.com/images/tshirt-classic.jpg'
);

INSERT INTO product_variant (name, display_name, display_order, is_required, product_id)
VALUES
  ('tshirt_classic_color', 'Color', 1, 0, 1),
  ('tshirt_classic_size',  'Size',  2, 1, 1),
  ('tshirt_classic_default', 'Default', 0, 1, 1);

  
-- Insert variant values (with stock quantities)
INSERT INTO product_variant_value (attribute_id, value, display_order, stock_quantity)
VALUES
  (1, 'White', 2, 30),
  (1, 'Black', 1, 50);

INSERT INTO product_variant_value (attribute_id, value, display_order, stock_quantity)
VALUES
  (2, 'S', 1, 20),
  (2, 'M', 2, 30),
  (2, 'L', 3, 25);


INSERT INTO price (
  variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to
) VALUES (1, 19.99, 24.99, 8.50, 1, CURRENT_TIMESTAMP, NULL),
(2, 29.99, 14.99, 18.50, 1, CURRENT_TIMESTAMP, NULL),
(3, 49.99, 34.99, 12.50, 1, CURRENT_TIMESTAMP, NULL);

INSERT INTO customer (first_name, last_name, email, address, phone_number)
VALUES
  ('Alice', 'Nguyen', 'alice.nguyen@example.com', '123 Main St, Hanoi', '0901234567'),
  ('Bob', 'Tran', 'bob.tran@example.com', '456 Le Loi, Ho Chi Minh City', '0912345678'),
  ('Carol', 'Pham', 'carol.pham@example.com', '789 Nguyen Hue, Da Nang', '0923456789'),
  ('David', 'Le', 'david.le@example.com', '321 Tran Phu, Can Tho', '0934567890'),
  ('Eve', 'Hoang', 'eve.hoang@example.com', '654 Bach Dang, Hai Phong', '0945678901')
ON DUPLICATE KEY UPDATE
  first_name = VALUES(first_name),
  last_name = VALUES(last_name),
  email = VALUES(email),
  address = VALUES(address);