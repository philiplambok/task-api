package profile

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/philiplambok/task-api/internal/pkg/ctx"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_servicer_test.go -package=profile -source=handler.go
type Servicer interface {
	GetProfile(ctx context.Context, email string) (*ProfileResultDTO, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(querier Querier) *Handler {
	return &Handler{
		service: NewService(querier),
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user email from context (set by JWT middleware)
	email, ok := ctx.ExtractUserEmail(r.Context())
	if !ok {
		httpErr := httperror.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
		httperror.Handle(w, r, httpErr)
		return
	}

	// Get profile via service
	result, err := h.service.GetProfile(r.Context(), email)
	if err != nil {
		// Map domain errors to HTTP errors
		if errors.Is(err, ErrUserNotFound) {
			httpErr := httperror.NewHTTPError(http.StatusNotFound, "User not found")
			httperror.Handle(w, r, httpErr)
			return
		}
		httperror.Handle(w, r, err)
		return
	}

	var resp v1.GetUserResponse
	resp.Data.User.Email = result.Email
	resp.Data.User.CreatedAt = result.CreatedAt

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
