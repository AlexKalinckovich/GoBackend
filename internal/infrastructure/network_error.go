package infrastructure

import (
	errconsts "github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
)

import infrastructureErrorCodes "github.com/brota/gobackend/internal/constants/errors/infrastructure"

type NetworkError struct {
	reason    string
	retryable bool
}

func NewNetworkError(reason string, retryable bool) *NetworkError {
	return &NetworkError{reason: reason, retryable: retryable}
}

func (e *NetworkError) Error() string {
	return e.reason
}

func (e *NetworkError) Code() errconsts.ErrorCode {
	return infrastructureErrorCodes.NetworkErrorCode
}

func (e *NetworkError) IsRetryable() bool {
	return e.retryable
}
