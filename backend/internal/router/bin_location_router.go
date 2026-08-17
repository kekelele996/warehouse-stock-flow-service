package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// RegisterBinLocationRoutes 注册库位路由。
func RegisterBinLocationRoutes(g *gin.RouterGroup, h *handler.BinLocationHandler, logger *slog.Logger) {
	bins := g.Group("/bin-locations")
	{
		bins.GET("", h.List)
		bins.POST("", middleware.RequirePermission("bin:manage"), h.Create)
		bins.POST("/batch", middleware.RequirePermission("bin:manage"), h.BatchCreate)
		bins.POST("/recommend", middleware.RequirePermission("bin:manage"), h.Recommend)
		bins.GET("/:id", h.Get)
		bins.GET("/:id/contents", h.Contents)
		bins.PUT("/:id", middleware.RequirePermission("bin:manage"), h.Update)
	}
}
