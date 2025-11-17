package create

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
	"github.com/philiplambok/task-api/internal/pkg/validator"
	"github.com/philiplambok/task-api/internal/user/shared/domain"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_servicer_test.go -package=create -source=handler.go
type Servicer interface {
	CreateUser(ctx context.Context, dto *CreateUserDTO) (*CreateUserResultDTO, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(querier Querier) *Handler {
	return &Handler{
		service: NewService(querier),
	}
}

// parseCreateUserRequestDTO parses and validates an HTTP request to CreateUserDTO.
func parseCreateUserRequestDTO(r *http.Request) (*CreateUserDTO, error) {
	var req v1.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	createReq := &CreateUserDTO{
		Email:                req.Email,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}

	if err := validator.ValidateStruct(createReq); err != nil {
		return nil, err
	}

	return createReq, nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	req, err := parseCreateUserRequestDTO(r)
	if err != nil {
		httperror.Handle(w, r, err)
		return
	}

	// Create user via service
	result, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		// Map domain errors to HTTP errors
		if errors.Is(err, domain.ErrDuplicateEmail) {
			httpErr := httperror.NewHTTPError(http.StatusConflict, "Email already exists")
			httperror.Handle(w, r, httpErr)
			return
		}
		httperror.Handle(w, r, err)
		return
	}

	var resp v1.CreateUserResponse
	resp.Data.User.CreatedAt = result.CreatedAt
	resp.Data.User.Email = result.Email

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}
