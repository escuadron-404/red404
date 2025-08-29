package models

import (
	"time"
)

type Media struct {
	ID                   string    `json:"id"`
	UserID               int       `json:"user_id"`
	OriginalFilename     string    `json:"original_filename"`
	StoredFilename       string    `json:"stored_filename"`
	FileType             string    `json:"file_type"`
	FileSize             int64     `json:"file_size_bytes"`
	Status               string    `json:"status"`
	TempLocalPath        string    `json:"temp_local_path"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ExpiresAt            time.Time `json:"expires_at"`
	PostID               *int      `json:"post_id"`
	ProfilePictureUserID *int      `json:"profile_picture_user_id"`
}
