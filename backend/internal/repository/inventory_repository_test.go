package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
)

func TestInventoryRepository_AddDecrement(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	ownerRepo := NewOwnerRepository(db)
	owner := newTestOwner(fmt.Sprintf("货主V%d", dbTxSeq()))
	if err := ownerRepo.Create(ctx, owner); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	productRepo := NewProductRepository(db)
	product := newTestProduct(owner.ID, fmt.Sprintf("SKU-V-%d", dbTxSeq()))
	if err := productRepo.Create(ctx, product); err != nil {
		t.Fatalf("create product: %v", err)
	}
	binRepo := NewBinLocationRepository(db)
	bin := newTestBin("A", fmtRack(dbTxSeq()), 2, 2)
	if err := binRepo.Create(ctx, bin); err != nil {
		t.Fatalf("create bin: %v", err)
	}
	repo := NewInventoryRepository(db)
	inv := &model.Inventory{
		ProductID: product.ID, OwnerID: owner.ID, BinLocationID: bin.ID, BatchNo: "VB1", Quantity: 100,
	}
	if err := repo.Create(ctx, inv); err != nil {
		t.Fatalf("create inventory: %v", err)
	}
	got, err := repo.FindByProductAndBin(ctx, product.ID, bin.ID, "VB1")
	if err != nil {
		t.Fatalf("find inventory: %v", err)
	}
	if got.Quantity != 100 {
		t.Fatalf("unexpected quantity: %d", got.Quantity)
	}
	if err := repo.DecrementQuantity(ctx, inv.ID, 30); err != nil {
		t.Fatalf("decrement: %v", err)
	}
	if err := repo.DecrementQuantity(ctx, inv.ID, 1000); err == nil {
		t.Fatalf("expected insufficient stock error")
	}
	list, total, err := repo.List(ctx, InventoryFilter{OwnerID: owner.ID, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list inventory: %v", err)
	}
	if total < 1 || len(list) < 1 {
		t.Fatalf("expected inventory in list")
	}
	rows, err := repo.SumValueByOwner(ctx)
	if err != nil {
		t.Fatalf("sum value by owner: %v", err)
	}
	_ = rows
	_ = db.Unscoped().Delete(&inv)
	_ = db.Unscoped().Delete(&bin)
	_ = db.Unscoped().Delete(&product)
	_ = db.Unscoped().Delete(&owner)
}

var _ = constants.StorageNormal
