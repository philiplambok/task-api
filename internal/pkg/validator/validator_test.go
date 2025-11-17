package validator

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestStruct struct {
	Email string `validate:"required,email"`
	Name  string `validate:"required"`
	Age   int    `validate:"min=0,max=150"`
}

var _ = Describe("ValidateStruct", func() {
	When("given a valid struct", func() {
		It("should return nil for struct with all valid fields", func() {
			s := TestStruct{
				Email: "test@example.com",
				Name:  "John Doe",
				Age:   25,
			}

			err := ValidateStruct(&s)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil for struct with minimum valid age", func() {
			s := TestStruct{
				Email: "test@example.com",
				Name:  "John Doe",
				Age:   0,
			}

			err := ValidateStruct(&s)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil for struct with maximum valid age", func() {
			s := TestStruct{
				Email: "test@example.com",
				Name:  "John Doe",
				Age:   150,
			}

			err := ValidateStruct(&s)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("given a struct with validation errors", func() {
		It("should return ValidationErrors for missing required field", func() {
			s := TestStruct{
				Email: "test@example.com",
				Name:  "",
				Age:   25,
			}

			err := ValidateStruct(&s)

			Expect(err).To(HaveOccurred())
			var validationErrs ValidationErrors
			Expect(errors.As(err, &validationErrs)).To(BeTrue())
			Expect(validationErrs.Errors).To(HaveLen(1))
			Expect(validationErrs.Errors[0].Field).To(Equal("Name"))
			Expect(validationErrs.Errors[0].Message).To(Equal("Name is required"))
		})

		It("should return ValidationErrors for invalid email", func() {
			s := TestStruct{
				Email: "invalid-email",
				Name:  "John Doe",
				Age:   25,
			}

			err := ValidateStruct(&s)

			Expect(err).To(HaveOccurred())
			var validationErrs ValidationErrors
			Expect(errors.As(err, &validationErrs)).To(BeTrue())
			Expect(validationErrs.Errors).To(HaveLen(1))
			Expect(validationErrs.Errors[0].Field).To(Equal("Email"))
			Expect(validationErrs.Errors[0].Message).To(Equal("Email must be a valid email address"))
		})

		It("should return ValidationErrors for multiple validation failures", func() {
			s := TestStruct{
				Email: "",
				Name:  "",
				Age:   25,
			}

			err := ValidateStruct(&s)

			Expect(err).To(HaveOccurred())
			var validationErrs ValidationErrors
			Expect(errors.As(err, &validationErrs)).To(BeTrue())
			Expect(validationErrs.Errors).To(HaveLen(2))
		})

		It("should return ValidationErrors with correct field names and messages", func() {
			s := TestStruct{
				Email: "invalid",
				Name:  "",
				Age:   25,
			}

			err := ValidateStruct(&s)

			Expect(err).To(HaveOccurred())
			var validationErrs ValidationErrors
			Expect(errors.As(err, &validationErrs)).To(BeTrue())
			Expect(validationErrs.Errors).To(HaveLen(2))

			// Check that both fields are present in errors
			fields := []string{validationErrs.Errors[0].Field, validationErrs.Errors[1].Field}
			Expect(fields).To(ContainElement("Email"))
			Expect(fields).To(ContainElement("Name"))
		})
	})
})

var _ = Describe("ValidationErrors", func() {
	Describe("Error", func() {
		It("should return first error message with single error", func() {
			ve := ValidationErrors{
				Errors: []ValidationError{
					{Field: "Email", Message: "Email is required"},
				},
			}

			Expect(ve.Error()).To(Equal("Email is required"))
		})

		It("should return first error message with multiple errors", func() {
			ve := ValidationErrors{
				Errors: []ValidationError{
					{Field: "Email", Message: "Email is required"},
					{Field: "Name", Message: "Name is required"},
					{Field: "Age", Message: "Age must be between 0 and 150"},
				},
			}

			Expect(ve.Error()).To(Equal("Email is required"))
		})

		It("should return empty string when no errors present", func() {
			ve := ValidationErrors{
				Errors: []ValidationError{},
			}

			Expect(ve.Error()).To(Equal(""))
		})
	})
})

var _ = Describe("getValidationMessage", func() {
	It("should return correct message for required tag", func() {
		// This is tested implicitly through ValidateStruct tests
		// but we can verify the behavior through struct validation
		s := TestStruct{
			Email: "test@example.com",
			Name:  "",
			Age:   25,
		}

		err := ValidateStruct(&s)
		var validationErrs ValidationErrors
		errors.As(err, &validationErrs)

		Expect(validationErrs.Errors[0].Message).To(Equal("Name is required"))
	})

	It("should return correct message for email tag", func() {
		s := TestStruct{
			Email: "invalid-email",
			Name:  "John Doe",
			Age:   25,
		}

		err := ValidateStruct(&s)
		var validationErrs ValidationErrors
		errors.As(err, &validationErrs)

		Expect(validationErrs.Errors[0].Message).To(Equal("Email must be a valid email address"))
	})
})
