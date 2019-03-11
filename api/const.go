package api

const (
	DEFAULT_LIMIT  = 100
	DEFAULT_OFFSET = 0
)

// Define all error codes send back to clients
// This helps debug easier
const (
	TOKEN_EXPIRED  = 1
	ERR_INTERNAL   = 2
	NO_TOKEN       = 3
	TOKEN_INVALID  = 4
	INVALID_PARAMS = 5
	DUPLICATE_NAME = 6
)
