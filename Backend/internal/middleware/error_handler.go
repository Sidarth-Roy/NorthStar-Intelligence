package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

		// Logic to handle errors attached via c.Error(err)
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err
			requestID := c.GetString("RequestID")
			
			// Industry Standard: Map error types to HTTP Status Codes
			statusCode := http.StatusInternalServerError

			// 1. Check GORM specific errors
			if errors.Is(lastErr, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			} else if errors.Is(lastErr, gorm.ErrDuplicatedKey) {
				statusCode = http.StatusConflict // 409
			}

			// 2. Check Postgres/Driver specific errors (e.g., Not Null violations)
			var pgErr *pgconn.PgError
			if errors.As(lastErr, &pgErr) {
				switch pgErr.Code {
				case "23502": // not_null_violation
					statusCode = http.StatusBadRequest
				case "23505": // unique_violation
					statusCode = http.StatusConflict
				case "23503": // foreign_key_violation
					statusCode = http.StatusBadRequest
				}
			}

			// 3. Log the error with full context
			logger.Get().Error("Request Error",
				zap.String("request_id", requestID),
				zap.Int("status_code", statusCode),
				zap.Int64("duration_ms", time.Since(start).Milliseconds()),
				zap.Error(lastErr),
			)

			// 4. CRITICAL: Stop the 200 OK and send the mapped error
			if !c.Writer.Written() {
				c.AbortWithStatusJSON(statusCode, gin.H{
					"error":      lastErr.Error(),
					"request_id": requestID,
				})
			}
		}
	}
}
