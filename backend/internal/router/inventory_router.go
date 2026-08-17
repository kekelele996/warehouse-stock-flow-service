package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
)

// RegisterInventoryRoutes 注册库存路由。
func RegisterInventoryRoutes(g *gin.RouterGroup, h *handler.InventoryHandler, logger *slog.Logger) {
	inv := g.Group("/inventory")
	{
		inv.GET("", h.List)
		inv.GET("/summary", h.Summary)
	}
}
