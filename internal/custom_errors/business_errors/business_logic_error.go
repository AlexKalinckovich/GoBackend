package business_errors

import (
	"github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

const (
	BusinessLogicErrorCode abstract_error_code.ErrorCode = "BUSINESS_LOGIC_ERROR"
)

type BusinessLogicError struct {
	reason string
}

func NewBusinessLogicError(reason string) *BusinessLogicError {
	return &BusinessLogicError{reason: reason}
}

func (e *BusinessLogicError) Error() string {
	return e.reason
}

func (e *BusinessLogicError) Code() abstract_error_code.ErrorCode {
	return BusinessLogicErrorCode
}
