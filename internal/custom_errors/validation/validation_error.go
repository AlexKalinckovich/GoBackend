package validation

import (
	errCode "github.com/brota/gobackend/internal/constants/errors"
	"github.com/brota/gobackend/internal/constants/errors/domain"
)

type ValidationError struct {
	field   string
	message string
}

func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{field: field, message: message}
}

func (e *ValidationError) Error() string {
	return e.message
}

func (e *ValidationError) Code() errCode.ErrorCode {
	return domain.ValidationErrorCode
}

func (e *ValidationError) Field() string {
	return e.field
}

func (e *ValidationError) ContextData() map[string]any {
	return map[string]any{"field": e.field, "message": e.message}
}

type ValidationAggregateError struct {
	errors map[string]string
}

func NewValidationAggregateError() *ValidationAggregateError {
	return &ValidationAggregateError{errors: make(map[string]string)}
}

func (e *ValidationAggregateError) AddField(field string, message string) {
	e.errors[field] = message
}

func (e *ValidationAggregateError) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ValidationAggregateError) Error() string {
	return "validation failed"
}

func (e *ValidationAggregateError) Code() errCode.ErrorCode {
	return domain.ValidationAggregateErrorCode
}

func (e *ValidationAggregateError) ContextData() map[string]any {
	return map[string]any{"details": e.errors}
}
