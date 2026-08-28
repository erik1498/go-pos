package domain

import "errors"

var (
	ErrInternalServer          = errors.New("the server was unable to complete your request. please try again later")
	ErrBadRequest              = errors.New("request body could not be read properly")
	ErrIDInvalid               = errors.New("id is not valid uuidv7")
	ErrChiperTextToShort       = errors.New("chipertext noncesize to short")
	ErrIdempotencyKeyDuplicate = errors.New("idempotency-key is already exist")
	ErrRateLimitExceeded       = errors.New("too many request, please try again later")

	SuccessGetDataByID = "get data by ID success"
	SuccessGetData     = "get data success"
	SuccessCreateData  = "create data success"
	SuccessUpdateData  = "update data success"

	AlertEnvNotFound       = ".env not found"
	AlertAESNot32Character = "aes key length is not 32 character"
	AlertBlindIndexEmpty   = "blind index is required"
)
