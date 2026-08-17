package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// BinLocationHandler 库位接口处理器。
type BinLocationHandler struct {
	binSvc service.BinLocationService
	logger *slog.Logger
}

// NewBinLocationHandler 构造库位接口处理器。
func NewBinLocationHandler(binSvc service.BinLocationService, logger *slog.Logger) *BinLocationHandler {
	return &BinLocationHandler{binSvc: binSvc, logger: logger}
}

// List 库位列表。
func (h *BinLocationHandler) List(c *gin.Context) {
	var query dto.BinQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	bins, total, err := h.binSvc.List(c.Request.Context(), query, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": bins, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// Get 库位详情。
func (h *BinLocationHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "库位ID参数无效")
		return
	}
	bin, err := h.binSvc.Get(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, bin)
}

// Contents 库位存放商品列表。
func (h *BinLocationHandler) Contents(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "库位ID参数无效")
		return
	}
	contents, err := h.binSvc.Contents(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, contents)
}

// Create 创建库位。
func (h *BinLocationHandler) Create(c *gin.Context) {
	var req dto.BinCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "库位创建参数校验失败，请检查区域/货架/层/列/容量字段")
		return
	}
	bin, err := h.binSvc.Create(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgBinCreated, bin)
}

// BatchCreate 批量创建库位。
func (h *BinLocationHandler) BatchCreate(c *gin.Context) {
	var req dto.BinBatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "批量创建库位参数校验失败，请检查区域/货架/层范围/列范围/容量字段")
		return
	}
	count, err := h.binSvc.BatchCreate(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgBinBatchCreated, gin.H{"created": count})
}

// Update 更新库位。
func (h *BinLocationHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "库位ID参数无效")
		return
	}
	var req dto.BinUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "库位更新参数校验失败，请检查容量/状态字段")
		return
	}
	bin, err := h.binSvc.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, bin)
}

// Recommend 推荐库位。
func (h *BinLocationHandler) Recommend(c *gin.Context) {
	var req dto.BinRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "推荐库位参数校验失败，请检查商品ID/数量字段")
		return
	}
	bin, err := h.binSvc.Recommend(c.Request.Context(), req.ProductID, req.Quantity)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, bin)
}
