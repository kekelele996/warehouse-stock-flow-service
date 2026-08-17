package service

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

func newOutboundServiceForTest(outboundRepo repository.OutboundRepository, invSvc InventoryService) OutboundService {
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: id, OwnerID: 1, Name: "五金件", SKU: "HW1", Unit: "盒", StorageRequirement: constants.StorageNormal}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.findByIDFn = func(ctx context.Context, id uint) (*model.BinLocation, error) {
		return &model.BinLocation{ID: id, Area: "B", RackNo: "R01", LayerNo: 1, ColumnNo: 1, Capacity: 8, Status: constants.BinStatusAvailable, StorageRequirement: constants.StorageNormal}, nil
	}
	userRepo := &fakeUserRepo{}
	userRepo.findByIDFn = func(ctx context.Context, id uint) (*model.User, error) {
		return &model.User{ID: id, Name: "操作员", Role: constants.RoleWarehouseManager}, nil
	}
	return NewOutboundService(fakeTx{}, outboundRepo, ownerRepo, productRepo, binRepo, userRepo, invSvc, &fakeAuditSvc{}, testLogger())
}

func TestOutboundService_Create_BinNotExists(t *testing.T) {
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: 1, OwnerID: 1, Name: "五金件", SKU: "HW1", Unit: "盒", StorageRequirement: constants.StorageNormal}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.findByIDFn = func(ctx context.Context, id uint) (*model.BinLocation, error) { return nil, repository.ErrNotFound }
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	svc := NewOutboundService(fakeTx{}, &fakeOutboundRepo{}, ownerRepo, productRepo, binRepo, &fakeUserRepo{}, &fakeInventorySvc{}, &fakeAuditSvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "manager", Role: constants.RoleWarehouseManager})
	_, err := svc.Create(ctx, &dto.OutboundCreateRequest{
		OwnerID: 1, ReceiverName: "收货方",
		Items: []dto.OutboundItemRequest{{ProductID: 1, BinLocationID: 404, ExpectedQty: 5}},
	})
	if err == nil {
		t.Fatalf("expected not found error")
	}
	if util.ToAppError(err).Code != constants.CodeNotFound {
		t.Fatalf("expected not found code, got %d", util.ToAppError(err).Code)
	}
}

func TestOutboundService_Picking_ReducesStock(t *testing.T) {
	order := &model.OutboundOrder{
		ID: 1, OrderNo: "OUT-20260801-0001", OwnerID: 1, Status: constants.OutboundStatusPending,
		Items: []model.OutboundItem{{ID: 21, ProductID: 1, BinLocationID: uintPtr(100), ExpectedQty: 5}},
	}
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.OutboundOrder, error) { return order, nil }
	outboundRepo.updateFn = func(ctx context.Context, o *model.OutboundOrder) error { return nil }
	outboundRepo.updateItemFn = func(ctx context.Context, item *model.OutboundItem) error { return nil }
	var reduceCalls int
	invSvc := &fakeInventorySvc{}
	invSvc.reduceStockFn = func(ctx context.Context, tx *gorm.DB, productID, binID uint, qty int) error {
		reduceCalls++
		return nil
	}
	svc := newOutboundServiceForTest(outboundRepo, invSvc)
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 3, Username: "picker", Role: constants.RolePicker})
	view, err := svc.Picking(ctx, 1, &dto.OutboundPickingRequest{Items: []dto.OutboundPickingItemRequest{{ItemID: 21, ActualQty: 5}}})
	if err != nil {
		t.Fatalf("picking error: %v", err)
	}
	if view.Status != constants.OutboundStatusPicking {
		t.Fatalf("expected Picking, got %s", view.Status)
	}
	if reduceCalls != 1 {
		t.Fatalf("expected 1 reduce stock call, got %d", reduceCalls)
	}
}

func TestOutboundService_Checking_InvalidTransition(t *testing.T) {
	order := &model.OutboundOrder{
		ID: 1, OrderNo: "OUT-20260801-0001", OwnerID: 1, Status: constants.OutboundStatusPending,
		Items: []model.OutboundItem{{ID: 21, ProductID: 1, BinLocationID: uintPtr(100), ExpectedQty: 5}},
	}
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.OutboundOrder, error) { return order, nil }
	svc := newOutboundServiceForTest(outboundRepo, &fakeInventorySvc{})
	_, err := svc.Checking(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected transition error")
	}
	if util.ToAppError(err).Code != constants.CodeBadStateTransition {
		t.Fatalf("expected transition code, got %d", util.ToAppError(err).Code)
	}
}

func uintPtr(v uint) *uint { return &v }
