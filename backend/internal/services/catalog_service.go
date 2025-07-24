package services

import (
	"fmt"
	"log"
	"time"

	"meesho-clone/configs"
	"meesho-clone/internal/models"

	"gorm.io/gorm"
)

// CatalogService handles catalog-related operations
type CatalogService struct {
	db *gorm.DB
}

// NewCatalogService creates a new catalog service
func NewCatalogService() *CatalogService {
	return &CatalogService{
		db: configs.DB,
	}
}

// GetCatalogData fetches catalog data for a user
func (s *CatalogService) GetCatalogData(userID string) (*models.CatalogResponse, error) {
	// Step 1: Get catalog IDs (same list as before)
	catalogIDs := s.getCatalogIDs()

	// Step 2: Query product_info table for these catalog IDs
	catalogProducts, err := s.getProductInfoByCatalogIDs(catalogIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
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
			Source:        "product_info_table",
		},
	}

	return response, nil
}

// getCatalogIDs returns the list of catalog IDs (same as before)
func (s *CatalogService) getCatalogIDs() []string {
	return []string{
		"48", "108", "540", "553", "857",
		"887", "1611", "1805", "2373", "2607",
	}
}

// getProductInfoByCatalogIDs queries the product_info table for given catalog IDs
func (s *CatalogService) getProductInfoByCatalogIDs(catalogIDs []string) ([]models.CatalogProduct, error) {
	var productInfos []models.ProductInfo

	// Query product_info table for the given catalog IDs
	result := s.db.Where("catalog_id IN ?", catalogIDs).Find(&productInfos)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query product_info table: %w", result.Error)
	}

	log.Printf("Found %d products in product_info table for %d catalog IDs", len(productInfos), len(catalogIDs))

	// Convert ProductInfo to CatalogProduct
	var catalogProducts []models.CatalogProduct
	for _, productInfo := range productInfos {
		catalogProduct := s.convertProductInfoToCatalogProduct(productInfo)
		catalogProducts = append(catalogProducts, catalogProduct)
	}

	return catalogProducts, nil
}

// convertProductInfoToCatalogProduct converts ProductInfo to CatalogProduct
func (s *CatalogService) convertProductInfoToCatalogProduct(productInfo models.ProductInfo) models.CatalogProduct {
	// Generate image URL from the images field or catalog_id
	imageURL := s.generateImageURL(productInfo.CatalogID, productInfo.Images)

	// Use the name field as title, fallback to category + subcategory
	title := productInfo.Name
	if title == "" {
		title = fmt.Sprintf("%s - %s", productInfo.Category, productInfo.Sscat)
	}

	// Generate price (you might want to add a price field to the table later)
	price := s.generatePriceFromProductInfo(productInfo)

	return models.CatalogProduct{
		CatalogID:   productInfo.CatalogID,
		ProductID:   productInfo.ProductID,
		ImageURL:    imageURL,
		Category:    productInfo.Category,
		SubCategory: productInfo.Sscat,
		Title:       title,
		Price:       price,
	}
}

// generateImageURL generates image URL from catalog_id and images field
func (s *CatalogService) generateImageURL(catalogID, images string) string {
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
func (s *CatalogService) generatePriceFromProductInfo(productInfo models.ProductInfo) string {
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

// GetCatalogDataByIDs fetches catalog data for specific catalog IDs
func (s *CatalogService) GetCatalogDataByIDs(catalogIDs []string) (*models.CatalogResponse, error) {
	// Validate catalog IDs
	validIDs := s.ValidateCatalogIDs(catalogIDs)
	if len(validIDs) == 0 {
		return nil, fmt.Errorf("no valid catalog IDs provided")
	}

	// Query product_info table for the valid catalog IDs
	catalogProducts, err := s.getProductInfoByCatalogIDs(validIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
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
			Source:        "product_info_table_direct",
		},
	}

	return response, nil
}

// ValidateCatalogIDs validates the provided catalog IDs
func (s *CatalogService) ValidateCatalogIDs(catalogIDs []string) []string {
	var validIDs []string

	for _, id := range catalogIDs {
		// Basic validation: check if it's not empty and has reasonable length
		if id != "" && len(id) >= 4 && len(id) <= 20 {
			validIDs = append(validIDs, id)
		}
	}

	return validIDs
}

// TruncateCaption truncates text to specified length
func (s *CatalogService) TruncateCaption(caption string, length int) string {
	if len(caption) > length {
		return caption[:length] + "..."
	}
	return caption
}
