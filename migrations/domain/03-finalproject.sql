-- MySQL dump 10.13  Distrib 8.0.23, for Win64 (x86_64)
--
-- Host: localhost    Database: bootcamp_final_project_product_service
-- ------------------------------------------------------
-- Server version	8.0.23

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `atc_cart`
--

DROP TABLE IF EXISTS `atc_cart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `atc_cart` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `atc_cart`
--

LOCK TABLES `atc_cart` WRITE;
/*!40000 ALTER TABLE `atc_cart` DISABLE KEYS */;
INSERT INTO `atc_cart` VALUES ('5029b6dc-5a01-45cc-b6e8-613c027473c7','c192d20b-10c1-4e29-8d86-56af8f774193','2023-08-09 01:36:12','c192d20b-10c1-4e29-8d86-56af8f774193','2023-08-11 09:49:29','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL);
/*!40000 ALTER TABLE `atc_cart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `atc_cart_item`
--

DROP TABLE IF EXISTS `atc_cart_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `atc_cart_item` (
  `cart_id` varchar(36) NOT NULL,
  `product_id` varchar(36) NOT NULL,
  `quantity` int NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`cart_id`,`product_id`),
  KEY `composite_idx_1` (`cart_id`,`product_id`),
  CONSTRAINT `atc_cart_item_ibfk_1` FOREIGN KEY (`cart_id`) REFERENCES `atc_cart` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `atc_cart_item`
--

LOCK TABLES `atc_cart_item` WRITE;
/*!40000 ALTER TABLE `atc_cart_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `atc_cart_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `atc_order`
--

DROP TABLE IF EXISTS `atc_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `atc_order` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `status` enum('pending','shipping','delivered','completed','cancelled') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `atc_order`
--

LOCK TABLES `atc_order` WRITE;
/*!40000 ALTER TABLE `atc_order` DISABLE KEYS */;
INSERT INTO `atc_order` VALUES ('089730a0-41bd-47e9-a6d9-99005eab5d20','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 10:26:42','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('0af57b3f-936a-49dd-a00e-a05b3763b4cb','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 10:31:23','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('158b8d60-c862-4b32-b802-d0ff32212c49','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 10:25:30','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('2904bc1b-44c1-46af-adec-9a68c354e9d2','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 04:31:22','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('2be794bf-2155-4960-970a-29268c6f0e54','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-10 02:57:17','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('5c6e3f36-357d-11ee-a136-0a0027000018','c192d20b-10c1-4e29-8d86-56af8f774193','123 Main St','pending','2023-08-08 06:51:43','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('6fd4527c-10eb-4014-8f4b-4cae82fb9ae0','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-10 06:25:40','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('709ba641-44d2-4e16-844f-d7e6a8097c74','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-10 03:28:27','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('8b780eef-49c7-4425-ae58-77eea572e423','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-10 02:59:51','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('ac236ae5-5638-4ebe-ac66-e78e50b34380','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-10 03:41:06','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('bb7da091-923f-464a-a9cc-c3813b8bd9b0','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 10:28:46','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('cce4127f-bd51-4d6c-bd1d-3fb11819031e','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 10:25:28','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('efd9908e-234d-4e2c-9380-3a0c54e8e46d','c192d20b-10c1-4e29-8d86-56af8f774193','Bandung','pending','2023-08-11 10:23:20','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `atc_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `atc_order_item`
--

DROP TABLE IF EXISTS `atc_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `atc_order_item` (
  `order_id` varchar(36) NOT NULL,
  `product_id` varchar(36) NOT NULL,
  `quantity` int NOT NULL,
  `unit_price` decimal(10,2) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`order_id`,`product_id`),
  KEY `product_id` (`product_id`),
  KEY `composite_idx_2` (`order_id`,`product_id`),
  CONSTRAINT `atc_order_item_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `atc_order` (`id`),
  CONSTRAINT `atc_order_item_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `atc_product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `atc_order_item`
--

LOCK TABLES `atc_order_item` WRITE;
/*!40000 ALTER TABLE `atc_order_item` DISABLE KEYS */;
INSERT INTO `atc_order_item` VALUES ('0af57b3f-936a-49dd-a00e-a05b3763b4cb','fd23e6c0-357c-11ee-a136-0a0027000018',2,1499.99,'2023-08-11 10:31:23','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('2904bc1b-44c1-46af-adec-9a68c354e9d2','fd235640-357c-11ee-a136-0a0027000018',2,999.99,'2023-08-11 04:31:22','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('2904bc1b-44c1-46af-adec-9a68c354e9d2','fd23e6c0-357c-11ee-a136-0a0027000018',4,1499.99,'2023-08-11 04:31:22','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('5c6e3f36-357d-11ee-a136-0a0027000018','fd235640-357c-11ee-a136-0a0027000018',1,NULL,'2023-08-08 06:55:13','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('5c6e3f36-357d-11ee-a136-0a0027000018','fd23e6c0-357c-11ee-a136-0a0027000018',2,NULL,'2023-08-08 06:55:13','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('6fd4527c-10eb-4014-8f4b-4cae82fb9ae0','fd23e6c0-357c-11ee-a136-0a0027000018',1,1499.99,'2023-08-10 06:25:40','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('6fd4527c-10eb-4014-8f4b-4cae82fb9ae0','fd23f10e-357c-11ee-a136-0a0027000018',1,199.99,'2023-08-10 06:25:40','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('709ba641-44d2-4e16-844f-d7e6a8097c74','379a0303-c15f-4199-a751-e47bc8cb12b9',2,29.00,'2023-08-10 03:28:27','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('8b780eef-49c7-4425-ae58-77eea572e423','379a0303-c15f-4199-a751-e47bc8cb12b9',2,29.00,'2023-08-10 02:59:51','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('ac236ae5-5638-4ebe-ac66-e78e50b34380','379a0303-c15f-4199-a751-e47bc8cb12b9',1,29.00,'2023-08-10 03:41:06','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `atc_order_item` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `after_insert_order_item` AFTER INSERT ON `atc_order_item` FOR EACH ROW BEGIN

    UPDATE atc_product

    SET stock = stock - NEW.quantity

    WHERE id = NEW.product_id;

END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `atc_product`
--

DROP TABLE IF EXISTS `atc_product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `atc_product` (
  `id` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `category` varchar(50) NOT NULL,
  `brand` varchar(50) NOT NULL,
  `stock` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(36) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `atc_product_index_0` (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `atc_product`
--

LOCK TABLES `atc_product` WRITE;
/*!40000 ALTER TABLE `atc_product` DISABLE KEYS */;
INSERT INTO `atc_product` VALUES ('379a0303-c15f-4199-a751-e47bc8cb12b9','vyatta airboor pro','true wireless earphone','gadget','vyatta',25,29.00,'2023-08-08 03:48:29','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('fd235640-357c-11ee-a136-0a0027000018','Smartphone','High-end smartphone','Electronics','Samsung',48,999.99,'2023-08-08 06:49:03','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('fd23e6c0-357c-11ee-a136-0a0027000018','Laptop','Powerful laptop','Electronics','Dell',44,1499.99,'2023-08-08 06:49:03','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL),('fd23f10e-357c-11ee-a136-0a0027000018','Headphones','Noise-cancelling headphones','Audio','Sony',99,199.99,'2023-08-08 06:49:03','c192d20b-10c1-4e29-8d86-56af8f774193',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `atc_product` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-08-11 21:00:54
