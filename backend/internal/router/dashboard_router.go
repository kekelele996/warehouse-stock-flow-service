package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
)

// RegisterDashboardRoutes 注册仓库总览路由。
func RegisterDashboardRoutes(g *gin.RouterGroup, h *handler.DashboardHandler, logger *slog.Logger) {
	g.GET("/dashboard/summary", h.Summary)
}
