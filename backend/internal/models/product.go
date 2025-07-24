package models

import "time"

// ProductDetailsRequest represents the request for product details
type ProductDetailsRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

// ProductDetailsResponse represents the response for product details
type ProductDetailsResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    ProductDetails `json:"data"`
	Meta    ProductMeta    `json:"meta"`
}

// ProductDetails represents detailed product information
type ProductDetails struct {
	ProductID       string            `json:"product_id"`
	CatalogID       string            `json:"catalog_id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Category        string            `json:"category"`
	SubCategory     string            `json:"sub_category"`
	Price           string            `json:"price"`
	OriginalPrice   string            `json:"original_price"`
	Discount        string            `json:"discount"`
	DiscountPercent int               `json:"discount_percent"`
	Images          []string          `json:"images"`
	MainImage       string            `json:"main_image"`
	Rating          float64           `json:"rating"`
	Reviews         int               `json:"reviews"`
	Stock           int               `json:"stock"`
	Brand           string            `json:"brand"`
	Seller          string            `json:"seller"`
	DeliveryInfo    string            `json:"delivery_info"`
	ReturnPolicy    string            `json:"return_policy"`
	Warranty        string            `json:"warranty"`
	Specifications  map[string]string `json:"specifications"`
	Variants        []ProductVariant  `json:"variants"`
	ReviewsList     []ProductReview   `json:"reviews_list"`
	SimilarProducts []SimilarProduct  `json:"similar_products"`
}

// ProductVariant represents product variants (size, color, etc.)
type ProductVariant struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Price    string `json:"price"`
	Stock    int    `json:"stock"`
	Selected bool   `json:"selected"`
}

// ProductReview represents a product review
type ProductReview struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	UserName string    `json:"user_name"`
	Rating   int       `json:"rating"`
	Title    string    `json:"title"`
	Comment  string    `json:"comment"`
	Date     time.Time `json:"date"`
	Verified bool      `json:"verified"`
	Helpful  int       `json:"helpful"`
}

// SimilarProduct represents a similar product recommendation
type SimilarProduct struct {
	ProductID string  `json:"product_id"`
	Title     string  `json:"title"`
	Image     string  `json:"image"`
	Price     string  `json:"price"`
	Rating    float64 `json:"rating"`
	Reviews   int     `json:"reviews"`
}

// ProductMeta represents metadata for product details
type ProductMeta struct {
	ProductID    string    `json:"product_id"`
	UserID       string    `json:"user_id"`
	GeneratedAt  time.Time `json:"generated_at"`
	Source       string    `json:"source"`
	CacheHit     bool      `json:"cache_hit"`
	ResponseTime int64     `json:"response_time_ms"`
}
