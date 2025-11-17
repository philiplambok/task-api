package login

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
	"github.com/philiplambok/task-api/internal/auth/shared/datamodel"
	"github.com/philiplambok/task-api/internal/auth/shared/domain"
)

var _ = Describe("Handler.Login", func() {
	var (
		ctrl        *gomock.Controller
		handler     *Handler
		mockQuerier *MockQuerier
		jwtSecret   string
		tokenExpHrs int
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockQuerier = NewMockQuerier(ctrl)
		jwtSecret = "test-secret-key"
		tokenExpHrs = 24
		handler = NewHandler(mockQuerier, jwtSecret, tokenExpHrs)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("login with valid credentials", func() {
		It("should return 200 OK with token", func() {
			// Create a user with hashed password
			hashedPassword, err := domain.HashPassword("SecurePass123!")
			Expect(err).NotTo(HaveOccurred())

			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil)

			// Create request
			body := `{"email":"test@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute handler
			handler.Login(w, req)

			// Assert response
			Expect(w.Code).To(Equal(http.StatusOK))

			var resp v1.LoginResponse
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Data.User.Email).To(Equal("test@example.com"))
			Expect(resp.Data.Token).NotTo(BeEmpty())
			Expect(resp.Data.ExpiresAt).To(BeTemporally(">", time.Now()))
		})

		It("should pass context from request to service", func() {
			hashedPassword, _ := domain.HashPassword("SecurePass123!")
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				DoAndReturn(func(ctx context.Context, email string) (*datamodel.User, error) {
					Expect(ctx).NotTo(BeNil())
					return expectedUser, nil
				})

			body := `{"email":"test@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should handle different valid email formats", func() {
			testCases := []string{
				"user+tag@example.com",
				"first.last@example.com",
				"user@mail.example.com",
			}

			for _, email := range testCases {
				hashedPassword, _ := domain.HashPassword("SecurePass123!")
				expectedUser := &datamodel.User{
					ID:             1,
					Email:          email,
					PasswordDigest: hashedPassword,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}

				mockQuerier.EXPECT().
					FindUserByEmail(gomock.Any(), email).
					Return(expectedUser, nil)

				body := `{"email":"` + email + `","password":"SecurePass123!"}`
				req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
				w := httptest.NewRecorder()

				handler.Login(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			}
		})
	})

	When("login with invalid credentials", func() {
		It("should return 401 for non-existent user", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "nonexistent@example.com").
				Return(nil, nil)

			body := `{"email":"nonexistent@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(Equal("Invalid credentials"))
		})

		It("should return 401 for wrong password", func() {
			hashedPassword, _ := domain.HashPassword("CorrectPassword123!")
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil)

			body := `{"email":"test@example.com","password":"WrongPassword123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(Equal("Invalid credentials"))
		})

		It("should not return token when credentials are invalid", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "nonexistent@example.com").
				Return(nil, nil)

			body := `{"email":"nonexistent@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))

			// Response should be error, not login response
			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("request validation fails", func() {
		It("should return 422 for invalid email format", func() {
			// No mock expectation - validation should fail before calling querier
			body := `{"email":"not-an-email","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("Email"))
		})

		It("should return 422 for missing email field", func() {
			body := `{"password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should return 422 for empty email string", func() {
			body := `{"email":"","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should return 422 for missing password", func() {
			body := `{"email":"test@example.com"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should return 422 for empty password", func() {
			body := `{"email":"test@example.com","password":""}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(ContainSubstring("required"))
		})

		It("should not call querier when validation fails", func() {
			// No mock expectation - test will fail if mock is called
			body := `{"email":"invalid","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	When("request has invalid JSON", func() {
		It("should return 400 for malformed JSON", func() {
			body := `{"email":}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).NotTo(BeEmpty())
		})

		It("should return 400 for empty body", func() {
			body := ``
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).NotTo(BeEmpty())
		})

		It("should not call querier when JSON parsing fails", func() {
			// No mock expectation - test will fail if mock is called
			body := `invalid json`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("service returns unexpected error", func() {
		It("should return 500 for generic database error", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(nil, errors.New("database connection failed"))

			body := `{"email":"test@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).To(Equal("Internal server error"))
		})

		It("should return 500 for context cancelled error", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(nil, context.Canceled)

			body := `{"email":"test@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("handling edge cases", func() {
		It("should handle very long email addresses", func() {
			longEmail := strings.Repeat("a", 64) + "@" + strings.Repeat("b", 63) + ".com"
			hashedPassword, _ := domain.HashPassword("SecurePass123!")
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          longEmail,
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), longEmail).
				Return(expectedUser, nil)

			body := `{"email":"` + longEmail + `","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should handle request with extra JSON fields", func() {
			hashedPassword, _ := domain.HashPassword("SecurePass123!")
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil)

			body := `{"email":"test@example.com","password":"SecurePass123!","extra":"field","another":"value"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should set correct Content-Type header", func() {
			hashedPassword, _ := domain.HashPassword("SecurePass123!")
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil)

			body := `{"email":"test@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("application/json"))
		})
	})

	When("testing response format", func() {
		It("should return response with correct structure", func() {
			hashedPassword, _ := domain.HashPassword("SecurePass123!")
			createdAt := time.Date(2025, 11, 16, 12, 0, 0, 0, time.UTC)
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      createdAt,
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil)

			body := `{"email":"test@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var resp v1.LoginResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			// Verify structure
			Expect(resp.Data.User.Email).To(Equal("test@example.com"))
			Expect(resp.Data.User.CreatedAt).To(BeTemporally("==", createdAt))
			Expect(resp.Data.Token).NotTo(BeEmpty())
			Expect(resp.Data.ExpiresAt).NotTo(BeZero())
		})

		It("should return error response with correct structure for invalid credentials", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "nonexistent@example.com").
				Return(nil, nil)

			body := `{"email":"nonexistent@example.com","password":"SecurePass123!"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))

			var resp v1.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Error.Message).NotTo(BeEmpty())
		})
	})

	When("testing handler initialization", func() {
		It("should create handler with correct configuration", func() {
			newHandler := NewHandler(mockQuerier, "secret-key", 48)

			Expect(newHandler).NotTo(BeNil())
			Expect(newHandler.service).NotTo(BeNil())
		})
	})
})
