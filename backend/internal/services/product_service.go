package services

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"meesho-clone/configs"
	"meesho-clone/internal/models"

	"gorm.io/gorm"
)

// ProductService handles product-related operations
type ProductService struct {
	db *gorm.DB
}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{
		db: configs.DB,
	}
}

// GetProductDetails retrieves detailed product information
func (s *ProductService) GetProductDetails(productID, userID string) (*models.ProductDetailsResponse, error) {
	startTime := time.Now()

	// Try to get product details from product_info table
	productDetails, err := s.getProductDetailsFromDatabase(productID)
	if err == nil && productDetails != nil {
		// Add mock data for fields not in product_info table
		s.enrichProductDetails(productDetails, productID, userID)

		responseTime := time.Since(startTime).Milliseconds()
		return &models.ProductDetailsResponse{
			Success: true,
			Message: "Product details retrieved successfully",
			Data:    *productDetails,
			Meta: models.ProductMeta{
				ProductID:    productID,
				UserID:       userID,
				GeneratedAt:  time.Now(),
				Source:       "product_info_table",
				CacheHit:     false,
				ResponseTime: responseTime,
			},
		}, nil
	}

	log.Printf("Failed to get product details from database: %v", err)

	// Fall back to mock data
	mockProductDetails := s.generateMockProductDetails(productID, userID)
	responseTime := time.Since(startTime).Milliseconds()

	return &models.ProductDetailsResponse{
		Success: true,
		Message: "Product details retrieved successfully (mock data)",
		Data:    mockProductDetails,
		Meta: models.ProductMeta{
			ProductID:    productID,
			UserID:       userID,
			GeneratedAt:  time.Now(),
			Source:       "Mock Data",
			CacheHit:     false,
			ResponseTime: responseTime,
		},
	}, nil
}

// getProductDetailsFromDatabase retrieves product details from product_info table
func (s *ProductService) getProductDetailsFromDatabase(productID string) (*models.ProductDetails, error) {
	var productInfo models.ProductInfo

	// Query product_info table for the product
	result := s.db.Where("product_id = ?", productID).First(&productInfo)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query product_info table: %w", result.Error)
	}

	log.Printf("Found product in database: %s", productInfo.ProductID)

	// Convert ProductInfo to ProductDetails
	productDetails := s.convertProductInfoToProductDetails(productInfo)

	return &productDetails, nil
}

// convertProductInfoToProductDetails converts ProductInfo to ProductDetails
func (s *ProductService) convertProductInfoToProductDetails(productInfo models.ProductInfo) models.ProductDetails {
	// Generate image URL from the images field or catalog_id
	imageURL := s.generateImageURL(productInfo.CatalogID, productInfo.Images)

	// Use the name field as title, fallback to category + subcategory
	title := productInfo.Name
	if title == "" {
		title = fmt.Sprintf("%s - %s", productInfo.Category, productInfo.Sscat)
	}

	// Generate price based on product info
	price := s.generatePriceFromProductInfo(productInfo)

	// Generate basic images array
	images := []string{imageURL}

	return models.ProductDetails{
		ProductID:   productInfo.ProductID,
		CatalogID:   productInfo.CatalogID,
		Title:       title,
		Category:    productInfo.Category,
		SubCategory: productInfo.Sscat,
		Price:       price,
		MainImage:   imageURL,
		Images:      images,
		Brand:       productInfo.BrandName,
		Description: fmt.Sprintf("High-quality %s %s. Perfect for your needs with excellent durability and style.",
			strings.ToLower(productInfo.Category), strings.ToLower(productInfo.Sscat)),
	}
}

// generateImageURL generates image URL from catalog_id and images field
func (s *ProductService) generateImageURL(catalogID, images string) string {
	// If images field has data, use it
	if images != "" {
		// You might need to parse the images field based on your data format
		// For now, assuming it's a direct URL or needs to be prefixed
		if len(images) > 0 && images[:4] == "http" {
			return images
		}
		// If it's a relative path, prefix with Meesho CDN
		return fmt.Sprintf("https://images.meesho.com%s", images)
	}

	// Fallback to generating URL from catalog_id
	return fmt.Sprintf("https://images.meesho.com/images/products/%s/1_256.jpg", catalogID)
}

