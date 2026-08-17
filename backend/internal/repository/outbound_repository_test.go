package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
)

func TestOutboundRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	ownerRepo := NewOwnerRepository(db)
	owner := newTestOwner(fmt.Sprintf("货主O%d", dbTxSeq()))
	if err := ownerRepo.Create(ctx, owner); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	repo := NewOutboundRepository(db)
	order := &model.OutboundOrder{
		OrderNo: fmt.Sprintf("OUT-TEST-%d", dbTxSeq()), OwnerID: owner.ID,
		ReceiverName: "测试收货方", Status: constants.OutboundStatusPending,
		Items: []model.OutboundItem{{ProductID: 1, BinLocationID: nil, ExpectedQty: 5}},
	}
	if err := repo.Create(ctx, order); err != nil {
		t.Fatalf("create outbound: %v", err)
	}
	got, err := repo.FindByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("find outbound: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("expected preloaded items")
	}
	got.Status = constants.OutboundStatusPicking
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update outbound: %v", err)
	}
	list, total, err := repo.List(ctx, OutboundFilter{OwnerID: owner.ID, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list outbound: %v", err)
	}
	if total < 1 || len(list) < 1 {
		t.Fatalf("expected outbound in list")
	}
	inc, err := repo.ListIncomplete(ctx, 0, 5)
	if err != nil {
		t.Fatalf("list incomplete: %v", err)
	}
	_ = inc
	_ = db.Unscoped().Delete(&order)
	_ = db.Unscoped().Delete(&owner)
}
