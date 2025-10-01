package repositories

import (
	"context"
	"fmt"

	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaRepository interface {
	CreateMedia(ctx context.Context, media *models.Media) error
	GetMediaByID(ctx context.Context, mediaID string) (*models.Media, error)
	GetMediaByIDAndUserID(ctx context.Context, mediaID string, userID int) (*models.Media, error)
	UpdateMediaStatusAndSize(ctx context.Context, mediaID string, status string, size int64, fileType string, storedPath string) error
	LinkMediaToPost(ctx context.Context, mediaID string, postID int) error
	LinkMediaToProfilePicture(ctx context.Context, mediaID string, userID int) error
	MarkMediaAsUsed(ctx context.Context, mediaID string, newStatus string) error
	DeleteMedia(ctx context.Context, mediaID string) error
	GetExpiredUnusedMedia(ctx context.Context, limit int) ([]models.Media, error)
}

type mediaRepository struct {
	db dbExecutor
}

func NewMediaRepository(db *pgxpool.Pool) MediaRepository {
	return &mediaRepository{db: db}
}

// scanMedia scans a pgx.Row into a models.Media struct.
func scanMedia(row pgx.Row) (*models.Media, error) {
	media := &models.Media{}
	err := row.Scan(
		&media.ID,
		&media.UserID,
		&media.OriginalFilename,
		&media.StoredFilename,
		&media.FileType,
		&media.FileSize,
		&media.Status,
		&media.TempLocalPath,
		&media.CreatedAt,
		&media.UpdatedAt,
		&media.ExpiresAt,
		&media.PostID,
		&media.ProfilePictureUserID,
	)
	if err != nil {
		return nil, err
	}
	return media, nil
}

func (r *mediaRepository) CreateMedia(ctx context.Context, media *models.Media) error {
	query := `
		INSERT INTO media (id, user_id, original_filename, stored_filename, file_type, file_size_bytes, status, temp_local_path, expires_at, created_at, updated_at, post_id, profile_picture_user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.Exec(ctx, query,
		media.ID,
		media.UserID,
		media.OriginalFilename,
		media.StoredFilename,
		media.FileType,
		media.FileSize,
		media.Status,
		media.TempLocalPath,
		media.ExpiresAt,
		media.CreatedAt,
		media.UpdatedAt,
		media.PostID,
		media.ProfilePictureUserID,
	)
	if err != nil {
		return fmt.Errorf("failed to create media record: %w", err)
	}
	return nil
}

func (r *mediaRepository) GetMediaByID(ctx context.Context, mediaID string) (*models.Media, error) {
	query := `SELECT id, user_id, original_filename, stored_filename, file_type, file_size_bytes, status, temp_local_path, created_at, updated_at, expires_at, post_id, profile_picture_user_id FROM media WHERE id = $1`
	media, err := scanMedia(r.db.QueryRow(ctx, query, mediaID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get media by ID: %w", err)
	}
	return media, nil
}

func (r *mediaRepository) GetMediaByIDAndUserID(ctx context.Context, mediaID string, userID int) (*models.Media, error) {
	query := `SELECT id, user_id, original_filename, stored_filename, file_type, file_size_bytes, status, temp_local_path, created_at, updated_at, expires_at, post_id, profile_picture_user_id FROM media WHERE id = $1 AND user_id = $2`
	media, err := scanMedia(r.db.QueryRow(ctx, query, mediaID, userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found or not owned by user
		}
		return nil, fmt.Errorf("failed to get media by ID and user ID: %w", err)
	}
	return media, nil
}

func (r *mediaRepository) UpdateMediaStatusAndSize(ctx context.Context, mediaID string, status string, size int64, fileType string, storedPath string) error {
	query := `
		UPDATE media
		SET status = $2, file_size_bytes = $3, file_type = $4, temp_local_path = $5, updated_at = NOW()
		WHERE id = $1
	`
	cmdTag, err := r.db.Exec(ctx, query, mediaID, status, size, fileType, storedPath)
	if err != nil {
		return fmt.Errorf("failed to update media status and size: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("media not found for update")
	}
	return nil
}

func (r *mediaRepository) LinkMediaToPost(ctx context.Context, mediaID string, postID int) error {
	query := `
		UPDATE media
		SET status = 'used', post_id = $2, updated_at = NOW(), expires_at = NULL
		WHERE id = $1
	`
	cmdTag, err := r.db.Exec(ctx, query, mediaID, postID)
	if err != nil {
		return fmt.Errorf("failed to link media to post: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("media not found to link to post")
	}
	return nil
}

func (r *mediaRepository) LinkMediaToProfilePicture(ctx context.Context, mediaID string, userID int) error {
	query := `
        UPDATE media
        SET status = 'used', profile_picture_user_id = $2, updated_at = NOW(), expires_at = NULL
        WHERE id = $1
    `
	cmdTag, err := r.db.Exec(ctx, query, mediaID, userID)
	if err != nil {
		return fmt.Errorf("failed to link media to profile picture: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("media not found to link to profile picture")
	}
	return nil
}

func (r *mediaRepository) MarkMediaAsUsed(ctx context.Context, mediaID string, newStatus string) error {
	query := `
		UPDATE media
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`
	cmdTag, err := r.db.Exec(ctx, query, mediaID, newStatus)
	if err != nil {
		return fmt.Errorf("failed to mark media as used: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("media not found to mark as used")
	}
	return nil
}

func (r *mediaRepository) DeleteMedia(ctx context.Context, mediaID string) error {
	query := `DELETE FROM media WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, mediaID)
	if err != nil {
		return fmt.Errorf("failed to delete media record: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("media record not found for deletion")
	}
	return nil
}

func (r *mediaRepository) GetExpiredUnusedMedia(ctx context.Context, limit int) ([]models.Media, error) {
	var mediaList []models.Media
	query := `
		SELECT id, user_id, original_filename, stored_filename, file_type, file_size_bytes, status, temp_local_path, created_at, updated_at, expires_at, post_id, profile_picture_user_id
		FROM media
		WHERE status = 'uploaded' AND expires_at < NOW()
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired unused media: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		media, err := scanMedia(rows) // rows also implements pgx.Row
		if err != nil {
			return nil, fmt.Errorf("failed to scan expired media row: %w", err)
		}
		mediaList = append(mediaList, *media)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expired media rows: %w", err)
	}

	return mediaList, nil
}