// generatePriceFromProductInfo generates price based on product info
func (s *ProductService) generatePriceFromProductInfo(productInfo models.ProductInfo) string {
	// You can implement price generation logic based on:
	// - Category (different categories have different price ranges)
	// - Brand (branded vs non-branded)
	// - Scale (large vs small scale)
	// - Other fields in the product_info table

	// For now, generating a basic price based on catalog_id
	// In a real implementation, you might want to add a price field to the table
	// or implement a pricing algorithm based on the available fields

	// Simple hash-based price generation
	price := 0
	for _, char := range productInfo.CatalogID {
		price += int(char)
	}
	price = (price % 4500) + 500 // Price between 500-5000

	return fmt.Sprintf("₹%d", price)
}

// enrichProductDetails adds mock data to enrich the product details
func (s *ProductService) enrichProductDetails(product *models.ProductDetails, productID, userID string) {
	// Generate random price data
	priceValue := rand.Intn(1000) + 100
	originalPrice := priceValue + rand.Intn(500) + 100
	discountPercent := int(float64(originalPrice-priceValue) / float64(originalPrice) * 100)

	product.OriginalPrice = fmt.Sprintf("₹%d", originalPrice)
	product.Discount = fmt.Sprintf("₹%d OFF", originalPrice-priceValue)
	product.DiscountPercent = discountPercent

	// Generate mock data for other fields
	if product.Description == "" {
		product.Description = fmt.Sprintf("High-quality %s %s. Perfect for your needs with excellent durability and style.",
			strings.ToLower(product.Category), strings.ToLower(product.SubCategory))
	}

	product.Rating = 3.5 + rand.Float64()*1.5 // 3.5 to 5.0
	product.Reviews = rand.Intn(10000) + 100
	product.Stock = rand.Intn(50) + 10
	if product.Brand == "" {
		product.Brand = "Meesho Brand"
	}
	product.Seller = "Meesho Seller"
	product.DeliveryInfo = "Free delivery by tomorrow"
	product.ReturnPolicy = "7 days return policy"
	product.Warranty = "1 year warranty"

	// Generate specifications
	product.Specifications = map[string]string{
		"Material": "Premium Quality",
		"Color":    "Multiple Options",
		"Size":     "Standard",
		"Weight":   "Lightweight",
		"Care":     "Easy to maintain",
	}

	// Generate variants
	product.Variants = []models.ProductVariant{
		{ID: "1", Name: "Size", Value: "Small", Price: product.Price, Stock: 15, Selected: true},
		{ID: "2", Name: "Size", Value: "Medium", Price: product.Price, Stock: 20, Selected: false},
		{ID: "3", Name: "Size", Value: "Large", Price: product.Price, Stock: 10, Selected: false},
	}

	// Generate reviews
	product.ReviewsList = s.generateMockReviews(productID)

	// Generate similar products
	product.SimilarProducts = s.generateSimilarProducts(productID)
}

