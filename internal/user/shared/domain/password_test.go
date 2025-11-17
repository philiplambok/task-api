package domain_test

import (
	"testing"

	"github.com/philiplambok/task-api/internal/user/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "SecurePassword123!"

		hashedPassword, err := domain.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword, "hashed password should not equal plain password")
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		password := "SecurePassword123!"

		hash1, err1 := domain.HashPassword(password)
		hash2, err2 := domain.HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "bcrypt should generate different salts")
	})

	t.Run("should hash empty password", func(t *testing.T) {
		password := ""

		hashedPassword, err := domain.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("should return error for password exceeding 72 bytes", func(t *testing.T) {
		// bcrypt has a 72 byte limit and should return an error
		password := "ThisIsAVeryLongPasswordThatExceedsNormalLengthButShouldStillWorkWithBcrypt123456789!"

		hashedPassword, err := domain.HashPassword(password)

		assert.Error(t, err)
		assert.Empty(t, hashedPassword)
		assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
	})

	t.Run("should hash password with special characters", func(t *testing.T) {
		password := "P@$$w0rd!#%&*()[]{}|\\:;\"'<>,.?/~`"

		hashedPassword, err := domain.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("should hash password with unicode characters", func(t *testing.T) {
		password := "パスワード123!@#"

		hashedPassword, err := domain.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("should return nil when password matches", func(t *testing.T) {
		password := "SecurePassword123!"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, password)

		assert.NoError(t, err)
	})

	t.Run("should return error when password does not match", func(t *testing.T) {
		password := "SecurePassword123!"
		wrongPassword := "WrongPassword123!"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, wrongPassword)

		assert.Error(t, err)
	})

	t.Run("should be case sensitive", func(t *testing.T) {
		password := "SecurePassword123!"
		differentCase := "securepassword123!"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, differentCase)

		assert.Error(t, err)
	})

	t.Run("should fail with empty password when hash is not empty", func(t *testing.T) {
		password := "SecurePassword123!"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, "")

		assert.Error(t, err)
	})

	t.Run("should fail with invalid hash format", func(t *testing.T) {
		invalidHash := "not-a-valid-bcrypt-hash"
		password := "SecurePassword123!"

		err := domain.ComparePassword(invalidHash, password)

		assert.Error(t, err)
	})

	t.Run("should work with empty password hash", func(t *testing.T) {
		emptyPassword := ""
		hashedPassword, err := domain.HashPassword(emptyPassword)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, emptyPassword)

		assert.NoError(t, err)
	})

	t.Run("should work with special characters", func(t *testing.T) {
		password := "P@$$w0rd!#%&*()[]{}|\\:;\"'<>,.?/~`"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, password)

		assert.NoError(t, err)
	})

	t.Run("should work with unicode characters", func(t *testing.T) {
		password := "パスワード123!@#"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, password)

		assert.NoError(t, err)
	})

	t.Run("should fail when single character differs", func(t *testing.T) {
		password := "SecurePassword123!"
		slightlyDifferent := "SecurePassword123?"
		hashedPassword, err := domain.HashPassword(password)
		require.NoError(t, err)

		err = domain.ComparePassword(hashedPassword, slightlyDifferent)

		assert.Error(t, err)
	})
}

func TestHashPasswordAndCompareIntegration(t *testing.T) {
	t.Run("should hash and verify multiple different passwords", func(t *testing.T) {
		passwords := []string{
			"Password1!",
			"DifferentPassword2@",
			"AnotherOne3#",
			"短いパス",
			"",
		}

		for _, password := range passwords {
			hashedPassword, err := domain.HashPassword(password)
			require.NoError(t, err, "failed to hash password: %s", password)

			// Should match original password
			err = domain.ComparePassword(hashedPassword, password)
			assert.NoError(t, err, "password should match for: %s", password)

			// Should not match different password
			err = domain.ComparePassword(hashedPassword, password+"wrong")
			assert.Error(t, err, "wrong password should not match for: %s", password)
		}
	})

	t.Run("should not allow cross-password matching", func(t *testing.T) {
		password1 := "FirstPassword123!"
		password2 := "SecondPassword456@"

		hash1, err := domain.HashPassword(password1)
		require.NoError(t, err)
		hash2, err := domain.HashPassword(password2)
		require.NoError(t, err)

		// Password 1 should not match hash 2
		err = domain.ComparePassword(hash2, password1)
		assert.Error(t, err)

		// Password 2 should not match hash 1
		err = domain.ComparePassword(hash1, password2)
		assert.Error(t, err)
	})
}
