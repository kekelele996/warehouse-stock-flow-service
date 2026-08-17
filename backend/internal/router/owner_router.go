package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// RegisterOwnerRoutes 注册货主路由。
func RegisterOwnerRoutes(g *gin.RouterGroup, h *handler.OwnerHandler, logger *slog.Logger) {
	owners := g.Group("/owners")
	{
		owners.GET("", h.List)
		owners.POST("", middleware.RequirePermission("owner:manage"), h.Create)
		owners.GET("/:id", h.Get)
		owners.GET("/:id/stats", h.Stats)
		owners.PUT("/:id", middleware.RequirePermission("owner:manage"), h.Update)
		owners.PUT("/:id/credit", middleware.RequirePermission("owner:manage"), h.AdjustCredit)
		owners.PUT("/:id/suspend", middleware.RequirePermission("owner:manage"), h.Suspend)
		owners.PUT("/:id/activate", middleware.RequirePermission("owner:manage"), h.Activate)
	}
}
