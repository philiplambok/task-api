package endpoint

import (
	"github.com/go-chi/chi/v5"
	"github.com/philiplambok/task-api/internal/task/common/repository"
	"github.com/philiplambok/task-api/internal/task/create"
	"github.com/philiplambok/task-api/internal/task/list"
	"gorm.io/gorm"
)

// NewEndpoint creates and configures the task module router
// This is the single source of truth for all task module endpoints
// Note: Authentication middleware should be applied at the transport layer
func NewEndpoint(db *gorm.DB) *chi.Mux {
	router := chi.NewMux()

	// Shared repository for all features
	repo := repository.NewRepository(db)

	// List feature
	listHandler := list.NewHandler(repo)
	router.Get("/", listHandler.ListTasks)

	// Create feature
	createHandler := create.NewHandler(repo)
	router.Post("/", createHandler.CreateTask)

	// Future features will be registered here:

	// Show feature
	// showRepo := repository.NewRepository(db)
	// showHandler := show.NewHandler(showRepo)
	// router.Get("/{id}", showHandler.ShowTask)

	// Update feature
	// updateRepo := repository.NewRepository(db)
	// updateHandler := update.NewHandler(updateRepo)
	// router.Put("/{id}", updateHandler.UpdateTask)

	// Delete feature
	// deleteRepo := repository.NewRepository(db)
	// deleteHandler := delete.NewHandler(deleteRepo)
	// router.Delete("/{id}", deleteHandler.DeleteTask)

	return router
}
