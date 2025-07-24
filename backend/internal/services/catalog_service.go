package services

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"meesho-clone/internal/models"
)

// CatalogService handles catalog-related operations
type CatalogService struct {
	httpClient        *http.Client
	databricksService *DatabricksService
}

// NewCatalogService creates a new catalog service
func NewCatalogService() *CatalogService {
	// Initialize Databricks service (primary)
	databricksService, err := NewDatabricksService()
	if err != nil {
		log.Printf("Warning: Failed to initialize Databricks service: %v", err)
		log.Printf("Falling back to mock data mode")
		databricksService = nil
	} else {
		// Test the connection
		if err := databricksService.TestConnection(); err != nil {
			log.Printf("Warning: Databricks connection test failed: %v", err)
			databricksService = nil
		} else {
			log.Printf("✅ Successfully connected to Databricks SQL")
		}
	}

	return &CatalogService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		databricksService: databricksService,
	}
}

// GetCatalogData fetches catalog data for a user
func (s *CatalogService) GetCatalogData(userID string) (*models.CatalogResponse, error) {
	// Step 1: Call external API to get catalog IDs
	catalogIDs, err := s.fetchCatalogIDsFromExternalAPI(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog IDs: %w", err)
	}

	// Step 2: Process catalog data (similar to Python script)
	catalogProducts, err := s.processCatalogData(catalogIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to process catalog data: %w", err)
	}

	// Step 3: Create response
	response := &models.CatalogResponse{
		Success: true,
		Message: "Catalog data retrieved successfully",
		Data:    catalogProducts,
		Meta: models.CatalogMeta{
			TotalProducts: len(catalogProducts),
			UserID:        userID,
			GeneratedAt:   time.Now(),
			Source:        "catalog_service",
		},
	}

	return response, nil
}

// fetchCatalogIDsFromExternalAPI calls the external API to get catalog IDs
func (s *CatalogService) fetchCatalogIDsFromExternalAPI(userID string) ([]string, error) {
	// For now, we'll mock the external API response
	// In a real implementation, this would make an HTTP call to the external service

	// Mock response based on user ID
	mockResponse := s.generateMockCatalogIDs(userID)

	return mockResponse.CatalogIDs, nil
}

// generateMockCatalogIDs generates mock catalog IDs based on user ID
func (s *CatalogService) generateMockCatalogIDs(userID string) *models.ExternalCatalogAPIResponse {
	catalogIDs := []string{
		"1119458", "3461848", "4663520", "1284242", "3582567",
		"5923065", "1822776", "6459587", "1485802", "2386881",
		"2392204", "3372467", "3372251", "6286760", "5476755",
		"2310673", "3467450", "1118824", "288006", "2384938",
		"702697", "1811163", "1117063", "3303733", "4814236",
		"6210856", "1356739", "2550124", "6007551", "948590",
		"3791774",
	}

	return &models.ExternalCatalogAPIResponse{
		Success:    true,
		CatalogIDs: catalogIDs,
		Message:    "Hardcoded catalog IDs returned successfully",
	}
}

// processCatalogData processes catalog IDs into product data using real Databricks queries
func (s *CatalogService) processCatalogData(catalogIDs []string) ([]models.CatalogProduct, error) {
	// If Databricks service is available, use real Databricks queries
	if s.databricksService != nil {
		// Try the complex query first (original Python script logic)
		products, err := s.databricksService.ExecuteComplexQuery(catalogIDs)
		if err == nil && len(products) > 0 {
			return products, nil
		}

		// Fall back to simple query if complex query fails
		products, err = s.databricksService.GetCatalogDataByIDs(catalogIDs)
		if err == nil && len(products) > 0 {
			return products, nil
		}

		log.Printf("Databricks queries failed, falling back to mock data: %v", err)
	}

	// Fall back to mock data if Databricks is not available or queries fail
	return s.processCatalogDataWithMockData(catalogIDs)
}

// processCatalogDataWithMockData processes catalog IDs using mock data (fallback)
func (s *CatalogService) processCatalogDataWithMockData(catalogIDs []string) ([]models.CatalogProduct, error) {
	var products []models.CatalogProduct

	for _, catalogID := range catalogIDs {
		// Generate product ID (similar to catalog ID)
		productID := s.generateProductIDFromCatalog(catalogID)

		// Get image URL (from silver.supply__products equivalent)
		imageURL := s.generateImageURL(catalogID)

		// Get taxonomy data (from gold.product_info equivalent)
		category, subCategory := s.getTaxonomyData(catalogID)

		// Get price data (from catalog__attributes_agg equivalent)
		price := s.generatePriceData(catalogID)

		// Create product (this is what the original query was returning)
		product := models.CatalogProduct{
			CatalogID:   catalogID,
			ProductID:   productID,
			ImageURL:    imageURL,
			Category:    category,
			SubCategory: subCategory,
			Title:       fmt.Sprintf("%s - %s", category, subCategory),
			Price:       price,
		}

		products = append(products, product)
	}

	return products, nil
}

