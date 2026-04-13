package transport

import (
	"github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

const (
	UnknownErrorCode  abstract_error_code.ErrorCode = "UNKNOWN_ERROR"
	InternalErrorCode abstract_error_code.ErrorCode = "INTERNAL_ERROR"
)

type UnknownError struct {
	message string
	cause   error
}

func NewUnknownError(message string) *UnknownError {
	return &UnknownError{message: message}
}

func NewUnknownErrorWithCause(message string, cause error) *UnknownError {
	return &UnknownError{message: message, cause: cause}
}

func (e *UnknownError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *UnknownError) Code() abstract_error_code.ErrorCode {
	return UnknownErrorCode
}

func (e *UnknownError) ContextData() map[string]any {
	data := map[string]any{"message": e.message}
	if e.cause != nil {
		data["cause"] = e.cause.Error()
	}
	return data
}

type InternalError struct {
	message string
}

func NewInternalError(message string) *InternalError {
	return &InternalError{message: message}
}

func (e *InternalError) Error() string {
	return e.message
}

func (e *InternalError) Code() abstract_error_code.ErrorCode {
	return InternalErrorCode
}

func (e *InternalError) ContextData() map[string]any {
	return map[string]any{"message": e.message}
}
