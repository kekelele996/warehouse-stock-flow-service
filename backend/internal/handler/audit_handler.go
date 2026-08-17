package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// AuditHandler 操作日志接口处理器。
type AuditHandler struct {
	auditSvc service.AuditService
	logger   *slog.Logger
}

// NewAuditHandler 构造操作日志接口处理器。
func NewAuditHandler(auditSvc service.AuditService, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc, logger: logger}
}

// List 操作日志列表。
func (h *AuditHandler) List(c *gin.Context) {
	var query dto.AuditQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	logs, total, err := h.auditSvc.List(c.Request.Context(), repository.OperationLogFilter{
		Module: query.Module, Action: query.Action, UserID: query.UserID, Page: page.Page, PageSize: page.PageSize,
	})
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgAuditListLoaded, gin.H{"list": logs, "total": total, "page": page.Page, "page_size": page.PageSize})
}
