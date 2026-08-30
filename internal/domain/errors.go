package domain

import "net/http"

type CustomError struct {
	HTTPCode int
	Message  string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(httpCode int, message string) *CustomError {
	return &CustomError{
		HTTPCode: httpCode,
		Message:  message,
	}
}

var (
	ErrInternalServer    = NewCustomError(http.StatusInternalServerError, "The server was unable to complete your request. Please try again later.")
	ErrRateLimitExceeded = NewCustomError(http.StatusTooManyRequests, "Too many requests. Please try again later.")

	ErrBadRequest          = NewCustomError(http.StatusBadRequest, "Invalid request payload or parameters.")
	ErrIDInvalid           = NewCustomError(http.StatusBadRequest, "Invalid ID format.")
	ErrIdempotencyRequired = NewCustomError(http.StatusBadRequest, "The 'X-Idempotency-Key' header is required.")

	ErrUsernameOrPasswordInvalid = NewCustomError(http.StatusUnauthorized, "Invalid username or password.")
	ErrUnauthorized              = NewCustomError(http.StatusUnauthorized, "Unauthorized access.")
	ErrSessionExpired            = NewCustomError(http.StatusUnauthorized, "User session is invalid or has expired.")
	ErrForbidden                 = NewCustomError(http.StatusForbidden, "Access forbidden. You do not have the required permissions.")

	ErrTaxNotFound      = NewCustomError(http.StatusNotFound, "Tax data not found.")
	ErrCategoryNotFound = NewCustomError(http.StatusNotFound, "Category data not found.")
	ErrMemberNotFound   = NewCustomError(http.StatusNotFound, "Member data not found.")
	ErrOrderNotFound    = NewCustomError(http.StatusNotFound, "Order not found.")
	ErrProductNotFound  = NewCustomError(http.StatusNotFound, "Product data not found.")

	ErrIdempotencyKeyDuplicate    = NewCustomError(http.StatusConflict, "The provided X-Idempotency-Key has already been used.")
	ErrRequestProcessed           = NewCustomError(http.StatusConflict, "The request is currently being processed. Please wait.")
	ErrOrderNoIsAlreadyRegistered = NewCustomError(http.StatusConflict, "Order number is already registered.")
	ErrProductStockIsNotEnough    = NewCustomError(http.StatusConflict, "Insufficient product stock.")

	ErrPhoneAlreadyRegistered        = NewCustomError(http.StatusConflict, "Phone number is already registered.")
	ErrEmailAlreadyRegistered        = NewCustomError(http.StatusConflict, "Email address is already registered.")
	ErrProductSKUIsAlreadyRegistered = NewCustomError(http.StatusConflict, "Product SKU is already registered.")
)
