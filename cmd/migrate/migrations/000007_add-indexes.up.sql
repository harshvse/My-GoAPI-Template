CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);

CREATE INDEX idx_posts_title ON posts using gin (title gin_trgm_ops);
CREATE INDEX idx_posts_tags ON posts using gin (tags);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_post_user_id ON posts (user_id);
CREATE INDEX idx_comments_posts_id on comments (post_id);