// generateMockProductDetails creates complete mock product details
func (s *ProductService) generateMockProductDetails(productID, userID string) models.ProductDetails {
	// Generate random price data
	priceValue := rand.Intn(1000) + 100
	originalPrice := priceValue + rand.Intn(500) + 100
	discountPercent := int(float64(originalPrice-priceValue) / float64(originalPrice) * 100)

	// Generate random category and subcategory
	categories := []string{"Electronics", "Fashion", "Home", "Beauty", "Sports"}
	subCategories := []string{"Smartphones", "Clothing", "Furniture", "Skincare", "Fitness"}

	category := categories[rand.Intn(len(categories))]
	subCategory := subCategories[rand.Intn(len(subCategories))]

	// Generate mock images
	images := []string{
		"https://images.meesho.com/images/products/1234567/1_400.jpg",
		"https://images.meesho.com/images/products/1234567/2_400.jpg",
		"https://images.meesho.com/images/products/1234567/3_400.jpg",
	}

	return models.ProductDetails{
		ProductID: productID,
		CatalogID: fmt.Sprintf("CAT%d", rand.Intn(1000000)),
		Title:     fmt.Sprintf("Premium %s %s", category, subCategory),
		Description: fmt.Sprintf("High-quality %s %s with excellent features. Perfect for your daily needs with premium quality and durability.",
			strings.ToLower(category), strings.ToLower(subCategory)),
		Category:        category,
		SubCategory:     subCategory,
		Price:           fmt.Sprintf("₹%d", priceValue),
		OriginalPrice:   fmt.Sprintf("₹%d", originalPrice),
		Discount:        fmt.Sprintf("₹%d OFF", originalPrice-priceValue),
		DiscountPercent: discountPercent,
		Images:          images,
		MainImage:       images[0],
		Rating:          3.5 + rand.Float64()*1.5,
		Reviews:         rand.Intn(10000) + 100,
		Stock:           rand.Intn(50) + 10,
		Brand:           "Meesho Brand",
		Seller:          "Meesho Seller",
		DeliveryInfo:    "Free delivery by tomorrow",
		ReturnPolicy:    "7 days return policy",
		Warranty:        "1 year warranty",
		Specifications: map[string]string{
			"Material": "Premium Quality",
			"Color":    "Multiple Options",
			"Size":     "Standard",
			"Weight":   "Lightweight",
			"Care":     "Easy to maintain",
		},
		Variants: []models.ProductVariant{
			{ID: "1", Name: "Size", Value: "Small", Price: fmt.Sprintf("₹%d", priceValue), Stock: 15, Selected: true},
			{ID: "2", Name: "Size", Value: "Medium", Price: fmt.Sprintf("₹%d", priceValue), Stock: 20, Selected: false},
			{ID: "3", Name: "Size", Value: "Large", Price: fmt.Sprintf("₹%d", priceValue), Stock: 10, Selected: false},
		},
		ReviewsList:     s.generateMockReviews(productID),
		SimilarProducts: s.generateSimilarProducts(productID),
	}
}

// generateMockReviews creates mock product reviews
func (s *ProductService) generateMockReviews(productID string) []models.ProductReview {
	reviews := []models.ProductReview{}
	userNames := []string{"Rahul K.", "Priya S.", "Amit M.", "Neha P.", "Vikram R."}
	reviewTitles := []string{
		"Great product!",
		"Excellent quality",
		"Worth the money",
		"Good value for money",
		"Happy with purchase",
	}
	reviewComments := []string{
		"Really happy with this purchase. Quality is excellent!",
		"Good product, fast delivery. Would recommend!",
		"Value for money. Meets all expectations.",
		"Great quality and perfect fit. Very satisfied!",
		"Excellent product with good features.",
	}

	for i := 0; i < 5; i++ {
		review := models.ProductReview{
			ID:       fmt.Sprintf("review_%s_%d", productID, i+1),
			UserID:   fmt.Sprintf("user_%d", rand.Intn(1000)),
			UserName: userNames[rand.Intn(len(userNames))],
			Rating:   rand.Intn(3) + 3, // 3-5 stars
			Title:    reviewTitles[rand.Intn(len(reviewTitles))],
			Comment:  reviewComments[rand.Intn(len(reviewComments))],
			Date:     time.Now().AddDate(0, 0, -rand.Intn(30)),
			Verified: rand.Float64() > 0.3, // 70% verified
			Helpful:  rand.Intn(50),
		}
		reviews = append(reviews, review)
	}

	return reviews
}

// generateSimilarProducts creates mock similar products
func (s *ProductService) generateSimilarProducts(productID string) []models.SimilarProduct {
	similarProducts := []models.SimilarProduct{}

	for i := 0; i < 4; i++ {
		price := rand.Intn(1000) + 100
		similarProduct := models.SimilarProduct{
			ProductID: fmt.Sprintf("similar_%s_%d", productID, i+1),
			Title:     fmt.Sprintf("Similar Product %d", i+1),
			Image:     fmt.Sprintf("https://images.meesho.com/images/products/%d/1_400.jpg", rand.Intn(1000000)),
			Price:     fmt.Sprintf("₹%d", price),
			Rating:    3.0 + rand.Float64()*2.0,
			Reviews:   rand.Intn(5000) + 50,
		}
		similarProducts = append(similarProducts, similarProduct)
	}

	return similarProducts
}
