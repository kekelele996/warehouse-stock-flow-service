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

func TestProductService_Create_DuplicateSKU(t *testing.T) {
	productRepo := &fakeProductRepo{}
	productRepo.createFn = func(ctx context.Context, p *model.Product) error { return repository.ErrDuplicate }
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	svc := NewProductService(productRepo, ownerRepo, &fakeAuditSvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "admin", Role: constants.RoleAdmin})
	_, err := svc.Create(ctx, &dto.ProductCreateRequest{
		OwnerID: 1, Name: "水泥", SKU: "SKU-001", Unit: "袋", StorageRequirement: constants.StorageNormal,
	})
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
	if util.ToAppError(err).Code != constants.CodeDuplicate {
		t.Fatalf("expected duplicate code, got %d", util.ToAppError(err).Code)
	}
}

func TestProductService_Create_OwnerNotExists(t *testing.T) {
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) { return nil, repository.ErrNotFound }
	svc := NewProductService(&fakeProductRepo{}, ownerRepo, &fakeAuditSvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "admin", Role: constants.RoleAdmin})
	_, err := svc.Create(ctx, &dto.ProductCreateRequest{OwnerID: 999, Name: "水泥", SKU: "SKU-002", Unit: "袋", StorageRequirement: constants.StorageNormal})
	if err == nil {
		t.Fatalf("expected not found error")
	}
	if util.ToAppError(err).Code != constants.CodeNotFound {
		t.Fatalf("expected not found code, got %d", util.ToAppError(err).Code)
	}
}

func TestProductService_List_OwnerScope(t *testing.T) {
	myOwnerID := uint(3)
	productRepo := &fakeProductRepo{}
	productRepo.listFn = func(ctx context.Context, filter repository.ProductFilter) ([]model.Product, int64, error) {
		if filter.OwnerID != myOwnerID {
			t.Fatalf("expected owner scope filter %d, got %d", myOwnerID, filter.OwnerID)
		}
		return []model.Product{{ID: 1, OwnerID: myOwnerID, Name: "水泥", SKU: "S1", Unit: "袋", StorageRequirement: constants.StorageNormal}}, 1, nil
	}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: myOwnerID, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	svc := NewProductService(productRepo, ownerRepo, &fakeAuditSvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	list, total, err := svc.List(ctx, dto.ProductQuery{}, 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].SKU != "S1" {
		t.Fatalf("unexpected list: %+v total=%d", list, total)
	}
}
