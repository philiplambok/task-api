package endpoint

import (
	"github.com/go-chi/chi/v5"
	"github.com/philiplambok/task-api/internal/task/common/repository"
	"github.com/philiplambok/task-api/internal/task/create"
	"gorm.io/gorm"
)

// NewEndpoint creates and configures the task module router
// This is the single source of truth for all task module endpoints
// Note: Authentication middleware should be applied at the transport layer
func NewEndpoint(db *gorm.DB) *chi.Mux {
	router := chi.NewMux()

	// Create feature
	createRepo := repository.NewRepository(db)
	createHandler := create.NewHandler(createRepo)
	router.Post("/", createHandler.CreateTask)

	// Future features will be registered here:
	// List feature
	// listRepo := repository.NewRepository(db)
	// listHandler := list.NewHandler(listRepo)
	// router.Get("/", listHandler.ListTasks)

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
