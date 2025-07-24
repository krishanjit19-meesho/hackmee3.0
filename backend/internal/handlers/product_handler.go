package handlers

import (
	"net/http"

	"meesho-clone/internal/services"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productService *services.ProductService
	userService    *services.UserService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *services.ProductService, userService *services.UserService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		userService:    userService,
	}
}

// GetProductDetails handles GET requests for product details
func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	// Get product_id from query parameter
	productID := c.Query("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "product_id is required as query parameter",
		})
		return
	}

	// Get user_id from query parameter
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "user_id is required as query parameter",
		})
		return
	}

	// Validate that user exists
	_, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
			"details": "user not found",
		})
		return
	}

	// Get product details
	productResponse, err := h.productService.GetProductDetails(productID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get product details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, productResponse)
}

// GetProductDetailsByID handles GET requests for product details by ID in URL
func (h *ProductHandler) GetProductDetailsByID(c *gin.Context) {
	// Get product_id from URL parameter
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "product_id is required in URL path",
		})
		return
	}

	// Get user_id from query parameter
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "user_id is required as query parameter",
		})
		return
	}

	// Validate that user exists
	_, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
			"details": "user not found",
		})
		return
	}

	// Get product details
	productResponse, err := h.productService.GetProductDetails(productID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get product details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, productResponse)
}

// HealthCheck handles health check requests for product service
func (h *ProductHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product service is healthy",
		"service": "product",
		"status":  "running",
	})
}
