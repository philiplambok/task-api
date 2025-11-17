package httperror

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/render"
	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/pkg/validator"
)

// HTTPError represents a custom HTTP error with a status code and message.
// Use NewHTTPError to create instances of this error type.
type HTTPError struct {
	statusCode int
	message    string
}

func (e HTTPError) Error() string {
	return e.message
}

// NewHTTPError creates a new HTTP error with the given status code and message.
// Use this in handlers to map domain errors to HTTP responses.
//
// Example usage:
//
//	err := httperror.NewHTTPError(http.StatusConflict, "Email already exists")
//	httperror.Handle(w, r, err)
func NewHTTPError(statusCode int, message string) error {
	return HTTPError{
		statusCode: statusCode,
		message:    message,
	}
}

// Handle processes an error and writes the appropriate HTTP response.
// It handles different error types and returns the correct status code and error message:
//   - HTTPError: returns status code and message from the error itself
//   - validator.ValidationErrors: returns 422 Unprocessable Entity
//   - json.SyntaxError or json.UnmarshalTypeError: returns 400 Bad Request with "Invalid request format"
//   - Other errors: returns 500 Internal Server Error with "Internal server error"
//
// Example usage:
//
//	func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
//	    params, err := NewCreateUserRequest(r)
//	    if err != nil {
//	        httperror.Handle(w, r, err)
//	        return
//	    }
//	    // ... continue with handler logic
//	}
func Handle(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	var message string

	// Check if error is HTTPError (highest priority)
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		statusCode = httpErr.statusCode
		message = httpErr.message
		renderError(w, r, statusCode, message)
		return
	}

	// Handle validation errors
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		statusCode = http.StatusUnprocessableEntity
		message = validationErrs.Error()
		renderError(w, r, statusCode, message)
		return
	}

	// Handle JSON parsing errors
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) {
		statusCode = http.StatusBadRequest
		message = "Invalid request format"
		renderError(w, r, statusCode, message)
		return
	}

	// Handle EOF errors (empty/incomplete request body)
	if errors.Is(err, io.EOF) {
		statusCode = http.StatusBadRequest
		message = "Bad request"
		renderError(w, r, statusCode, message)
		return
	}

	// Handle all other errors with a generic 500 error
	statusCode = http.StatusInternalServerError
	message = "Internal server error"
	renderError(w, r, statusCode, message)
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	render.Status(r, statusCode)
	var resp v1.ErrorResponse
	resp.Error.Message = message
	render.JSON(w, r, resp)
}
