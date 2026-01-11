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

INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  2, 'Slim Fit Jeans',
  'Comfortable slim fit jeans made from high-quality denim. Perfect for everyday wear.',
  '90 x 60 x 2',
  (SELECT `id` FROM `category` WHERE `name` = 'Jeans'),
  1,
  0.50,
  'JNS-SLM-FT-0002',
  '1234567890124',
  '100% Cotton Denim',
  'Vietnam',
  'DenimCo',
  'https://example.com/images/jeans-slim.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (4, 'jeans_color', 'Color', 1, 0, 2),
  (5, 'jeans_size', 'Size', 2, 1, 2),
  (6, 'jeans_default', 'Default', 0, 1, 2);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (6, 4, 'Blue', 1, 40),
  (7, 4, 'Black', 2, 35);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (8, 5, '28', 1, 15),
  (9, 5, '30', 2, 20),
  (10, 5, '32', 3, 25);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (4, 4, 59.99, 69.99, 25.00, 1, CURRENT_TIMESTAMP, NULL),
  (5, 5, 59.99, 69.99, 25.00, 1, CURRENT_TIMESTAMP, NULL),
  (6, 6, 59.99, 69.99, 25.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 3: Oversized Hoodie
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  3, 'Oversized Hoodie',
  'Warm oversized hoodie with kangaroo pocket and adjustable drawstring hood.',
  '80 x 60 x 5',
  (SELECT `id` FROM `category` WHERE `name` = 'Hoodie'),
  1,
  0.60,
  'HD-OVSZD-0003',
  '1234567890127',
  '80% Cotton, 20% Polyester',
  'Vietnam',
  'UrbanWear',
  'https://example.com/images/hoodie-oversized.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (7, 'hoodie_color', 'Color', 1, 0, 3),
  (8, 'hoodie_size', 'Size', 2, 1, 3);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (11, 7, 'Black', 1, 20),
  (12, 7, 'Gray', 2, 15);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (13, 8, 'M', 1, 10),
  (14, 8, 'L', 2, 12),
  (15, 8, 'XL', 3, 8);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (7, 7, 34.99, 44.99, 18.00, 1, CURRENT_TIMESTAMP, NULL),
  (8, 8, 36.99, 46.99, 19.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 4: Denim Jacket
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  4, 'Denim Jacket',
  'Classic denim jacket with button closure and chest pockets.',
  '85 x 65 x 3',
  (SELECT `id` FROM `category` WHERE `name` = 'Jacket'),
  1,
  0.75,
  'JKT-DNM-0004',
  '1234567890128',
  '100% Cotton Denim',
  'Vietnam',
  'DenimCo',
  'https://example.com/images/jacket-denim.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (9, 'jacket_color', 'Color', 1, 0, 4),
  (10, 'jacket_size', 'Size', 2, 1, 4);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (16, 9, 'Blue', 1, 18),
  (17, 9, 'Black', 2, 14);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (18, 10, 'M', 1, 9),
  (19, 10, 'L', 2, 11),
  (20, 10, 'XL', 3, 7);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (9, 9, 54.99, 64.99, 28.00, 1, CURRENT_TIMESTAMP, NULL),
  (10, 10, 56.99, 66.99, 29.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 5: Pleated Skirt
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  5, 'Pleated Skirt',
  'Elegant pleated skirt with elastic waistband. Flows beautifully for any occasion.',
  '70 x 50 x 2',
  (SELECT `id` FROM `category` WHERE `name` = 'Skirt'),
  1,
  0.25,
  'SKT-PLT-0005',
  '1234567890129',
  'Polyester',
  'Vietnam',
  'ChicStyle',
  'https://example.com/images/skirt-pleated.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (11, 'skirt_color', 'Color', 1, 0, 5),
  (12, 'skirt_size', 'Size', 2, 1, 5);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (21, 11, 'Pink', 1, 13),
  (22, 11, 'Navy', 2, 10);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (23, 12, 'S', 1, 7),
  (24, 12, 'M', 2, 8),
  (25, 12, 'L', 3, 6);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (11, 11, 27.99, 37.99, 12.00, 1, CURRENT_TIMESTAMP, NULL),
  (12, 12, 29.99, 39.99, 13.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 6: Canvas Sneakers
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  6, 'Canvas Sneakers',
  'Classic canvas sneakers with rubber sole, perfect for casual wear.',
  '32 x 18 x 12',
  (SELECT `id` FROM `category` WHERE `name` = 'Sneakers'),
  1,
  0.65,
  'SNK-CNVS-0006',
  '1234567890130',
  'Canvas, Rubber',
  'Vietnam',
  'StepUp',
  'https://example.com/images/sneakers-canvas.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (13, 'canvas_color', 'Color', 1, 0, 6),
  (14, 'canvas_size', 'Size', 2, 1, 6);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (26, 13, 'White', 1, 20),
  (27, 13, 'Black', 2, 18);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (28, 14, '38', 1, 10),
  (29, 14, '39', 2, 12),
  (30, 14, '40', 3, 8);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (13, 13, 39.99, 49.99, 18.00, 1, CURRENT_TIMESTAMP, NULL),
  (14, 14, 41.99, 51.99, 19.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 7: Cropped Sweater
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  7, 'Cropped Sweater',
  'Trendy cropped sweater with ribbed hem and cuffs.',
  '60 x 45 x 2',
  (SELECT `id` FROM `category` WHERE `name` = 'Sweater'),
  1,
  0.35,
  'SWT-CRPD-0007',
  '1234567890131',
  'Acrylic Blend',
  'Vietnam',
  'CozyWear',
  'https://example.com/images/sweater-cropped.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (15, 'cropped_color', 'Color', 1, 0, 7),
  (16, 'cropped_size', 'Size', 2, 1, 7);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (31, 15, 'Pink', 1, 12),
  (32, 15, 'Beige', 2, 10);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (33, 16, 'S', 1, 7),
  (34, 16, 'M', 2, 8),
  (35, 16, 'L', 3, 6);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (15, 15, 29.99, 39.99, 13.00, 1, CURRENT_TIMESTAMP, NULL),
  (16, 16, 31.99, 41.99, 14.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 8: Summer Shorts
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  8, 'Summer Shorts',
  'Lightweight shorts for summer activities. Elastic waistband and side pockets.',
  '50 x 40 x 2',
  (SELECT `id` FROM `category` WHERE `name` = 'Shorts'),
  1,
  0.18,
  'SRT-SMMR-0008',
  '1234567890132',
  'Cotton Blend',
  'Vietnam',
  'ActiveLife',
  'https://example.com/images/shorts-summer.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (17, 'shorts_color', 'Color', 1, 0, 8),
  (18, 'shorts_size', 'Size', 2, 1, 8);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (36, 17, 'Khaki', 1, 9),
  (37, 17, 'Navy', 2, 11);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (38, 18, 'M', 1, 6),
  (39, 18, 'L', 2, 8),
  (40, 18, 'XL', 3, 5);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (17, 17, 24.99, 34.99, 11.00, 1, CURRENT_TIMESTAMP, NULL),
  (18, 18, 26.99, 36.99, 12.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 9: Midi Dress
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  9, 'Midi Dress',
  'Elegant midi dress with a fitted waist and flowing skirt.',
  '110 x 70',
  (SELECT `id` FROM `category` WHERE `name` = 'Dress'),
  1,
  0.28,
  'DRS-MIDI-0009',
  '1234567890133',
  'Rayon',
  'Vietnam',
  'SummerStyle',
  'https://example.com/images/dress-midi.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (19, 'midi_color', 'Color', 1, 0, 9),
  (20, 'midi_size', 'Size', 2, 1, 9);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (41, 19, 'Red', 1, 8),
  (42, 19, 'Blue', 2, 7);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (43, 20, 'S', 1, 4),
  (44, 20, 'M', 2, 5),
  (45, 20, 'L', 3, 3);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (19, 19, 49.99, 59.99, 22.00, 1, CURRENT_TIMESTAMP, NULL),
  (20, 20, 51.99, 61.99, 23.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 10: Leather Shoes
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  10, 'Leather Shoes',
  'Premium leather shoes for formal occasions. Durable and stylish.',
  '30 x 15 x 10',
  (SELECT `id` FROM `category` WHERE `name` = 'Shoes'),
  1,
  0.80,
  'SHO-LTHR-0010',
  '1234567890134',
  'Genuine Leather',
  'Vietnam',
  'ClassicStep',
  'https://example.com/images/shoes-leather.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (21, 'leather_color', 'Color', 1, 0, 10),
  (22, 'leather_size', 'Size', 2, 1, 10);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (46, 21, 'Brown', 1, 6),
  (47, 21, 'Black', 2, 8);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (48, 22, '41', 1, 3),
  (49, 22, '42', 2, 4),
  (50, 22, '43', 3, 2);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (21, 21, 89.99, 109.99, 45.00, 1, CURRENT_TIMESTAMP, NULL),
  (22, 22, 91.99, 111.99, 46.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 11: Bomber Jacket
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  11, 'Bomber Jacket',
  'Classic bomber jacket with ribbed collar and cuffs. Perfect for cool weather.',
  '80 x 60 x 4',
  (SELECT `id` FROM `category` WHERE `name` = 'Jacket'),
  1,
  0.68,
  'JKT-BMBR-0011',
  '1234567890135',
  'Polyester, Nylon',
  'Vietnam',
  'UrbanWear',
  'https://example.com/images/jacket-bomber.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (23, 'bomber_color', 'Color', 1, 0, 11),
  (24, 'bomber_size', 'Size', 2, 1, 11);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (51, 23, 'Olive', 1, 11),
  (52, 23, 'Black', 2, 14);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (53, 24, 'M', 1, 6),
  (54, 24, 'L', 2, 9),
  (55, 24, 'XL', 3, 5);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (23, 23, 64.99, 79.99, 32.00, 1, CURRENT_TIMESTAMP, NULL),
  (24, 24, 66.99, 81.99, 33.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 12: Ripped Jeans
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  12, 'Ripped Jeans',
  'Trendy ripped jeans with distressed detailing. Slim fit design.',
  '92 x 62 x 2',
  (SELECT `id` FROM `category` WHERE `name` = 'Jeans'),
  1,
  0.55,
  'JNS-RPPD-0012',
  '1234567890136',
  'Cotton Denim',
  'Vietnam',
  'DenimCo',
  'https://example.com/images/jeans-ripped.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (25, 'ripped_color', 'Color', 1, 0, 12),
  (26, 'ripped_size', 'Size', 2, 1, 12);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (56, 25, 'Light Blue', 1, 17),
  (57, 25, 'Black', 2, 13);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (58, 26, '29', 1, 8),
  (59, 26, '30', 2, 10),
  (60, 26, '32', 3, 7);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (25, 25, 54.99, 69.99, 26.00, 1, CURRENT_TIMESTAMP, NULL),
  (26, 26, 56.99, 71.99, 27.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 13: V-Neck T-Shirt
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  13, 'V-Neck T-Shirt',
  'Casual v-neck t-shirt with a relaxed fit. Soft and breathable fabric.',
  '95 x 95',
  (SELECT `id` FROM `category` WHERE `name` = 'T-Shirt'),
  1,
  0.19,
  'TSH-VNCK-0013',
  '1234567890137',
  'Cotton Jersey',
  'Vietnam',
  'SimpleWear',
  'https://example.com/images/tshirt-vneck.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (27, 'vneck_color', 'Color', 1, 0, 13),
  (28, 'vneck_size', 'Size', 2, 1, 13);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (61, 27, 'White', 1, 22),
  (62, 27, 'Gray', 2, 19);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (63, 28, 'S', 1, 11),
  (64, 28, 'M', 2, 15),
  (65, 28, 'L', 3, 12);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (27, 27, 18.99, 24.99, 8.00, 1, CURRENT_TIMESTAMP, NULL),
  (28, 28, 20.99, 26.99, 9.00, 1, CURRENT_TIMESTAMP, NULL);

-- Product 14: Mini Skirt
INSERT INTO `product` (
  `id`, `name`, `description`, `dimension`, `category_id`, `status`, `weight`, `sku`, `barcode`, `material`, `origin`, `brand`, `img_url`
) VALUES (
  14, 'Mini Skirt',
  'Trendy mini skirt with A-line silhouette. Perfect for casual outings.',
  '60 x 40 x 1',
  (SELECT `id` FROM `category` WHERE `name` = 'Skirt'),
  1,
  0.15,
  'SKT-MINI-0014',
  '1234567890138',
  'Cotton Blend',
  'Vietnam',
  'ChicStyle',
  'https://example.com/images/skirt-mini.jpg'
);

INSERT INTO product_variant (id, name, display_name, display_order, is_required, product_id)
VALUES
  (29, 'mini_color', 'Color', 1, 0, 14),
  (30, 'mini_size', 'Size', 2, 1, 14);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (66, 29, 'Black', 1, 14),
  (67, 29, 'Red', 2, 9);

INSERT INTO product_variant_value (id, attribute_id, value, display_order, stock_quantity)
VALUES
  (68, 30, 'S', 1, 8),
  (69, 30, 'M', 2, 10),
  (70, 30, 'L', 3, 5);

INSERT INTO price (id, variant_id, price, compare_at_price, cost_price, status, effective_from, effective_to)
VALUES 
  (29, 29, 22.99, 32.99, 10.00, 1, CURRENT_TIMESTAMP, NULL),
  (30, 30, 24.99, 34.99, 11.00, 1, CURRENT_TIMESTAMP, NULL);
