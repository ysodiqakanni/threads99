package models

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors,omitempty"`
}

// NewSuccessResponse creates a success response with data
func NewSuccessResponse(data interface{}, message string) BaseResponse {
	return BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(errors []string, message string) BaseResponse {
	return BaseResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
