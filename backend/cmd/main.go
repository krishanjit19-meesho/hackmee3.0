package main

import (
	"log"
	"meesho-clone/configs"
	"meesho-clone/internal/handlers"
	"meesho-clone/internal/middleware"
	"meesho-clone/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Connect to database
	configs.ConnectDatabase()

	// Set Gin mode (release for production)
	gin.SetMode(gin.DebugMode) // Change to gin.ReleaseMode for production

	// Create Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())

	// Initialize services
	userService := services.NewUserService()
	productService := services.NewProductService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	homescreenHandler := handlers.NewHomescreenHandler()
	catalogHandler := handlers.NewCatalogHandler()
	productHandler := handlers.NewProductHandler(productService, userService)
	orderHandler := handlers.NewOrderHandler()

	// Health check endpoint
	router.GET("/health", authHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/validate", authHandler.ValidateUser)
			auth.GET("/profile/:user_id", authHandler.GetUserProfile)
			auth.PUT("/profile/:user_id", authHandler.UpdateUserProfile)
		}

		// Homescreen routes
		home := v1.Group("/home")
		{
			home.POST("/", homescreenHandler.GetHomescreen)
			home.GET("/", homescreenHandler.GetHomescreen) // Support both GET and POST
			home.GET("/categories", homescreenHandler.GetCategories)
			home.GET("/products", homescreenHandler.GetProducts)
			home.GET("/refresh", homescreenHandler.RefreshHomescreen)
		}

		// Search and product routes
		products := v1.Group("/products")
		{
			products.GET("/search", homescreenHandler.SearchProducts)
			products.GET("/:product_id", homescreenHandler.GetProductDetails)
		}

		// Product details routes
		product := v1.Group("/product")
		{
			product.GET("/details", productHandler.GetProductDetails)
			product.GET("/:id", productHandler.GetProductDetailsByID)
			product.GET("/health", productHandler.HealthCheck)
		}

		// Catalog routes
		catalog := v1.Group("/catalog")
		{
			catalog.GET("/", catalogHandler.GetCatalogData) // Single API for banner widget click
			catalog.GET("/health", catalogHandler.HealthCheck)
		}

		// Order routes
		order := v1.Group("/order")
		{
			order.POST("/place", orderHandler.PlaceOrder)
			order.GET("/health", orderHandler.HealthCheck)
		}
	}

	// Add a simple route to test CORS
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Meesho Clone API is running",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
			"endpoints": gin.H{
				"health":   "/health",
				"auth":     "/api/v1/auth/*",
				"home":     "/api/v1/home/*",
				"products": "/api/v1/products/*",
				"catalog":  "/api/v1/catalog/*",
				"order":    "/api/v1/order/*",
			},
		})
	})

	// Start server
	port := "8080"
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📱 Meesho Clone API is ready")
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("📊 API docs: http://localhost:%s/", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
