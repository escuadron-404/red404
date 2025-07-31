DROP INDEX idx_users_email_deleted;
DROP INDEX idx_posts_user_id;
DROP INDEX idx_likes_user_id;
DROP INDEX idx_likes_post_id;
DROP INDEX idx_comments_post_id;
DROP INDEX idx_comments_user_id;
DROP INDEX idx_followers_follower_id;
DROP INDEX idx_followers_followed_id;
DROP INDEX idx_follow_request_sender_id;
DROP INDEX idx_follow_request_receiver_id;
DROP INDEX idx_hashtags_name;
DROP INDEX idx_post_tags_post_id;
DROP INDEX idx_post_tags_tag_id;

-- Remove foreign key constraints from post_tags
ALTER TABLE post_tags DROP CONSTRAINT IF EXISTS post_tags_post_id_fkey;
ALTER TABLE post_tags DROP CONSTRAINT IF EXISTS post_tags_tag_id_fkey;

DROP TABLE posts;
DROP TABLE likes;
DROP TABLE comments;
DROP TABLE followers;
DROP TABLE follow_request;
DROP TABLE hashtags;
DROP TABLE post_tags;