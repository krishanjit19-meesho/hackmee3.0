package handlers

import (
	"meesho-clone/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CatalogHandler handles catalog-related requests
type CatalogHandler struct {
	catalogService *services.CatalogService
	userService    *services.UserService
}

// NewCatalogHandler creates a new catalog handler
func NewCatalogHandler() *CatalogHandler {
	return &CatalogHandler{
		catalogService: services.NewCatalogService(),
		userService:    services.NewUserService(),
	}
}

// GetCatalogData handles the banner widget click and returns catalog data
func (h *CatalogHandler) GetCatalogData(c *gin.Context) {
	// Get user_id from query parameter (simplified approach)
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
			"details": err.Error(),
		})
		return
	}

	// Fetch catalog data using hardcoded catalog IDs
	catalogResponse, err := h.catalogService.GetCatalogData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch catalog data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, catalogResponse)
}

// HealthCheck provides health check for catalog service
func (h *CatalogHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "catalog",
		"status":  "healthy",
		"message": "Catalog service is running",
	})
}
