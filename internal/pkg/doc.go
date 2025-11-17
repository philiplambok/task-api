// Package pkg contains reusable utilities and infrastructure code
// that can be shared across all domains in the application.
//
// This package provides common functionality that is not specific to any
// particular business domain, making it available for use throughout the
// entire application.
//
// # Available Packages
//
// api/v1: HTTP API request and response schemas following JSON:API specification
//
// httperror: HTTP error handling utilities with automatic error type detection
// and response formatting
//
// validator: Struct validation utilities using go-playground/validator
//
// # Usage Guidelines
//
// Packages in internal/pkg should:
//   - Be domain-agnostic and reusable
//   - Have no dependencies on domain-specific code
//   - Provide clear, well-documented interfaces
//   - Follow Go idioms and best practices
//
// Example import:
//
//	import (
//	    v1 "github.com/philiplambok/task-api/internal/pkg/api/v1"
//	    "github.com/philiplambok/task-api/internal/pkg/httperror"
//	    "github.com/philiplambok/task-api/internal/pkg/validator"
//	)
package pkg
