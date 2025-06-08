package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware for handling errors consistently
func ErrorHandler() gin.HandlerFunc {
	return gin.ErrorLogger()
}

// HandleError is a utility function to handle errors in handlers
func HandleError(c *gin.Context, err error, message string, statusCode int) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error": gin.H{
			"code":    getErrorCode(statusCode),
			"details": err.Error(),
		},
	})
}

// HandleValidationError handles validation errors specifically
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "Validation failed",
		"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"details": err.Error(),
		},
	})
}

// getErrorCode returns appropriate error code based on HTTP status
func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}
