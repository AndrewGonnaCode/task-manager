package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger middleware logs incoming requests with correlation IDs
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or generate correlation ID
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Set correlation ID in context
		c.Set("correlation_id", correlationID)
		c.Writer.Header().Set("X-Correlation-ID", correlationID)

		// Start timer
		start := time.Now()

		// Log request
		log.Printf("[%s] %s %s - Started", correlationID, c.Request.Method, c.Request.URL.Path)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		log.Printf("[%s] %s %s - Completed in %v (status: %d)",
			correlationID,
			c.Request.Method,
			c.Request.URL.Path,
			duration,
			c.Writer.Status(),
		)
	}
}
