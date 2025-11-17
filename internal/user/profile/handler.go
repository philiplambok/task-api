package profile

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/philiplambok/task-api/internal/pkg/ctx"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
	"github.com/philiplambok/task-api/internal/user/common/datamodel"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_querier_test.go -package=profile -source=handler.go

// Querier defines the interface for querying user data
type Querier interface {
	FindUserByEmail(ctx context.Context, email string) (*datamodel.User, error)
}

type Handler struct {
	querier Querier
}

func NewHandler(querier Querier) *Handler {
	return &Handler{
		querier: querier,
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

	// Get user from repository
	user, err := h.querier.FindUserByEmail(r.Context(), email)
	if err != nil {
		httperror.Handle(w, r, err)
		return
	}

	// Check if user exists
	if user == nil {
		httpErr := httperror.NewHTTPError(http.StatusNotFound, "User not found")
		httperror.Handle(w, r, httpErr)
		return
	}

	var resp v1.GetUserResponse
	resp.Data.User.Email = user.Email
	resp.Data.User.CreatedAt = user.CreatedAt

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
