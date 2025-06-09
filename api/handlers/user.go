package handlers

import (
	"net/http"

	"vdt-dashboard-backend/api/middleware"
	"vdt-dashboard-backend/models"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct{}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetCurrentUser handles GET /user/me
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	// Return user info (excluding sensitive data)
	userResponse := gin.H{
		"id":              user.ID,
		"clerkUserId":     user.ClerkUserID,
		"email":           user.Email,
		"firstName":       user.FirstName,
		"lastName":        user.LastName,
		"profileImageUrl": user.ProfileImageURL,
		"fullName":        user.GetFullName(),
		"createdAt":       user.CreatedAt,
		"updatedAt":       user.UpdatedAt,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User retrieved successfully", userResponse))
} 