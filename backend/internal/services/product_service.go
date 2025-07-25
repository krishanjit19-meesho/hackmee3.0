package services

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

// getProductDetailsFromDatabase retrieves product details from price_product_info table
func (s *ProductService) getProductDetailsFromDatabase(productID string) (*models.ProductDetails, error) {
	var priceProductInfo models.PriceProductInfo

	// Query price_product_info table for the product
	result := s.db.Where("product_id = ?", productID).First(&priceProductInfo)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query price_product_info table: %w", result.Error)
	}

	log.Printf("Found product in database: %s", priceProductInfo.ProductID)

	// Convert PriceProductInfo to ProductDetails
	productDetails := s.convertPriceProductInfoToProductDetails(priceProductInfo)

	return &productDetails, nil
}

// convertPriceProductInfoToProductDetails converts PriceProductInfo to ProductDetails
func (s *ProductService) convertPriceProductInfoToProductDetails(priceProductInfo models.PriceProductInfo) models.ProductDetails {
	// Generate image URL from the images field or catalog_id
	imageURL := s.generateImageURL(priceProductInfo.CatalogID, priceProductInfo.Images)

	// Use the name field for both title and description
	title := priceProductInfo.Name
	if title == "" {
		title = fmt.Sprintf("%s - %s", priceProductInfo.Category, priceProductInfo.Sscat)
	}

	// Use the name field as description as well, with fallback
	description := priceProductInfo.Name
	if description == "" {
		description = fmt.Sprintf("High-quality %s %s. Perfect for your needs with excellent durability and style.",
			strings.ToLower(priceProductInfo.Category), strings.ToLower(priceProductInfo.Sscat))
	}

	// Generate price and discount based on meesho_price_with_shipping and supplier_listed_price
	price, originalPrice, discount, discountPercent := s.generatePriceFromPriceProductInfo(priceProductInfo)

	// Generate images array with counter-based approach
	images := s.generateProductImages(priceProductInfo.ProductID, imageURL)

	return models.ProductDetails{
		ProductID:       priceProductInfo.ProductID,
		CatalogID:       priceProductInfo.CatalogID,
		Title:           title,
		Category:        priceProductInfo.Category,
		SubCategory:     priceProductInfo.Sscat,
		Price:           price,
		OriginalPrice:   originalPrice,
		Discount:        discount,
		DiscountPercent: discountPercent,
		MainImage:       imageURL,
		Images:          images,
		Brand:           priceProductInfo.BrandName,
		Description:     description,
	}
}

// generateImageURL generates image URL from catalog_id and images field
func (s *ProductService) generateImageURL(catalogID, images string) string {
	var imageURL string

	// If images field has data, use it
	if images != "" {
		// Split by comma to get individual image paths
		imagePaths := strings.Split(images, ",")
		if len(imagePaths) > 0 {
			// Take the first image as the main image
			mainImage := strings.TrimSpace(imagePaths[0])

			// Check if it's already a full URL
			if len(mainImage) > 0 && mainImage[:4] == "http" {
				imageURL = mainImage
			} else if len(mainImage) > 0 {
				// If it's a relative path, prefix with Meesho CDN
				imageURL = fmt.Sprintf("https://images.meesho.com%s", mainImage)
			}
		}
	}

	// If no image URL was set, fallback to generating URL from catalog_id
	if imageURL == "" {
		imageURL = fmt.Sprintf("https://images.meesho.com/images/products/%s/1_256.jpg", catalogID)
	}

	// Check if the generated image exists, if not use a fallback
	if !s.imageExists(imageURL) {
		log.Printf("Main image not found, using fallback: %s", imageURL)
		// Use a default fallback image
		imageURL = "https://images.meesho.com/images/products/default/1_256.jpg"
	}

	return imageURL
}

// generatePriceFromPriceProductInfo generates price and discount based on meesho_price_with_shipping and supplier_listed_price
func (s *ProductService) generatePriceFromPriceProductInfo(priceProductInfo models.PriceProductInfo) (string, string, string, int) {
	// Use supplier_listed_price as the actual price (what customers pay)
	actualPrice := priceProductInfo.SupplierListedPrice

	// Use meesho_price_with_shipping as the original price (strikethrough price)
	originalPrice := priceProductInfo.MeeshoPriceWithShipping

	// Calculate discount amount
	discountAmount := originalPrice - actualPrice

	// Calculate discount percentage
	var discountPercent int
	if originalPrice > 0 {
		discountPercent = int((discountAmount / originalPrice) * 100)
	}

	// Format prices as strings
	priceStr := fmt.Sprintf("₹%.0f", actualPrice)
	originalPriceStr := fmt.Sprintf("₹%.0f", originalPrice)
	discountStr := fmt.Sprintf("₹%.0f OFF", discountAmount)

	return priceStr, originalPriceStr, discountStr, discountPercent
}

// generateProductImages generates images array with counter-based approach
func (s *ProductService) generateProductImages(productID, mainImage string) []string {
	images := []string{mainImage}

	// Generate additional images from 1-4 and check if they exist concurrently
	imageURLs := make([]string, 4)
	for i := 1; i <= 4; i++ {
		imageURLs[i-1] = fmt.Sprintf("https://images.meesho.com/images/products/%s/%d_256.jpg", productID, i)
	}

	// Check images concurrently for better performance
	existingImages := s.checkImagesConcurrently(imageURLs)
	images = append(images, existingImages...)

	return images
}

// checkImagesConcurrently checks multiple images concurrently
func (s *ProductService) checkImagesConcurrently(imageURLs []string) []string {
	results := make(chan string, len(imageURLs))

	// Start goroutines for each image check
	for _, url := range imageURLs {
		go func(imageURL string) {
			if s.imageExists(imageURL) {
				results <- imageURL
			} else {
				results <- "" // Empty string for non-existent images
			}
		}(url)
	}

	// Collect results with timeout
	var existingImages []string
	timeout := time.After(10 * time.Second) // 10 second timeout for all image checks

	for i := 0; i < len(imageURLs); i++ {
		select {
		case result := <-results:
			if result != "" {
				existingImages = append(existingImages, result)
			}
		case <-timeout:
			log.Printf("Timeout reached while checking images for product")
			return existingImages
		}
	}

	return existingImages
}

// imageExists checks if an image URL exists using HTTP HEAD request
func (s *ProductService) imageExists(imageURL string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second, // 5 second timeout
	}

	req, err := http.NewRequest("HEAD", imageURL, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", imageURL, err)
		return false
	}

	// Set user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error checking image %s: %v", imageURL, err)
		return false
	}
	defer resp.Body.Close()

	exists := resp.StatusCode >= 200 && resp.StatusCode < 300
	return exists
}

// enrichProductDetails adds mock data to enrich the product details
func (s *ProductService) enrichProductDetails(product *models.ProductDetails, productID, userID string) {
	// Note: Pricing data is now handled by the database, so we don't override it here

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

	// Generate a mock name that will be used for both title and description
	mockName := fmt.Sprintf("Premium %s %s with excellent features", category, subCategory)

	// Generate mock images with counter-based approach and existence check
	mainImage := "https://images.meesho.com/images/products/1234567/1_256.jpg"
	images := s.generateProductImages("1234567", mainImage)

	return models.ProductDetails{
		ProductID:       productID,
		CatalogID:       fmt.Sprintf("CAT%d", rand.Intn(1000000)),
		Title:           mockName,
		Description:     mockName,
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
		ReviewsList: s.generateMockReviews(productID),
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
