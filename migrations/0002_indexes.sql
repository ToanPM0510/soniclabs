CREATE INDEX IF NOT EXISTS idx_courses_diff_created_id
  ON courses (difficulty, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_courses_diff_created_id
  ON courses (created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_enrollments_email_enrolled_id
  ON enrollments (student_email, enrolled_at DESC, id DESC);

CREATE UNIQUE INDEX IF NOT EXISTS idx_enrollments_unique_email_course
  ON enrollments (lower(student_email), course_id);
  
CREATE INDEX IF NOT EXISTS idx_outbox_unprocessed
  ON outbox (processed_at) WHERE processed_at IS NULL;
