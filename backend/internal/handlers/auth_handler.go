package handlers

import (
	"meesho-clone/internal/models"
	"meesho-clone/internal/services"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(),
	}
}

// Login handles user login with phone number
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate phone number format (Indian mobile numbers)
	if !h.isValidIndianPhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid phone number format. Please enter a valid 10-digit Indian mobile number",
		})
		return
	}

	// Create or get user
	user, err := h.userService.CreateOrGetUser(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process login",
			"details": err.Error(),
		})
		return
	}

	// Prepare response
	response := models.LoginResponse{
		UserID:      user.UserID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		Message:     "Login successful",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetUserProfile retrieves user profile information
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateUserProfile updates user profile information
func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Remove fields that shouldn't be updated directly
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "phone_number")
	delete(updates, "created_at")

	err := h.userService.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user profile",
			"details": err.Error(),
		})
		return
	}

	// Get updated user data
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch updated user data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    user,
	})
}

// ValidateUser validates if a user exists
func (h *AuthHandler) ValidateUser(c *gin.Context) {
	userID := c.Query("user_id")
	phoneNumber := c.Query("phone_number")

	if userID == "" && phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Either user_id or phone_number is required",
		})
		return
	}

	var user *models.User
	var err error

	if userID != "" {
		user, err = h.userService.GetUserByID(userID)
	} else {
		user, err = h.userService.GetUserByPhoneNumber(phoneNumber)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"exists":  true,
		"data": gin.H{
			"user_id":      user.UserID,
			"phone_number": user.PhoneNumber,
			"name":         user.Name,
		},
	})
}

// isValidIndianPhoneNumber validates Indian mobile number format
func (h *AuthHandler) isValidIndianPhoneNumber(phone string) bool {
	// Indian mobile numbers: 10 digits starting with 6, 7, 8, or 9
	phoneRegex := regexp.MustCompile(`^[6-9][0-9]{9}$`)
	return phoneRegex.MatchString(phone)
}

// Health check endpoint
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "meesho-clone-auth",
		"timestamp": gin.H{
			"unix": c.GetInt64("timestamp"),
		},
	})
}
