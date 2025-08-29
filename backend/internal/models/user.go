// Package models defines the shape of data in memory
package models

import (
	"time"
)

type User struct {
	ID             int        `json:"id"`
	FullName       *string    `json:"full_name"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	Bio            *string    `json:"bio"`
	ProfilePicture *string    `json:"profile_picture"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Deleted        bool       `json:"deleted"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type UserWithoutPassword struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
