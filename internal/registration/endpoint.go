package registration

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// NewEndpoint creates and configures the registration module router
// This is the single source of truth for all registration endpoints
func NewEndpoint(db *gorm.DB) *chi.Mux {
	router := chi.NewMux()

	// Registration feature - public endpoint
	registrationRepo := NewRepository(db)
	registrationHandler := NewHandler(registrationRepo)
	router.Post("/", registrationHandler.Register)

	return router
}
