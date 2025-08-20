// Package routes defines the HTTP verbs that the application supports
package routes

import (
	"net/http"

	"github.com/escuadron-404/red404/backend/internal/handlers"
	"github.com/escuadron-404/red404/backend/pkg/middleware"
)

// SetupRoutes configures all application routes.
func SetupRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	mediaHandler *handlers.MediaHandler,
	postHandler *handlers.PostHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	RegisterFrontendHandlers(mux)

	// Register authentication-related routes
	AuthRoutes(mux, authHandler, postHandler, mediaHandler, authMiddleware) // <-- UPDATED: Pass new handlers and middleware

	// Register user-related routes
	UserRoutes(mux, userHandler, authMiddleware)

	return mux
}

// AuthRoutes registers authentication and related routes that might involve new features like posts and media.
func AuthRoutes(
	mux *http.ServeMux,
	authHandler *handlers.AuthHandler,
	postHandler *handlers.PostHandler,
	mediaHandler *handlers.MediaHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/register", authHandler.Register)

	// The post creation handler will use the pre-uploaded media ID from the request body
	mux.HandleFunc("POST /api/post", authMiddleware.Auth(postHandler.CreatePost))

	// New Media Upload Route (protected)
	mux.HandleFunc("POST /api/media", authMiddleware.Auth(mediaHandler.Upload))
}

// UserRoutes remains mostly the same, just keeping it here for completeness
func UserRoutes(mux *http.ServeMux, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	mux.HandleFunc("GET /api/users/{id}", authMiddleware.Auth(userHandler.GetUserByID))
	mux.HandleFunc("GET /api/users", authMiddleware.Auth(userHandler.GetAllUsers))
	// Add other user routes if you have them, e.g., UpdateUser, DeleteUser
	mux.HandleFunc("PUT /api/users/{id}", authMiddleware.Auth(userHandler.UpdateUser))    // Example
	mux.HandleFunc("DELETE /api/users/{id}", authMiddleware.Auth(userHandler.DeleteUser)) // Example
}
