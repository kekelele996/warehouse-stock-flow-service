package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/util"
)

// RequireRoles 角色校验中间件：允许列表中的角色访问。
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := util.CurrentUser(c.Request.Context())
		if !ok {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		for _, role := range roles {
			if claims.Role == role {
				c.Next()
				return
			}
		}
		util.Fail(c, http.StatusForbidden, constants.CodeRoleForbidden,
			"角色 "+claims.Role+" 无权执行该操作，需要角色："+util.JoinNames(roles))
		c.Abort()
	}
}

// RequirePermission 权限点校验中间件。
func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := util.CurrentUser(c.Request.Context())
		if !ok {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		if !constants.HasPermission(claims.Role, perm) {
			util.Fail(c, http.StatusForbidden, constants.CodeRoleForbidden,
				"角色 "+claims.Role+" 缺少权限点 "+perm)
			c.Abort()
			return
		}
		c.Next()
	}
}
