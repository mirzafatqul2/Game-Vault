package response

// SuccessResponse is standard success response format
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is standard error response format
type ErrorResponse struct {
	Message string `json:"message"`
}

// TokenResponse is response format for authentication token
type TokenResponse struct {
	Token string `json:"token"`
}

// NewSuccessResponse creates formatted success response
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates formatted error response
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}
