package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vdt-dashboard-backend/models"
	"vdt-dashboard-backend/repositories"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthConfig holds Clerk configuration
type AuthConfig struct {
	SecretKey string
}

// AuthMiddleware handles Clerk JWT authentication using Clerk SDK
func AuthMiddleware(userRepo repositories.UserRepository, clerkSecretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Authorization header is required", models.ErrUnauthorized, "Missing Authorization header"))
			c.Abort()
			return
		}

		// Extract the token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Invalid authorization header format", models.ErrUnauthorized, "Use Bearer <token>"))
			c.Abort()
			return
		}

		sessionToken := parts[1]

		// Set Clerk API key
		clerk.SetKey(clerkSecretKey)

		// Verify the token using Clerk SDK v2
		ctx := context.Background()
		
		// First decode the token to get the key ID
		decoded, err := jwt.Decode(ctx, &jwt.DecodeParams{Token: sessionToken})
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Invalid token", models.ErrUnauthorized, err.Error()))
			c.Abort()
			return
		}

		// Fetch the JSON web key for verification
		jwk, err := jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
			KeyID: decoded.KeyID,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Failed to get JWT key", models.ErrUnauthorized, err.Error()))
			c.Abort()
			return
		}

		// Verify the token with the retrieved key
		claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
			Token: sessionToken,
			JWK:   jwk,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Invalid token", models.ErrUnauthorized, err.Error()))
			c.Abort()
			return
		}

		// Get user info from Clerk using the SDK
		clerkUser, err := user.Get(ctx, claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Failed to fetch user from Clerk", models.ErrUnauthorized, err.Error()))
			c.Abort()
			return
		}

		// Get or create user in our database
		user, err := getOrCreateUserFromClerk(userRepo, clerkUser, claims.Subject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to authenticate user", models.ErrInternalError, err.Error()))
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("clerkUserID", user.ClerkUserID)

		c.Next()
	}
}

// getOrCreateUserFromClerk retrieves or creates a user in our database based on Clerk user data
func getOrCreateUserFromClerk(userRepo repositories.UserRepository, clerkUser *clerk.User, clerkUserID string) (*models.User, error) {
	// Try to find existing user by Clerk ID
	user, err := userRepo.GetByClerkID(clerkUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Extract user info from Clerk user object
	var email, firstName, lastName, profileImageURL string
	
	// Get primary email
	if len(clerkUser.EmailAddresses) > 0 {
		for _, emailAddr := range clerkUser.EmailAddresses {
			// Handle the pointer comparison correctly
			if clerkUser.PrimaryEmailAddressID != nil && emailAddr.ID == *clerkUser.PrimaryEmailAddressID {
				email = emailAddr.EmailAddress
				break
			}
		}
		// Fallback to first email if primary not found
		if email == "" {
			email = clerkUser.EmailAddresses[0].EmailAddress
		}
	}

	// Get name info
	if clerkUser.FirstName != nil {
		firstName = *clerkUser.FirstName
	}
	if clerkUser.LastName != nil {
		lastName = *clerkUser.LastName
	}
	if clerkUser.ImageURL != nil {
		profileImageURL = *clerkUser.ImageURL
	}

	// If user doesn't exist, create a new one
	if err == gorm.ErrRecordNotFound {
		user = &models.User{
			ID:              uuid.New(),
			ClerkUserID:     clerkUserID,
			Email:           email,
			FirstName:       firstName,
			LastName:        lastName,
			ProfileImageURL: profileImageURL,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := userRepo.Create(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update existing user info
		user.Email = email
		user.FirstName = firstName
		user.LastName = lastName
		user.ProfileImageURL = profileImageURL
		user.UpdatedAt = time.Now()

		if err := userRepo.Update(user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return user, nil
}

// GetUserFromContext extracts the authenticated user from gin context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// GetUserIDFromContext extracts the authenticated user ID from gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	return userID.(uuid.UUID), true
} 