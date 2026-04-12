package domain

import "github.com/brota/gobackend/internal/constants/errors"

const (
	ValidationErrorCode    errors.ErrorCode = "VALIDATION_ERROR"
	BusinessLogicErrorCode errors.ErrorCode = "BUSINESS_LOGIC_ERROR"
	ConstraintErrorCode    errors.ErrorCode = "CONSTRAINT_ERROR"
)
