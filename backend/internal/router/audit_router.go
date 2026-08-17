package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// RegisterAuditRoutes 注册操作日志路由。
func RegisterAuditRoutes(g *gin.RouterGroup, h *handler.AuditHandler, logger *slog.Logger) {
	g.GET("/audit", middleware.RequirePermission("audit:view"), h.List)
}
