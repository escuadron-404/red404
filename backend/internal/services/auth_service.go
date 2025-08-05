// Package services defines business-level logic for use in handlers
package services

import (
	"context"
	"fmt"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/escuadron-404/red404/backend/internal/repositories"
	"github.com/escuadron-404/red404/backend/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	validator *validator.Validate
	jwtUtil   *utils.JWTUtil
}

func NewAuthService(userRepo repositories.UserRepository, authValidator *validator.Validate, jwtUtil *utils.JWTUtil) AuthService {
	return &authService{
		userRepo:  userRepo,
		validator: authValidator,
		jwtUtil:   jwtUtil,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user model
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Generate JWT token
	//  not necessary here i believe
	// token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate token: %v", err)
	// }

	// Return response
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &dto.AuthResponse{
		// Token: token,
		User: userResponse,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	// Return response
	userResponse := dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	return &dto.AuthResponse{
		Token: token,
		User:  userResponse,
	}, nil
}
