-- Insert 5 users
INSERT INTO users (name, email, password) VALUES
('Alice Johnson', 'alice@example.com', '$2a$10$TbKnsnI8pB/KovbdnbvbbOybf1SESd0o8nB7y/iCwkYtoLa2vhjiu'),
('Bob Smith', 'bob@example.com', 'password2'),
('Charlie Brown', 'charlie@example.com', 'password3'),
('David Wilson', 'david@example.com', 'password4'),
('Eva Green', 'eva@example.com', 'password5');

-- Insert 5 posts (one for each user)
INSERT INTO posts (title, excerpt, body, user_id) VALUES
('First Post', 'This is the first post.', 'This is the full content of the first post.', 1),
('Second Post', 'A brief about the second post.', 'The full content of the second post.', 2),
('Third Post', 'A short description for the third.', 'Details of the third post go here.', 3),
('Fourth Post', 'Fourth post summary.', 'Complete content of the fourth post.', 4),
('Fifth Post', 'This is the last post.', 'Here is the body of the fifth post.', 5);
