package list

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	pkgctx "github.com/philiplambok/task-api/internal/pkg/ctx"
	"github.com/philiplambok/task-api/internal/pkg/httperror"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_servicer_test.go -package=list -source=handler.go
type Servicer interface {
	ListTasks(ctx context.Context, userID int64, dto *ListTasksDTO) (*ListTasksResultDTO, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(querier Querier) *Handler {
	return &Handler{
		service: NewService(querier),
	}
}

// buildTaskResponse converts a TaskDTO to v1.Task API response
func buildTaskResponse(task TaskDTO) (v1.Task, error) {
	taskUUID, err := uuid.Parse(task.UUID)
	if err != nil {
		return v1.Task{}, err
	}

	apiTask := v1.Task{
		Id:          taskUUID,
		Title:       task.Title,
		Description: task.Description,
		IsDone:      task.IsDone,
		CreatedAt:   task.CreatedAt,
	}

	if task.Deadline != nil {
		deadline := types.Date{Time: *task.Deadline}
		apiTask.Deadline = &deadline
	}

	return apiTask, nil
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by JWT middleware)
	userID, ok := pkgctx.ExtractUserID(r.Context())
	if !ok {
		httpErr := httperror.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
		httperror.Handle(w, r, httpErr)
		return
	}

	// Parse query parameter
	dto := &ListTasksDTO{}
	isDoneStr := r.URL.Query().Get("is_done")
	if isDoneStr != "" {
		isDone, err := strconv.ParseBool(isDoneStr)
		if err != nil {
			httpErr := httperror.NewHTTPError(http.StatusBadRequest, "Invalid is_done parameter")
			httperror.Handle(w, r, httpErr)
			return
		}
		dto.IsDone = &isDone
	}

	// List tasks via service
	result, err := h.service.ListTasks(r.Context(), userID, dto)
	if err != nil {
		httperror.Handle(w, r, err)
		return
	}

	// Build response
	var resp v1.ListTasksResponse
	resp.Data.Tasks = make([]v1.Task, len(result.Tasks))
	for i, task := range result.Tasks {
		apiTask, err := buildTaskResponse(task)
		if err != nil {
			httperror.Handle(w, r, err)
			return
		}
		resp.Data.Tasks[i] = apiTask
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
