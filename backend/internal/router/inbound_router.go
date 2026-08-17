package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// RegisterInboundRoutes 注册入库单路由。
func RegisterInboundRoutes(g *gin.RouterGroup, h *handler.InboundHandler, logger *slog.Logger) {
	inbound := g.Group("/inbound")
	{
		inbound.GET("", h.List)
		inbound.POST("", h.Create)
		inbound.GET("/:id", h.Get)
		inbound.POST("/:id/receive", middleware.RequireRoles("Admin", "WarehouseManager"), h.Receive)
		inbound.POST("/:id/qc", middleware.RequireRoles("Admin", "WarehouseManager", "QCInspector"), h.QC)
		inbound.POST("/:id/shelve", middleware.RequireRoles("Admin", "WarehouseManager"), h.Shelve)
		inbound.POST("/:id/complete", middleware.RequireRoles("Admin", "WarehouseManager"), h.Complete)
	}
}
