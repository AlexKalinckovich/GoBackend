package custom_errors

import errCode "github.com/brota/gobackend/internal/constants/errors"
import domainErrors "github.com/brota/gobackend/internal/constants/errors/domain"

type ConstraintError struct {
	constraint string
}

func NewConstraintError(constraint string) *ConstraintError {
	return &ConstraintError{constraint: constraint}
}

func (e *ConstraintError) Error() string {
	return e.constraint
}

func (e *ConstraintError) Code() errCode.ErrorCode {
	return domainErrors.ConstraintErrorCode
}
