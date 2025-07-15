CREATE TABLE IF NOT EXISTS assignments(
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `name` VARCHAR(255),
  `due_date` DATETIME,
  `points` INT,
  `published` BOOLEAN,
  `course_id` BIGINT,

  CONSTRAINT fk_assignment_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE SET NULL
);
