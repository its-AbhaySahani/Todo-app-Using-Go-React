CREATE TABLE users (
  id varchar(36) NOT NULL,
  username varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY username (username)
);

CREATE TABLE todos (
  id varchar(36) NOT NULL,
  task varchar(255) NOT NULL,
  description text,
  done tinyint(1) NOT NULL,
  important tinyint(1) NOT NULL DEFAULT 0,
  user_id varchar(36) DEFAULT NULL,
  date date DEFAULT NULL,
  time time DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE shared_todos (
  id varchar(36) PRIMARY KEY,
  task varchar(255),
  description text,
  done BOOLEAN,
  important BOOLEAN DEFAULT 0,
  user_id varchar(36),
  date DATE,
  time TIME,
  shared_by varchar(36),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (shared_by) REFERENCES users(id)
);

CREATE TABLE teams (
  id varchar(36) NOT NULL,
  name varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  admin_id varchar(36) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (admin_id) REFERENCES users(id)
);

CREATE TABLE team_members (
  team_id varchar(36) NOT NULL,
  user_id varchar(36) NOT NULL,
  is_admin BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (team_id, user_id),
  FOREIGN KEY (team_id) REFERENCES teams(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE team_todos (
  id varchar(36) NOT NULL,
  task varchar(255) NOT NULL,
  description text,
  done BOOLEAN NOT NULL,
  important BOOLEAN DEFAULT FALSE,
  team_id varchar(36) NOT NULL,
  assigned_to varchar(36),
  date DATE DEFAULT NULL,
  time TIME DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (team_id) REFERENCES teams(id),
  FOREIGN KEY (assigned_to) REFERENCES users(id)
);

CREATE TABLE routines (
  id varchar(36) NOT NULL,
  day ENUM('sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday') NOT NULL,
  scheduleType ENUM('morning', 'noon', 'evening', 'night') NOT NULL,
  taskId varchar(36) NOT NULL,
  userId varchar(36) NOT NULL,
  createdAt DATE NOT NULL,
  updatedAt DATE NOT NULL,
  isActive BOOLEAN DEFAULT TRUE,
  PRIMARY KEY (id),
  FOREIGN KEY (taskId) REFERENCES todos(id) ON DELETE CASCADE,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);