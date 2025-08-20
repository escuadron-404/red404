package models

import (
	"time"
)

type Post struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	ImageURL    *string    `json:"image_url"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deleted     bool       `json:"deleted"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
