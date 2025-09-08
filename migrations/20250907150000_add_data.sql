-- +migrate Up
-- Add your up migration here

INSERT INTO users (name, email) VALUES
  ('John Doe', 'john@example.com'),
  ('Jane Smith', 'jane@example.com'),
  ('Bob Johnson', 'bob@example.com');




-- +migrate Down
-- Add your down migration here

DELETE FROM users;