package services

import (
	"fmt"
	"strings"
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
	var catalogIDs []string
	var source string

	// Step 1: Try to get user's code from user_mapping table
	userCode, err := s.getUserCode(userID)
	if err != nil {
		fmt.Printf("Warning: Failed to get user code: %v. Using fallback catalog IDs.\n", err)
		catalogIDs = s.getCatalogIDs()
		source = "fallback_catalog_ids_with_ranking"
	} else {
		// Step 2: Try to get catalog IDs from RTO API based on user's code
		rtoService := NewRTOService()
		catalogIDs = rtoService.GetCatalogIDsFromRTOWithFallback(userCode)

		if len(catalogIDs) == 0 {
			fmt.Printf("Warning: No catalog IDs from RTO API. Using fallback catalog IDs.\n")
			catalogIDs = s.getCatalogIDs()
			source = "fallback_catalog_ids_with_ranking"
		} else {
			source = "rto_api_with_ranking"
		}
	}

	// Step 3: Get ranked catalog IDs from ranking service
	rankingService := NewRankingService()
	rankedCatalogIDs := rankingService.GetRankedCatalogIDsWithFallback(catalogIDs, userID)

	// Step 4: Query price_product_info table for the ranked catalog IDs
	catalogProducts, err := s.getProductInfoByCatalogIDs(rankedCatalogIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	// Step 4.5: Sort catalog products based on the ranked order
	sortedCatalogProducts := s.sortCatalogProductsByRanking(catalogProducts, rankedCatalogIDs)

	// Step 5: Create response
	response := &models.CatalogResponse{
		Success: true,
		Message: "Catalog data retrieved successfully",
		Data:    sortedCatalogProducts,
		Meta: models.CatalogMeta{
			TotalProducts: len(sortedCatalogProducts),
			UserID:        userID,
			GeneratedAt:   time.Now(),
			Source:        source,
		},
	}

	return response, nil
}

// getUserCode fetches the user's code from user_mapping table
func (s *CatalogService) getUserCode(userID string) (string, error) {
	var userMapping models.UserMapping

	result := s.db.Where("user_id = ?", userID).First(&userMapping)
	if result.Error != nil {
		return "", fmt.Errorf("failed to find user mapping for user_id %s: %w", userID, result.Error)
	}

	if userMapping.Code == "" {
		return "", fmt.Errorf("no code found for user_id %s", userID)
	}

	fmt.Printf("Found code '%s' for user %s\n", userMapping.Code, userID)
	return userMapping.Code, nil
}

// getCatalogIDs returns the list of top 100 catalog IDs (fallback method)
func (s *CatalogService) getCatalogIDs() []string {
	return []string{
		"647628", "3446585", "279116", "1610750", "1929554",
		"1222068", "3948449", "329972", "3310138", "381394",
		"258058", "2566266", "874707", "794925", "340768",
		"588251", "1918153", "1235372", "173571", "339856",
		"2764020", "2356994", "4057058", "313798", "373563",
		"276440", "89085", "271993", "1426675", "689640",
		"2597747", "3224910", "582131", "721729", "1701047",
		"2489323", "692563", "2309360", "3764919", "535354",
		"278790", "3011754", "633737", "454653", "491521",
		"508493", "600594", "505419", "505884", "933930",
		"340457", "1161355", "939873", "1034004", "939873",
		"251580", "2260285", "998799", "1291857", "1318580",
		"1087236", "205522", "152045", "2107126", "1039513",
		"3546354", "467602", "368361", "3943714", "3873515",
		"1311926", "2811468", "870959", "2321705", "642439",
		"79349", "999423", "506203", "713681", "1606309",
		"636532", "57790", "761239", "756322", "781576",
		"801433", "754541", "170989", "801666", "800622",
		"801161", "534900", "4106562", "83457", "1822307",
		"536092", "564252", "3195474", "29187", "4114942",
	}
}

// getProductInfoByCatalogIDs queries the price_product_info table for given catalog IDs
func (s *CatalogService) getProductInfoByCatalogIDs(catalogIDs []string) ([]models.CatalogProduct, error) {
	var priceProductInfos []models.PriceProductInfo

	// Debug logging
	fmt.Printf("Querying price_product_info table for %d catalog IDs\n", len(catalogIDs))
	
	// Query price_product_info table for the given catalog IDs
	result := s.db.Where("catalog_id IN ?", catalogIDs).Find(&priceProductInfos)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to query price_product_info table: %w", result.Error)
	}

	fmt.Printf("Found %d products in price_product_info table for the given catalog IDs\n", len(priceProductInfos))

	// Convert PriceProductInfo to CatalogProduct
	var catalogProducts []models.CatalogProduct
	for _, priceProductInfo := range priceProductInfos {
		catalogProduct := s.convertPriceProductInfoToCatalogProduct(priceProductInfo)
		catalogProducts = append(catalogProducts, catalogProduct)
	}

	return catalogProducts, nil
}

