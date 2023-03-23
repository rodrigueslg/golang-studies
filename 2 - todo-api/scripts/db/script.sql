CREATE DATABASE `go-challenges`;

CREATE TABLE `go-challenges`.task (
	ID INT auto_increment NOT NULL,
	Name varchar(100) NOT NULL,
	Completed BOOL NULL,
	CONSTRAINT task_PK PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `go-challenges`.task (Name, Completed) VALUES ('Task 1', 0);
INSERT INTO `go-challenges`.task (Name, Completed) VALUES ('Task 2', 1);
INSERT INTO `go-challenges`.task (Name, Completed) VALUES ('Task 3', 0);
INSERT INTO `go-challenges`.task (Name, Completed) VALUES ('Task 4', 1);
