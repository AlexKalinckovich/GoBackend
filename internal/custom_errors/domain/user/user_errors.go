package user

import (
	errconsts "github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

const (
	AlreadyExistsErrorCode = "ALREADY_EXISTS"
)

type AlreadyExistsError struct {
	email string
}

func NewUserAlreadyExistsError(email string) *AlreadyExistsError {
	return &AlreadyExistsError{email: email}
}

func (e *AlreadyExistsError) Error() string {
	return "user already exists"
}

func (e *AlreadyExistsError) Code() errconsts.ErrorCode {
	return AlreadyExistsErrorCode
}

func (e *AlreadyExistsError) ContextData() map[string]any {
	return map[string]any{"email": e.email}
}
