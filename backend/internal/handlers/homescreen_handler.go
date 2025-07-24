package handlers

import (
	"meesho-clone/internal/models"
	"meesho-clone/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomescreenHandler handles homescreen related requests
type HomescreenHandler struct {
	userService   *services.UserService
	meeshoService *services.MeeshoService
}

// NewHomescreenHandler creates a new homescreen handler
func NewHomescreenHandler() *HomescreenHandler {
	return &HomescreenHandler{
		userService:   services.NewUserService(),
		meeshoService: services.NewMeeshoService(),
	}
}

// GetHomescreen fetches homescreen data for a user
func (h *HomescreenHandler) GetHomescreen(c *gin.Context) {
	var req models.HomescreenRequest

	// Try to get user_id from JSON body first, then from query params
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON binding fails, try query params
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user_id is required in request body or query parameter",
			})
			return
		}
		req.UserID = userID
	}

	// Validate that user exists
	user, err := h.userService.GetUserByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	// Fetch homescreen data from Meesho API
	meeshoData, err := h.meeshoService.FetchHomescreenData(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch homescreen data",
			"details": err.Error(),
		})
		return
	}

	// Format response for frontend consumption
	response := h.meeshoService.FormatHomescreenResponse(meeshoData, req.UserID)

	// Add user information to response
	response["user_info"] = gin.H{
		"user_id": user.UserID,
		"name":    user.Name,
		"phone":   user.PhoneNumber,
	}

	c.JSON(http.StatusOK, response)
}

// GetCategories fetches just the categories from navigation
func (h *HomescreenHandler) GetCategories(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	// Validate user exists
	_, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Fetch data from Meesho API
	meeshoData, err := h.meeshoService.FetchHomescreenData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch categories",
			"details": err.Error(),
		})
		return
	}

	// Extract just categories
	categories := make([]map[string]interface{}, 0)
	for _, tile := range meeshoData.TopNavBar.Tiles {
		category := map[string]interface{}{
			"id":          tile.ID,
			"title":       tile.Title,
			"image":       tile.Image,
			"destination": tile.DestinationData,
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"user_id":    userID,
		"categories": categories,
		"total":      len(categories),
	})
}

// GetProducts fetches just the products/widgets
func (h *HomescreenHandler) GetProducts(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	// Validate user exists
	_, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Fetch data from Meesho API
	meeshoData, err := h.meeshoService.FetchHomescreenData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
		return
	}

	// Extract products from widget groups
	products := make([]map[string]interface{}, 0)
	for _, group := range meeshoData.WidgetGroups {
		for _, widget := range group.Widgets {
			product := map[string]interface{}{
				"id":          widget.ID,
				"title":       widget.Title,
				"image":       widget.Image,
				"screen":      widget.Screen,
				"group_title": group.Title,
				"group_id":    group.ID,
				"destination": widget.Data,
				"priority":    widget.Priority,
			}
			products = append(products, product)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"user_id":  userID,
		"products": products,
		"total":    len(products),
	})
}

// SearchProducts handles product search
func (h *HomescreenHandler) SearchProducts(c *gin.Context) {
	userID := c.Query("user_id")
	query := c.Query("q")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "search query 'q' is required",
		})
		return
	}

	// Validate user exists
	_, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Use Meesho service to search (placeholder implementation)
	results, err := h.meeshoService.SearchProducts(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Search failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"query":   query,
		"user_id": userID,
		"results": results,
	})
}

// GetProductDetails handles product detail requests
func (h *HomescreenHandler) GetProductDetails(c *gin.Context) {
	userID := c.Query("user_id")
	productID := c.Param("product_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "product_id is required",
		})
		return
	}

	// Validate user exists
	_, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Get product details from Meesho service
	details, err := h.meeshoService.GetProductDetails(productID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch product details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"user_id":    userID,
		"product_id": productID,
		"details":    details,
	})
}

// RefreshHomescreen forces a fresh fetch from Meesho API
func (h *HomescreenHandler) RefreshHomescreen(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	// Validate user exists
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Force fresh fetch from Meesho API
	meeshoData, err := h.meeshoService.FetchHomescreenData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to refresh homescreen data",
			"details": err.Error(),
		})
		return
	}

	// Format response
	response := h.meeshoService.FormatHomescreenResponse(meeshoData, userID)
	response["refreshed"] = true
	response["user_info"] = gin.H{
		"user_id": user.UserID,
		"name":    user.Name,
		"phone":   user.PhoneNumber,
	}

	c.JSON(http.StatusOK, response)
}
