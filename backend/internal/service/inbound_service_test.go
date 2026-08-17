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

// newInboundServiceForTest 构造入库服务并预置 toView 所需的关联查询。
func newInboundServiceForTest(inboundRepo repository.InboundRepository, invSvc InventoryService) InboundService {
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: id, OwnerID: 1, Name: "水泥", SKU: "S1", Unit: "袋", StorageRequirement: constants.StorageNormal, Volume: 0.05}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.findByIDFn = func(ctx context.Context, id uint) (*model.BinLocation, error) {
		return &model.BinLocation{ID: id, Area: "A", RackNo: "R01", LayerNo: 1, ColumnNo: 1, Capacity: 8, Status: constants.BinStatusAvailable, StorageRequirement: constants.StorageNormal}, nil
	}
	userRepo := &fakeUserRepo{}
	userRepo.findByIDFn = func(ctx context.Context, id uint) (*model.User, error) {
		return &model.User{ID: id, Name: "操作员", Role: constants.RoleWarehouseManager}, nil
	}
	audit := &fakeAuditSvc{}
	return NewInboundService(fakeTx{}, inboundRepo, ownerRepo, productRepo, binRepo, userRepo, invSvc, audit, testLogger())
}

func TestInboundService_Create_ProductOwnerMismatch(t *testing.T) {
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: 1, OwnerID: 999, Name: "别人的商品", SKU: "S-OTHER", Unit: "箱", StorageRequirement: constants.StorageNormal}, nil
	}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	svc := NewInboundService(fakeTx{}, &fakeInboundRepo{}, ownerRepo, productRepo, &fakeBinRepo{}, &fakeUserRepo{}, &fakeInventorySvc{}, &fakeAuditSvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "manager", Role: constants.RoleWarehouseManager})
	_, err := svc.Create(ctx, &dto.InboundCreateRequest{
		OwnerID: 1, SupplierName: "供应商",
		Items: []dto.InboundItemRequest{{ProductID: 1, BatchNo: "B1", ExpectedQty: 10}},
	})
	if err == nil {
		t.Fatalf("expected mismatch error")
	}
	if util.ToAppError(err).Code != constants.CodeForbidden {
		t.Fatalf("expected forbidden code, got %d", util.ToAppError(err).Code)
	}
}

func TestInboundService_Receive_Success(t *testing.T) {
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0001", OwnerID: 1, Status: constants.InboundStatusPending,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 10}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	inboundRepo.updateFn = func(ctx context.Context, o *model.InboundOrder) error { return nil }
	svc := newInboundServiceForTest(inboundRepo, &fakeInventorySvc{})
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "manager", Role: constants.RoleWarehouseManager})
	view, err := svc.Receive(ctx, 1)
	if err != nil {
		t.Fatalf("receive error: %v", err)
	}
	if view.Status != constants.InboundStatusReceived {
		t.Fatalf("expected Received, got %s", view.Status)
	}
	if view.ActualArrivalDate == nil {
		t.Fatalf("expected actual arrival date set")
	}
}

func TestInboundService_QC_InvalidTransition(t *testing.T) {
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0001", OwnerID: 1, Status: constants.InboundStatusPending,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 10}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	svc := newInboundServiceForTest(inboundRepo, &fakeInventorySvc{})
	_, err := svc.QC(context.Background(), 1, &dto.InboundQCRequest{
		Items: []dto.InboundQCItemRequest{{ItemID: 11, ActualQty: 10, QCResult: constants.QCResultPass}},
	})
	if err == nil {
		t.Fatalf("expected transition error")
	}
	if util.ToAppError(err).Code != constants.CodeBadStateTransition {
		t.Fatalf("expected transition code %d, got %d", constants.CodeBadStateTransition, util.ToAppError(err).Code)
	}
}

func TestInboundService_Shelve_AddsStock(t *testing.T) {
	order := &model.InboundOrder{
		ID: 1, OrderNo: "INB-20260801-0001", OwnerID: 1, Status: constants.InboundStatusQCInProgress,
		Items: []model.InboundItem{{ID: 11, ProductID: 1, BatchNo: "B1", ExpectedQty: 10, ActualQty: 10, QCResult: constants.QCResultPass}},
	}
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.findByIDFn = func(ctx context.Context, id uint) (*model.InboundOrder, error) { return order, nil }
	inboundRepo.updateFn = func(ctx context.Context, o *model.InboundOrder) error { return nil }
	inboundRepo.updateItemFn = func(ctx context.Context, item *model.InboundItem) error { return nil }
	var addCalls int
	invSvc := &fakeInventorySvc{}
	invSvc.addStockFn = func(ctx context.Context, tx *gorm.DB, productID, ownerID, binID uint, batchNo string, qty int) error {
		addCalls++
		return nil
	}
	svc := newInboundServiceForTest(inboundRepo, invSvc)
	view, err := svc.Shelve(context.Background(), 1, &dto.InboundShelveRequest{Items: []dto.InboundShelveItemRequest{{ItemID: 11, BinLocationID: 100}}})
	if err != nil {
		t.Fatalf("shelve error: %v", err)
	}
	if view.Status != constants.InboundStatusShelved {
		t.Fatalf("expected Shelved, got %s", view.Status)
	}
	if addCalls != 1 {
		t.Fatalf("expected 1 add stock call, got %d", addCalls)
	}
}
