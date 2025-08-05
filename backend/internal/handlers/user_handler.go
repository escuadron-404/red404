package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/services"
	"github.com/escuadron-404/red404/backend/pkg/common"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService services.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService services.UserService, userValidator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   userValidator,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.handleValidationErrors(w, err)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	common.CreatedResponse(w, user, "User created successfully")
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		common.ErrorResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	common.SuccessResponse(w, user, "User retrieved successfully")
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, totalCount, err := h.userService.GetAllUsers(ctx, limit, offset)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	response := struct {
		Data       []dto.UserResponse `json:"data"`
		TotalCount int                `json:"total_count"`
		Limit      int                `json:"limit"`
		Offset     int                `json:"offset"`
	}{
		Data:       users,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding GetPaginatedUsers response: %v", err)
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.handleValidationErrors(w, err)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, req)
	if err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	common.SuccessResponse(w, user, "User updated successfully")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	err = h.userService.DeleteUser(r.Context(), id)
	if err != nil {
		common.ErrorResponse(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	common.SuccessResponse(w, nil, "User deleted successfully")
}

func (h *UserHandler) handleValidationErrors(w http.ResponseWriter, err error) {
	var validationErrors = make([]dto.ValidationError, len(err.(validator.ValidationErrors)))

	for _, err := range err.(validator.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", err.Field())
		case "email":
			message = "Invalid email format"
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", err.Field())
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