// generateProductIDFromCatalog simulates the product_id generation from the SQL query
func (s *CatalogService) generateProductIDFromCatalog(catalogID string) string {
	// In the original query, this was getting the top product by orders for each catalog
	// We'll generate a deterministic product ID based on catalog ID
	seed := int64(0)
	for _, char := range catalogID {
		seed += int64(char)
	}
	rand.Seed(seed)

	// Generate a product ID that's related to the catalog ID
	productID := rand.Intn(9000000) + 1000000
	return strconv.Itoa(productID)
}

// generateImageURL simulates the image URL generation from silver.supply__products
func (s *CatalogService) generateImageURL(catalogID string) string {
	// Original query was: concat(lit("https://images.meesho.com"), col("image_url"))
	// This is the actual Meesho CDN format
	return fmt.Sprintf("https://images.meesho.com/images/products/%s/1_256.jpg", catalogID)
}

// getTaxonomyData simulates the taxonomy data from gold.product_info
func (s *CatalogService) getTaxonomyData(catalogID string) (string, string) {
	// Original query was getting scat and sscat from gold.product_info
	// We'll use deterministic categories based on catalog ID

	// Categories and subcategories (realistic Meesho categories)
	categories := []string{"Electronics", "Fashion", "Home & Living", "Beauty", "Sports", "Books", "Toys", "Automotive"}
	subCategories := map[string][]string{
		"Electronics":   {"Smartphones", "Laptops", "Accessories", "Gaming"},
		"Fashion":       {"Men's Clothing", "Women's Clothing", "Kids' Wear", "Footwear"},
		"Home & Living": {"Furniture", "Kitchen", "Decor", "Bedding"},
		"Beauty":        {"Skincare", "Makeup", "Hair Care", "Fragrances"},
		"Sports":        {"Fitness", "Outdoor", "Team Sports", "Yoga"},
		"Books":         {"Fiction", "Non-Fiction", "Academic", "Children's"},
		"Toys":          {"Educational", "Action Figures", "Board Games", "Puzzles"},
		"Automotive":    {"Car Accessories", "Motorcycle", "Tools", "Maintenance"},
	}

	// Use catalog ID to deterministically select category
	seed := int64(0)
	for _, char := range catalogID {
		seed += int64(char)
	}
	rand.Seed(seed)

	category := categories[seed%int64(len(categories))]
	subCategoryList := subCategories[category]
	subCategory := subCategoryList[seed%int64(len(subCategoryList))]

	return category, subCategory
}

// generatePriceData simulates the price data from catalog__attributes_agg
func (s *CatalogService) generatePriceData(catalogID string) string {
	// Original query was getting price_sscat_decile from catalog__attributes_agg
	// We'll generate realistic price ranges based on catalog ID

	seed := int64(0)
	for _, char := range catalogID {
		seed += int64(char)
	}
	rand.Seed(seed)

	// Generate price between 500-5000 based on catalog ID
	price := rand.Intn(4500) + 500
	return fmt.Sprintf("₹%d", price)
}

// GetCatalogDataByIDs fetches catalog data for specific catalog IDs
func (s *CatalogService) GetCatalogDataByIDs(catalogIDs []string) (*models.CatalogResponse, error) {
	// Process catalog data for specific IDs
	catalogProducts, err := s.processCatalogData(catalogIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to process catalog data: %w", err)
	}

	// Create response
	response := &models.CatalogResponse{
		Success: true,
		Message: "Catalog data retrieved successfully",
		Data:    catalogProducts,
		Meta: models.CatalogMeta{
			TotalProducts: len(catalogProducts),
			UserID:        "direct_request",
			GeneratedAt:   time.Now(),
			Source:        "catalog_service_direct",
		},
	}

	return response, nil
}

// ValidateCatalogIDs validates the provided catalog IDs
func (s *CatalogService) ValidateCatalogIDs(catalogIDs []string) []string {
	var validIDs []string

	for _, id := range catalogIDs {
		// Basic validation: check if it's a numeric string with 6-8 digits
		if len(id) >= 6 && len(id) <= 8 {
			if _, err := strconv.Atoi(id); err == nil {
				validIDs = append(validIDs, id)
			}
		}
	}

	return validIDs
}

// TruncateCaption truncates text to specified length (similar to Python function)
func (s *CatalogService) TruncateCaption(caption string, length int) string {
	if len(caption) > length {
		return caption[:length] + "..."
	}
	return caption
}
