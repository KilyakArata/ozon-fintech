CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    body VARCHAR(2000) NOT NULL,
    post_id BIGINT REFERENCES posts (id) ON DELETE CASCADE,
    parent_comment_id BIGINT REFERENCES comments (id) ON DELETE CASCADE
);