package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestBinLocationRepository_CRUD(t *testing.T) {
	db := testDB(t)
	repo := NewBinLocationRepository(db)
	ctx := context.Background()
	area := "A"
	bin := newTestBin(area, fmtRack(dbTxSeq()), 1, 1)
	if err := repo.Create(ctx, bin); err != nil {
		t.Fatalf("create bin: %v", err)
	}
	got, err := repo.FindByCode(ctx, area, bin.RackNo, bin.LayerNo, bin.ColumnNo)
	if err != nil {
		t.Fatalf("find bin by code: %v", err)
	}
	got.OccupancyRate = 50
	got.Status = "Occupied"
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update bin: %v", err)
	}
	avail, err := repo.FindAvailable(ctx, "Normal", 5)
	if err != nil {
		t.Fatalf("find available bins: %v", err)
	}
	if len(avail) < 1 {
		t.Fatalf("expected available bins")
	}
	all, err := repo.CountAll(ctx)
	if err != nil {
		t.Fatalf("count all: %v", err)
	}
	if all < 1 {
		t.Fatalf("expected count >= 1")
	}
	_, err = repo.SumOccupancyRate(ctx)
	if err != nil {
		t.Fatalf("sum occupancy: %v", err)
	}
	_ = db.Unscoped().Delete(&bin)
}

func fmtRack(n int64) string {
	return fmt.Sprintf("R%02d", n%100)
}
