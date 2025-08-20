package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *models.Post) (int, error)
}

type postRepository struct {
	db dbExecutor
}

func NewPostRepository(db *pgxpool.Pool) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *models.Post) (int, error) {
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.Deleted = false
	post.DeletedAt = nil

	query := `
		INSERT INTO posts (user_id, image_url, description, created_at, updated_at, deleted, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var postID int
	err := r.db.QueryRow(ctx, query,
		post.UserID,
		post.ImageURL,
		post.Description,
		post.CreatedAt,
		post.UpdatedAt,
		post.Deleted,
		post.DeletedAt,
	).Scan(&postID)

	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}
	return postID, nil
}
