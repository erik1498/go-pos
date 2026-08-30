package domain

const (
	SuccessLogin       = "Login success."
	SuccessGetDataByID = "Get data by ID success."
	SuccessGetData     = "Get data success."
	SuccessCreateData  = "Create data success."
	SuccessUpdateData  = "Update data success."

	AlertEnvNotFound       = ".env file not found. Falling back to system environment variables."
	AlertAESNot32Character = "AES-256 key length is strictly required to be 32 characters."
	AlertBlindIndexEmpty   = "Blind index key is strictly required for cryptography."
)
