package middleware

import (
    "github.com/rs/cors"
)

// NewCORS returns a configured CORS middleware
func NewCORS() *cors.Cors {
    return cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"}, // your frontend dev URL
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })
}

