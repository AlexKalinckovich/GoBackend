package domainerrors

import errconsts "github.com/brota/gobackend/internal/constants/errors"
import domainconsts "github.com/brota/gobackend/internal/constants/errors/domain"

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

func (e *ValidationError) Code() errconsts.ErrorCode {
	return domainconsts.ValidationErrorCode
}

func (e *ValidationError) Field() string {
	return e.field
}
