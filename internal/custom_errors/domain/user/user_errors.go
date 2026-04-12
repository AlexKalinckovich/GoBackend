package user

import (
	errconsts "github.com/brota/gobackend/internal/constants/errors"
	"github.com/brota/gobackend/internal/constants/errors/domain"
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
	return domain.AlreadyExistsErrorCode
}

func (e *AlreadyExistsError) ContextData() map[string]any {
	return map[string]any{"email": e.email}
}
