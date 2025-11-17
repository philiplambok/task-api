package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/philiplambok/task-api/internal"
	"github.com/philiplambok/task-api/internal/auth/login"
	authendpoint "github.com/philiplambok/task-api/internal/auth/shared/endpoint"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/registration"
	taskendpoint "github.com/philiplambok/task-api/internal/task/shared/endpoint"
	userendpoint "github.com/philiplambok/task-api/internal/user/shared/endpoint"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(config internal.Config, db *gorm.DB) *HTTPServer {
	routes := chi.NewRouter()

	routes.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		s, err := v1.GetSwagger()
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to load OpenAPI spec", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "Failed to load OpenAPI spec")
			return
		}

		// Downgrade OpenAPI version to 3.0.0 for swagger-ui compatibility
		// oapi-codegen upgrades the spec to 3.1.0, but swagger-ui doesn't support it well
		s.OpenAPI = "3.0.0"

		b, err := s.MarshalJSON()
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to marshal OpenAPI spec", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "Failed to marshal OpenAPI spec")
			return
		}

		render.Status(r, http.StatusOK)
		render.Data(w, r, b)
	})

	routes.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	authEndpoint := authendpoint.NewEndpoint(db, config.JWT.Secret, config.JWT.ExpirationHours)
	registrationEndpoint := registration.NewEndpoint(db)
	userEndpoint := userendpoint.NewEndpoint(db)
	taskEndpoint := taskendpoint.NewEndpoint(db)

	routes.Route("/v1", func(v1 chi.Router) {
		v1.Mount("/auth", authEndpoint)
		v1.Mount("/registration", registrationEndpoint)

		v1.Group(func(protected chi.Router) {
			protected.Use(login.JWTAuth(config.JWT.Secret))

			protected.Mount("/users", userEndpoint)
			protected.Mount("/tasks", taskEndpoint)
		})
	})

	return &HTTPServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.HTTPServer.Port),
			Handler: routes,
		},
	}
}

func (h *HTTPServer) ListenAndServe() error {
	return h.server.ListenAndServe()
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
