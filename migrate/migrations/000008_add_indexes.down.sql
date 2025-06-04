DROP index IF EXISTS idx_comments_content ON POSTS;
DROP index IF EXISTS dx_posts_title ON posts;
DROP index IF EXISTS idx_posts_tags ON posts;
DROP index IF EXISTS idx_users_username ON users;
DROP index IF EXISTS idx_posts_user_id ON posts;
DROP index IF EXISTS idx_comments_post_id ON comments;