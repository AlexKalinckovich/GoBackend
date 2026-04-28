package test_errors

import (
	"github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

const InsufficientFundsErrorCode abstract_error_code.ErrorCode = "INSUFFICIENT_FUNDS"

type InsufficientFundsError struct {
	accountID string
	required  float64
	available float64
}

func NewInsufficientFundsError(accountID string, required, available float64) *InsufficientFundsError {
	return &InsufficientFundsError{
		accountID: accountID,
		required:  required,
		available: available,
	}
}

func (e *InsufficientFundsError) Error() string {
	return "insufficient funds"
}

func (e *InsufficientFundsError) Code() abstract_error_code.ErrorCode {
	return InsufficientFundsErrorCode
}

func (e *InsufficientFundsError) ContextData() map[string]any {
	return map[string]any{
		"accountID": e.accountID,
		"required":  e.required,
		"available": e.available,
	}
}
