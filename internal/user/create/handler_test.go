package create

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
	"github.com/philiplambok/task-api/internal/user/shared/datamodel"
	"github.com/philiplambok/task-api/internal/user/shared/domain"
)

var _ = Describe("Handler.CreateUser", func() {
	var (
		ctrl        *gomock.Controller
		handler     *Handler
		mockQuerier *MockQuerier
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockQuerier = NewMockQuerier(ctrl)
		handler = NewHandler(mockQuerier)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("creating a user with valid request", func() {
		It("should return 201 Created with user data", func() {
			// Setup mock
			expectedTime := time.Now()
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(expectedTime, nil)

			// Create request
			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute handler
			handler.CreateUser(w, req)

			// Assert response
			Expect(w.Code).To(Equal(http.StatusCreated))

			var resp v1.CreateUserResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Data.User.Email).To(Equal("test@example.com"))
			Expect(resp.Data.User.CreatedAt).To(BeTemporally("==", expectedTime))
		})

		It("should pass context from request to querier", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, params datamodel.CreateUser) (time.Time, error) {
					Expect(ctx).NotTo(BeNil())
					return time.Now(), nil
				})

			body := `{"email":"context@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should handle different valid email formats", func() {
			testCases := []string{
				"user+tag@example.com",
				"first.last@example.com",
				"user@mail.example.com",
			}

			for _, email := range testCases {
				mockQuerier.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(time.Now(), nil)

				body := `{"email":"` + email + `","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
				req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
				w := httptest.NewRecorder()

				handler.CreateUser(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))
			}
		})
	})

	When("creating a user with duplicate email", func() {
		It("should return 409 Conflict", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Time{}, domain.ErrDuplicateEmail)

			body := `{"email":"duplicate@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusConflict))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(Equal("Email already exists"))
		})

		It("should not call querier for subsequent operations after duplicate error", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Time{}, domain.ErrDuplicateEmail).
				Times(1)

			body := `{"email":"duplicate@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusConflict))
		})
	})

	When("request validation fails", func() {
		It("should return 422 for invalid email format", func() {
			// No mock expectation - validation should fail before calling querier
			body := `{"email":"not-an-email","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("Email"))
		})

		It("should return 422 for missing email field", func() {
			body := `{"password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should return 422 for empty email string", func() {
			body := `{"email":"","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should return 422 for missing password", func() {
			body := `{"email":"test@example.com"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should return 422 for password too short", func() {
			body := `{"email":"test@example.com","password":"short","passwordConfirmation":"short"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("Password"))
		})

		It("should return 422 for password mismatch", func() {
			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"DifferentPass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("Password"))
		})

		It("should not call querier when validation fails", func() {
			// No mock expectation - test will fail if mock is called
			body := `{"email":"invalid","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	When("request has invalid JSON", func() {
		It("should return 400 for malformed JSON", func() {
			body := `{"email":}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).NotTo(BeEmpty())
		})

		It("should return 400 for empty body", func() {
			body := ``
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).NotTo(BeEmpty())
		})

		It("should not call querier when JSON parsing fails", func() {
			// No mock expectation - test will fail if mock is called
			body := `invalid json`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("querier returns unexpected error", func() {
		It("should return 500 for generic database error", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Time{}, errors.New("database connection failed"))

			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(Equal("Internal server error"))
		})

		It("should return 500 for context cancelled error", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Time{}, context.Canceled)

			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("handling edge cases", func() {
		It("should handle very long email addresses", func() {
			longEmail := strings.Repeat("a", 64) + "@" + strings.Repeat("b", 63) + ".com"
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Now(), nil)

			body := `{"email":"` + longEmail + `","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should handle request with extra JSON fields", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Now(), nil)

			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!","extra":"field","another":"value"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should set correct Content-Type header", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Now(), nil)

			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("application/json"))
		})
	})

	When("testing response format", func() {
		It("should return response with correct structure", func() {
			createdAt := time.Date(2025, 11, 16, 12, 0, 0, 0, time.UTC)
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(createdAt, nil)

			body := `{"email":"test@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))

			var resp v1.CreateUserResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			// Verify structure
			Expect(resp.Data.User.Email).To(Equal("test@example.com"))
			Expect(resp.Data.User.CreatedAt).To(BeTemporally("==", createdAt))
		})

		It("should return error response with correct structure for duplicate", func() {
			mockQuerier.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(time.Time{}, domain.ErrDuplicateEmail)

			body := `{"email":"duplicate@example.com","password":"SecurePass123!","passwordConfirmation":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			Expect(w.Code).To(Equal(http.StatusConflict))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).NotTo(BeEmpty())
		})
	})
})
