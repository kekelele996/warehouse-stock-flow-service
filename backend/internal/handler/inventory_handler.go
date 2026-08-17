package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// InventoryHandler 库存接口处理器。
type InventoryHandler struct {
	invSvc service.InventoryService
	logger *slog.Logger
}

// NewInventoryHandler 构造库存接口处理器。
func NewInventoryHandler(invSvc service.InventoryService, logger *slog.Logger) *InventoryHandler {
	return &InventoryHandler{invSvc: invSvc, logger: logger}
}

// List 库存列表。
func (h *InventoryHandler) List(c *gin.Context) {
	var query dto.InventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	list, total, err := h.invSvc.List(c.Request.Context(), query, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": list, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// Summary 货主库存金额汇总。
func (h *InventoryHandler) Summary(c *gin.Context) {
	items, err := h.invSvc.Summary(c.Request.Context())
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, items)
}
