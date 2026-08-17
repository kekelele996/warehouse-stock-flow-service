package service

import (
	"context"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

func TestInventoryService_Summary(t *testing.T) {
	invRepo := &fakeInvRepo{}
	invRepo.sumOwnerFn = func(ctx context.Context) ([]repository.InventorySummaryRow, error) {
		return []repository.InventorySummaryRow{
			{OwnerID: 1, OwnerName: "货主A", TotalQty: 100, TotalValue: 5000},
			{OwnerID: 2, OwnerName: "货主B", TotalQty: 50, TotalValue: 3000},
		}, nil
	}
	svc := NewInventoryService(invRepo, &fakeBinRepo{}, &fakeProductRepo{}, &fakeOwnerRepo{}, testLogger())
	items, err := svc.Summary(context.Background())
	if err != nil {
		t.Fatalf("summary error: %v", err)
	}
	if len(items) != 2 || items[0].OwnerName != "货主A" || items[0].TotalValue != 5000 {
		t.Fatalf("unexpected summary: %+v", items)
	}
}

func TestInventoryService_List_OwnerScope(t *testing.T) {
	myOwnerID := uint(6)
	invRepo := &fakeInvRepo{}
	invRepo.listFn = func(ctx context.Context, filter repository.InventoryFilter) ([]model.Inventory, int64, error) {
		if filter.OwnerID != myOwnerID {
			t.Fatalf("expected scoped owner %d, got %d", myOwnerID, filter.OwnerID)
		}
		return []model.Inventory{{ID: 1, ProductID: 1, OwnerID: myOwnerID, BinLocationID: 1, BatchNo: "B1", Quantity: 10}}, 1, nil
	}
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: 1, OwnerID: myOwnerID, Name: "水泥", SKU: "S1", Unit: "袋", StorageRequirement: constants.StorageNormal, Price: 26.5}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.findByIDFn = func(ctx context.Context, id uint) (*model.BinLocation, error) {
		return &model.BinLocation{ID: 1, Area: "A", RackNo: "R01", LayerNo: 1, ColumnNo: 1, Capacity: 8, Status: constants.BinStatusOccupied, StorageRequirement: constants.StorageNormal}, nil
	}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: myOwnerID, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	svc := NewInventoryService(invRepo, binRepo, productRepo, ownerRepo, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	list, total, err := svc.List(ctx, dto.InventoryQuery{}, 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].TotalValue != 265 {
		t.Fatalf("unexpected list: %+v total=%d", list, total)
	}
}
