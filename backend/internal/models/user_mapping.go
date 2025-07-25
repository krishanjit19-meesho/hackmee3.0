package models

import (
	"time"
)

// UserMapping represents the user_mapping table structure
type UserMapping struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(50);index"`
	Code      string    `json:"code" gorm:"column:code;type:varchar(20);index"`
	City      string    `json:"city" gorm:"column:city;type:varchar(100);index"`
	State     string    `json:"state" gorm:"column:state;type:varchar(100)"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for UserMapping
func (UserMapping) TableName() string {
	return "user_mapping"
}
