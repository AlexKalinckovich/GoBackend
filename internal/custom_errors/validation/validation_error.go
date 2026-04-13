package validation

import (
	"github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
	"strings"
)

const (
	ErrorCode          abstract_error_code.ErrorCode = "VALIDATION_ERROR"
	AggregateErrorCode abstract_error_code.ErrorCode = "VALIDATION_AGGREGATE_ERROR"
)

type Error struct {
	field   string
	message string
}

func NewValidationError(field string, message string) *Error {
	return &Error{field: field, message: message}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() abstract_error_code.ErrorCode {
	return ErrorCode
}

func (e *Error) Field() string {
	return e.field
}

func (e *Error) ContextData() map[string]any {
	return map[string]any{"field": e.field, "message": e.message}
}

type FieldError struct {
	Field   string `json:"field"`
	Value   any    `json:"value,omitempty"`
	Message string `json:"message"`
}

type AggregateError struct {
	Errors []FieldError
}

func NewAggregateError() *AggregateError {
	return &AggregateError{}
}

func (e *AggregateError) Add(field string, value any, message string) {
	e.Errors = append(e.Errors, FieldError{
		Field:   field,
		Value:   value,
		Message: message,
	})
}

func (e *AggregateError) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *AggregateError) Error() string {
	if len(e.Errors) == 0 {
		return "validation passed"
	}
	msgs := make([]string, len(e.Errors))
	for i, fe := range e.Errors {
		msgs[i] = fe.Field + ": " + fe.Message
	}
	return strings.Join(msgs, "; ")
}

func (e *AggregateError) Code() abstract_error_code.ErrorCode {
	return AggregateErrorCode
}

func (e *AggregateError) ContextData() map[string]any {
	return map[string]any{
		"errors": e.Errors,
	}
}
