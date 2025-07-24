package services

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"meesho-clone/internal/models"
)

// ProductService handles product-related operations
type ProductService struct {
	databricksService *DatabricksService
}

// NewProductService creates a new product service
func NewProductService(databricksService *DatabricksService) *ProductService {
	return &ProductService{
		databricksService: databricksService,
	}
}

// GetProductDetails retrieves detailed product information
func (s *ProductService) GetProductDetails(productID, userID string) (*models.ProductDetailsResponse, error) {
	startTime := time.Now()

	// Try to get product details from Databricks first
	if s.databricksService != nil {
		productDetails, err := s.getProductDetailsFromDatabricks(productID)
		if err == nil && productDetails != nil {
			// Add mock data for fields not in Databricks
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
					Source:       "Databricks",
					CacheHit:     false,
					ResponseTime: responseTime,
				},
			}, nil
		}
		log.Printf("Failed to get product details from Databricks: %v", err)
	}

	// Fall back to mock data
	productDetails := s.generateMockProductDetails(productID, userID)
	responseTime := time.Since(startTime).Milliseconds()

	return &models.ProductDetailsResponse{
		Success: true,
		Message: "Product details retrieved successfully (mock data)",
		Data:    productDetails,
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

// getProductDetailsFromDatabricks retrieves product details from Databricks
func (s *ProductService) getProductDetailsFromDatabricks(productID string) (*models.ProductDetails, error) {
	// Query Databricks for product details
	query := `
		SELECT 
			p.product_id,
			p.catalog_id,
			p.scat as category,
			p.sscat as sub_category,
			sp.image_url,
			COALESCE(ca.price_sscat_decile, '₹999') as price
		FROM gold.product_info p
		LEFT JOIN silver.supply__products sp ON p.product_id = sp.product_id
		LEFT JOIN catalog__attributes_agg ca ON p.catalog_id = ca.catalog_id
		WHERE p.product_id = ?
		LIMIT 1
	`

	rows, err := s.databricksService.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query Databricks: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var product models.ProductDetails
		var imageURL string
		var price string

		err := rows.Scan(
			&product.ProductID,
			&product.CatalogID,
			&product.Category,
			&product.SubCategory,
			&imageURL,
			&price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product data: %w", err)
		}

		// Set basic fields
		product.Title = fmt.Sprintf("%s - %s", product.Category, product.SubCategory)
		product.Price = price
		product.MainImage = fmt.Sprintf("https://images.meesho.com%s", imageURL)
		product.Images = []string{product.MainImage}

		return &product, nil
	}

	return nil, fmt.Errorf("product not found: %s", productID)
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
	product.Description = fmt.Sprintf("High-quality %s %s. Perfect for your needs with excellent durability and style.",
		strings.ToLower(product.Category), strings.ToLower(product.SubCategory))

	product.Rating = 3.5 + rand.Float64()*1.5 // 3.5 to 5.0
	product.Reviews = rand.Intn(10000) + 100
	product.Stock = rand.Intn(50) + 10
	product.Brand = "Meesho Brand"
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
