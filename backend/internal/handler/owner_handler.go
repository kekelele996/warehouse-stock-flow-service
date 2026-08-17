package handler

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// OwnerHandler 货主接口处理器。
type OwnerHandler struct {
	ownerSvc service.OwnerService
	logger   *slog.Logger
}

// NewOwnerHandler 构造货主接口处理器。
func NewOwnerHandler(ownerSvc service.OwnerService, logger *slog.Logger) *OwnerHandler {
	return &OwnerHandler{ownerSvc: ownerSvc, logger: logger}
}

// List 货主列表。
func (h *OwnerHandler) List(c *gin.Context) {
	var query dto.OwnerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, constants.MsgInvalidParams)
		return
	}
	page := parsePage(c)
	owners, total, err := h.ownerSvc.List(c.Request.Context(), query, page.Page, page.PageSize)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, gin.H{"list": owners, "total": total, "page": page.Page, "page_size": page.PageSize})
}

// Get 货主详情。
func (h *OwnerHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	owner, err := h.ownerSvc.Get(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, owner)
}

// Stats 货主详情统计（商品数/单据数/库存金额）。
func (h *OwnerHandler) Stats(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	stats, err := h.ownerSvc.GetStats(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, stats)
}

// Create 创建货主。
func (h *OwnerHandler) Create(c *gin.Context) {
	var req dto.OwnerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主创建参数校验失败，请检查货主名称/联系方式/结算方式字段")
		return
	}
	owner, err := h.ownerSvc.Create(c.Request.Context(), &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOwnerCreated, owner)
}

// Update 更新货主。
func (h *OwnerHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	var req dto.OwnerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主更新参数校验失败，请检查联系人/电话字段")
		return
	}
	owner, err := h.ownerSvc.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOwnerUpdated, owner)
}

// AdjustCredit 调整信用额度。
func (h *OwnerHandler) AdjustCredit(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	var req dto.OwnerCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "信用额度参数校验失败，请检查信用额度字段")
		return
	}
	owner, err := h.ownerSvc.AdjustCredit(c.Request.Context(), id, &req)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOK, owner)
}

// Suspend 暂停合作。
func (h *OwnerHandler) Suspend(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	owner, err := h.ownerSvc.Suspend(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOwnerSuspended, owner)
}

// Activate 恢复合作。
func (h *OwnerHandler) Activate(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		util.Fail(c, 400, constants.CodeInvalidParams, "货主ID参数无效")
		return
	}
	owner, err := h.ownerSvc.Activate(c.Request.Context(), id)
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgOwnerActivated, owner)
}

// parseID 解析路径 ID。
func parseID(c *gin.Context) (uint, error) {
	raw := c.Param("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}
