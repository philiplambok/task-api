package endpoint

import (
	"github.com/go-chi/chi/v5"
	"github.com/philiplambok/task-api/internal/auth/login"
	"github.com/philiplambok/task-api/internal/auth/shared/repository"
	"gorm.io/gorm"
)

// NewEndpoint creates and configures the auth module router
func NewEndpoint(db *gorm.DB, jwtSecret string, tokenExpHours int) *chi.Mux {
	router := chi.NewMux()

	// Login feature
	loginRepo := repository.NewRepository(db)
	loginHandler := login.NewHandler(loginRepo, jwtSecret, tokenExpHours)
	router.Post("/login", loginHandler.Login)

	return router
}
