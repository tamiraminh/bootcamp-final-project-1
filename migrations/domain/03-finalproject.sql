CREATE TABLE `atc_product` (
  `id` varchar(36) PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `description` varchar(255),
  `category` varchar(50) NOT NULL,
  `brand` varchar(50) NOT NULL,
  `stock` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime,
  `updated_by` varchar(36),
  `deleted_at` datetime,
  `deleted_by` varchar(36)
);

CREATE TABLE `atc_cart` (
  `id` varchar(36) PRIMARY KEY NOT NULL,
  `user_id` varchar(36),
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime,
  `updated_by` varchar(36),
  `deleted_at` datetime,
  `deleted_by` varchar(36)
);

CREATE TABLE `atc_cart_item` (
  `cart_id` varchar(36) NOT NULL,
  `product_id` varchar(36) NOT NULL,
  `quantity` int NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime,
  `updated_by` varchar(36),
  `deleted_at` datetime,
  `deleted_by` varchar(36),
  PRIMARY KEY (`cart_id`, `product_id`)
);

CREATE TABLE `atc_order` (
  `id` varchar(36) PRIMARY KEY NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `address` varchar(36),
  `status` ENUM ('pending', 'shipping', 'delivered', 'completed', 'cancelled') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime,
  `updated_by` varchar(36),
  `deleted_at` datetime,
  `deleted_by` varchar(36)
);

CREATE TABLE `atc_order_item` (
  `order_id` varchar(36) NOT NULL,
  `product_id` varchar(36) NOT NULL,
  `quantity` int NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime,
  `updated_by` varchar(36),
  `deleted_at` datetime,
  `deleted_by` varchar(36),
  PRIMARY KEY (`order_id`, `product_id`)
);

CREATE INDEX `atc_product_index_0` ON `atc_product` (`id`) USING BTREE;

ALTER TABLE `atc_cart_item` ADD FOREIGN KEY (`cart_id`) REFERENCES `atc_cart` (`id`);

ALTER TABLE `atc_cart_item` ADD FOREIGN KEY (`product_id`) REFERENCES `atc_product` (`id`);

ALTER TABLE `atc_order_item` ADD FOREIGN KEY (`order_id`) REFERENCES `atc_order` (`id`);

ALTER TABLE `atc_order_item` ADD FOREIGN KEY (`product_id`) REFERENCES `atc_product` (`id`);


INSERT INTO `atc_product` (`id`, `name`, `description`, `category`, `brand`, `stock`, `price`, `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`)
VALUES
(UUID(), 'Smartphone', 'High-end smartphone', 'Electronics', 'Samsung', 50, 999.99, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL),
(UUID(), 'Laptop', 'Powerful laptop', 'Electronics', 'Dell', 30, 1499.99, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL),
(UUID(), 'Headphones', 'Noise-cancelling headphones', 'Audio', 'Sony', 100, 199.99, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL);


INSERT INTO `atc_cart` (`id`, `user_id`, `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`)
VALUES
(UUID(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL);


INSERT INTO `atc_cart_item` (`cart_id`, `product_id`, `quantity`, `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`)
VALUES
('026812a3-357d-11ee-a136-0a0027000018', 'fd235640-357c-11ee-a136-0a0027000018', 2, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL),
('026812a3-357d-11ee-a136-0a0027000018', 'fd23e6c0-357c-11ee-a136-0a0027000018', 1, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL);


INSERT INTO `atc_order` (`id`, `user_id`, `address`, `status`, `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`)
VALUES
(UUID(), 'c192d20b-10c1-4e29-8d86-56af8f774193', '123 Main St', 'pending', NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL);


INSERT INTO `atc_order_item` (`order_id`, `product_id`, `quantity`, `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`)
VALUES
('5c6e3f36-357d-11ee-a136-0a0027000018', 'fd235640-357c-11ee-a136-0a0027000018', 1, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL),
('5c6e3f36-357d-11ee-a136-0a0027000018', 'fd23e6c0-357c-11ee-a136-0a0027000018', 2, NOW(), 'c192d20b-10c1-4e29-8d86-56af8f774193', NULL, NULL, NULL, NULL);



