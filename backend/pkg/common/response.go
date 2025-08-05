// Package common defines useful interfaces for the whole application
package common

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// TODO: maybe actually return the error?
		log.Printf("Error encoding JSON response: %v", err)
		// TODO: maybe write a 500?
	}
}

func SuccessResponse(w http.ResponseWriter, data any, message string) {
	JSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(w http.ResponseWriter, data any, message string) {
	JSONResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err any) {
	JSONResponse(w, statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}
