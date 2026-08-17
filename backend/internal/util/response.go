package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应体：{ "code": 0, "message": "ok", "data": ... }。
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Success 返回成功响应。
func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: message, Data: data})
}

// Fail 返回错误响应，HTTP 状态由调用方决定。
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Response{Code: code, Message: message, Data: nil})
}

// FailWithAppError 将 AppError 渲染为统一错误响应。
func FailWithAppError(c *gin.Context, err error) {
	appErr := ToAppError(err)
	Fail(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
}
