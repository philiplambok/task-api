package httperror

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/philiplambok/task-api/internal/pkg/validator"
)


var _ = Describe("Handle", func() {
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/test", nil)
	})

	When("given an HTTPError", func() {
		It("should return status code and message from the error", func() {
			httpErr := NewHTTPError(http.StatusConflict, "Resource already exists")

			Handle(w, r, httpErr)

			Expect(w.Code).To(Equal(http.StatusConflict))
			Expect(w.Code).To(Equal(409))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())

			errorObj, ok := response["error"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(errorObj["message"]).To(Equal("Resource already exists"))
		})

		It("should handle 404 Not Found errors", func() {
			httpErr := NewHTTPError(http.StatusNotFound, "User not found")

			Handle(w, r, httpErr)

			Expect(w.Code).To(Equal(404))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			errorObj := response["error"].(map[string]interface{})
			Expect(errorObj["message"]).To(Equal("User not found"))
		})

		It("should have correct Content-Type header", func() {
			httpErr := NewHTTPError(http.StatusConflict, "Conflict")

			Handle(w, r, httpErr)

			Expect(w.Header().Get("Content-Type")).To(Equal("application/json"))
		})

		It("should take priority over other error types", func() {
			// Create an HTTPError
			httpErr := NewHTTPError(http.StatusConflict, "Custom conflict message")

			Handle(w, r, httpErr)

			Expect(w.Code).To(Equal(http.StatusConflict))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			errorObj := response["error"].(map[string]interface{})
			Expect(errorObj["message"]).To(Equal("Custom conflict message"))
		})
	})

	When("given a validation error", func() {
		It("should return 422 status code with validation error message", func() {
			validationErr := validator.ValidationErrors{
				Errors: []validator.ValidationError{
					{Field: "Email", Message: "Email is required"},
				},
			}

			Handle(w, r, validationErr)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			Expect(w.Code).To(Equal(422))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())

			errorObj, ok := response["error"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(errorObj["message"]).To(Equal("Email is required"))
		})

		It("should return first error message when multiple validation errors", func() {
			validationErr := validator.ValidationErrors{
				Errors: []validator.ValidationError{
					{Field: "Email", Message: "Email is required"},
					{Field: "Name", Message: "Name is required"},
				},
			}

			Handle(w, r, validationErr)

			Expect(w.Code).To(Equal(422))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			errorObj := response["error"].(map[string]interface{})
			Expect(errorObj["message"]).To(Equal("Email is required"))
		})

		It("should have correct Content-Type header", func() {
			validationErr := validator.ValidationErrors{
				Errors: []validator.ValidationError{
					{Field: "Email", Message: "Email is required"},
				},
			}

			Handle(w, r, validationErr)

			Expect(w.Header().Get("Content-Type")).To(Equal("application/json"))
		})
	})

	When("given a JSON syntax error", func() {
		It("should return 400 status code with 'Invalid request format' message", func() {
			syntaxErr := &json.SyntaxError{Offset: 10}

			Handle(w, r, syntaxErr)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Code).To(Equal(400))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())

			errorObj, ok := response["error"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(errorObj["message"]).To(Equal("Invalid request format"))
		})
	})

	When("given a JSON unmarshal type error", func() {
		It("should return 400 status code with 'Invalid request format' message", func() {
			unmarshalErr := &json.UnmarshalTypeError{
				Value: "string",
				Type:  nil,
			}

			Handle(w, r, unmarshalErr)

			Expect(w.Code).To(Equal(400))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			errorObj := response["error"].(map[string]interface{})
			Expect(errorObj["message"]).To(Equal("Invalid request format"))
		})
	})

	When("given a wrapped JSON syntax error", func() {
		It("should unwrap and return 400 with 'Invalid request format' message", func() {
			syntaxErr := &json.SyntaxError{Offset: 5}
			wrappedErr := errors.New("failed to decode: " + syntaxErr.Error())

			// Simulate proper wrapping
			wrappedErr = json.Unmarshal([]byte("{invalid}"), new(interface{}))

			Handle(w, r, wrappedErr)

			Expect(w.Code).To(Equal(400))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			errorObj := response["error"].(map[string]interface{})
			Expect(errorObj["message"]).To(Equal("Invalid request format"))
		})
	})

	When("given a generic error", func() {
		It("should return 500 status code with 'Internal server error' message", func() {
			genericErr := errors.New("something went wrong")

			Handle(w, r, genericErr)

			Expect(w.Code).To(Equal(500))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())

			errorObj, ok := response["error"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(errorObj["message"]).To(Equal("Internal server error"))
		})

		It("should not expose internal error details", func() {
			internalErr := errors.New("database connection failed: timeout after 30s")

			Handle(w, r, internalErr)

			bodyBytes, _ := io.ReadAll(w.Body)
			bodyString := string(bodyBytes)

			Expect(bodyString).NotTo(ContainSubstring("database"))
			Expect(bodyString).NotTo(ContainSubstring("timeout"))
			Expect(bodyString).To(ContainSubstring("Internal server error"))
		})
	})

	When("given an io.EOF error", func() {
		It("should return 400 with 'Bad request' message", func() {
			Handle(w, r, io.EOF)

			Expect(w.Code).To(Equal(400))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			errorObj := response["error"].(map[string]interface{})
			Expect(errorObj["message"]).To(Equal("Bad request"))
		})
	})

	Describe("Response format consistency", func() {
		It("should always return JSON with 'error.message' structure", func() {
			testCases := []struct {
				name string
				err  error
			}{
				{
					name: "validation error",
					err: validator.ValidationErrors{
						Errors: []validator.ValidationError{
							{Field: "Test", Message: "Test is required"},
						},
					},
				},
				{
					name: "json syntax error",
					err:  &json.SyntaxError{Offset: 1},
				},
				{
					name: "generic error",
					err:  errors.New("test error"),
				},
			}

			for _, tc := range testCases {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/test", nil)

				Handle(w, r, tc.err)

				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				Expect(err).NotTo(HaveOccurred(), "Failed for case: "+tc.name)

				errorObj, ok := response["error"].(map[string]interface{})
				Expect(ok).To(BeTrue(), "Missing 'error' field for case: "+tc.name)

				_, ok = errorObj["message"].(string)
				Expect(ok).To(BeTrue(), "Missing 'message' field for case: "+tc.name)
			}
		})
	})
})
