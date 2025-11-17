package login

import (
	"net/http"
	"strings"

	"github.com/philiplambok/task-api/internal/auth/common/domain"
	"github.com/philiplambok/task-api/internal/pkg/ctx"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
)

// JWTAuth creates a middleware that validates JWT tokens
func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				err := httperror.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
				httperror.Handle(w, r, err)
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				err := httperror.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
				httperror.Handle(w, r, err)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := domain.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				httpErr := httperror.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
				httperror.Handle(w, r, httpErr)
				return
			}

			// Add user data to context
			newCtx := ctx.InjectUserEmail(r.Context(), claims.Email)
			newCtx = ctx.InjectUserID(newCtx, claims.UserID)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
