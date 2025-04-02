-- +migrate up
CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       title VARCHAR(200) NOT NULL,
                       content TEXT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_posts_user_id ON posts(user_id);

-- +migrate down
DROP TABLE IF EXISTS posts;
