USE Todo_app;

CREATE TABLE `todos` (
  `id` varchar(36) NOT NULL,
  `task` varchar(255) NOT NULL,
  `done` tinyint(1) NOT NULL,
  `user_id` varchar(36) DEFAULT NULL,
  `date` date DEFAULT NULL,
  `time` time DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci



CREATE TABLE shared_todos (
    id VARCHAR(36) PRIMARY KEY,
    task VARCHAR(255),
    done BOOLEAN,
    user_id VARCHAR(36),
    date DATE,
    time TIME,
    shared_by VARCHAR(36),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (shared_by) REFERENCES users(id)
);