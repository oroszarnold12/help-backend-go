CREATE TABLE IF NOT EXISTS quizzes(
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `name` VARCHAR(255),
  `due_date` DATETIME,
  `points` DOUBLE,
  `published` BOOLEAN,
  `course_id` BIGINT,

  CONSTRAINT fk_quiz_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE SET NULL
);
