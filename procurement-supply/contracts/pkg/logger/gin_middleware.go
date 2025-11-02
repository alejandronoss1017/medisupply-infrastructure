package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDKey = "request_id"

// GinLogger returns Gin middleware for structured logging with Zap
func GinLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Build log entry
		fields := []interface{}{
			"request_id", requestID,
			"status", statusCode,
			"method", method,
			"path", path,
			"query", query,
			"ip", clientIP,
			"latency_ms", latency.Milliseconds(),
			"user_agent", c.Request.UserAgent(),
		}

		// Add error if present
		if errorMessage != "" {
			fields = append(fields, "error", errorMessage)
		}

		// Log based on status code
		switch {
		case statusCode >= 500:
			logger.Errorw("Server error", fields...)
		case statusCode >= 400:
			logger.Warnw("Client error", fields...)
		case statusCode >= 300:
			logger.Infow("Redirection", fields...)
		default:
			logger.Infow("Request completed", fields...)
		}
	}
}

// GinRecovery returns Gin middleware for panic recovery with structured logging
func GinRecovery(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get request ID if available
				requestID, _ := c.Get(RequestIDKey)

				logger.Errorw("Panic recovered",
					"request_id", requestID,
					"error", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"ip", c.ClientIP(),
				)

				// Abort with 500
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// GetRequestID extracts the request ID from the Gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
