CREATE TABLE `rooms` (
  `room_id` varchar(255) NOT NULL,
  `room_name` varchar(255) DEFAULT NULL,
  `room_cap` int DEFAULT NULL,
  `room_access` tinyint(1) DEFAULT NULL,
  `room_owner` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;