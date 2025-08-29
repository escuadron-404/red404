package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/escuadron-404/red404/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	mediaFolderName     = "uploads"
	tempFileExpiry      = time.Minute * 30
	maxCleanupBatchSize = 100
)

type MediaService interface {
	PrepareFileUpload(ctx context.Context, originalFilename string, contentType string) (io.WriteCloser, *dto.MediaUploadResponse, error)
	FinalizeFileUpload(ctx context.Context, mediaID string, bytesWritten int64, actualContentType string) error
	GetMediaForPostCreation(ctx context.Context, mediaID string, userID int) (*models.Media, error)
	LinkMediaToPost(ctx context.Context, mediaID string, postID int) error
	LinkMediaToProfilePicture(ctx context.Context, mediaID string, userID int) error
	CleanupExpiredUnusedMedia(ctx context.Context)
}

type mediaService struct {
	validator       *validator.Validate
	mediaRepository repositories.MediaRepository
}

// NewMediaService takes pgxpool.Pool to pass to repositories
func NewMediaService(mediaValidator *validator.Validate, mediaRepo repositories.MediaRepository) MediaService {
	return &mediaService{
		validator:       mediaValidator,
		mediaRepository: mediaRepo,
	}
}

// ensureUploadDir checks if the upload directory exists and creates it if not.
func ensureUploadDir() error {
	if _, err := os.Stat(mediaFolderName); os.IsNotExist(err) {
		slog.Info("Upload directory does not exist, creating it.", "path", mediaFolderName)
		return os.Mkdir(mediaFolderName, 0755)
	}
	return nil
}

// getUniqueFilename generates a unique filename (UUID) while preserving the original file extension.
func getUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return uuid.New().String() + ext
}

// isAllowedMediaType checks if the content type is allowed for uploads.
func isAllowedMediaType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "video/")
}

// PrepareFileUpload implements MediaService.PrepareFileUpload.
func (s *mediaService) PrepareFileUpload(ctx context.Context, originalFilename string, contentType string) (io.WriteCloser, *dto.MediaUploadResponse, error) {
	claims := GetUserFromContext(ctx)
	if claims == nil {
		slog.Error("Unauthorized attempt to upload file: user claims missing from context.")
		return nil, nil, fmt.Errorf("unauthorized")
	}

	if !isAllowedMediaType(contentType) {
		slog.Warn("Attempted upload with disallowed media type", "user_id", claims.UserID, "content_type", contentType)
		return nil, nil, fmt.Errorf("unsupported media type: %s", contentType)
	}

	if err := ensureUploadDir(); err != nil {
		slog.Error("Couldn't ensure upload directory", "err", err, "user_id", claims.UserID)
		return nil, nil, fmt.Errorf("couldn't prepare upload directory: %w", err)
	}

	storedFilename := getUniqueFilename(originalFilename)
	dstPath := filepath.Join(mediaFolderName, storedFilename)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		slog.Error("Couldn't create local file for upload", "err", err, "filename", originalFilename, "stored_filename", storedFilename, "user_id", claims.UserID)
		return nil, nil, fmt.Errorf("couldn't create local file for upload: %w", err)
	}

	now := time.Now()
	expiresAt := now.Add(tempFileExpiry)

	mediaRecord := &models.Media{ // Use models.Media
		ID:                   uuid.New().String(),
		UserID:               claims.UserID, // Direct access to UserID
		OriginalFilename:     originalFilename,
		StoredFilename:       storedFilename,
		FileType:             contentType,
		FileSize:             0,
		Status:               "pending",
		TempLocalPath:        dstPath,
		CreatedAt:            now,
		UpdatedAt:            now,
		ExpiresAt:            expiresAt,
		PostID:               nil, // Initially no post linked
		ProfilePictureUserID: nil, // Initially no profile picture linked
	}

	if err := s.mediaRepository.CreateMedia(ctx, mediaRecord); err != nil {
		dstFile.Close()
		os.Remove(dstPath)
		slog.Error("Failed to create media record in DB", "err", err, "user_id", claims.UserID)
		return nil, nil, fmt.Errorf("failed to record media: %w", err)
	}

	response := &dto.MediaUploadResponse{
		ID:               mediaRecord.ID,
		OriginalFilename: originalFilename,
		FileType:         contentType,
		FileURL:          fmt.Sprintf("/%s/%s", mediaFolderName, storedFilename),
		ExpiresAt:        expiresAt,
		Message:          "File upload started. Keep this ID for post creation.",
	}

	slog.Info("Prepared file for upload", "media_id", mediaRecord.ID, "user_id", claims.UserID, "original_filename", originalFilename, "stored_filename", storedFilename)
	return dstFile, response, nil
}

