package errors

// ErrorCode represents application-specific error codes
type ErrorCode int

const (
	// Unknown errors
	Unknown ErrorCode = 10000

	// Input validation errors (10xxx)
	InvalidInput ErrorCode = 10001
	NotFound    ErrorCode = 10002

	// Price validation errors (20xxx)
	ErrPriceBelowMin      ErrorCode = 20001
	ErrPriceDropExceeded  ErrorCode = 20002
	ErrPriceRiseExceeded  ErrorCode = 20003

	// Platform API errors (30xxx)
	ErrPlatformAPI     ErrorCode = 30001
	ErrPlatformTimeout ErrorCode = 30002
	ErrPlatformAuth    ErrorCode = 30003

	// Database errors (40xxx)
	ErrDatabase      ErrorCode = 40001
	ErrDuplicate     ErrorCode = 40002
	ErrForeignKey     ErrorCode = 40003

	// Notification errors (50xxx)
	ErrNotificationSend ErrorCode = 50001
)

// HTTPStatus returns the appropriate HTTP status code for an error code
func (c ErrorCode) HTTPStatus() int {
	switch c {
	case InvalidInput:
		return 400
	case NotFound:
		return 404
	case ErrPlatformAuth:
		return 401
	case ErrDatabase, ErrDuplicate, ErrForeignKey:
		return 500
	case ErrPlatformAPI, ErrPlatformTimeout:
		return 502
	default:
		return 500
	}
}

// Message returns a human-readable message for the error code
func (c ErrorCode) Message() string {
	switch c {
	case InvalidInput:
		return "Invalid input"
	case NotFound:
		return "Resource not found"
	case ErrPriceBelowMin:
		return "Price is below minimum threshold"
	case ErrPriceDropExceeded:
		return "Price drop exceeded threshold"
	case ErrPriceRiseExceeded:
		return "Price rise exceeded threshold"
	case ErrPlatformAPI:
		return "Platform API error"
	case ErrPlatformTimeout:
		return "Platform API timeout"
	case ErrPlatformAuth:
		return "Platform authentication failed"
	case ErrDatabase:
		return "Database error"
	case ErrDuplicate:
		return "Duplicate entry"
	case ErrForeignKey:
		return "Foreign key constraint violation"
	case ErrNotificationSend:
		return "Failed to send notification"
	default:
		return "Unknown error"
	}
}
