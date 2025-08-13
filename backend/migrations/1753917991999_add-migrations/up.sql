CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    post_id INTEGER,
    created_at TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    image_url VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted BOOLEAN,
    deleted_at TIMESTAMP
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    post_id INTEGER,
    text TEXT,
    created_at TIMESTAMP,
    deleted BOOLEAN,
    deleted_at TIMESTAMP,
    deleted_by TIMESTAMP
);

CREATE TABLE followers (
    id SERIAL PRIMARY KEY,
    follower_id INTEGER,
    followed_id INTEGER,
    created_at TIMESTAMP
);

CREATE TABLE follow_request (
    id UUID PRIMARY KEY,
    sender_id UUID,
    receiver_id UUID,
    status VARCHAR(8) CHECK (status IN ('pending', 'accepted', 'rejected')),
    UNIQUE (sender_id, receiver_id, status)
);

CREATE TABLE hashtags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    created_at TIMESTAMP,
    deleted BOOLEAN,
    deleted_at TIMESTAMP
);

CREATE TABLE post_tags (
    post_id INTEGER,
    tag_id INTEGER,
    created_at TIMESTAMP,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (tag_id) REFERENCES hashtags(id)
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  full_name VARCHAR(255),
  email VARCHAR(255),
  password VARCHAR(255),
  bio TEXT,
  profile_picture TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN DEFAULT FALSE,
  deleted_at TIMESTAMP
);


CREATE INDEX idx_users_email_deleted ON users(email,deleted);

CREATE INDEX idx_posts_user_id ON posts(user_id);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_post_id ON likes(post_id);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);

CREATE INDEX idx_followers_follower_id ON followers(follower_id);
CREATE INDEX idx_followers_followed_id ON followers(followed_id);

CREATE INDEX idx_follow_request_sender_id ON follow_request(sender_id);
CREATE INDEX idx_follow_request_receiver_id ON follow_request(receiver_id);

CREATE INDEX idx_hashtags_name ON hashtags(name);

CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
