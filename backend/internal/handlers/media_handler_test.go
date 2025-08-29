package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
)

const testUploadDir = "../../test_uploads"

// SetupTestDB creates a connection pool to a test database.
// Ensure you have a 'red404_test' database created and accessible.
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	connStr := "postgres://user:password@localhost:5432/red404_test?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	// clear existing data before each test run for isolation
	_, err = pool.Exec(context.Background(), `
		DELETE FROM media;
		DELETE FROM posts;
		DELETE FROM users;
	`)
	if err != nil {
		t.Fatalf("Failed to clear test database: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
		// clean up test_uploads directory
		os.RemoveAll(testUploadDir)
	})
	return pool
}

// Helper to create a dummy user and return its ID
func CreateTestUser(t *testing.T, pool *pgxpool.Pool) int {
	userRepo := repositories.NewUserRepository(pool)
	user := &models.User{
		Email:          fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
		Password:       "password123",
		FullName:       new(string),
		Bio:            new(string),
		ProfilePicture: new(string),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Deleted:        false,
		DeletedAt:      nil,
	}
	if err := userRepo.Create(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user.ID
}

func TestMediaUploadHandler_Upload(t *testing.T) {
	pool := SetupTestDB(t)
	validate := validator.New()
	testLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(testLogger)

	// Create test_uploads directory if it doesn't exist
	if err := os.MkdirAll(testUploadDir, 0755); err != nil {
		t.Fatalf("Failed to create test upload directory: %v", err)
	}

	mediaRepo := repositories.NewMediaRepository(pool)
	mediaService := services.NewMediaService(validate, mediaRepo)
	mediaHandler := handlers.NewMediaUploadHandler(mediaService, validate)

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
		fileName       string
		fileContent    string
		contentType    string
		expectedStatus int
		expectMediaID  bool
	}{
		{
			name:           "Successful Image Upload",
			fileName:       "test_image.jpg",
			fileContent:    "a beautiful redhead damsel sitting alone in a forest",
			contentType:    "image/jpeg",
			expectedStatus: http.StatusCreated,
			expectMediaID:  true,
		},
		{
			name:           "Too Large File",
			fileName:       "large_file.txt",
			fileContent:    strings.Repeat("iaaaa", 11<<20),
			contentType:    "text/plain",
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectMediaID:  false,
		},
		{
			name:           "Disallowed File Type",
			fileName:       "script.sh",
			fileContent:    "echo 'hello'",
			contentType:    "application/x-sh",
			expectedStatus: http.StatusBadRequest,
			expectMediaID:  false,
		},
		{
			name:           "No File Uploaded",
			fileName:       "",
			fileContent:    "",
			contentType:    "",
			expectedStatus: http.StatusBadRequest,
			expectMediaID:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.fileName != "" {
				part, err := writer.CreateFormFile("file_upload", tt.fileName)
				if err != nil {
					t.Fatalf("Failed to create form file: %v", err)
				}
				_, err = io.WriteString(part, tt.fileContent)
				if err != nil {
					t.Fatalf("Failed to write file content: %v", err)
				}
			}
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			ctx := context.WithValue(req.Context(), middleware.UserContextKey, testClaims)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			mediaHandler.Upload(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.expectMediaID {
				var resp dto.MediaUploadResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.ID == "" {
					t.Errorf("Expected media ID to be present, but got empty")
				}
				if resp.OriginalFilename != tt.fileName {
					t.Errorf("Expected original filename %s, got %s", tt.fileName, resp.OriginalFilename)
				}

				mediaInfo, err := mediaRepo.GetMediaByID(context.Background(), resp.ID)
				if err != nil {
					t.Fatalf("Failed to retrieve media from DB: %v", err)
				}
				if mediaInfo == nil {
					t.Fatalf("Media record not found in DB for ID: %s", resp.ID)
				}

				expectedPath := filepath.Join(testUploadDir, mediaInfo.StoredFilename)
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("Expected file to exist at %s, but it does not.", expectedPath)
				}

				os.Remove(expectedPath)
			}
		})
	}
}
