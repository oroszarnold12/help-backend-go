CREATE TABLE IF NOT EXISTS announcements (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `name` VARCHAR(255),
  `date` DATETIME,
  `content` TEXT,
  `course_id` BIGINT,
  `creator_id` BIGINT,

  CONSTRAINT fk_announcement_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE SET NULL,
  CONSTRAINT fk_announcement_creator FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL
);
