package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if Request ID already exists in header (from a Load Balancer / Gateway)
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set it in the Gin Context so other functions can access it
		c.Set("RequestID", requestID)

		// Set it in the Response Header so the client can report it on error
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}