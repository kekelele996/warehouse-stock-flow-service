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

func newBinServiceForTest(binRepo repository.BinLocationRepository) BinLocationService {
	return NewBinLocationService(fakeTx{}, binRepo, &fakeProductRepo{}, &fakeInvRepo{}, &fakeOwnerRepo{}, &fakeAuditSvc{}, testLogger())
}

func TestBinLocationService_Create_DuplicateCode(t *testing.T) {
	binRepo := &fakeBinRepo{}
	binRepo.findByCodeFn = func(ctx context.Context, area, rackNo string, layerNo, columnNo int) (*model.BinLocation, error) {
		return &model.BinLocation{ID: 1, Area: area, RackNo: rackNo, LayerNo: layerNo, ColumnNo: columnNo, Capacity: 8}, nil
	}
	svc := newBinServiceForTest(binRepo)
	_, err := svc.Create(context.Background(), &dto.BinCreateRequest{Area: "A", RackNo: "R01", LayerNo: 1, ColumnNo: 1, Capacity: 8, StorageRequirement: constants.StorageNormal})
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
	if util.ToAppError(err).Code != constants.CodeDuplicate {
		t.Fatalf("expected duplicate code, got %d", util.ToAppError(err).Code)
	}
}

func TestBinLocationService_Create_Success(t *testing.T) {
	binRepo := &fakeBinRepo{}
	binRepo.findByCodeFn = func(ctx context.Context, area, rackNo string, layerNo, columnNo int) (*model.BinLocation, error) {
		return nil, repository.ErrNotFound
	}
	binRepo.createFn = func(ctx context.Context, b *model.BinLocation) error {
		b.ID = 55
		return nil
	}
	svc := newBinServiceForTest(binRepo)
	view, err := svc.Create(context.Background(), &dto.BinCreateRequest{Area: "B", RackNo: "R02", LayerNo: 2, ColumnNo: 3, Capacity: 10, StorageRequirement: constants.StorageColdChain})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if view.Code != "B-R02-2-3" {
		t.Fatalf("unexpected bin code: %s", view.Code)
	}
	if view.Status != constants.BinStatusAvailable {
		t.Fatalf("expected available status")
	}
}

func TestBinLocationService_Recommend_NoAvailable(t *testing.T) {
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: 1, Name: "水泥", SKU: "S1", Unit: "袋", StorageRequirement: constants.StorageNormal}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.findAvailFn = func(ctx context.Context, req string, limit int) ([]model.BinLocation, error) {
		return []model.BinLocation{}, nil
	}
	svc := NewBinLocationService(fakeTx{}, binRepo, productRepo, &fakeInvRepo{}, &fakeOwnerRepo{}, &fakeAuditSvc{}, testLogger())
	_, err := svc.Recommend(context.Background(), 1, 10)
	if err == nil {
		t.Fatalf("expected conflict error")
	}
	if util.ToAppError(err).Code != constants.CodeConflict {
		t.Fatalf("expected conflict code, got %d", util.ToAppError(err).Code)
	}
}

func TestBinLocationService_Recommend_ReturnsLowestOccupancy(t *testing.T) {
	productRepo := &fakeProductRepo{}
	productRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Product, error) {
		return &model.Product{ID: 1, Name: "水泥", SKU: "S1", Unit: "袋", StorageRequirement: constants.StorageFragile}, nil
	}
	binRepo := &fakeBinRepo{}
	binRepo.findAvailFn = func(ctx context.Context, req string, limit int) ([]model.BinLocation, error) {
		if req != constants.StorageFragile {
			t.Fatalf("expected fragile requirement filter, got %s", req)
		}
		return []model.BinLocation{
			{ID: 2, Area: "D", RackNo: "R01", LayerNo: 2, ColumnNo: 3, Capacity: 8, OccupancyRate: 5},
			{ID: 1, Area: "D", RackNo: "R01", LayerNo: 1, ColumnNo: 3, Capacity: 8, OccupancyRate: 30},
		}, nil
	}
	svc := NewBinLocationService(fakeTx{}, binRepo, productRepo, &fakeInvRepo{}, &fakeOwnerRepo{}, &fakeAuditSvc{}, testLogger())
	view, err := svc.Recommend(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("recommend error: %v", err)
	}
	if view.ID != 2 {
		t.Fatalf("expected bin 2 (lowest occupancy), got %d", view.ID)
	}
}
