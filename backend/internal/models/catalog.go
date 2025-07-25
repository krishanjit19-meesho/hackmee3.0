package models

import (
	"time"
)

// CatalogRequest represents the request for catalog data
type CatalogRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// CatalogResponse represents the response from catalog API
type CatalogResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []CatalogProduct `json:"data"`
	Meta    CatalogMeta      `json:"meta"`
}

// CatalogProduct represents a product in the catalog
type CatalogProduct struct {
	CatalogID       string `json:"catalog_id"`
	ProductID       string `json:"product_id"`
	ImageURL        string `json:"image_url"`
	Category        string `json:"category"`
	SubCategory     string `json:"sub_category"`
	Title           string `json:"title"`
	Price           string `json:"price,omitempty"`
	OriginalPrice   string `json:"original_price,omitempty"`
	Discount        string `json:"discount,omitempty"`
	DiscountPercent int    `json:"discount_percent,omitempty"`
}

// CatalogMeta represents metadata about the catalog response
type CatalogMeta struct {
	TotalProducts int       `json:"total_products"`
	UserID        string    `json:"user_id"`
	GeneratedAt   time.Time `json:"generated_at"`
	Source        string    `json:"source"`
}

// ExternalCatalogAPIRequest represents request to external catalog service
type ExternalCatalogAPIRequest struct {
	UserID string `json:"user_id"`
}

// ExternalCatalogAPIResponse represents response from external catalog service
type ExternalCatalogAPIResponse struct {
	Success    bool     `json:"success"`
	CatalogIDs []string `json:"catalog_ids"`
	Message    string   `json:"message"`
}

// MockCatalogData represents mock catalog data structure
type MockCatalogData struct {
	CatalogID   string `json:"catalog_id"`
	ProductID   string `json:"product_id"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
}
