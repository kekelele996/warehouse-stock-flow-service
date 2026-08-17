package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// OutboundHandler 出库单接口处理器。
type OutboundHandler struct {
	outboundSvc service.OutboundService
	logger      *slog.Logger
}

// NewOutboundHandler 构造出库单接口处理器。
func NewOutboundHandler(outboundSvc service.OutboundService, logger *slog.Logger) *OutboundHandler {
	return &OutboundHandler{outboundSvc: outboundSvc, logger: logger}
}

// List 出库单列表。
func (h *OutboundHandler) List(c *gin.Context) {
	var query dto.OutboundQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	orders, total, err := h.outboundSvc.List(c.Request.Context(), query, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": orders, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// Get 出库单详情。
func (h *OutboundHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单ID参数无效")
		return
	}
	order, err := h.outboundSvc.Get(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, order)
}

// Create 创建出库单。
func (h *OutboundHandler) Create(c *gin.Context) {
	var req dto.OutboundCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单创建参数校验失败，请检查货主/收货方/出库明细字段")
		return
	}
	order, err := h.outboundSvc.Create(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOutboundCreated, order)
}

// Picking 拣货。
func (h *OutboundHandler) Picking(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单ID参数无效")
		return
	}
	var req dto.OutboundPickingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "拣货参数校验失败，请检查明细ID/实拣数量字段")
		return
	}
	order, err := h.outboundSvc.Picking(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOutboundPicked, order)
}

// Checking 复核。
func (h *OutboundHandler) Checking(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单ID参数无效")
		return
	}
	order, err := h.outboundSvc.Checking(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOutboundChecked, order)
}

// Packing 打包。
func (h *OutboundHandler) Packing(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单ID参数无效")
		return
	}
	order, err := h.outboundSvc.Packing(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOutboundPacked, order)
}

// Ship 发货。
func (h *OutboundHandler) Ship(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单ID参数无效")
		return
	}
	var req dto.OutboundShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "发货参数校验失败，请检查快递单号字段")
		return
	}
	order, err := h.outboundSvc.Ship(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOutboundShipped, order)
}

// Complete 完成出库。
func (h *OutboundHandler) Complete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "出库单ID参数无效")
		return
	}
	order, err := h.outboundSvc.Complete(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOutboundCompleted, order)
}
