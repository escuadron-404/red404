package dto

import (
	"time"
)

type MediaUploadRequest struct {
	// TODO: nani wo suru?
}

type MediaUploadResponse struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	FileType         string    `json:"file_type"`
	Size             int64     `json:"size"`
	FileURL          string    `json:"file_url"`
	ExpiresAt        time.Time `json:"expires_at,omitempty"`
	Message          string    `json:"message"`
}

type MediaInfo struct {
	ID                   string    `db:"id"`
	UserID               int       `db:"user_id"`
	OriginalFilename     string    `db:"original_filename"`
	StoredFilename       string    `db:"stored_filename"`
	FileType             string    `db:"file_type"`
	FileSize             int64     `db:"file_size_bytes"`
	Status               string    `db:"status"`
	TempLocalPath        string    `db:"temp_local_path"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
	ExpiresAt            time.Time `db:"expires_at"`
	PostID               *int      `db:"post_id"`
	ProfilePictureUserID *int      `db:"profile_picture_user_id"`
}
