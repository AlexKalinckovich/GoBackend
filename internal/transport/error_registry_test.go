package transport

import (
	errconsts "github.com/brota/gobackend/internal/constants/errors"
	"github.com/brota/gobackend/internal/constants/errors/domain"
	"github.com/brota/gobackend/internal/constants/errors/transport"
	"github.com/brota/gobackend/internal/custom_errors"
	"github.com/brota/gobackend/internal/custom_errors/domain/user"
	"github.com/brota/gobackend/internal/custom_errors/validation"
	"reflect"
	"testing"
)

func TestErrorRegistry_Translate(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		registerCode    errconsts.ErrorCode
		registerHandler ContextualErrorHandler
		expectedStatus  int
		expectedCode    string
		expectedDetail  string
	}{
		{
			name:            "UserAlreadyExists with context",
			err:             user.NewUserAlreadyExistsError("john@gmail.com"),
			registerCode:    domain.AlreadyExistsErrorCode,
			registerHandler: handleUserAlreadyExistsTest,
			expectedStatus:  409,
			expectedCode:    "USER_ALREADY_EXISTS",
			expectedDetail:  "user john@gmail.com exists",
		},
		{
			name:            "FieldValidationError with context",
			err:             user.NewFieldValidationError("name", "Al", "min length 3"),
			registerCode:    domain.ValidationErrorCode,
			registerHandler: handleFieldValidationTest,
			expectedStatus:  422,
			expectedCode:    "VALIDATION_ERROR",
			expectedDetail:  "name: min length 3",
		},
		{
			name:            "Unregistered error uses fallback",
			err:             custom_errors.NewBusinessLogicError("unknown"),
			registerCode:    domain.BusinessLogicErrorCode,
			registerHandler: nil,
			expectedStatus:  500,
			expectedCode:    "INTERNAL_ERROR",
			expectedDetail:  "unexpected error",
		},
		{
			name:            "ContextCarrier with missing context key",
			err:             user.NewFieldValidationError("age", -1, "must be positive"),
			registerCode:    domain.ValidationErrorCode,
			registerHandler: handleValidationWithMissingKeyTest,
			expectedStatus:  422,
			expectedCode:    "VALIDATION_ERROR",
			expectedDetail:  "age invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, executeTestCase(tt))
	}
}

func executeTestCase(tt struct {
	name            string
	err             error
	registerCode    errconsts.ErrorCode
	registerHandler ContextualErrorHandler
	expectedStatus  int
	expectedCode    string
	expectedDetail  string
}) func(t *testing.T) {
	return func(t *testing.T) {
		registry := NewErrorRegistry()
		registerHandlerIfPresent(registry, tt.registerCode, tt.registerHandler)
		result := registry.Translate(tt.err)
		verifyResponse(t, result, tt.expectedStatus, tt.expectedCode, tt.expectedDetail)
	}
}

func registerHandlerIfPresent(registry *ErrorRegistry, code errconsts.ErrorCode, handler ContextualErrorHandler) {
	if handler != nil {
		registry.Register(code, handler)
	}
}

func verifyResponse(t *testing.T, result HTTPResponse, expectedStatus int, expectedCode string, expectedDetail string) {
	if result.Status != expectedStatus {
		t.Errorf("Status = %d, want %d", result.Status, expectedStatus)
	}
	if result.Code != expectedCode {
		t.Errorf("Code = %s, want %s", result.Code, expectedCode)
	}
	if result.Detail != expectedDetail {
		t.Errorf("Detail = %s, want %s", result.Detail, expectedDetail)
	}
}

func handleUserAlreadyExistsTest(err error, ctx map[string]any) HTTPResponse {
	email := extractStringContext(ctx, "email")
	return HTTPResponse{
		Status: 409,
		Code:   string(domain.AlreadyExistsErrorCode),
		Detail: "user " + email + " exists",
	}
}

func handleFieldValidationTest(err error, ctx map[string]any) HTTPResponse {
	field := extractStringContext(ctx, "field")
	rule := extractStringContext(ctx, "rule")
	return HTTPResponse{
		Status: 422,
		Code:   string(domain.ValidationErrorCode),
		Detail: field + ": " + rule,
	}
}

func handleValidationWithMissingKeyTest(err error, ctx map[string]any) HTTPResponse {
	field := extractStringContext(ctx, "field")
	if field == "" {
		field = "unknown"
	}
	return HTTPResponse{
		Status: 422,
		Code:   string(domain.ValidationErrorCode),
		Detail: field + " invalid",
	}
}

func extractStringContext(ctx map[string]any, key string) string {
	value, ok := ctx[key].(string)
	if !ok {
		return ""
	}
	return value
}

func TestErrorRegistry_ContextExtraction(t *testing.T) {
	t.Run("Non-ContextCarrier returns empty map", testNonContextCarrier)
	t.Run("ContextCarrier returns populated map", testContextCarrier)
}

func testNonContextCarrier(t *testing.T) {
	err := custom_errors.NewBusinessLogicError("test")
	ctx := extractContext(err)
	if len(ctx) != 0 {
		t.Errorf("Expected empty context, got %v", ctx)
	}
}

func testContextCarrier(t *testing.T) {
	err := user.NewUserAlreadyExistsError("test@example.com")
	ctx := extractContext(err)
	email := extractStringContext(ctx, "email")
	if email != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got '%s'", email)
	}
}

func TestErrorRegistry_ResolveHandler(t *testing.T) {
	t.Run("Registered code returns handler", testRegisteredHandler)
	t.Run("Unregistered code returns fallback", testUnregisteredHandler)
}

func testRegisteredHandler(t *testing.T) {
	registry := NewErrorRegistry()
	registry.Register(domain.ValidationErrorCode, handleFieldValidationTest)
	handler := registry.ResolveHandler(domain.ValidationErrorCode)
	if handler == nil {
		t.Error("Expected handler, got nil")
	}
}

func testUnregisteredHandler(t *testing.T) {
	registry := NewErrorRegistry()
	handler := registry.ResolveHandler(domain.AlreadyExistsErrorCode)

	p1 := reflect.ValueOf(handler).Pointer()
	p2 := reflect.ValueOf(registry.fallback).Pointer()

	if p1 != p2 {
		t.Error("Expected fallback handler")
	}
}

func TestErrorRegistry_ExtractErrorCode(t *testing.T) {
	t.Run("ErrorCodeProvider returns code", testExtractFromProvider)
	t.Run("Non-provider returns unknown", testExtractFromNonProvider)
}

func testExtractFromProvider(t *testing.T) {
	err := validation.NewValidationError("email", "required")
	code := ExtractErrorCode(err)
	if code != domain.ValidationErrorCode {
		t.Errorf("Expected ValidationErrorCode, got %s", code)
	}
}

func testExtractFromNonProvider(t *testing.T) {
	err := &standardError{message: "plain error"}
	code := ExtractErrorCode(err)
	if code != transport.UnknownErrorCode {
		t.Errorf("Expected UnknownErrorCode, got %s", code)
	}
}

type standardError struct {
	message string
}

func (e *standardError) Error() string {
	return e.message
}
