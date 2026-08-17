package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
)

func TestInboundRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	ownerRepo := NewOwnerRepository(db)
	owner := newTestOwner(fmt.Sprintf("货主I%d", dbTxSeq()))
	if err := ownerRepo.Create(ctx, owner); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	repo := NewInboundRepository(db)
	order := &model.InboundOrder{
		OrderNo: fmt.Sprintf("INB-TEST-%d", dbTxSeq()), OwnerID: owner.ID,
		SupplierName: "测试供应商", Status: constants.InboundStatusPending,
		Items: []model.InboundItem{{ProductID: 1, BatchNo: "BT1", ExpectedQty: 10}},
	}
	if err := repo.Create(ctx, order); err != nil {
		t.Fatalf("create inbound: %v", err)
	}
	if len(order.Items) == 0 || order.Items[0].InboundOrderID == 0 {
		t.Fatalf("expected items created with order id")
	}
	got, err := repo.FindByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("find inbound: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("expected preloaded items")
	}
	got.Status = constants.InboundStatusReceived
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update inbound: %v", err)
	}
	list, total, err := repo.List(ctx, InboundFilter{OwnerID: owner.ID, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list inbound: %v", err)
	}
	if total < 1 || len(list) < 1 {
		t.Fatalf("expected inbound in list")
	}
	n, err := repo.CountByOwner(ctx, owner.ID)
	if err != nil {
		t.Fatalf("count by owner: %v", err)
	}
	if n < 1 {
		t.Fatalf("expected count >= 1")
	}
	_ = db.Unscoped().Delete(&order)
	_ = db.Unscoped().Delete(&owner)
}
