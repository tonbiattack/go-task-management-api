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
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tasks` (
  `id` char(36) NOT NULL COMMENT 'uuid',
  `title` varchar(255) NOT NULL COMMENT 'タスクのタイトル',
  `description` text DEFAULT NULL COMMENT 'タスクの詳細説明',
  `status` varchar(50) NOT NULL DEFAULT 'PENDING' COMMENT 'タスクの状態。デフォルトは「PENDING」(未完了)',
  `created_at` datetime DEFAULT current_timestamp() COMMENT 'タスクの作成日時',
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'タスクの最終更新日時',
  `deadline` datetime DEFAULT NULL COMMENT 'タスク期限日',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='システム内の個々のタスクを表します。タスクは、プロジェクトまたはワークフロー内の特定のアクションや活動を指します。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tasks`
--

LOCK TABLES `tasks` WRITE;
/*!40000 ALTER TABLE `tasks` DISABLE KEYS */;
INSERT INTO `tasks` VALUES ('117ab453-044d-46da-8085-1e686616d7ec','新しいタスク','これは新しいタスクの詳細説明です。','','2023-12-28 22:04:00','2023-12-28 22:04:00','0000-00-00 00:00:00'),('2516e3bf-6f51-4ad2-b9b1-2e8cfc47df23','新しいタスク','これは新しいタスクの詳細説明です。','PENDING','2023-12-28 09:53:37','2023-12-28 09:53:37',NULL),('5f79b7e7-1394-4802-a096-98538298a2da','新しいタスク','これは新しいタスクの詳細説明です。','PENDING','2023-12-26 21:59:08','2023-12-26 21:59:08',NULL),('65be153d-5758-4f6c-8e91-67ed405d4817','新しいタスク','これは新しいタスクの詳細説明です。','','2023-12-28 22:14:48','2023-12-28 22:14:48',NULL),('667804ff-8363-486c-889b-3e5f6bce49fb','新しいタスク','これは新しいタスクの詳細説明です。','PENDING','2023-12-28 09:53:48','2023-12-28 09:53:48',NULL),('6cd5a240-f97c-493f-8196-23c8715013ec','新しいタスク','これは新しいタスクの詳細説明です。','','2023-12-28 22:34:39','2023-12-28 22:34:39','2024-01-01 08:59:59'),('75c8cac8-e202-46d6-b6b8-653ce1289941','新しいタスク','これは新しいタスクの詳細説明です。','PENDING','2023-12-26 21:58:55','2023-12-26 21:58:55',NULL),('7b7789e8-7806-4651-a77d-9e6d10c73353','ワークフロータスク','これは新しいタスクの詳細説明です。','pending','2024-01-13 12:45:50','2024-01-13 12:45:50',NULL),('b9b22431-b7c9-4cd9-a19c-82dc1027ab02','新しいタスク','これは新しいタスクの詳細説明です。','','2023-12-28 22:35:34','2023-12-28 22:35:34','2023-12-28 22:34:39'),('ba9c3787-70e9-48cb-a22e-f5e3598240f9','新しいタスク','これは新しいタスクの詳細説明です。','','2023-12-28 22:42:40','2023-12-28 22:42:40',NULL),('ce0e5d99-8112-400d-aae3-95fbecdca36c','更新したタスク','これは新しいタスクの詳細説明です。','','0000-00-00 00:00:00','2023-12-28 22:05:16','0000-00-00 00:00:00'),('ee4f36f1-b98e-4fc8-8256-8e74a87a6d23','新しいタスク','これは新しいタスクの詳細説明です。','pending','2023-12-28 22:47:23','2023-12-28 22:47:23',NULL);
/*!40000 ALTER TABLE `tasks` ENABLE KEYS */;
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
