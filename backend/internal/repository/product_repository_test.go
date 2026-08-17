package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestProductRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	ownerRepo := NewOwnerRepository(db)
	owner := newTestOwner(fmt.Sprintf("货主P%d", dbTxSeq()))
	if err := ownerRepo.Create(ctx, owner); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	repo := NewProductRepository(db)
	sku := fmt.Sprintf("SKU-T-%d", dbTxSeq())
	product := newTestProduct(owner.ID, sku)
	if err := repo.Create(ctx, product); err != nil {
		t.Fatalf("create product: %v", err)
	}
	if err := repo.Create(ctx, newTestProduct(owner.ID, sku)); err == nil {
		t.Fatalf("expected duplicate sku error")
	}
	got, err := repo.FindBySKU(ctx, sku)
	if err != nil {
		t.Fatalf("find product by sku: %v", err)
	}
	got.Price = 99.9
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update product: %v", err)
	}
	list, total, err := repo.List(ctx, ProductFilter{OwnerID: owner.ID, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list products: %v", err)
	}
	if total < 1 || len(list) < 1 {
		t.Fatalf("expected products in list, total=%d", total)
	}
	count, err := repo.CountByOwner(ctx, owner.ID)
	if err != nil {
		t.Fatalf("count by owner: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected count >= 1")
	}
	_ = db.Unscoped().Delete(&product)
	_ = db.Unscoped().Delete(&owner)
}
