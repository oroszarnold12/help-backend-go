CREATE TABLE IF NOT EXISTS assignment_grades (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `submitter_id` BIGINT,
  `assignment_id` BIGINT,
  `grade` DOUBLE,

  CONSTRAINT fk_assignment_grade_assignment FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE SET NULL,
  CONSTRAINT fk_assignment_grade_submitter FOREIGN KEY (submitter_id) REFERENCES users(id) ON DELETE SET NULL
);
