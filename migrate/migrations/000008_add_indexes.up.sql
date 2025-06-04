CREATE extension if not exists pg_trgm;

CREATE index IF NOT EXISTS idx_comments_content on POSTS USING gin(content gin_trgm_ops);
CREATE index IF NOT EXISTS dx_posts_title ON posts USING gin(title gin_trgm_ops);
CREATE index IF NOT EXISTS idx_posts_tags ON posts USING gin(tags);
CREATE index IF NOT EXISTS idx_users_username ON users(username);
CREATE index IF NOT EXISTS idx_posts_user_id on posts(user_id);
CREATE index IF NOT EXISTS idx_comments_post_id ON comments(post_id);