package dto

import "time"

// PostCreateRequest represents the data required to create a new post.
type PostCreateRequest struct {
	Description string `json:"description" validate:"required,min=1,max=1000"`
	MediaID     string `json:"media_id,omitempty" validate:"uuid"`
}

// PostResponse represents the data returned after a post is created or retrieved.
type PostResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ImageURL    string    `json:"image_url,omitempty"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
