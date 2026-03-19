// This middleware catches panics and errors, formatting them into your required JSON log structure.
package middleware

import (
	"net/http"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		start := time.Now()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Get().Error().
				Str("request_id", requestID).
				Str("error_type", "API_ERROR").
				Int64("duration_ms", time.Since(start).Milliseconds()).
				Msg(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_id": requestID,
				"message":  "An unexpected error occurred",
			})
		}
	}
}