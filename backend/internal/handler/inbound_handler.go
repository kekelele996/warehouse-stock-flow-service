package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// InboundHandler 入库单接口处理器。
type InboundHandler struct {
	inboundSvc service.InboundService
	logger     *slog.Logger
}

// NewInboundHandler 构造入库单接口处理器。
func NewInboundHandler(inboundSvc service.InboundService, logger *slog.Logger) *InboundHandler {
	return &InboundHandler{inboundSvc: inboundSvc, logger: logger}
}

// List 入库单列表。
func (h *InboundHandler) List(c *gin.Context) {
	var query dto.InboundQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	orders, total, err := h.inboundSvc.List(c.Request.Context(), query, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": orders, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// Get 入库单详情。
func (h *InboundHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "入库单ID参数无效")
		return
	}
	order, err := h.inboundSvc.Get(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, order)
}

// Create 创建入库单。
func (h *InboundHandler) Create(c *gin.Context) {
	var req dto.InboundCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "入库单创建参数校验失败，请检查货主/入库明细字段")
		return
	}
	order, err := h.inboundSvc.Create(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgInboundCreated, order)
}

// Receive 收货。
func (h *InboundHandler) Receive(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "入库单ID参数无效")
		return
	}
	order, err := h.inboundSvc.Receive(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgInboundReceived, order)
}

// QC 收货质检。
func (h *InboundHandler) QC(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "入库单ID参数无效")
		return
	}
	var req dto.InboundQCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "质检参数校验失败，请检查明细ID/实收数量/质检结果字段")
		return
	}
	order, err := h.inboundSvc.QC(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgInboundQCCompleted, order)
}

// Shelve 上架。
func (h *InboundHandler) Shelve(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "入库单ID参数无效")
		return
	}
	var req dto.InboundShelveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "上架参数校验失败，请检查明细ID/库位ID字段")
		return
	}
	order, err := h.inboundSvc.Shelve(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgInboundShelved, order)
}

// Complete 完成入库。
func (h *InboundHandler) Complete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "入库单ID参数无效")
		return
	}
	order, err := h.inboundSvc.Complete(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgInboundCompleted, order)
}
