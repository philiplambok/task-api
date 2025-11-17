package repository

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	commondatamodel "github.com/philiplambok/task-api/internal/common/datamodel"
	"github.com/philiplambok/task-api/internal/user/common/datamodel"
	"github.com/philiplambok/task-api/internal/user/common/domain"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("CreateUser", func() {
	var (
		ctx              context.Context
		pgContainer      *postgrescontainer.PostgresContainer
		db               *gorm.DB
		repository       *Repository
		cleanupContainer func()
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Start PostgreSQL container
		var err error
		pgContainer, err = postgrescontainer.Run(ctx,
			"postgres:16-alpine",
			postgrescontainer.WithDatabase("testdb"),
			postgrescontainer.WithUsername("testuser"),
			postgrescontainer.WithPassword("testpass"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second)),
		)
		Expect(err).NotTo(HaveOccurred())

		// Get connection string
		connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
		Expect(err).NotTo(HaveOccurred())

		// Connect to database with GORM
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		// Run migrations using Goose
		migrationDB, err := goose.OpenDBWithDriver("pgx", connStr)
		Expect(err).NotTo(HaveOccurred())
		defer migrationDB.Close()

		goose.SetTableName("schema_migrations")
		err = goose.Up(migrationDB, "../../../../db/migrations")
		Expect(err).NotTo(HaveOccurred())

		repository = NewRepository(db)

		// Setup cleanup function
		cleanupContainer = func() {
			if pgContainer != nil {
				_ = pgContainer.Terminate(ctx)
			}
		}
	})

	AfterEach(func() {
		cleanupContainer()
	})

	When("creating a user with valid email", func() {
		It("should create user and return created_at timestamp", func() {
			params := datamodel.CreateUser{
				Email: "test@example.com",
			}

			createdAt, err := repository.CreateUser(ctx, params)

			Expect(err).NotTo(HaveOccurred())
			Expect(createdAt).NotTo(BeZero())
			Expect(createdAt).To(BeTemporally("~", time.Now(), 5*time.Second))
		})

		It("should insert the user into the database", func() {
			params := datamodel.CreateUser{
				Email: "test2@example.com",
			}

			_, err := repository.CreateUser(ctx, params)
			Expect(err).NotTo(HaveOccurred())

			// Verify user exists in database
			var count int64
			err = db.Raw("SELECT COUNT(*) FROM users WHERE email = ?", params.Email).Scan(&count).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})

		It("should create a default list for the user", func() {
			params := datamodel.CreateUser{
				Email: "test3@example.com",
			}

			_, err := repository.CreateUser(ctx, params)
			Expect(err).NotTo(HaveOccurred())

			// Get the user ID
			var userID int64
			err = db.Raw("SELECT id FROM users WHERE email = ?", params.Email).Scan(&userID).Error
			Expect(err).NotTo(HaveOccurred())

			// Verify default list exists for the user
			var listCount int64
			err = db.Raw("SELECT COUNT(*) FROM lists WHERE user_id = ? AND name = ?",
				userID, commondatamodel.DefaultListName).Scan(&listCount).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(listCount).To(Equal(int64(1)))
		})

		It("should create default list with correct timestamps", func() {
			params := datamodel.CreateUser{
				Email: "test4@example.com",
			}

			createdAt, err := repository.CreateUser(ctx, params)
			Expect(err).NotTo(HaveOccurred())

			// Get the user ID
			var userID int64
			err = db.Raw("SELECT id FROM users WHERE email = ?", params.Email).Scan(&userID).Error
			Expect(err).NotTo(HaveOccurred())

			// Verify list timestamps
			var listResult struct {
				CreatedAt time.Time
				UpdatedAt time.Time
			}
			err = db.Raw("SELECT created_at, updated_at FROM lists WHERE user_id = ?", userID).Scan(&listResult).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(listResult.CreatedAt).NotTo(BeZero())
			Expect(listResult.UpdatedAt).NotTo(BeZero())
			Expect(listResult.CreatedAt).To(BeTemporally("~", createdAt, 1*time.Second))
		})
	})

	When("creating a user with duplicate email", func() {
		It("should return ErrDuplicateEmail", func() {
			params := datamodel.CreateUser{
				Email: "duplicate@example.com",
			}

			// Create first user
			_, err := repository.CreateUser(ctx, params)
			Expect(err).NotTo(HaveOccurred())

			// Try to create second user with same email
			createdAt, err := repository.CreateUser(ctx, params)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(domain.ErrDuplicateEmail))
			Expect(createdAt).To(BeZero())
		})
	})

	When("database operation fails", func() {
		It("should return error when table doesn't exist", func() {
			// Drop the lists and tasks tables first due to foreign key constraints
			err := db.Exec("DROP TABLE IF EXISTS tasks").Error
			Expect(err).NotTo(HaveOccurred())
			err = db.Exec("DROP TABLE IF EXISTS lists").Error
			Expect(err).NotTo(HaveOccurred())
			// Drop the users table to simulate error
			err = db.Exec("DROP TABLE users").Error
			Expect(err).NotTo(HaveOccurred())

			params := datamodel.CreateUser{
				Email: "test@example.com",
			}
			createdAt, err := repository.CreateUser(ctx, params)

			Expect(err).To(HaveOccurred())
			Expect(createdAt).To(BeZero())
			Expect(err).NotTo(MatchError(domain.ErrDuplicateEmail))
		})

		It("should return error when lists table doesn't exist", func() {
			// Drop the lists table to simulate error
			err := db.Exec("DROP TABLE lists").Error
			Expect(err).NotTo(HaveOccurred())

			params := datamodel.CreateUser{
				Email: "test-no-list@example.com",
			}
			createdAt, err := repository.CreateUser(ctx, params)

			Expect(err).To(HaveOccurred())
			Expect(createdAt).To(BeZero())
		})
	})

	When("context is canceled", func() {
		It("should return context error", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel() // Cancel immediately

			params := datamodel.CreateUser{
				Email: "test@example.com",
			}
			createdAt, err := repository.CreateUser(cancelCtx, params)

			Expect(err).To(HaveOccurred())
			Expect(createdAt).To(BeZero())
		})
	})

	When("verifying constraint name detection", func() {
		It("should specifically detect users_email_key constraint", func() {
			params := datamodel.CreateUser{
				Email: "constraint-test@example.com",
			}

			// Create first user
			_, err := repository.CreateUser(ctx, params)
			Expect(err).NotTo(HaveOccurred())

			// Verify the constraint name in database
			var constraintName string
			err = db.Raw(`
				SELECT constraint_name
				FROM information_schema.table_constraints
				WHERE table_name = 'users'
				AND constraint_type = 'UNIQUE'
			`).Scan(&constraintName).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(constraintName).To(Equal("users_email_key"))

			// Try duplicate - should trigger our specific constraint check
			_, err = repository.CreateUser(ctx, params)
			Expect(err).To(MatchError(domain.ErrDuplicateEmail))
		})
	})

	When("handling concurrent user creation", func() {
		It("should handle concurrent requests correctly", func() {
			params := datamodel.CreateUser{
				Email: "concurrent@example.com",
			}
			errChan := make(chan error, 2)

			// Try to create the same user concurrently
			go func() {
				_, err := repository.CreateUser(ctx, params)
				errChan <- err
			}()

			go func() {
				_, err := repository.CreateUser(ctx, params)
				errChan <- err
			}()

			// Collect results
			err1 := <-errChan
			err2 := <-errChan

			// One should succeed, one should fail with duplicate error
			if err1 == nil {
				Expect(err2).To(MatchError(domain.ErrDuplicateEmail))
			} else if err2 == nil {
				Expect(err1).To(MatchError(domain.ErrDuplicateEmail))
			} else {
				// Both failed - at least one should be duplicate error
				Expect(fmt.Sprintf("%v %v", err1, err2)).To(ContainSubstring("email already exists"))
			}

			// Verify only one user was created
			var count int64
			err := db.Raw("SELECT COUNT(*) FROM users WHERE email = ?", params.Email).Scan(&count).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})
	})
})
