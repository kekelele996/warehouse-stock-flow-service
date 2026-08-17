package service

import (
	"context"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/util"
)

func TestDashboardService_Summary(t *testing.T) {
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.countTodayFn = func(ctx context.Context, ownerID uint) (int64, error) { return 3, nil }
	inboundRepo.countStatusFn = func(ctx context.Context, status string, ownerID uint) (int64, error) { return 1, nil }
	inboundRepo.listIncFn = func(ctx context.Context, ownerID uint, limit int) ([]model.InboundOrder, error) {
		return []model.InboundOrder{{ID: 1, OrderNo: "INB-1", Status: constants.InboundStatusPending}}, nil
	}
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.countTodayFn = func(ctx context.Context, ownerID uint) (int64, error) { return 2, nil }
	outboundRepo.countStatusFn = func(ctx context.Context, status string, ownerID uint) (int64, error) { return 0, nil }
	outboundRepo.listIncFn = func(ctx context.Context, ownerID uint, limit int) ([]model.OutboundOrder, error) {
		return []model.OutboundOrder{{ID: 2, OrderNo: "OUT-1", Status: constants.OutboundStatusPending}}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.countAllFn = func(ctx context.Context) (int64, error) { return 40, nil }
	binRepo.countOccFn = func(ctx context.Context) (int64, error) { return 10, nil }
	binRepo.sumRateFn = func(ctx context.Context) (float64, error) { return 400, nil }
	invSvc := &fakeInventorySvc{}
	invSvc.summaryFn = func(ctx context.Context) ([]dto.InventorySummaryItem, error) {
		return []dto.InventorySummaryItem{
			{OwnerID: 1, OwnerName: "货主A", TotalQty: 100, TotalValue: 5000},
			{OwnerID: 2, OwnerName: "货主B", TotalQty: 50, TotalValue: 9000},
		}, nil
	}
	svc := NewDashboardService(inboundRepo, outboundRepo, binRepo, invSvc, testLogger())
	summary, err := svc.Summary(context.Background())
	if err != nil {
		t.Fatalf("summary error: %v", err)
	}
	if summary.TodayInboundCount != 3 || summary.TodayOutboundCount != 2 {
		t.Fatalf("unexpected today counts: %+v", summary)
	}
	if summary.PendingInbound != 4 || summary.PendingOutbound != 0 {
		t.Fatalf("unexpected pending counts: %+v", summary)
	}
	if summary.BinOccupancyRate != 10 {
		t.Fatalf("expected avg occupancy 10, got %v", summary.BinOccupancyRate)
	}
	if len(summary.OwnerTop10) != 2 || summary.OwnerTop10[0].OwnerName != "货主B" {
		t.Fatalf("expected top10 sorted by value desc, got %+v", summary.OwnerTop10)
	}
	if len(summary.PendingTasks) != 2 {
		t.Fatalf("expected 2 pending tasks, got %d", len(summary.PendingTasks))
	}
}

func TestDashboardService_Summary_OwnerScope(t *testing.T) {
	myOwnerID := uint(9)
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.countTodayFn = func(ctx context.Context, ownerID uint) (int64, error) {
		if ownerID != myOwnerID {
			t.Fatalf("expected scoped owner %d, got %d", myOwnerID, ownerID)
		}
		return 1, nil
	}
	inboundRepo.countStatusFn = func(ctx context.Context, status string, ownerID uint) (int64, error) { return 0, nil }
	inboundRepo.listIncFn = func(ctx context.Context, ownerID uint, limit int) ([]model.InboundOrder, error) { return nil, nil }
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.countTodayFn = func(ctx context.Context, ownerID uint) (int64, error) { return 0, nil }
	outboundRepo.countStatusFn = func(ctx context.Context, status string, ownerID uint) (int64, error) { return 0, nil }
	outboundRepo.listIncFn = func(ctx context.Context, ownerID uint, limit int) ([]model.OutboundOrder, error) { return nil, nil }
	binRepo := &fakeBinRepo{}
	binRepo.countAllFn = func(ctx context.Context) (int64, error) { return 40, nil }
	binRepo.countOccFn = func(ctx context.Context) (int64, error) { return 0, nil }
	binRepo.sumRateFn = func(ctx context.Context) (float64, error) { return 0, nil }
	svc := NewDashboardService(inboundRepo, outboundRepo, binRepo, &fakeInventorySvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	if _, err := svc.Summary(ctx); err != nil {
		t.Fatalf("summary error: %v", err)
	}
}
