package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"meesho-clone/configs"
	"meesho-clone/internal/models"

	"gorm.io/gorm"
)

// UserService handles user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		db: configs.DB,
	}
}

// generateUserID generates a unique user ID
func (s *UserService) generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "user_" + hex.EncodeToString(bytes)[:12]
}

// CreateOrGetUser creates a new user or returns existing user by phone number
func (s *UserService) CreateOrGetUser(phoneNumber string) (*models.User, error) {
	var user models.User

	// Check if user already exists
	err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err == nil {
		// User exists, return existing user
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		// Database error
		return nil, fmt.Errorf("database error: %v", err)
	}

	// User doesn't exist, create new user
	newUser := models.User{
		UserID:      s.generateUserID(),
		PhoneNumber: phoneNumber,
		Name:        fmt.Sprintf("User %s", phoneNumber[len(phoneNumber)-4:]), // Use last 4 digits as name
	}

	err = s.db.Create(&newUser).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &newUser, nil
}

// GetUserByID retrieves user by user ID
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &user, nil
}

// GetUserByPhoneNumber retrieves user by phone number
func (s *UserService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID string, updates map[string]interface{}) error {
	err := s.db.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}
