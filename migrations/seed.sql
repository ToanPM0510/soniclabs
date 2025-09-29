INSERT INTO courses (id, title, description, difficulty)
VALUES
  (gen_random_uuid(), 'Go Basics', 'Introduction to Go programming language', 'Beginner'),
  (gen_random_uuid(), 'RESTful APIs with Gin', 'Build robust REST APIs using Gin', 'Intermediate'),
  (gen_random_uuid(), 'PostgreSQL for Developers', 'Indexes, transactions, and query plans', 'Intermediate'),
  (gen_random_uuid(), 'High-Performance Go', 'Concurrency, memory, and profiling', 'Advanced')
ON CONFLICT DO NOTHING;
