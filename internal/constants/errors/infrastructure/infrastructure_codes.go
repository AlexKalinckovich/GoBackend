package infrastructure

import "github.com/brota/gobackend/internal/constants/errors"

const (
	NetworkErrorCode            errors.ErrorCode = "NETWORK_ERROR"
	ExternalServiceErrorCode    errors.ErrorCode = "EXTERNAL_SERVICE_ERROR"
	DatabaseConstraintErrorCode errors.ErrorCode = "DATABASE_CONSTRAINT_ERROR"
)
