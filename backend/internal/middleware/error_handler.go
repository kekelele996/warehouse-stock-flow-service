package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/util"
)

// ErrorHandler 统一异常处理：panic 恢复 + 错误渲染。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, _ := r.(error)
				logger.ErrorContext(c.Request.Context(), constants.LogPanicRecovered,
					"path", c.Request.URL.Path, "request_id", GetRequestID(c), "err", err)
				util.Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
				c.Abort()
			}
		}()
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *util.AppError
			if errors.As(err, &appErr) {
				util.Fail(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
				return
			}
			util.Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
	}
}
