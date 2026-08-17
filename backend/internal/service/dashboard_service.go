package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

// DashboardService 仓库总览服务。
type DashboardService interface {
	Summary(ctx context.Context) (*dto.DashboardSummary, error)
}

type dashboardService struct {
	inboundRepo  repository.InboundRepository
	outboundRepo repository.OutboundRepository
	binRepo      repository.BinLocationRepository
	inventorySvc InventoryService
	logger       *slog.Logger
}

// NewDashboardService 构造仓库总览服务。
func NewDashboardService(
	inboundRepo repository.InboundRepository,
	outboundRepo repository.OutboundRepository,
	binRepo repository.BinLocationRepository,
	inventorySvc InventoryService,
	logger *slog.Logger,
) DashboardService {
	return &dashboardService{inboundRepo: inboundRepo, outboundRepo: outboundRepo, binRepo: binRepo, inventorySvc: inventorySvc, logger: logger}
}

// Summary 汇总仓库总览数据；Owner 角色仅统计自己的单据与库存。
func (s *dashboardService) Summary(ctx context.Context) (*dto.DashboardSummary, error) {
	ownerID := uint(0)
	if claims, ok := util.CurrentUser(ctx); ok && claims.Role == constants.RoleOwner && claims.OwnerID != nil {
		ownerID = *claims.OwnerID
	}
	todayInbound, err := s.inboundRepo.CountToday(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary today inbound: %w", err)
	}
	todayOutbound, err := s.outboundRepo.CountToday(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary today outbound: %w", err)
	}
	var pendingInbound, pendingOutbound int64
	for _, status := range []string{
		constants.InboundStatusPending, constants.InboundStatusReceived,
		constants.InboundStatusQCInProgress, constants.InboundStatusShelved,
	} {
		n, err := s.inboundRepo.CountByStatus(ctx, status, ownerID)
		if err != nil {
			return nil, fmt.Errorf("dashboard summary pending inbound %s: %w", status, err)
		}
		pendingInbound += n
	}
	for _, status := range []string{
		constants.OutboundStatusPending, constants.OutboundStatusPicking,
		constants.OutboundStatusChecking, constants.OutboundStatusPacking,
	} {
		n, err := s.outboundRepo.CountByStatus(ctx, status, ownerID)
		if err != nil {
			return nil, fmt.Errorf("dashboard summary pending outbound %s: %w", status, err)
		}
		pendingOutbound += n
	}
	totalBins, err := s.binRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary bin count: %w", err)
	}
	occupiedBins, err := s.binRepo.CountOccupied(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary occupied bin: %w", err)
	}
	sumRate, err := s.binRepo.SumOccupancyRate(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary bin rate: %w", err)
	}
	avgRate := 0.0
	if totalBins > 0 {
		avgRate = sumRate / float64(totalBins)
	}
	ownerTop, err := s.inventorySvc.Summary(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary owner top: %w", err)
	}
	topItems := make([]dto.OwnerTopItem, 0, len(ownerTop))
	for _, item := range ownerTop {
		topItems = append(topItems, dto.OwnerTopItem{OwnerID: item.OwnerID, OwnerName: item.OwnerName, TotalValue: item.TotalValue})
	}
	topItems = sortOwnerTop(topItems, 10)
	pendingTasks := make([]dto.PendingTaskItem, 0, 10)
	inboundPending, err := s.inboundRepo.ListIncomplete(ctx, ownerID, 5)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary inbound tasks: %w", err)
	}
	for i := range inboundPending {
		pendingTasks = append(pendingTasks, dto.PendingTaskItem{
			Type:    "inbound",
			Title:   "入库单待处理：" + util.InboundStatusText(inboundPending[i].Status),
			OrderNo: inboundPending[i].OrderNo,
			Status:  inboundPending[i].Status,
		})
	}
	outboundPending, err := s.outboundRepo.ListIncomplete(ctx, ownerID, 5)
	if err != nil {
		return nil, fmt.Errorf("dashboard summary outbound tasks: %w", err)
	}
	for i := range outboundPending {
		pendingTasks = append(pendingTasks, dto.PendingTaskItem{
			Type:    "outbound",
			Title:   "出库单待处理：" + util.OutboundStatusText(outboundPending[i].Status),
			OrderNo: outboundPending[i].OrderNo,
			Status:  outboundPending[i].Status,
		})
	}
	s.logger.InfoContext(ctx, constants.LogDashboardLoaded, "owner_id", ownerID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return &dto.DashboardSummary{
		TodayInboundCount:  todayInbound,
		TodayOutboundCount: todayOutbound,
		PendingInbound:     pendingInbound,
		PendingOutbound:    pendingOutbound,
		BinOccupancyRate:   avgRate,
		TotalBinCount:      totalBins,
		OccupiedBinCount:   occupiedBins,
		OwnerTop10:         topItems,
		PendingTasks:       pendingTasks,
	}, nil
}

// sortOwnerTop 按库存金额降序排序并截取 TopN。
func sortOwnerTop(items []dto.OwnerTopItem, n int) []dto.OwnerTopItem {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].TotalValue > items[i].TotalValue {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	if len(items) > n {
		items = items[:n]
	}
	return items
}
