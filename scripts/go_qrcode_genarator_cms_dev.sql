-- MySQL dump 10.13  Distrib 8.0.34, for Win64 (x86_64)
--
-- Host: localhost    Database: go_qrcode_genarator_cms_dev
-- ------------------------------------------------------
-- Server version	8.0.34

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `qrcodes`
--

DROP TABLE IF EXISTS `qrcodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `qrcodes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `uuid` longtext NOT NULL,
  `content` longtext NOT NULL,
  `type` varchar(191) DEFAULT 'text',
  `background` varchar(191) DEFAULT '#FFFFFF',
  `foreground` varchar(191) DEFAULT '#000000',
  `border_width` bigint DEFAULT '20',
  `circle_shape` tinyint(1) NOT NULL DEFAULT '0',
  `transparent_background` tinyint(1) NOT NULL DEFAULT '0',
  `version` bigint DEFAULT '2',
  `error_level` bigint DEFAULT '2',
  `public_url` longtext NOT NULL,
  `encode_content` longtext NOT NULL,
  `file_path` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_qrcodes_deleted_at` (`deleted_at`),
  KEY `fk_users_qr_codes` (`user_id`),
  CONSTRAINT `fk_users_qr_codes` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `qrcodes`
--

LOCK TABLES `qrcodes` WRITE;
/*!40000 ALTER TABLE `qrcodes` DISABLE KEYS */;
/*!40000 ALTER TABLE `qrcodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `uuid` longtext NOT NULL,
  `name` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'2023-11-07 16:53:03.214','2023-11-07 16:53:03.214',NULL,'dfdbd98a-d9c3-46ed-ae0e-e0aae0de8cad','Administrator'),(2,'2023-11-07 16:53:13.963','2023-11-07 16:53:13.963',NULL,'d43b2a6e-9943-456d-b0db-d4d004fb5f44','Normal User'),(3,'2023-11-07 16:53:22.572','2023-11-07 16:53:22.572',NULL,'fd27487f-bb57-44d5-a8e7-748f1e3d7c4f','Payment User');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT NULL,
  `uuid` longtext NOT NULL,
  `first_name` longtext NOT NULL,
  `last_name` longtext NOT NULL,
  `gender` tinyint(1) DEFAULT '0',
  `email` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `activate` tinyint(1) DEFAULT '0',
  `activation_code` longtext NOT NULL,
  `avatar_url` longtext NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `fk_roles_user` (`role_id`),
  CONSTRAINT `fk_roles_user` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2023-11-07 16:54:03.436','2023-11-07 16:54:03.436',NULL,2,'461351ba-f393-427a-b327-be30949f9a25','Lê','Xuân Hòa',0,'dso.intern.xuanhoa@gmail.com','$2a$05$xmhIvt6FQdBRXoGRDvyn8OW9IGa.ajy2XyaLMpCodgO4iVU6AFMQi',1,'183a203e-0117-4597-8508-a7405eea9bfb','https://lh3.googleusercontent.com/a/ACg8ocI6SpcNl5mcPoaKeujESCctTfj6w-agbQXqYvZQf8y-TQ=s96-c');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-11-07 16:56:24
