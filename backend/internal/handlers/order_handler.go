package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"meesho-clone/configs"
	"meesho-clone/internal/models"
	"meesho-clone/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	userService *services.UserService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		userService: services.NewUserService(),
	}
}

// PlaceOrderRequest represents the request for placing an order
type PlaceOrderRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	CatalogID string `json:"catalog_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// PlaceOrderResponse represents the response from place order API
type PlaceOrderResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// RTODeleteRequest represents the request to external RTO delete API
type RTODeleteRequest struct {
	Code      string `json:"code"`
	ProductID int    `json:"product_id"`
	CatalogID int    `json:"catalog_id"`
}

// PlaceOrder handles the place order request
func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	fmt.Printf("=== PLACE ORDER DEBUG START ===\n")

	var req PlaceOrderRequest

	// Log the raw request body
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("Raw request body: %s\n", string(body))

	// Reset the body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("Parsed request: %+v\n", req)

	// Validate that user exists
	fmt.Printf("Validating user: %s\n", req.UserID)
	_, err := h.userService.GetUserByID(req.UserID)
	if err != nil {
		fmt.Printf("User validation failed: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}
	fmt.Printf("User validation successful\n")

	// Generate order ID
	orderID := h.generateOrderID()
	fmt.Printf("Generated order ID: %s\n", orderID)

	// Call external RTO delete API
	fmt.Printf("Calling RTO drop API...\n")
	rtoSuccess := h.callRTODropAPI(req.UserID, req.ProductID, req.CatalogID)
	fmt.Printf("RTO drop API result: %v\n", rtoSuccess)

	// Create response
	response := PlaceOrderResponse{
		Success:   true,
		Message:   "Order placed successfully",
		OrderID:   orderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	// Add RTO API status to response
	if rtoSuccess {
		response.Message += " - RTO drop API called successfully"
	} else {
		response.Message += " - RTO drop API call failed, but order placed"
	}

	fmt.Printf("Sending response: %+v\n", response)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})

	fmt.Printf("=== PLACE ORDER DEBUG END ===\n")
}

// callRTODropAPI calls the external RTO delete API
func (h *OrderHandler) callRTODropAPI(userID, productID, catalogID string) bool {
	fmt.Printf("=== RTO API DEBUG START ===\n")
	fmt.Printf("Input parameters: userID=%s, productID=%s, catalogID=%s\n", userID, productID, catalogID)

	// Get user code from user_mapping table
	userCode, err := h.getUserCode(userID)
	if err != nil {
		fmt.Printf("Warning: Failed to get user code for RTO API call: %v\n", err)
		return false
	}

	fmt.Printf("Retrieved user code: %s\n", userCode)

	// Convert product_id and catalog_id to integers
	productIDInt := h.parseProductID(productID)
	catalogIDInt := h.parseCatalogID(catalogID)

	fmt.Printf("Parsed IDs: productIDInt=%d, catalogIDInt=%d\n", productIDInt, catalogIDInt)

	// Prepare request body
	requestBody := RTODeleteRequest{
		Code:      userCode,
		ProductID: productIDInt,
		CatalogID: catalogIDInt,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshaling RTO delete request: %v\n", err)
		return false
	}

	fmt.Printf("RTO API request body: %s\n", string(jsonData))

	// Call external RTO delete API
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("DELETE", "http://localhost:3001/rto/delete-by-product", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating RTO delete request: %v\n", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error calling RTO delete API: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading RTO delete response: %v\n", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("RTO delete API call successful: %s\n", string(body))
		fmt.Printf("=== RTO API DEBUG END (SUCCESS) ===\n")
		return true
	} else {
		fmt.Printf("RTO delete API call failed with status %d: %s\n", resp.StatusCode, string(body))
		fmt.Printf("=== RTO API DEBUG END (FAILED) ===\n")
		return false
	}
}

// getUserCode fetches the user's code from user_mapping table
func (h *OrderHandler) getUserCode(userID string) (string, error) {
	fmt.Printf("Getting user code for userID: %s\n", userID)

	var userMapping models.UserMapping

	// Use the database directly from configs
	result := configs.DB.Where("user_id = ?", userID).First(&userMapping)
	if result.Error != nil {
		fmt.Printf("Database error: %v\n", result.Error)
		return "", fmt.Errorf("failed to find user mapping for user_id %s: %w", userID, result.Error)
	}

	fmt.Printf("Found user mapping: %+v\n", userMapping)

	if userMapping.Code == "" {
		fmt.Printf("No code found for user %s\n", userID)
		return "", fmt.Errorf("no code found for user_id %s", userID)
	}

	fmt.Printf("Found code '%s' for user %s\n", userMapping.Code, userID)
	return userMapping.Code, nil
}

// generateOrderID generates a unique order ID
func (h *OrderHandler) generateOrderID() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("MEESH%d", timestamp)
}

// parseProductID converts product ID string to integer
func (h *OrderHandler) parseProductID(productID string) int {
	// Remove 's-' prefix if present
	if len(productID) > 2 && productID[:2] == "s-" {
		productID = productID[2:]
	}

	// Try to parse as integer, return fallback if failed
	var result int
	_, err := fmt.Sscanf(productID, "%d", &result)
	if err != nil {
		return 67890 // Fallback value
	}
	return result
}

// parseCatalogID converts catalog ID string to integer
func (h *OrderHandler) parseCatalogID(catalogID string) int {
	// Try to parse as integer, return fallback if failed
	var result int
	_, err := fmt.Sscanf(catalogID, "%d", &result)
	if err != nil {
		return 12345 // Fallback value
	}
	return result
}

// HealthCheck provides health check for order service
func (h *OrderHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "order",
		"status":  "healthy",
		"message": "Order service is running",
	})
}
