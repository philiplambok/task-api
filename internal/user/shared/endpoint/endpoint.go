package endpoint

import (
	"github.com/go-chi/chi/v5"
	"github.com/philiplambok/task-api/internal/auth/login"
	"github.com/philiplambok/task-api/internal/user/create"
	"github.com/philiplambok/task-api/internal/user/profile"
	"github.com/philiplambok/task-api/internal/user/common/repository"
	"gorm.io/gorm"
)

// NewEndpoint creates and configures the user module router
// This is the single source of truth for all user module endpoints
func NewEndpoint(db *gorm.DB, jwtSecret string) *chi.Mux {
	router := chi.NewMux()

	// Create feature
	createRepo := repository.NewRepository(db)
	createHandler := create.NewHandler(createRepo)
	router.Post("/", createHandler.CreateUser)

	// Profile feature (requires authentication)
	profileRepo := repository.NewRepository(db)
	profileHandler := profile.NewHandler(profileRepo)
	router.With(login.JWTAuth(jwtSecret)).Get("/profile", profileHandler.GetProfile)

	// Future features will be registered here:
	// List feature
	// listRepo := list.NewRepository(db)
	// listHandler := list.NewHandler(listRepo)
	// router.Get("/", listHandler.ListUsers)

	// Show feature
	// showRepo := show.NewRepository(db)
	// showHandler := show.NewHandler(showRepo)
	// router.Get("/{id}", showHandler.ShowUser)

	return router
}
