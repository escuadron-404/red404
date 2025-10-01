package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/handlers"
	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/escuadron-404/red404/backend/internal/repositories"
	"github.com/escuadron-404/red404/backend/internal/services"
	"github.com/escuadron-404/red404/backend/pkg/middleware"
	"github.com/escuadron-404/red404/backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Helper to create a test media record (simulating a pre-upload)
func CreateTestMedia(t *testing.T, pool *pgxpool.Pool, userID int, status string) string {
	mediaRepo := repositories.NewMediaRepository(pool)
	mediaID := uuid.New().String()
	storedFilename := uuid.New().String() + ".jpg" // Simulate file on disk
	media := &models.Media{
		ID:               mediaID,
		UserID:           userID,
		OriginalFilename: "original.jpg",
		StoredFilename:   storedFilename,
		FileType:         "image/jpeg",
		FileSize:         1024,
		Status:           status,
		TempLocalPath:    filepath.Join(testUploadDir, storedFilename),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(time.Hour),
	}
	if err := mediaRepo.CreateMedia(context.Background(), media); err != nil {
		t.Fatalf("Failed to create test media: %v", err)
	}
	// Create a dummy file on disk for the service to find
	err := os.WriteFile(filepath.Join(testUploadDir, storedFilename), []byte("dummy image content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy file for test media: %v", err)
	}
	return mediaID
}

func TestPostHandler_CreatePost(t *testing.T) {
	pool := SetupTestDB(t) // Re-use SetupTestDB from media_handler_test
	validate := validator.New()
	testLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(testLogger)

	// Ensure test_uploads exists
	if err := os.MkdirAll(testUploadDir, 0755); err != nil {
		t.Fatalf("Failed to create test upload directory: %v", err)
	}

	// userRepo := repositories.NewUserRepository(pool)
	mediaRepo := repositories.NewMediaRepository(pool)
	postRepo := repositories.NewPostRepository(pool)

	mediaService := services.NewMediaService(validate, mediaRepo)
	postService := services.NewPostService(pool, validate, postRepo, mediaService)
	postHandler := handlers.NewPostHandler(postService, validate)

	// User ID for context
	testUserID := CreateTestUser(t, pool)
	testClaims := &utils.Claims{
		UserID: testUserID,
		Email:  "testuser@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	tests := []struct {
		name           string
		requestBody    dto.PostCreateRequest
		setupMediaID   string
		mediaStatus    string
		expectedStatus int
		expectPostID   bool
	}{
		{
			name: "Successful Post with Pre-uploaded Media",
			requestBody: dto.PostCreateRequest{
				Description: "This is a test post with an image.",
			},
			mediaStatus:    "uploaded",
			expectedStatus: http.StatusCreated,
			expectPostID:   true,
		},
		{
			name: "Successful Post without Media",
			requestBody: dto.PostCreateRequest{
				Description: "This is a text-only test post.",
			},
			setupMediaID:   "", // No media for this test
			expectedStatus: http.StatusCreated,
			expectPostID:   true,
		},
		{
			name: "Post with Non-existent Media ID",
			requestBody: dto.PostCreateRequest{
				Description: "Post with missing image.",
				MediaID:     uuid.New().String(), // A random, non-existent UUID
			},
			setupMediaID:   "", // No media created for this
			expectedStatus: http.StatusBadRequest,
			expectPostID:   false,
		},
		{
			name: "Post with Media in Pending Status",
			requestBody: dto.PostCreateRequest{
				Description: "Post with pending image.",
			},
			mediaStatus:    "pending", // Media exists but isn't finalized
			expectedStatus: http.StatusBadRequest,
			expectPostID:   false,
		},
		{
			name: "Invalid Request - Missing Description",
			requestBody: dto.PostCreateRequest{
				Description: "", // Missing description
			},
			expectedStatus: http.StatusBadRequest,
			expectPostID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup media if needed for the test case
			if tt.mediaStatus != "" {
				tt.requestBody.MediaID = CreateTestMedia(t, pool, testUserID, tt.mediaStatus)
			}

			reqBodyBytes, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Inject claims into context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, testClaims)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			postHandler.CreatePost(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.expectPostID {
				var resp dto.PostResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.ID == 0 {
					t.Errorf("Expected post ID to be present, but got 0")
				}
				if resp.Description != tt.requestBody.Description {
					t.Errorf("Expected description %q, got %q", tt.requestBody.Description, resp.Description)
				}
				if tt.requestBody.MediaID != "" && resp.ImageURL == "" {
					t.Errorf("Expected ImageURL to be present for post with media, but got empty")
				} else if tt.requestBody.MediaID == "" && resp.ImageURL != "" {
					t.Errorf("Expected ImageURL to be empty for post without media, but got %q", resp.ImageURL)
				}

				// Verify post created in DB
				var dbPost models.Post
				err := pool.QueryRow(context.Background(), "SELECT id, user_id, description, image_url FROM posts WHERE id = $1", resp.ID).
					Scan(&dbPost.ID, &dbPost.UserID, &dbPost.Description, &dbPost.ImageURL)
				if err != nil {
					t.Fatalf("Failed to retrieve post from DB: %v", err)
				}
				if dbPost.ID != resp.ID {
					t.Errorf("DB post ID mismatch: expected %d, got %d", resp.ID, dbPost.ID)
				}

				// If media was used, verify its status in DB
				if tt.requestBody.MediaID != "" {
					mediaInfo, err := mediaRepo.GetMediaByID(context.Background(), tt.requestBody.MediaID)
					if err != nil {
						t.Fatalf("Failed to retrieve media from DB after post creation: %v", err)
					}
					if mediaInfo == nil || mediaInfo.Status != "used" {
						t.Errorf("Expected media status to be 'used', got %s", mediaInfo.Status)
					}
					if mediaInfo.PostID == nil || *mediaInfo.PostID != resp.ID {
						t.Errorf("Expected media to be linked to post %d, but got %v", resp.ID, mediaInfo.PostID)
					}
					// Clean up the associated file on disk. The CleanupExpiredUnusedMedia job will also do this eventually.
					os.Remove(mediaInfo.TempLocalPath)
				}
			}
		})
	}
}
