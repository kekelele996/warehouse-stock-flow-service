package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// RegisterOutboundRoutes 注册出库单路由。
func RegisterOutboundRoutes(g *gin.RouterGroup, h *handler.OutboundHandler, logger *slog.Logger) {
	outbound := g.Group("/outbound")
	{
		outbound.GET("", h.List)
		outbound.POST("", h.Create)
		outbound.GET("/:id", h.Get)
		outbound.POST("/:id/picking", middleware.RequireRoles("Admin", "WarehouseManager", "Picker"), h.Picking)
		outbound.POST("/:id/checking", middleware.RequireRoles("Admin", "WarehouseManager", "Checker"), h.Checking)
		outbound.POST("/:id/packing", middleware.RequireRoles("Admin", "WarehouseManager", "Checker"), h.Packing)
		outbound.POST("/:id/ship", middleware.RequireRoles("Admin", "WarehouseManager"), h.Ship)
		outbound.POST("/:id/complete", middleware.RequireRoles("Admin", "WarehouseManager"), h.Complete)
	}
}
