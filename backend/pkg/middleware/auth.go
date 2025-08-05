// Package middleware defines intermediate processing for requests such as auth and cpntext retrieval
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/escuadron-404/red404/backend/pkg/common"
	"github.com/escuadron-404/red404/backend/pkg/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	jwtUtil *utils.JWTUtil
}

func NewAuthMiddleware(jwtUtil *utils.JWTUtil) *AuthMiddleware {
	return &AuthMiddleware{jwtUtil: jwtUtil}
}

func (am *AuthMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.ErrorResponse(w, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			common.ErrorResponse(w, http.StatusUnauthorized, "Bearer token required", nil)
			return
		}

		claims, err := am.jwtUtil.ValidateToken(tokenString)
		if err != nil {
			common.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserFromContext(ctx context.Context) *utils.Claims {
	if claims, ok := ctx.Value(UserContextKey).(*utils.Claims); ok {
		return claims
	}
	return nil
}
