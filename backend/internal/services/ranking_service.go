package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RankingService handles ranking-related operations
type RankingService struct {
	rankingAPIURL string
	httpClient    *http.Client
}

// RankingRequest represents the request to the ranking API
type RankingRequest struct {
	CatalogIDs []string `json:"catalog_ids"`
	UserID     string   `json:"user_id"`
}

// RankedCatalog represents a single ranked catalog item
type RankedCatalog struct {
	CatalogID string  `json:"catalog_id"`
	PctrScore float64 `json:"pctr_score"`
}

// RankingResponse represents the response from the ranking API
type RankingResponse struct {
	Success        bool            `json:"success"`
	RankedCatalogs []RankedCatalog `json:"ranked_catalogs"`
	TotalCatalogs  int             `json:"total_catalogs"`
	RawResponse    interface{}     `json:"raw_response,omitempty"`
}

// NewRankingService creates a new ranking service
func NewRankingService() *RankingService {
	// Get ranking API URL from environment variable with fallback
	rankingAPIURL := os.Getenv("RANKING_API_URL")
	if rankingAPIURL == "" {
		rankingAPIURL = "http://localhost:3000/rank"
	}

	return &RankingService{
		rankingAPIURL: rankingAPIURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRankedCatalogIDs calls the ranking API to get personalized catalog IDs
func (s *RankingService) GetRankedCatalogIDs(catalogIDs []string, userID string) ([]string, error) {
	// Prepare the request
	request := RankingRequest{
		CatalogIDs: catalogIDs,
		UserID:     userID,
	}

	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ranking request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", s.rankingAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create ranking request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ranking API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read ranking response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ranking API returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var rankingResponse RankingResponse
	if err := json.Unmarshal(responseBody, &rankingResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ranking response: %w", err)
	}

	// Check if the response indicates success
	if !rankingResponse.Success {
		return nil, fmt.Errorf("ranking API returned error: success field is false")
	}

	// Extract catalog IDs from the ranked catalogs
	var rankedCatalogIDs []string
	for _, rankedCatalog := range rankingResponse.RankedCatalogs {
		rankedCatalogIDs = append(rankedCatalogIDs, rankedCatalog.CatalogID)
	}

	// Log the ranking details for debugging
	if len(rankingResponse.RankedCatalogs) > 0 {
		fmt.Printf("Top ranked catalog: %s (PCTR: %.6f)\n",
			rankingResponse.RankedCatalogs[0].CatalogID,
			rankingResponse.RankedCatalogs[0].PctrScore)
	}

	return rankedCatalogIDs, nil
}

// GetRankedCatalogIDsWithFallback calls the ranking API with fallback to original catalog IDs
func (s *RankingService) GetRankedCatalogIDsWithFallback(catalogIDs []string, userID string) []string {
	fmt.Printf("Calling ranking API for user %s with %d catalog IDs\n", userID, len(catalogIDs))

	rankedIDs, err := s.GetRankedCatalogIDs(catalogIDs, userID)
	if err != nil {
		// Log the error but return original catalog IDs as fallback
		fmt.Printf("Warning: Failed to get ranked catalog IDs: %v. Using original catalog IDs.\n", err)
		return catalogIDs
	}

	// If ranking API returns empty result, use original catalog IDs
	if len(rankedIDs) == 0 {
		fmt.Printf("Warning: Ranking API returned empty result. Using original catalog IDs.\n")
		return catalogIDs
	}

	fmt.Printf("Successfully got %d ranked catalog IDs from ranking API\n", len(rankedIDs))
	return rankedIDs
}
