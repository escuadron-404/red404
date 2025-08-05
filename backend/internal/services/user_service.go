package services

import (
	"context"
	"fmt"
	"time"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/escuadron-404/red404/backend/internal/repositories"
	"github.com/escuadron-404/red404/backend/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]dto.UserResponse, int, error)
	UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo      repositories.UserRepository
	validator *validator.Validate
}

func NewUserService(repo repositories.UserRepository, userValidator *validator.Validate) UserService {
	return &userService{
		repo:      repo,
		validator: userValidator,
	}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user model
	now := time.Now()
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Return response
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]dto.UserResponse, int, error) {
	users, totalCount, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("service failed to get paginated users from repo: %w", err)
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return userResponses, totalCount, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get existing user
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Update fields
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		existingUser.Password = hashedPassword
	}
	existingUser.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Return response
	return &dto.UserResponse{
		ID:        existingUser.ID,
		Email:     existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Check if user exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Delete from database
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
