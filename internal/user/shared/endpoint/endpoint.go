package endpoint

import (
	"github.com/go-chi/chi/v5"
	"github.com/philiplambok/task-api/internal/user/common/repository"
	"github.com/philiplambok/task-api/internal/user/profile"
	"gorm.io/gorm"
)

// NewEndpoint creates and configures the user module router
// This is the single source of truth for all user module endpoints
// Authentication middleware should be applied at the transport layer via Group
func NewEndpoint(db *gorm.DB) *chi.Mux {
	router := chi.NewMux()

	// Profile feature
	profileRepo := repository.NewRepository(db)
	profileHandler := profile.NewHandler(profileRepo)
	router.Get("/profile", profileHandler.GetProfile)

	// Future features will be registered here:
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
