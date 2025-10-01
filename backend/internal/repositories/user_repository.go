package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, int, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db dbExecutor
}

// NewUserRepository initializes with a pgxpool.Pool which implements dbExecutor
func NewUserRepository(db dbExecutor) UserRepository {
	return &userRepository{db: db}
}

// scanUser scans a pgx.Row into a models.User struct.
func scanUser(row pgx.Row) (*models.User, error) {
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.ProfilePicture,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Deleted,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	// Ensure default values for new fields if they are nil
	if user.FullName == nil {
		user.FullName = new(string)
	}
	if user.Bio == nil {
		user.Bio = new(string)
	}
	if user.ProfilePicture == nil {
		user.ProfilePicture = new(string)
	}
	user.Deleted = false // Default to false
	user.DeletedAt = nil // Default to NULL

	query := `
		INSERT INTO users (full_name, email, password, bio, profile_picture, created_at, updated_at, deleted, deleted_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id`
	return r.db.QueryRow(ctx, query,
		user.FullName,
		user.Email,
		user.Password,
		user.Bio,
		user.ProfilePicture,
		user.CreatedAt,
		user.UpdatedAt,
		user.Deleted,
		user.DeletedAt,
	).Scan(&user.ID)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, full_name, email, password, bio, profile_picture, created_at, updated_at, deleted, deleted_at FROM users WHERE id = $1`
	user, err := scanUser(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, full_name, email, password, bio, profile_picture, created_at, updated_at, deleted, deleted_at FROM users WHERE email = $1`
	user, err := scanUser(r.db.QueryRow(ctx, query, email))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	} else if limit < 0 {
		limit = 0
	}

	var totalUsers int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted = FALSE` // Only count active users
	err := r.db.QueryRow(ctx, countQuery).Scan(&totalUsers)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users for pagination: %w", err)
	}

	if totalUsers == 0 {
		return []models.User{}, 0, nil
	}

	query := `SELECT id, full_name, email, password, bio, profile_picture, created_at, updated_at, deleted, deleted_at FROM users WHERE deleted = FALSE ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query paginated users: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0, limit)

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Password, &user.Bio, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt, &user.Deleted, &user.DeletedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user row during pagination: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during paginated rows iteration: %w", err)
	}

	return users, totalUsers, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now() // Update timestamp on modification

	query := `
		UPDATE users 
		SET full_name = $1, email = $2, password = $3, bio = $4, profile_picture = $5, updated_at = $6, deleted = $7, deleted_at = $8 
		WHERE id = $9
	`
	cmdTag, err := r.db.Exec(ctx, query,
		user.FullName,
		user.Email,
		user.Password,
		user.Bio,
		user.ProfilePicture,
		user.UpdatedAt,
		user.Deleted,
		user.DeletedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found for update")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE users SET deleted = TRUE, deleted_at = $1 WHERE id = $2`
	cmdTag, err := r.db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to soft delete user: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found for deletion")
	}
	return nil
}
