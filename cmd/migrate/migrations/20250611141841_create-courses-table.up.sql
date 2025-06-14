CREATE TABLE IF NOT EXISTS courses (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `name` VARCHAR(255),
  `long_name` VARCHAR(255),
  `description` TEXT,
  `teacher_id` BIGINT,

  CONSTRAINT fk_course_teacher FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE SET NULL
);
