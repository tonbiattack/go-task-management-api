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
-- Table structure for table `task_workflow_statuses`
--

DROP TABLE IF EXISTS `task_workflow_statuses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_workflow_statuses` (
  `id` char(36) NOT NULL COMMENT 'uuid',
  `task_id` char(36) NOT NULL COMMENT 'タスクID',
  `workflow_step_id` char(36) NOT NULL COMMENT 'ワークフローステップID',
  `status` varchar(50) NOT NULL COMMENT 'ステータス',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '作成日時',
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新日時',
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`),
  KEY `workflow_step_id` (`workflow_step_id`),
  CONSTRAINT `task_workflow_statuses_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`),
  CONSTRAINT `task_workflow_statuses_ibfk_2` FOREIGN KEY (`workflow_step_id`) REFERENCES `workflow_steps` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='タスクとワークフローステップとの関連を表し、タスクのワークフローにおける現在の状態を追跡します';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_workflow_statuses`
--

LOCK TABLES `task_workflow_statuses` WRITE;
/*!40000 ALTER TABLE `task_workflow_statuses` DISABLE KEYS */;
INSERT INTO `task_workflow_statuses` VALUES ('5eb5c4ce-bdce-4838-95ac-b32cb48c45fb','7b7789e8-7806-4651-a77d-9e6d10c73353','57ecd3f0-addc-4df7-b537-199db062fe63','COMPLETED','2024-01-13 13:50:53','2024-01-13 15:10:18');
/*!40000 ALTER TABLE `task_workflow_statuses` ENABLE KEYS */;
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
