package models

type BaseResponse struct {
	// shows that the api is live and was successfully called. The main action may still not succeed (eg a bad data)
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Errors    []string    `json:"errors,omitempty"`
	ErrorCode string      `json:"errorCode"`
}

// NewSuccessResponse creates a success response with data
func NewSuccessResponse(data interface{}, message string) BaseResponse {
	return BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// Todo: deprecate this!
// NewErrorResponse creates an error response
func NewErrorResponse(errors []string, message string, errorCode string) BaseResponse {
	return BaseResponse{
		Success:   true,
		Message:   message,
		Errors:    errors,
		ErrorCode: errorCode,
	}
}
func NewErrorOrFailureResponse(errors []string, message string, errorCode string, success bool) BaseResponse {
	return BaseResponse{
		Success:   success,
		Message:   message,
		Errors:    errors,
		ErrorCode: errorCode,
	}
}
