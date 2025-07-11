CREATE TABLE IF NOT EXISTS invitations (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `course_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,

  CONSTRAINT fk_invitation_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
  CONSTRAINT fk_invitation_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
