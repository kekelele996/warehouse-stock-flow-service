package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

// Audit 通用审计中间件：为写操作（POST/PUT/PATCH/DELETE）记录操作日志。
func Audit(logRepo repository.OperationLogRepository, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 400 {
			return
		}
		switch c.Request.Method {
		case "POST", "PUT", "PATCH", "DELETE":
		default:
			return
		}
		entry := &model.OperationLog{
			Module:     moduleFromPath(c.Request.URL.Path),
			Action:     actionFromRequest(c),
			EntityType: entityTypeFromPath(c.Request.URL.Path),
			EntityID:   c.Param("id"),
			Detail:     fmt.Sprintf("请求 %s %s", c.Request.Method, c.Request.URL.Path),
			IP:         c.ClientIP(),
		}
		if claims, ok := util.CurrentUser(c.Request.Context()); ok {
			entry.UserID = claims.UserID
			entry.Username = claims.Username
			entry.Role = claims.Role
		}
		if err := logRepo.Create(c.Request.Context(), entry); err != nil {
			logger.WarnContext(c.Request.Context(), "audit middleware create failed", "err", err)
		}
	}
}

func moduleFromPath(path string) string {
	for _, m := range []string{"inbound", "outbound", "owner", "product", "bin", "inventory", "auth", "audit", "dashboard"} {
		if containsPath(path, m) {
			return m
		}
	}
	return "misc"
}

func entityTypeFromPath(path string) string {
	for _, m := range []string{"inbound", "outbound", "owner", "product", "bin", "inventory"} {
		if containsPath(path, m) {
			return m
		}
	}
	return ""
}

func actionFromRequest(c *gin.Context) string {
	method := c.Request.Method
	switch method {
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "PATCH":
		return "patch"
	case "DELETE":
		return "delete"
	default:
		return "write"
	}
}

func containsPath(path, sub string) bool {
	for i := 0; i+len(sub) <= len(path); i++ {
		if path[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
