package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/services"
	"github.com/escuadron-404/red404/backend/pkg/common"
	"github.com/go-playground/validator/v10"
	"strings"
)

// Add new handler struct for posts
type PostHandler struct {
	postService services.PostService
	validator   *validator.Validate
}

func NewPostHandler(postService services.PostService, postValidator *validator.Validate) *PostHandler {
	return &PostHandler{
		postService: postService,
		validator:   postValidator,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req dto.PostCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Warn("Failed to decode post creation request", "err", err)
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	// Validate DTO
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("Post creation request validation failed", "err", err)
		common.ErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	postResponse, err := h.postService.CreatePost(r.Context(), &req)
	if err != nil {
		slog.Error("Failed to create post", "err", err)
		// Check for specific errors from service to return appropriate HTTP status
		if strings.Contains(err.Error(), "unauthorized") {
			common.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		} else if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "invalid or inaccessible media ID") {
			common.ErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		} else {
			common.ErrorResponse(w, http.StatusInternalServerError, "Failed to create post", nil)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(postResponse); err != nil {
		slog.Error("Failed to encode post response", "err", err, "post_id", postResponse.ID)
		common.ErrorResponse(w, http.StatusInternalServerError, "Failed to encode response", nil)
	}
}
