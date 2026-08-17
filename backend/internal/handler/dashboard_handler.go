package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

// DashboardHandler 仓库总览接口处理器。
type DashboardHandler struct {
	dashSvc service.DashboardService
	logger  *slog.Logger
}

// NewDashboardHandler 构造仓库总览接口处理器。
func NewDashboardHandler(dashSvc service.DashboardService, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{dashSvc: dashSvc, logger: logger}
}

// Summary 仓库总览汇总。
func (h *DashboardHandler) Summary(c *gin.Context) {
	data, err := h.dashSvc.Summary(c.Request.Context())
	if err != nil {
		util.FailWithAppError(c, err)
		return
	}
	util.Success(c, constants.MsgDashboardLoaded, data)
}
