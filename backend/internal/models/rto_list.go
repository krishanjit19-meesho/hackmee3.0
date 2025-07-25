package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// RTOProduct represents a single product in the RTO list
type RTOProduct struct {
	RTOCount    int    `json:"rto_count"`
	CatalogID   int64  `json:"catalog_id"`
	OrderDate   string `json:"order_date"`
	ProductID   int64  `json:"product_id"`
	SubOrderNum string `json:"sub_order_num"`
}

// RTOProducts is a custom type to handle JSON serialization/deserialization
type RTOProducts []RTOProduct

// Value implements the driver.Valuer interface
func (r RTOProducts) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface
func (r *RTOProducts) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan non-string/non-byte value into RTOProducts")
	}

	return json.Unmarshal(bytes, r)
}

// RTOList represents the rto_list table structure
type RTOList struct {
	ID        int         `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string      `json:"code" gorm:"column:code;type:varchar(20);index"`
	Products  RTOProducts `json:"products" gorm:"column:products;type:json"`
	CreatedAt time.Time   `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for RTOList
func (RTOList) TableName() string {
	return "rto_list"
}
