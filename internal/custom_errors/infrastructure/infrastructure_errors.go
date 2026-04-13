package infrastructure

import (
	"github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

const (
	NetworkErrorCode            abstract_error_code.ErrorCode = "NETWORK_ERROR"
	ExternalServiceErrorCode    abstract_error_code.ErrorCode = "EXTERNAL_SERVICE_ERROR"
	DatabaseConstraintErrorCode abstract_error_code.ErrorCode = "DATABASE_CONSTRAINT_ERROR"
)

type NetworkError struct {
	message string
	cause   error
}

func NewNetworkError(message string, cause error) *NetworkError {
	return &NetworkError{message: message, cause: cause}
}

func (e *NetworkError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *NetworkError) Code() abstract_error_code.ErrorCode {
	return NetworkErrorCode
}

func (e *NetworkError) ContextData() map[string]any {
	data := map[string]any{"message": e.message}
	if e.cause != nil {
		data["cause"] = e.cause.Error()
	}
	return data
}

type ExternalServiceError struct {
	service string
	message string
	cause   error
}

func NewExternalServiceError(service, message string, cause error) *ExternalServiceError {
	return &ExternalServiceError{service: service, message: message, cause: cause}
}

func (e *ExternalServiceError) Error() string {
	return e.message
}

func (e *ExternalServiceError) Code() abstract_error_code.ErrorCode {
	return ExternalServiceErrorCode
}

func (e *ExternalServiceError) ContextData() map[string]any {
	data := map[string]any{
		"service": e.service,
		"message": e.message,
	}
	if e.cause != nil {
		data["cause"] = e.cause.Error()
	}
	return data
}

type DatabaseConstraintError struct {
	constraint string
	message    string
	cause      error
}

func NewDatabaseConstraintError(constraint, message string, cause error) *DatabaseConstraintError {
	return &DatabaseConstraintError{constraint: constraint, message: message, cause: cause}
}

func (e *DatabaseConstraintError) Error() string {
	return e.message
}

func (e *DatabaseConstraintError) Code() abstract_error_code.ErrorCode {
	return DatabaseConstraintErrorCode
}

func (e *DatabaseConstraintError) ContextData() map[string]any {
	data := map[string]any{
		"constraint": e.constraint,
		"message":    e.message,
	}
	if e.cause != nil {
		data["cause"] = e.cause.Error()
	}
	return data
}
