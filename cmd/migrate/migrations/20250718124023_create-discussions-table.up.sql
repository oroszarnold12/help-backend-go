CREATE TABLE IF NOT EXISTS discussions(
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `name` VARCHAR(255),
  `date` DATETIME,
  `creator_id` BIGINT,
  `course_id` BIGINT,

  CONSTRAINT fk_discussion_creator FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_discussion_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE SET NULL
);
