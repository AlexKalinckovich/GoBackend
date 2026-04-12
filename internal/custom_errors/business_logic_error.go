package custom_errors

import errCode "github.com/brota/gobackend/internal/constants/errors"
import domainErrors "github.com/brota/gobackend/internal/constants/errors/domain"

type BusinessLogicError struct {
	reason string
}

func NewBusinessLogicError(reason string) *BusinessLogicError {
	return &BusinessLogicError{reason: reason}
}

func (e *BusinessLogicError) Error() string {
	return e.reason
}

func (e *BusinessLogicError) Code() errCode.ErrorCode {
	return domainErrors.BusinessLogicErrorCode
}
