package user

import (
	domainErrorCodes "github.com/brota/gobackend/internal/constants/errors/domain"
	"github.com/brota/gobackend/internal/constants/errors/test"
	"github.com/brota/gobackend/internal/custom_errors/test_errors"
	"testing"
)

func TestUserAlreadyExistsError(t *testing.T) {
	err := NewUserAlreadyExistsError("john@gmail.com")
	if err.Error() != "user already exists" {
		t.Errorf("Expected 'user already exists', got '%s'", err.Error())
	}
	if err.Code() != domainErrorCodes.AlreadyExistsErrorCode {
		t.Errorf("Expected UserAlreadyExistsErrorCode, got %s", err.Code())
	}
	ctx := err.ContextData()
	email, ok := ctx["email"].(string)
	if !ok || email != "john@gmail.com" {
		t.Errorf("Expected email context, got %v", ctx)
	}
}

func TestFieldValidationError(t *testing.T) {
	err := NewFieldValidationError("name", "Al", "min length 3")
	if err.Error() != "validation failed" {
		t.Errorf("Expected 'validation failed', got '%s'", err.Error())
	}
	if err.Code() != domainErrorCodes.ValidationErrorCode {
		t.Errorf("Expected ValidationErrorCode, got %s", err.Code())
	}
	ctx := err.ContextData()
	field, _ := ctx["field"].(string)
	rule, _ := ctx["rule"].(string)
	if field != "name" || rule != "min length 3" {
		t.Errorf("Expected context with field and rule, got %v", ctx)
	}
}

func TestNewErrorsAreExtensible(t *testing.T) {
	err := test_errors.NewInsufficientFundsError(100.0, 50.0, "USD")
	if err.Code() != test.InsufficientFundsErrorCode {
		t.Errorf("Expected InsufficientFundsErrorCode, got %s", err.Code())
	}
	ctx := err.ContextData()
	required, _ := ctx["required"].(float64)
	balance, _ := ctx["balance"].(float64)
	if required != 100.0 || balance != 50.0 {
		t.Errorf("Expected numeric context, got %v", ctx)
	}
}
