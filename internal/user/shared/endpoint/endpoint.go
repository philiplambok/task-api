package endpoint

import (
	"github.com/go-chi/chi/v5"
	"github.com/philiplambok/task-api/internal/user/common/repository"
	"github.com/philiplambok/task-api/internal/user/create"
	"github.com/philiplambok/task-api/internal/user/profile"
	"gorm.io/gorm"
)

// NewPublicEndpoint creates and configures the public user endpoints (no authentication required)
func NewPublicEndpoint(db *gorm.DB) *chi.Mux {
	router := chi.NewMux()

	// Create feature - public endpoint
	createRepo := repository.NewRepository(db)
	createHandler := create.NewHandler(createRepo)
	router.Post("/", createHandler.CreateUser)

	return router
}

// NewProtectedEndpoint creates and configures the protected user endpoints (authentication required)
// Authentication middleware should be applied at the transport layer
func NewProtectedEndpoint(db *gorm.DB) *chi.Mux {
	router := chi.NewMux()

	// Profile feature (requires authentication)
	profileRepo := repository.NewRepository(db)
	profileHandler := profile.NewHandler(profileRepo)
	router.Get("/profile", profileHandler.GetProfile)

	// Future protected features will be registered here:
	// Update feature
	// updateRepo := repository.NewRepository(db)
	// updateHandler := update.NewHandler(updateRepo)
	// router.Put("/profile", updateHandler.UpdateUser)

	// Delete feature
	// deleteRepo := repository.NewRepository(db)
	// deleteHandler := delete.NewHandler(deleteRepo)
	// router.Delete("/profile", deleteHandler.DeleteUser)

	return router
}
