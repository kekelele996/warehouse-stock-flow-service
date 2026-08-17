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

func TestInboundShelve_PassAddsActualQty(t *testing.T) {
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0001", OwnerID: 1, Status: constants.InboundStatusQCInProgress,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 12, ActualQty: 10, QCResult: constants.QCResultPass}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	inboundRepo.updateFn = func(ctx context.Context, o *model.InboundOrder) error { return nil }
	inboundRepo.updateItemFn = func(ctx context.Context, item *model.InboundItem) error { return nil }
	var addQty int
	var addCalls int
	invSvc := &fakeInventorySvc{}
	invSvc.addStockFn = func(ctx context.Context, tx *gorm.DB, productID, ownerID, binID uint, batchNo string, qty int) error {
		addCalls++
		addQty = qty
		return nil
	}
	svc := newInboundServiceForTest(inboundRepo, invSvc)
	view, err := svc.Shelve(context.Background(), 1, &dto.InboundShelveRequest{Items: []dto.InboundShelveItemRequest{{ItemID: 11, BinLocationID: 100}}})
	if err != nil {
		t.Fatalf("shelve error: %v", err)
	}
	if addCalls != 1 || addQty != 10 {
		t.Fatalf("expected add stock once with actual qty 10, got calls=%d qty=%d", addCalls, addQty)
	}
	if len(view.Items) != 1 || view.Items[0].ActualQty != 10 || view.Items[0].QCResultText != "合格" {
		t.Fatalf("unexpected item view: %+v", view.Items)
	}
}

func TestInboundShelve_FailRejected(t *testing.T) {
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0002", OwnerID: 1, Status: constants.InboundStatusQCInProgress,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 12, ActualQty: 10, QCResult: constants.QCResultFail}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	inboundRepo.updateFn = func(ctx context.Context, o *model.InboundOrder) error { return nil }
	inboundRepo.updateItemFn = func(ctx context.Context, item *model.InboundItem) error { return nil }
	svc := newInboundServiceForTest(inboundRepo, &fakeInventorySvc{})
	_, err := svc.Shelve(context.Background(), 1, &dto.InboundShelveRequest{Items: []dto.InboundShelveItemRequest{{ItemID: 11, BinLocationID: 100}}})
	if err == nil {
		t.Fatalf("expected fail item to be rejected")
	}
	if code := util.ToAppError(err).Code; code != constants.CodeConflict {
		t.Fatalf("expected conflict code, got %d", code)
	}
}

func TestInboundComplete_FromQCRejected(t *testing.T) {
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0003", OwnerID: 1, Status: constants.InboundStatusQCInProgress,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 12, ActualQty: 10, QCResult: constants.QCResultPass}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	inboundRepo.updateFn = func(ctx context.Context, o *model.InboundOrder) error { return nil }
	svc := newInboundServiceForTest(inboundRepo, &fakeInventorySvc{})
	_, err := svc.Complete(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected transition error")
	}
	if code := util.ToAppError(err).Code; code != constants.CodeBadStateTransition {
		t.Fatalf("expected bad state transition code, got %d", code)
	}
}
