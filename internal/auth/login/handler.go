package login

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
	"github.com/philiplambok/task-api/internal/pkg/validator"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_servicer_test.go -package=login -source=handler.go
type Servicer interface {
	CreateUserAuthToken(ctx context.Context, dto *LoginDTO) (*LoginResultDTO, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(querier Querier, jwtSecret string, tokenExpHours int) *Handler {
	return &Handler{
		service: NewService(querier, jwtSecret, tokenExpHours),
	}
}

// parseLoginRequestDTO parses and validates an HTTP request to LoginDTO.
func parseLoginRequestDTO(r *http.Request) (*LoginDTO, error) {
	var req v1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	loginReq := &LoginDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := validator.ValidateStruct(loginReq); err != nil {
		return nil, err
	}

	return loginReq, nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := parseLoginRequestDTO(r)
	if err != nil {
		httperror.Handle(w, r, err)
		return
	}

	// Create user auth token
	result, err := h.service.CreateUserAuthToken(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpErr := httperror.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
			httperror.Handle(w, r, httpErr)
			return
		}
		httperror.Handle(w, r, err)
		return
	}

	var resp v1.LoginResponse
	resp.Data.User.Email = result.User.Email
	resp.Data.User.CreatedAt = result.User.CreatedAt
	resp.Data.Token = result.Token
	resp.Data.ExpiresAt = result.ExpiresAt

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
