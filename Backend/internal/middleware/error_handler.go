package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GlobalExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("RequestID")
				
				logger.Get().Error("Panic Recovered",
					zap.String("request_id", requestID),
					zap.Any("error", err),
					zap.String("stack_trace", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal Server Error",
					"request_id": requestID,
				})
			}
		}()

		start := time.Now()
		c.Next()

		// Handle errors attached to context via c.Error(err)
		if len(c.Errors) > 0 {
			requestID := c.GetString("RequestID")
			duration := time.Since(start).Milliseconds()

			for _, e := range c.Errors {
				logger.Get().Error("Request Error",
					zap.String("request_id", requestID),
					zap.String("error_type", "BusinessLogicError"),
					zap.Int64("duration_ms", duration),
					zap.Error(e.Err),
				)
			}
		}
	}
}