package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      string         `json:"user_id" gorm:"uniqueIndex:idx_user_id,length:255;not null"`
	PhoneNumber string         `json:"phone_number" gorm:"uniqueIndex:idx_phone_number,length:20;not null"`
	Name        string         `json:"name"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=10"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Message     string `json:"message"`
}

// HomescreenRequest represents request for homescreen data
type HomescreenRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// MeeshoAPIResponse represents the response from Meesho API
type MeeshoAPIResponse struct {
	TopNavBar    TopNavBar     `json:"top_nav_bar"`
	WidgetGroups []WidgetGroup `json:"widget_groups"`
}

// TopNavBar represents the navigation structure
type TopNavBar struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Tiles []NavTile `json:"tiles"`
}

// NavTile represents individual navigation items
type NavTile struct {
	ID              int                    `json:"id"`
	Title           string                 `json:"title"`
	Image           string                 `json:"image"`
	NewImage        string                 `json:"new_image,omitempty"`
	DestinationData map[string]interface{} `json:"destination_data"`
}

// WidgetGroup represents grouped widgets/products
type WidgetGroup struct {
	ID              int      `json:"id"`
	MongoWGID       string   `json:"mongo_wg_id"`
	Position        int      `json:"position"`
	Title           string   `json:"title"`
	Tag             string   `json:"tag"`
	Type            int      `json:"type"`
	Widgets         []Widget `json:"widgets"`
	BackgroundColor string   `json:"background_color"`
	Dynamic         bool     `json:"dynamic"`
	DSEnabled       bool     `json:"ds_enabled"`
	AdsEnabled      bool     `json:"ads_enabled"`
	Priority        int      `json:"priority"`
	WidgetGroupType string   `json:"widget_group_type"`
}

// Widget represents individual product/content widgets
type Widget struct {
	ID               int                    `json:"id"`
	Title            string                 `json:"title"`
	Image            string                 `json:"image"`
	ImageAspectRatio float64                `json:"image_aspect_ratio"`
	Screen           string                 `json:"screen"`
	Type             int                    `json:"type"`
	DestinationID    int                    `json:"destination_id"`
	Data             map[string]interface{} `json:"data"`
	Fixed            bool                   `json:"fixed"`
	Priority         int                    `json:"priority"`
}
