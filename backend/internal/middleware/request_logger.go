package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
)

// RequestLogger 记录请求日志（method/path/status/latency_ms/request_id）。
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		rid := GetRequestID(c)
		logger.InfoContext(c.Request.Context(), constants.LogRequestIn, "method", c.Request.Method, "path", c.Request.URL.Path, "request_id", rid)
		c.Next()
		latency := time.Since(start).Milliseconds()
		logger.InfoContext(c.Request.Context(), constants.LogRequestOut,
			"method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "latency_ms", latency, "request_id", rid)
	}
}
