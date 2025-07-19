CREATE TABLE IF NOT EXISTS quiz_grades (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `uuid` VARCHAR(36) NOT NULL UNIQUE,
  `submitter_id` BIGINT,
  `quiz_id` BIGINT,
  `grade` DOUBLE,

  CONSTRAINT fk_quiz_grade_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE SET NULL,
  CONSTRAINT fk_quiz_grade_submitter FOREIGN KEY (submitter_id) REFERENCES users(id) ON DELETE SET NULL
);
