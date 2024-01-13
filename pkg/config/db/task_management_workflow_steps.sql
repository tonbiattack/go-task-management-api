-- MySQL dump 10.13  Distrib 8.0.34, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: task_management
-- ------------------------------------------------------
-- Server version	5.5.5-10.4.28-MariaDB

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
-- Table structure for table `workflow_steps`
--

DROP TABLE IF EXISTS `workflow_steps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workflow_steps` (
  `id` char(36) NOT NULL COMMENT 'uuid',
  `workflow_id` char(36) NOT NULL COMMENT 'ワークフローID',
  `name` varchar(255) NOT NULL COMMENT 'ステップの名前',
  `description` text DEFAULT NULL COMMENT 'ステップの説明',
  `order` int(11) NOT NULL COMMENT 'ステップの順序',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '作成日時',
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新日時',
  PRIMARY KEY (`id`),
  KEY `workflow_id` (`workflow_id`),
  CONSTRAINT `workflow_steps_ibfk_1` FOREIGN KEY (`workflow_id`) REFERENCES `workflows` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='このテーブルは、ワークフロー内の個々のステップを表します。各ステップは、ワークフローの特定の段階を指します。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workflow_steps`
--

LOCK TABLES `workflow_steps` WRITE;
/*!40000 ALTER TABLE `workflow_steps` DISABLE KEYS */;
INSERT INTO `workflow_steps` VALUES ('09d9d3f8-e0f4-4f47-87cd-31105772d33e','83060a3a-8708-430c-8472-d370ef707a30','企画','プロジェクトの初期企画段階',1,'2024-01-13 13:35:26','2024-01-13 13:35:26'),('57ecd3f0-addc-4df7-b537-199db062fe63','83060a3a-8708-430c-8472-d370ef707a30','完了','',4,'2024-01-13 13:36:45','2024-01-13 13:36:45'),('5b119c7e-22a6-4e9a-b041-0236b9acd40f','83060a3a-8708-430c-8472-d370ef707a30','実行','実際に実行',2,'2024-01-13 13:36:06','2024-01-13 13:36:06'),('a515cc51-1a92-4ab1-81c1-e67b7063fda9','83060a3a-8708-430c-8472-d370ef707a30','レビュー','上長レビュー',3,'2024-01-13 13:36:32','2024-01-13 13:36:32');
/*!40000 ALTER TABLE `workflow_steps` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-13 15:21:07
