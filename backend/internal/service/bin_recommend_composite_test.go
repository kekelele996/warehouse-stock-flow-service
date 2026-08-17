package service

import (
	"context"
	"testing"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
)

func TestBinRecommend_PicksLowestOccupancy(t *testing.T) {
	bins := []model.BinLocation{
		{ID: 1, Area: "A", RackNo: "R01", LayerNo: 12, ColumnNo: 34, Capacity: 100, OccupancyRate: 10, StorageRequirement: constants.StorageColdChain, Status: constants.BinStatusAvailable},
		{ID: 2, Area: "A", RackNo: "R01", LayerNo: 12, ColumnNo: 35, Capacity: 100, OccupancyRate: 80, StorageRequirement: constants.StorageColdChain, Status: constants.BinStatusAvailable},
	}
	binRepo := &fakeBinRepo{}
	binRepo.findAvailFn = func(ctx context.Context, req string, limit int) ([]model.BinLocation, error) {
		return bins, nil
	}
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: 1, OwnerID: 1, Name: "疫苗", SKU: "V1", StorageRequirement: constants.StorageColdChain}, nil
	}
	svc := NewBinLocationService(fakeTx{}, binRepo, productRepo, &fakeInvRepo{}, &fakeOwnerRepo{}, &fakeAuditSvc{}, testLogger())
	view, err := svc.Recommend(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("recommend error: %v", err)
	}
	if view.ID != 1 {
		t.Fatalf("expected lowest occupancy bin id 1, got %d (occupancy %.0f)", view.ID, view.OccupancyRate)
	}
}

func TestBinView_CodeAndStatusText(t *testing.T) {
	bin := &model.BinLocation{ID: 1, Area: "A", RackNo: "R01", LayerNo: 12, ColumnNo: 34, Capacity: 100, OccupancyRate: 10, StorageRequirement: constants.StorageColdChain, Status: constants.BinStatusAvailable}
	view := dto.ToBinView(bin)
	if view.Code != "A-R01-12-34" {
		t.Fatalf("unexpected bin code %q", view.Code)
	}
	if view.StatusText != "可用" {
		t.Fatalf("unexpected status text %q", view.StatusText)
	}
	if view.StorageText != "冷链" {
		t.Fatalf("unexpected storage text %q", view.StorageText)
	}
	if view.StatusText != "可用" {
		t.Fatalf("unexpected status text %q", view.StatusText)
	}
}
