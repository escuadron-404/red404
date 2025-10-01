CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; 

CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    user_id INTEGER NOT NULL REFERENCES users(id), 
    original_filename VARCHAR(255) NOT NULL, 
    stored_filename VARCHAR(255) NOT NULL UNIQUE, 
    file_type VARCHAR(50) NOT NULL, 
    file_size_bytes BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'uploaded', 'used', 'failed', 'deleted')),
    temp_local_path TEXT, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE, 
    post_id INTEGER REFERENCES posts(id) DEFAULT NULL, 
    profile_picture_user_id INTEGER REFERENCES users(id) DEFAULT NULL 
);

CREATE INDEX idx_media_user_id ON media(user_id);
CREATE INDEX idx_media_status ON media(status);
CREATE INDEX idx_media_expires_at ON media(expires_at);
CREATE UNIQUE INDEX idx_media_stored_filename ON media(stored_filename);

ALTER TABLE posts ALTER COLUMN image_url TYPE VARCHAR(512);

ALTER TABLE users ALTER COLUMN profile_picture TYPE VARCHAR(512);
