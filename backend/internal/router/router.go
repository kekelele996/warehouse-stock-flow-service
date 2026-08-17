package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// Register 注册全部路由与全局中间件。
func Register(
	r *gin.Engine,
	jwtSecret string,
	logger *slog.Logger,
	authHandler *handler.AuthHandler,
	ownerHandler *handler.OwnerHandler,
	productHandler *handler.ProductHandler,
	binHandler *handler.BinLocationHandler,
	inventoryHandler *handler.InventoryHandler,
	inboundHandler *handler.InboundHandler,
	outboundHandler *handler.OutboundHandler,
	dashboardHandler *handler.DashboardHandler,
	auditHandler *handler.AuditHandler,
	auditMiddleware gin.HandlerFunc,
) {
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(auditMiddleware)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "wmsflow"})
	})

	v1 := r.Group("/api/v1")
	{
		RegisterAuthRoutes(v1, authHandler, jwtSecret, logger)
		authed := v1.Group("")
		authed.Use(middleware.Auth(jwtSecret, logger))
		{
			RegisterOwnerRoutes(authed, ownerHandler, logger)
			RegisterProductRoutes(authed, productHandler, logger)
			RegisterBinLocationRoutes(authed, binHandler, logger)
			RegisterInventoryRoutes(authed, inventoryHandler, logger)
			RegisterInboundRoutes(authed, inboundHandler, logger)
			RegisterOutboundRoutes(authed, outboundHandler, logger)
			RegisterDashboardRoutes(authed, dashboardHandler, logger)
			RegisterAuditRoutes(authed, auditHandler, logger)
		}
	}
}
