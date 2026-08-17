package service

import (
	"context"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/util"
)

func TestOutboundGet_PendingDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Get on pending outbound panicked: %v", r)
		}
	}()
	order := &model.OutboundOrder{
		ID: 1, OrderNo: "OUT-20260801-0001", OwnerID: 1, Status: constants.OutboundStatusPending,
		Items: []model.OutboundItem{{ID: 21, ProductID: 1, BinLocationID: uintPtr(100), ExpectedQty: 5}},
	}
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.OutboundOrder, error) { return order, nil }
	svc := newOutboundServiceForTest(outboundRepo, &fakeInventorySvc{})
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "manager", Role: constants.RoleWarehouseManager})
	view, err := svc.Get(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view == nil || view.Status != constants.OutboundStatusPending {
		t.Fatalf("unexpected view: %+v", view)
	}
}

func TestInboundGet_QCNotDoneDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Get on pending inbound panicked: %v", r)
		}
	}()
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0001", OwnerID: 1, Status: constants.InboundStatusPending,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 10}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	svc := newInboundServiceForTest(inboundRepo, &fakeInventorySvc{})
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "manager", Role: constants.RoleWarehouseManager})
	view, err := svc.Get(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view == nil || view.Status != constants.InboundStatusPending {
		t.Fatalf("unexpected view: %+v", view)
	}
}
