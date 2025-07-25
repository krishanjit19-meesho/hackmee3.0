package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RTOService handles RTO-related operations
type RTOService struct {
	rtoAPIURL  string
	httpClient *http.Client
}

// RTOItem represents a single RTO item in the response
type RTOItem struct {
	CatalogID int64  `json:"catalog_id"`
	OrderDate string `json:"order_date"`
	ProductID int64  `json:"product_id"`
	RTOCount  int    `json:"rto_count"`
}

// RTOResponse represents the response from the RTO API
type RTOResponse struct {
	Code       string             `json:"code"`
	RTOList    map[string]RTOItem `json:"rto_list"`
	Success    bool               `json:"success"`
	TotalItems int                `json:"total_items"`
}

// NewRTOService creates a new RTO service
func NewRTOService() *RTOService {
	// Get RTO API URL from environment variable with fallback
	rtoAPIURL := os.Getenv("RTO_API_URL")
	if rtoAPIURL == "" {
		rtoAPIURL = "http://localhost:3001/rto/fetch"
	}

	return &RTOService{
		rtoAPIURL: rtoAPIURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCatalogIDsFromRTO fetches catalog IDs from RTO API based on user's code
func (s *RTOService) GetCatalogIDsFromRTO(userCode string) ([]string, error) {
	// Prepare the API URL
	apiURL := fmt.Sprintf("%s/%s", s.rtoAPIURL, userCode)

	fmt.Printf("Calling RTO API: %s\n", apiURL)

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RTO request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call RTO API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RTO response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RTO API returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var rtoResponse RTOResponse
	if err := json.Unmarshal(responseBody, &rtoResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RTO response: %w", err)
	}

	// Check if the response indicates success
	if !rtoResponse.Success {
		return nil, fmt.Errorf("RTO API returned error: success field is false")
	}

	// Extract catalog IDs from the RTO list
	var catalogIDs []string
	for _, rtoItem := range rtoResponse.RTOList {
		catalogIDs = append(catalogIDs, fmt.Sprintf("%d", rtoItem.CatalogID))
	}

	fmt.Printf("Successfully got %d catalog IDs from RTO API for code '%s'\n", len(catalogIDs), userCode)

	return catalogIDs, nil
}

// GetCatalogIDsFromRTOWithFallback calls the RTO API with fallback to empty list
func (s *RTOService) GetCatalogIDsFromRTOWithFallback(userCode string) []string {
	catalogIDs, err := s.GetCatalogIDsFromRTO(userCode)
	if err != nil {
		// Log the error but return empty list as fallback
		fmt.Printf("Warning: Failed to get catalog IDs from RTO API: %v. Using empty list.\n", err)
		return []string{}
	}

	// If RTO API returns empty result, return empty list
	if len(catalogIDs) == 0 {
		fmt.Printf("Warning: RTO API returned empty result. Using empty list.\n")
		return []string{}
	}

	return catalogIDs
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
