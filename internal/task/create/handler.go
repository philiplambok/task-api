package create

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/oapi-codegen/runtime/types"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	pkgctx "github.com/philiplambok/task-api/internal/pkg/ctx"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
	"github.com/philiplambok/task-api/internal/pkg/validator"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_servicer_test.go -package=create -source=handler.go
type Servicer interface {
	CreateTask(ctx context.Context, userID int64, dto *CreateTaskDTO) (*CreateTaskResultDTO, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(querier Querier) *Handler {
	return &Handler{
		service: NewService(querier),
	}
}

// parseCreateTaskRequestDTO parses and validates an HTTP request to CreateTaskDTO.
func parseCreateTaskRequestDTO(r *http.Request) (*CreateTaskDTO, error) {
	var req v1.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	var deadline *string
	if req.Deadline != nil {
		// Extract the date from the openapi Date type
		deadlineStr := req.Deadline.Time.Format("2006-01-02")
		deadline = &deadlineStr
	}

	createReq := &CreateTaskDTO{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    deadline,
	}

	if err := validator.ValidateStruct(createReq); err != nil {
		return nil, err
	}

	return createReq, nil
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by JWT middleware)
	userID, ok := pkgctx.ExtractUserID(r.Context())
	if !ok {
		httpErr := httperror.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
		httperror.Handle(w, r, httpErr)
		return
	}

	req, err := parseCreateTaskRequestDTO(r)
	if err != nil {
		httperror.Handle(w, r, err)
		return
	}

	// Create task via service
	result, err := h.service.CreateTask(r.Context(), userID, req)
	if err != nil {
		httperror.Handle(w, r, err)
		return
	}

	var resp v1.CreateTaskResponse
	resp.Data.Task.Id = strconv.FormatInt(result.ID, 10)
	resp.Data.Task.Title = result.Title
	resp.Data.Task.Description = result.Description
	if result.Deadline != nil {
		deadline := types.Date{Time: *result.Deadline}
		resp.Data.Task.Deadline = &deadline
	}
	resp.Data.Task.IsDone = result.IsDone
	resp.Data.Task.CreatedAt = result.CreatedAt

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}
