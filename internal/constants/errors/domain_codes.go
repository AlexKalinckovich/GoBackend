package errors

const (
	ValidationErrorCode        ErrorCode = "VALIDATION_ERROR"
	BusinessLogicErrorCode     ErrorCode = "BUSINESS_LOGIC_ERROR"
	ConstraintErrorCode        ErrorCode = "CONSTRAINT_ERROR"
	UserAlreadyExistsErrorCode ErrorCode = "USER_ALREADY_EXISTS"
	InsufficientFundsErrorCode ErrorCode = "INSUFFICIENT_FUNDS"
	UnknownErrorCode           ErrorCode = "UNKNOWN_ERROR"
)
