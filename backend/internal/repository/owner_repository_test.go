package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestOwnerRepository_CRUD(t *testing.T) {
	db := testDB(t)
	repo := NewOwnerRepository(db)
	ctx := context.Background()
	name := fmt.Sprintf("测试货主%d", dbTxSeq())
	owner := newTestOwner(name)

	if err := repo.Create(ctx, owner); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	if owner.ID == 0 {
		t.Fatalf("expected generated id")
	}
	// 重复名称应返回 ErrDuplicate
	if err := repo.Create(ctx, newTestOwner(name)); err == nil {
		t.Fatalf("expected duplicate error")
	}
	got, err := repo.FindByID(ctx, owner.ID)
	if err != nil {
		t.Fatalf("find owner: %v", err)
	}
	if got.Name != name {
		t.Fatalf("unexpected owner: %s", got.Name)
	}
	got.Phone = "13999999999"
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update owner: %v", err)
	}
	list, total, err := repo.List(ctx, OwnerFilter{Keyword: name, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list owner: %v", err)
	}
	if total < 1 || len(list) < 1 {
		t.Fatalf("expected owner in list, total=%d", total)
	}
	_ = db.Unscoped().Delete(&owner)
}
