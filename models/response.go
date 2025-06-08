package models

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents error information in API responses
type APIError struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Error      *APIError           `json:"error,omitempty"`
}

// SuccessResponse creates a successful API response
func SuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(message string, code string, details string) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    code,
			Details: details,
		},
	}
}

// PaginatedSuccessResponse creates a successful paginated API response
func PaginatedSuccessResponse(message string, data interface{}, pagination *PaginationResponse) *PaginatedResponse {
	return &PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

// Error codes constants
const (
	ErrValidation             = "VALIDATION_ERROR"
	ErrSchemaNotFound         = "SCHEMA_NOT_FOUND"
	ErrDatabaseError          = "DATABASE_ERROR"
	ErrDuplicateName          = "DUPLICATE_NAME"
	ErrInvalidJSON            = "INVALID_JSON"
	ErrMissingRequiredField   = "MISSING_REQUIRED_FIELD"
	ErrUnsupportedDataType    = "UNSUPPORTED_DATA_TYPE"
	ErrForeignKeyError        = "FOREIGN_KEY_ERROR"
	ErrDatabaseCreationFailed = "DATABASE_CREATION_FAILED"
	ErrInternalError          = "INTERNAL_ERROR"
)
