package login

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/philiplambok/task-api/internal/auth/common/datamodel"
	"github.com/philiplambok/task-api/internal/auth/common/domain"
)

var _ = Describe("Service.CreateUserAuthToken", func() {
	var (
		ctrl        *gomock.Controller
		service     *Service
		mockQuerier *MockQuerier
		jwtSecret   string
		tokenExpHrs int
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockQuerier = NewMockQuerier(ctrl)
		jwtSecret = "test-secret-key"
		tokenExpHrs = 24
		service = NewService(mockQuerier, jwtSecret, tokenExpHrs)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("login with valid credentials", func() {
		It("should return login result with token", func() {
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

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.User.Email).To(Equal("test@example.com"))
			Expect(result.Token).NotTo(BeEmpty())
			Expect(result.ExpiresAt).To(BeTemporally(">", time.Now()))
		})

		It("should pass context to querier", func() {
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

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			_, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should generate valid JWT token", func() {
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

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Token).To(MatchRegexp(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`))
		})

		It("should set correct expiration time", func() {
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

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			now := time.Now()
			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).NotTo(HaveOccurred())
			expectedExpiration := now.Add(time.Duration(tokenExpHrs) * time.Hour)
			Expect(result.ExpiresAt).To(BeTemporally("~", expectedExpiration, 5*time.Second))
		})
	})

	When("user does not exist", func() {
		It("should return ErrInvalidCredentials", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "nonexistent@example.com").
				Return(nil, nil)

			dto := &LoginDTO{
				Email:    "nonexistent@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(ErrInvalidCredentials))
			Expect(result).To(BeNil())
		})

		It("should not attempt password comparison when user not found", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "nonexistent@example.com").
				Return(nil, nil).
				Times(1)

			dto := &LoginDTO{
				Email:    "nonexistent@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(ErrInvalidCredentials))
			Expect(result).To(BeNil())
		})
	})

	When("password is incorrect", func() {
		It("should return ErrInvalidCredentials", func() {
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

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "WrongPassword123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(ErrInvalidCredentials))
			Expect(result).To(BeNil())
		})

		It("should not generate token when password is wrong", func() {
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
				Return(expectedUser, nil).
				Times(1)

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "WrongPassword123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(ErrInvalidCredentials))
			Expect(result).To(BeNil())
		})
	})

	When("database returns error", func() {
		It("should propagate the error", func() {
			dbError := errors.New("database connection failed")
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(nil, dbError)

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(dbError))
			Expect(result).To(BeNil())
		})

		It("should handle context cancellation", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(nil, context.Canceled)

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(context.Canceled))
			Expect(result).To(BeNil())
		})

		It("should handle context deadline exceeded", func() {
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(nil, context.DeadlineExceeded)

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(context.DeadlineExceeded))
			Expect(result).To(BeNil())
		})
	})

	When("testing edge cases", func() {
		It("should handle empty password field in user record", func() {
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: "",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil)

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			result, err := service.CreateUserAuthToken(context.Background(), dto)

			Expect(err).To(Equal(ErrInvalidCredentials))
			Expect(result).To(BeNil())
		})

		It("should handle different password formats", func() {
			testCases := []struct {
				password string
			}{
				{password: "SimplePass123"},
				{password: "Complex!@#$%Pass123"},
				{password: "Pass with spaces 123"},
				{password: "UTF8Password日本語123"},
			}

			for _, tc := range testCases {
				hashedPassword, _ := domain.HashPassword(tc.password)
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

				dto := &LoginDTO{
					Email:    "test@example.com",
					Password: tc.password,
				}

				result, err := service.CreateUserAuthToken(context.Background(), dto)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Token).NotTo(BeEmpty())
			}
		})

		It("should handle different email formats", func() {
			testEmails := []string{
				"user@example.com",
				"user+tag@example.com",
				"first.last@example.com",
				"user@subdomain.example.com",
			}

			for _, email := range testEmails {
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

				dto := &LoginDTO{
					Email:    email,
					Password: "SecurePass123!",
				}

				result, err := service.CreateUserAuthToken(context.Background(), dto)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.User.Email).To(Equal(email))
			}
		})
	})

	When("testing service initialization", func() {
		It("should create service with correct configuration", func() {
			newService := NewService(mockQuerier, "secret-key", 48)

			Expect(newService).NotTo(BeNil())
			Expect(newService.querier).To(Equal(mockQuerier))
			Expect(newService.jwtSecret).To(Equal("secret-key"))
			Expect(newService.tokenExpHours).To(Equal(48))
		})

		It("should handle different expiration hours", func() {
			testCases := []int{1, 24, 48, 168, 720}

			for _, expHours := range testCases {
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

				svc := NewService(mockQuerier, "test-secret", expHours)
				dto := &LoginDTO{
					Email:    "test@example.com",
					Password: "SecurePass123!",
				}

				now := time.Now()
				result, err := svc.CreateUserAuthToken(context.Background(), dto)

				Expect(err).NotTo(HaveOccurred())
				expectedExpiration := now.Add(time.Duration(expHours) * time.Hour)
				Expect(result.ExpiresAt).To(BeTemporally("~", expectedExpiration, 5*time.Second))
			}
		})
	})

	When("testing concurrent requests", func() {
		It("should handle multiple concurrent login attempts", func() {
			hashedPassword, _ := domain.HashPassword("SecurePass123!")
			expectedUser := &datamodel.User{
				ID:             1,
				Email:          "test@example.com",
				PasswordDigest: hashedPassword,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// Expect 5 concurrent calls
			mockQuerier.EXPECT().
				FindUserByEmail(gomock.Any(), "test@example.com").
				Return(expectedUser, nil).
				Times(5)

			dto := &LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			}

			done := make(chan bool)
			for i := 0; i < 5; i++ {
				go func() {
					result, err := service.CreateUserAuthToken(context.Background(), dto)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).NotTo(BeNil())
					done <- true
				}()
			}

			for i := 0; i < 5; i++ {
				<-done
			}
		})
	})
})
