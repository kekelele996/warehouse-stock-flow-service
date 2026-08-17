package service

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/util"
)

func TestOutboundPicking_Composite(t *testing.T) {
	order := &model.OutboundOrder{
		ID: 1, OrderNo: "OUT-20260801-0001", OwnerID: 1, Status: constants.OutboundStatusPending,
		Items: []model.OutboundItem{{ID: 21, ProductID: 1, BinLocationID: uintPtr(100), ExpectedQty: 5}},
	}
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.OutboundOrder, error) { return order, nil }
	outboundRepo.updateFn = func(ctx context.Context, o *model.OutboundOrder) error { return nil }
	outboundRepo.updateItemFn = func(ctx context.Context, item *model.OutboundItem) error { return nil }

	var reducedQty int
	var reduceCalls int
	invSvc := &fakeInventorySvc{}
	invSvc.reduceStockFn = func(ctx context.Context, tx *gorm.DB, productID, binID uint, qty int) error {
		reduceCalls++
		reducedQty = qty
		return nil
	}

	svc := newOutboundServiceForTest(outboundRepo, invSvc)
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 3, Username: "picker", Role: constants.RolePicker})
	view, err := svc.Picking(ctx, 1, &dto.OutboundPickingRequest{Items: []dto.OutboundPickingItemRequest{{ItemID: 21, ActualQty: 3}}})
	if err != nil {
		t.Fatalf("picking error: %v", err)
	}
	if view.StatusText != "拣货中" {
		t.Fatalf("unexpected status text %q", view.StatusText)
	}
	if reduceCalls != 1 || reducedQty != 3 {
		t.Fatalf("expected reduce once with actual qty 3, got calls=%d qty=%d", reduceCalls, reducedQty)
	}
	if len(view.Items) != 1 || view.Items[0].ActualQty != 3 || view.Items[0].ExpectedQty != 5 {
		t.Fatalf("unexpected item view: %+v", view.Items)
	}
}

func TestOutboundChecking_FromPendingRejected(t *testing.T) {
	order := &model.OutboundOrder{
		ID: 1, OrderNo: "OUT-20260801-0002", OwnerID: 1, Status: constants.OutboundStatusPending,
		Items: []model.OutboundItem{{ID: 21, ProductID: 1, BinLocationID: uintPtr(100), ExpectedQty: 5}},
	}
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.OutboundOrder, error) { return order, nil }
	outboundRepo.updateFn = func(ctx context.Context, o *model.OutboundOrder) error { return nil }
	svc := newOutboundServiceForTest(outboundRepo, &fakeInventorySvc{})
	_, err := svc.Checking(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected transition error")
	}
	if util.ToAppError(err).Code != constants.CodeBadStateTransition {
		t.Fatalf("expected transition code %d, got %d", constants.CodeBadStateTransition, util.ToAppError(err).Code)
	}
}
