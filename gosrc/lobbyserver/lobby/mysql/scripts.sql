-- MySQL dump 10.13  Distrib 8.0.13, for Win64 (x86_64)
--
-- Host: localhost    Database: game
-- ------------------------------------------------------
-- Server version	5.7.17-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `account`
--

DROP TABLE IF EXISTS `account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `account` (
  `account` varchar(64) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `password` varchar(64) DEFAULT NULL,
  `register_time` datetime DEFAULT NULL,
  PRIMARY KEY (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='账号表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `account`
--

LOCK TABLES `account` WRITE;
/*!40000 ALTER TABLE `account` DISABLE KEYS */;
INSERT INTO `account` VALUES ('',10000004,NULL,NULL),('0a5e82dc-d8ef-4103-9200-3cd3123f040a',10000009,NULL,NULL),('0f2177e2-f767-4da7-a50c-fe251b7f21c0',10000010,NULL,NULL),('123',10000003,NULL,NULL),('20908b78-1c1d-4a1f-843a-1f1334cc1e7f',10000019,'',NULL),('2cab2050-db31-4ba9-9f72-e47ee4cbfe17',10000016,'',NULL),('3493fcaa-9365-4755-aa39-7413e70b1193',10000011,NULL,NULL),('3823da3e-feb6-411c-b202-25b21ef0ba72',10000020,'',NULL),('527b61d4-1c9e-4dce-8f3f-4059f90d9461',10000005,NULL,NULL),('56161bf0-946f-4a35-b397-d432b97247b8',10000006,NULL,NULL),('810c90b1-7b5e-4f25-85b1-ae92079b3457',10000015,'',NULL),('910d3e32-8234-4c6f-9f88-bd75ef0c440e',10000014,'',NULL),('925a07c2-fd6e-4a58-bd8d-bf53a913dc3b',10000008,NULL,NULL),('9d11cbd1-6fab-41e3-a0f3-e43bca2feeca',10000013,'',NULL),('d842456f-3bbd-403b-952e-ff3fa0a80cea',10000017,'',NULL),('db860726-ea20-410c-8e2d-85e6ffcecdf4',10000007,NULL,NULL),('ea077104-ed69-40c2-bfdd-1e3dac0f4c83',10000012,NULL,NULL),('f8508a24-00c7-43d0-9061-342e9f13b1da',10000018,'',NULL);
/*!40000 ALTER TABLE `account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `diamond`
--

DROP TABLE IF EXISTS `diamond`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `diamond` (
  `user_id` int(11) NOT NULL,
  `num` int(11) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `diamond`
--

LOCK TABLES `diamond` WRITE;
/*!40000 ALTER TABLE `diamond` DISABLE KEYS */;
INSERT INTO `diamond` VALUES (123,0);
/*!40000 ALTER TABLE `diamond` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `diamond_log`
--

DROP TABLE IF EXISTS `diamond_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `diamond_log` (
  `user_id` int(11) NOT NULL,
  `change_num` int(11) DEFAULT NULL,
  `current_num` int(11) DEFAULT NULL,
  `descript` varchar(128) DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='钻石改变日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `diamond_log`
--

LOCK TABLES `diamond_log` WRITE;
/*!40000 ALTER TABLE `diamond_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `diamond_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `test`
--

DROP TABLE IF EXISTS `test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `test` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(45) NOT NULL,
  `phone` varchar(45) DEFAULT NULL,
  `open_id` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`,`user_id`),
  UNIQUE KEY `user_id_UNIQUE` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test`
--

LOCK TABLES `test` WRITE;
/*!40000 ALTER TABLE `test` DISABLE KEYS */;
INSERT INTO `test` VALUES (6,'1','2220',NULL),(7,'3','2220',NULL),(8,'4','2220',NULL);
/*!40000 ALTER TABLE `test` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `test1`
--

DROP TABLE IF EXISTS `test1`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `test1` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test1`
--

LOCK TABLES `test1` WRITE;
/*!40000 ALTER TABLE `test1` DISABLE KEYS */;
INSERT INTO `test1` VALUES (1,'1'),(2,'22');
/*!40000 ALTER TABLE `test1` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `open_id` varchar(45) DEFAULT NULL,
  `phone` varchar(11) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `mod_name` varchar(64) DEFAULT NULL,
  `mod_version` varchar(32) DEFAULT NULL,
  `core_version` varchar(32) DEFAULT NULL,
  `lobby_version` varchar(32) DEFAULT NULL,
  `operating_system` varchar(32) DEFAULT NULL,
  `system_family` varchar(45) DEFAULT NULL,
  `device_id` varchar(45) DEFAULT NULL,
  `device_name` varchar(64) DEFAULT NULL,
  `device_mode` varchar(45) DEFAULT NULL,
  `network_type` varchar(45) DEFAULT NULL,
  `register_type` varchar(16) DEFAULT NULL,
  `nick_name` varchar(64) DEFAULT NULL,
  `sex` int(1) DEFAULT NULL,
  `provice` varchar(32) DEFAULT NULL,
  `city` varchar(32) DEFAULT NULL,
  `country` varchar(32) DEFAULT NULL,
  `head_img_url` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`,`user_id`),
  UNIQUE KEY `user_id_UNIQUE` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,10000001,'1001',NULL,'riguang','login','1.0.0','1.0.0','1.0.0','ios','c#','00.1','IPHONE','IPHONE','4G',NULL,'abc',1,'guangDong','shenzhen','nanshan','http://baidu.com'),(2,10000002,'',NULL,'','loginMode','1.0.1','1.0.0','1.0.0','IOS','iphone5','111222','chen_phone','phone','4G',NULL,'',0,'','','',''),(3,10000003,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(4,10000004,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(5,10000005,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(6,10000006,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(7,10000007,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(8,10000008,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(9,10000009,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(10,10000010,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(11,10000011,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(12,10000012,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(13,10000013,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(14,10000014,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(15,10000015,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(16,10000016,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(17,10000017,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(18,10000018,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(19,10000019,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL),(20,10000020,NULL,'',NULL,'','','','','','','','','','',NULL,NULL,NULL,NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-04-24 12:08:23