// FinalizeFileUpload implements MediaService.FinalizeFileUpload.
func (s *mediaService) FinalizeFileUpload(ctx context.Context, mediaID string, bytesWritten int64, actualContentType string) error {
	claims := GetUserFromContext(ctx) // Now returns *utils.Claims
	if claims == nil {
		slog.Error("Unauthorized attempt to finalize file: user claims missing from context.")
		return fmt.Errorf("unauthorized")
	}

	media, err := s.mediaRepository.GetMediaByIDAndUserID(ctx, mediaID, claims.UserID)
	if err != nil {
		slog.Error("Failed to retrieve media record for finalization", "err", err, "media_id", mediaID, "user_id", claims.UserID)
		return fmt.Errorf("media record not found or not owned")
	}
	if media == nil {
		slog.Error("Media record not found for finalization", "media_id", mediaID, "user_id", claims.UserID)
		return fmt.Errorf("media record not found")
	}
	if media.Status != "pending" {
		slog.Warn("Attempted to finalize a media record not in 'pending' state", "media_id", mediaID, "status", media.Status)
		return fmt.Errorf("media already processed or invalid state")
	}

	if actualContentType == "" {
		actualContentType = media.FileType
	}
	if !isAllowedMediaType(actualContentType) {
		slog.Warn("Actual uploaded file content type is disallowed after upload", "media_id", mediaID, "content_type", actualContentType)
		go func() {
			if media.TempLocalPath != "" {
				os.Remove(media.TempLocalPath)
				slog.Info("Removed disallowed file post-upload", "path", media.TempLocalPath)
			}
		}()
		return fmt.Errorf("disallowed file content type detected after upload")
	}

	err = s.mediaRepository.UpdateMediaStatusAndSize(ctx, mediaID, "uploaded", bytesWritten, actualContentType, media.TempLocalPath)
	if err != nil {
		slog.Error("Failed to update media record after upload complete", "err", err, "media_id", mediaID, "user_id", claims.UserID)
		return fmt.Errorf("failed to finalize media record: %w", err)
	}

	slog.Info("File upload finalized successfully", "media_id", mediaID, "user_id", claims.UserID, "bytes_written", bytesWritten)
	return nil
}

// GetMediaForPostCreation retrieves an uploaded media item for linking to a post.
func (s *mediaService) GetMediaForPostCreation(ctx context.Context, mediaID string, userID int) (*models.Media, error) {
	media, err := s.mediaRepository.GetMediaByIDAndUserID(ctx, mediaID, userID)
	if err != nil {
		slog.Error("Failed to retrieve media record for post creation", "err", err, "media_id", mediaID, "user_id", userID)
		return nil, fmt.Errorf("media record lookup failed: %w", err)
	}
	if media == nil {
		slog.Warn("Media record not found or not owned by user", "media_id", mediaID, "user_id", userID)
		return nil, fmt.Errorf("media not found or unauthorized")
	}
	if media.Status != "uploaded" {
		slog.Warn("Media record is not in 'uploaded' state", "media_id", mediaID, "status", media.Status, "user_id", userID)
		return nil, fmt.Errorf("media not ready or already used")
	}
	return media, nil
}

// LinkMediaToPost implements MediaService.LinkMediaToPost.
func (s *mediaService) LinkMediaToPost(ctx context.Context, mediaID string, postID int) error {
	err := s.mediaRepository.LinkMediaToPost(ctx, mediaID, postID)
	if err != nil {
		slog.Error("Failed to link media to post", "err", err, "media_id", mediaID, "post_id", postID)
		return fmt.Errorf("failed to link media to post: %w", err)
	}
	slog.Info("Media linked to post successfully", "media_id", mediaID, "post_id", postID)
	return nil
}

// LinkMediaToProfilePicture implements MediaService.LinkMediaToProfilePicture.
func (s *mediaService) LinkMediaToProfilePicture(ctx context.Context, mediaID string, userID int) error {
	err := s.mediaRepository.LinkMediaToProfilePicture(ctx, mediaID, userID)
	if err != nil {
		slog.Error("Failed to link media to profile picture", "err", err, "media_id", mediaID, "user_id", userID)
		return fmt.Errorf("failed to link media to profile picture: %w", err)
	}
	slog.Info("Media linked to profile picture successfully", "media_id", mediaID, "user_id", userID)
	return nil
}

// CleanupExpiredUnusedMedia is a background job to clean up old, unused files.
func (s *mediaService) CleanupExpiredUnusedMedia(ctx context.Context) {
	slog.Info("Starting media cleanup job...")
	ticker := time.NewTicker(time.Hour * 6)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Media cleanup job shutting down.")
			return
		case <-ticker.C:
			slog.Info("Running media cleanup iteration...")
			mediaToClean, err := s.mediaRepository.GetExpiredUnusedMedia(context.Background(), maxCleanupBatchSize)
			if err != nil {
				slog.Error("Failed to retrieve expired unused media for cleanup", "err", err)
				continue
			}

			if len(mediaToClean) == 0 {
				slog.Info("No expired unused media found to clean up.")
				continue
			}

			for _, media := range mediaToClean {
				slog.Info("Cleaning up expired media", "media_id", media.ID, "stored_filename", media.StoredFilename, "status", media.Status)
				if media.TempLocalPath != "" {
					if err := os.Remove(media.TempLocalPath); err != nil {
						slog.Error("Failed to delete expired media file from disk", "err", err, "path", media.TempLocalPath)
						s.mediaRepository.MarkMediaAsUsed(context.Background(), media.ID, "failed_cleanup_file")
						continue
					}
					slog.Info("Successfully deleted expired media file from disk", "path", media.TempLocalPath)
				} else {
					slog.Warn("Media record has no local path to delete, skipping file removal", "media_id", media.ID)
				}

				if err := s.mediaRepository.DeleteMedia(context.Background(), media.ID); err != nil {
					slog.Error("Failed to delete expired media record from DB", "err", err, "media_id", media.ID)
					s.mediaRepository.MarkMediaAsUsed(context.Background(), media.ID, "failed_cleanup_db")
				} else {
					slog.Info("Successfully deleted expired media record from DB", "media_id", media.ID)
				}
			}
		}
	}
}
