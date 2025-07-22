CREATE TABLE IF NOT EXISTS course_files(
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `name` VARCHAR(255),
  `size` BIGINT,
  `creation_date` DATETIME,
  `uploader_id` BIGINT,
  `course_id` BIGINT,

  CONSTRAINT fk_course_file_uploader FOREIGN KEY (uploader_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_course_file_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE SET NULL
);
