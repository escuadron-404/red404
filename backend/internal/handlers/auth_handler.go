// Package handlers defines methods for initial request handling
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/services"
	"github.com/escuadron-404/red404/backend/pkg/common"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService services.AuthService, authValidator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   authValidator,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.handleValidationErrors(w, err)
		return
	}

	authResponse, err := h.authService.Register(r.Context(), req)
	if err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	common.CreatedResponse(w, authResponse, "User registered successfully")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.handleValidationErrors(w, err)
		return
	}

	authResponse, err := h.authService.Login(r.Context(), req)
	if err != nil {
		common.ErrorResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	common.SuccessResponse(w, authResponse, "Login successful")
}

func (h *AuthHandler) handleValidationErrors(w http.ResponseWriter, err error) {
	var validationErrors = make([]dto.ValidationError, len(err.(validator.ValidationErrors)))

	for _, err := range err.(validator.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = "This field is required"
		case "email":
			message = "Invalid email format"
		case "min":
			message = fmt.Sprintf("Must be at least %s characters long", err.Param())
		default:
			message = "Invalid value"
		}
		validationErrors = append(validationErrors, dto.ValidationError{
			Field:   err.Field(),
			Message: message,
		})
	}

	response := dto.ErrorResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  validationErrors,
	}

	common.JSONResponse(w, http.StatusBadRequest, common.Response{
		Success: false,
		Message: "Validation failed",
		Error:   response,
	})
}
