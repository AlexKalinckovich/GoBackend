package domainerrors

import errconsts "github.com/brota/gobackend/internal/constants/errors"

type UserAlreadyExistsError struct {
	email string
}

func NewUserAlreadyExistsError(email string) *UserAlreadyExistsError {
	return &UserAlreadyExistsError{email: email}
}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

func (e *UserAlreadyExistsError) Code() errconsts.ErrorCode {
	return errconsts.UserAlreadyExistsErrorCode
}

func (e *UserAlreadyExistsError) ContextData() map[string]any {
	return map[string]any{"email": e.email}
}

type FieldValidationError struct {
	field string
	value any
	rule  string
}

func NewFieldValidationError(field string, value any, rule string) *FieldValidationError {
	return &FieldValidationError{field: field, value: value, rule: rule}
}

func (e *FieldValidationError) Error() string {
	return "validation failed"
}

func (e *FieldValidationError) Code() errconsts.ErrorCode {
	return errconsts.ValidationErrorCode
}

func (e *FieldValidationError) ContextData() map[string]any {
	return map[string]any{"field": e.field, "value": e.value, "rule": e.rule}
}
