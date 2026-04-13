package common

import (
	"github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

const (
	AlreadyExistsErrorCode abstract_error_code.ErrorCode = "ALREADY_EXISTS_ERROR"
)

type EntityAlreadyExistsError struct {
	message string
	entity  any
}

func NewEntityAlreadyExistsNilEntityError(message string) *EntityAlreadyExistsError {
	return &EntityAlreadyExistsError{
		message: message,
		entity:  nil,
	}
}

func NewEntityAlreadyExistsError(message string, entity any) *EntityAlreadyExistsError {
	return &EntityAlreadyExistsError{
		message: message,
		entity:  entity,
	}
}

func (e *EntityAlreadyExistsError) Error() string {
	return e.message
}

func (e *EntityAlreadyExistsError) Code() abstract_error_code.ErrorCode {
	return AlreadyExistsErrorCode
}

func (e *EntityAlreadyExistsError) ContextData() map[string]any {

	return map[string]any{
		"entity":  e.entity,
		"message": e.message,
	}
}
