package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Something went wrong. Please try again later.",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// AddTimestampMiddleware adds timestamp to context
func AddTimestampMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("timestamp", time.Now().Unix())
		c.Next()
	}
}
