package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// AuthHandler 认证接口处理器。
type AuthHandler struct {
	authSvc service.AuthService
	logger  *slog.Logger
}

// NewAuthHandler 构造认证接口处理器。
func NewAuthHandler(authSvc service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, logger: logger}
}

// Register 注册（默认创建货主角色账号并同步创建货主档案）。
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	user, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgRegisterSuccess, user)
}

// Login 登录。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	resp, err := h.authSvc.Login(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgLoginSuccess, resp)
}

// Me 当前用户信息。
func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.authSvc.Me(c.Request.Context())
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, user)
}
