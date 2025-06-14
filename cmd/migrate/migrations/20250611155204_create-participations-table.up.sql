CREATE TABLE IF NOT EXISTS participations (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `course_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `show_on_dashboard` BOOLEAN NOT NULL,

  CONSTRAINT fk_participation_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
  CONSTRAINT fk_participation_person FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
