package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/util"
)

// Auth JWT 鉴权中间件。
func Auth(jwtSecret string, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(jwtSecret, tokenStr)
		if err != nil {
			util.Fail(c, http.StatusUnauthorized, constants.CodeTokenExpired, "登录已过期，请重新登录")
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(util.WithUser(c.Request.Context(), claims))
		c.Next()
	}
}
