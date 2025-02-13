USE Todo_app;

-- Drop existing tables if they exist
DROP TABLE IF EXISTS shared_todos;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create todos table
CREATE TABLE `todos` (
  `id` varchar(36) NOT NULL,
  `task` varchar(255) NOT NULL,
  `description` text,
  `done` tinyint(1) NOT NULL,
  `important` tinyint(1) NOT NULL DEFAULT 0,
  `user_id` varchar(36) DEFAULT NULL,
  `date` date DEFAULT NULL,
  `time` time DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create shared_todos table
CREATE TABLE `shared_todos` (
  `id` varchar(36) PRIMARY KEY,
  `task` varchar(255),
  `description` text,
  `done` BOOLEAN,
  `important` BOOLEAN DEFAULT 0,
  `user_id` varchar(36),
  `date` DATE,
  `time` TIME,
  `shared_by` varchar(36),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
  FOREIGN KEY (`shared_by`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;