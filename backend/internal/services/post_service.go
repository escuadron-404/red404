package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/escuadron-404/red404/backend/internal/repositories"
	"github.com/escuadron-404/red404/backend/pkg/middleware" // <--- ADD THIS IMPORT
	"github.com/escuadron-404/red404/backend/pkg/utils"      // <--- ADD THIS IMPORT
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostService interface {
	CreatePost(ctx context.Context, req *dto.PostCreateRequest) (*dto.PostResponse, error)
}

type postService struct {
	pool           *pgxpool.Pool
	validator      *validator.Validate
	postRepository repositories.PostRepository
	mediaService   MediaService
}

func NewPostService(
	pool *pgxpool.Pool,
	postValidator *validator.Validate,
	postRepo repositories.PostRepository,
	mediaSvc MediaService,
) PostService {
	return &postService{
		pool:           pool,
		validator:      postValidator,
		postRepository: postRepo,
		mediaService:   mediaSvc,
	}
}

func (s *postService) CreatePost(ctx context.Context, req *dto.PostCreateRequest) (*dto.PostResponse, error) {
	claims := GetUserFromContext(ctx)
	if claims == nil {
		slog.Error("Unauthorized attempt to create post: user claims missing from context.")
		return nil, fmt.Errorf("unauthorized")
	}

	if err := s.validator.Struct(req); err != nil {
		slog.Warn("Validation error for post creation", "err", err, "user_id", claims.UserID)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var media *models.Media
	var imageURL *string
	if req.MediaID != "" {
		var err error
		media, err = s.mediaService.GetMediaForPostCreation(ctx, req.MediaID, claims.UserID)
		if err != nil {
			slog.Warn("Failed to retrieve or validate pre-uploaded media for post", "err", err, "media_id", req.MediaID, "user_id", claims.UserID)
			return nil, fmt.Errorf("invalid or inaccessible media ID: %w", err)
		}
		if media.TempLocalPath == "" {
			slog.Error("Media record found but no local path specified. This indicates an internal issue.", "media_id", req.MediaID)
			return nil, fmt.Errorf("internal error: media path missing")
		}
		url := fmt.Sprintf("/%s/%s", mediaFolderName, media.StoredFilename)
		imageURL = &url
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Failed to begin transaction for post creation", "err", err, "user_id", claims.UserID)
		return nil, fmt.Errorf("failed to create post: internal error")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				slog.Error("Failed to commit transaction for post creation", "err", err, "user_id", claims.UserID)
			}
		}
	}()

	now := time.Now()
	post := &models.Post{
		UserID:      claims.UserID,
		Description: req.Description,
		ImageURL:    imageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
		Deleted:     false,
	}

	postID, err := s.postRepository.CreatePost(ctx, post)
	if err != nil {
		slog.Error("Failed to create post record in DB", "err", err, "user_id", claims.UserID)
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	if req.MediaID != "" {
		err = s.mediaService.LinkMediaToPost(ctx, req.MediaID, postID)
		if err != nil {
			slog.Error("Failed to link media to post after post creation", "err", err, "media_id", req.MediaID, "post_id", postID, "user_id", claims.UserID)
			return nil, fmt.Errorf("failed to finalize media for post: %w", err)
		}
	}

	slog.Info("Post created successfully", "post_id", postID, "user_id", claims.UserID, "media_id", req.MediaID)

	return &dto.PostResponse{
		ID:          postID,
		UserID:      claims.UserID,
		ImageURL:    *imageURL,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetUserFromContext helper for services to retrieve claims from context.
func GetUserFromContext(ctx context.Context) *utils.Claims {
	if claims, ok := ctx.Value(middleware.UserContextKey).(*utils.Claims); ok {
		return claims
	}
	return nil
}
