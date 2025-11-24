package models

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func NewErrorResponse(message, code string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Message: message,
			Code:    code,
			Details: details,
		},
	}
}

