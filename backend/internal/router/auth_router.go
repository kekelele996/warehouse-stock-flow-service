package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
)

// RegisterAuthRoutes 注册认证路由。
func RegisterAuthRoutes(g *gin.RouterGroup, h *handler.AuthHandler, jwtSecret string, logger *slog.Logger) {
	auth := g.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.GET("/me", middleware.Auth(jwtSecret, logger), h.Me)
	}
}