// convertPriceProductInfoToCatalogProduct converts PriceProductInfo to CatalogProduct
func (s *CatalogService) convertPriceProductInfoToCatalogProduct(priceProductInfo models.PriceProductInfo) models.CatalogProduct {
	// Generate image URL from the images field or catalog_id
	imageURL := s.generateImageURL(priceProductInfo.CatalogID, priceProductInfo.Images)

	// Use the name field as title, fallback to category + subcategory
	title := priceProductInfo.Name
	if title == "" {
		title = fmt.Sprintf("%s - %s", priceProductInfo.Category, priceProductInfo.Sscat)
	}

	// Generate price and discount based on meesho_price_with_shipping and supplier_listed_price
	price, originalPrice, discount, discountPercent := s.generatePriceFromPriceProductInfo(priceProductInfo)

	return models.CatalogProduct{
		CatalogID:       priceProductInfo.CatalogID,
		ProductID:       priceProductInfo.ProductID,
		ImageURL:        imageURL,
		Category:        priceProductInfo.Category,
		SubCategory:     priceProductInfo.Sscat,
		Title:           title,
		Price:           price,
		OriginalPrice:   originalPrice,
		Discount:        discount,
		DiscountPercent: discountPercent,
	}
}

// generateImageURL generates image URL from catalog_id and images field
func (s *CatalogService) generateImageURL(catalogID, images string) string {
	// If images field has data, use it
	if images != "" {
		// Split by comma to get individual image paths
		imagePaths := strings.Split(images, ",")
		if len(imagePaths) > 0 {
			// Take the first image as the main image
			mainImage := strings.TrimSpace(imagePaths[0])

			// Check if it's already a full URL
			if len(mainImage) > 0 && mainImage[:4] == "http" {
				return mainImage
			}

			// If it's a relative path, prefix with Meesho CDN
			if len(mainImage) > 0 {
				return fmt.Sprintf("https://images.meesho.com%s", mainImage)
			}
		}
	}

	// Fallback to generating URL from catalog_id
	return fmt.Sprintf("https://images.meesho.com/images/products/%s/1_256.jpg", catalogID)
}

// generatePriceFromPriceProductInfo generates price and discount based on meesho_price_with_shipping and supplier_listed_price
func (s *CatalogService) generatePriceFromPriceProductInfo(priceProductInfo models.PriceProductInfo) (string, string, string, int) {
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

// GetCatalogDataByIDs fetches catalog data for specific catalog IDs
func (s *CatalogService) GetCatalogDataByIDs(catalogIDs []string, userID string) (*models.CatalogResponse, error) {
	// Validate catalog IDs
	validIDs := s.ValidateCatalogIDs(catalogIDs)
	if len(validIDs) == 0 {
		return nil, fmt.Errorf("no valid catalog IDs provided")
	}

	// Get ranked catalog IDs from ranking service
	rankingService := NewRankingService()
	rankedCatalogIDs := rankingService.GetRankedCatalogIDsWithFallback(validIDs, userID)

	// Query price_product_info table for the ranked catalog IDs
	catalogProducts, err := s.getProductInfoByCatalogIDs(rankedCatalogIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	// Sort catalog products based on the ranked order
	sortedCatalogProducts := s.sortCatalogProductsByRanking(catalogProducts, rankedCatalogIDs)

	// Create response
	response := &models.CatalogResponse{
		Success: true,
		Message: "Catalog data retrieved successfully",
		Data:    sortedCatalogProducts,
		Meta: models.CatalogMeta{
			TotalProducts: len(sortedCatalogProducts),
			UserID:        userID,
			GeneratedAt:   time.Now(),
			Source:        "price_product_info_table_direct_with_ranking",
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

// sortCatalogProductsByRanking sorts catalog products based on the ranked catalog IDs order
func (s *CatalogService) sortCatalogProductsByRanking(catalogProducts []models.CatalogProduct, rankedCatalogIDs []string) []models.CatalogProduct {
	// Create a map for quick lookup of catalog ID positions
	rankedPositions := make(map[string]int)
	for i, catalogID := range rankedCatalogIDs {
		rankedPositions[catalogID] = i
	}

	// Create a copy of the catalog products to sort
	sortedProducts := make([]models.CatalogProduct, len(catalogProducts))
	copy(sortedProducts, catalogProducts)

	// Sort the products based on their position in the ranked catalog IDs
	// Use bubble sort for simplicity (can be optimized with sort.Slice if needed)
	for i := 0; i < len(sortedProducts)-1; i++ {
		for j := i + 1; j < len(sortedProducts); j++ {
			// Get positions of current and next products
			posI := len(rankedCatalogIDs) // Default high position for products not in ranking
			posJ := len(rankedCatalogIDs) // Default high position for products not in ranking

			if position, exists := rankedPositions[sortedProducts[i].CatalogID]; exists {
				posI = position
			}
			if position, exists := rankedPositions[sortedProducts[j].CatalogID]; exists {
				posJ = position
			}

			// Swap if current product has higher position (lower rank)
			if posI > posJ {
				sortedProducts[i], sortedProducts[j] = sortedProducts[j], sortedProducts[i]
			}
		}
	}

	fmt.Printf("Sorted %d catalog products based on ranking order\n", len(sortedProducts))
	return sortedProducts
}

// TruncateCaption truncates text to specified length
func (s *CatalogService) TruncateCaption(caption string, length int) string {
	if len(caption) > length {
		return caption[:length] + "..."
	}
	return caption
}
