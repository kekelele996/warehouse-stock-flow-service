package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// ProductHandler 商品接口处理器。
type ProductHandler struct {
	productSvc service.ProductService
	logger     *slog.Logger
}

// NewProductHandler 构造商品接口处理器。
func NewProductHandler(productSvc service.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{productSvc: productSvc, logger: logger}
}

// List 商品列表。
func (h *ProductHandler) List(c *gin.Context) {
	var query dto.ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	products, total, err := h.productSvc.List(c.Request.Context(), query, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": products, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// ListByOwner 按货主查询商品（复用商品列表逻辑）。
func (h *ProductHandler) ListByOwner(c *gin.Context) {
	ownerID, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	page := parsePage(c)
	products, total, err := h.productSvc.ListByOwner(c.Request.Context(), ownerID, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": products, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// Get 商品详情。
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "商品ID参数无效")
		return
	}
	product, err := h.productSvc.Get(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, product)
}

// Create 创建商品。
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "商品创建参数校验失败，请检查商品名称/SKU/单位/存储要求字段")
		return
	}
	product, err := h.productSvc.Create(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgProductCreated, product)
}

// Update 更新商品。
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "商品ID参数无效")
		return
	}
	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "商品更新参数校验失败，请检查商品名称/单位字段")
		return
	}
	product, err := h.productSvc.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOwnerUpdated, product)
}
