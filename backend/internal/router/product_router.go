package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
)

// RegisterProductRoutes 注册商品路由。
func RegisterProductRoutes(g *gin.RouterGroup, h *handler.ProductHandler, logger *slog.Logger) {
	products := g.Group("/products")
	{
		products.GET("", h.List)
		products.POST("", h.Create)
		products.GET("/by-owner/:ownerId", h.ListByOwner)
		products.GET("/:id", h.Get)
		products.PUT("/:id", h.Update)
	}
}
