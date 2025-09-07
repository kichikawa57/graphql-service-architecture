-- +migrate Up
-- Add your up migration here

INSERT INTO posts (user_id, title, content) VALUES
  (1, 'Post 1', 'Content 1'),
  (2, 'Post 2', 'Content 2'),
  (3, 'Post 3', 'Content 3'),
  (1, 'Post 4', 'Content 4'),
  (2, 'Post 5', 'Content 5'),
  (3, 'Post 6', 'Content 6'),
  (1, 'Post 7', 'Content 7'),
  (2, 'Post 8', 'Content 8'),
  (3, 'Post 9', 'Content 9');




-- +migrate Down
-- Add your down migration here

DELETE FROM posts;