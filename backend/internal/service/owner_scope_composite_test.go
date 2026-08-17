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

func TestInventoryList_OwnerScopeApplied(t *testing.T) {
	myOwnerID := uint(6)
	invRepo := &fakeInvRepo{}
	invRepo.listFn = func(ctx context.Context, filter repository.InventoryFilter) ([]model.Inventory, int64, error) {
		if filter.OwnerID != myOwnerID {
			t.Fatalf("expected scoped owner %d, got %d", myOwnerID, filter.OwnerID)
		}
		return []model.Inventory{}, 0, nil
	}
	svc := NewInventoryService(invRepo, &fakeBinRepo{}, &fakeProductRepo{}, &fakeOwnerRepo{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	_, _, err := svc.List(ctx, dto.InventoryQuery{}, 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
}

func TestInboundList_OwnerScopeApplied(t *testing.T) {
	myOwnerID := uint(6)
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.listFn = func(ctx context.Context, filter repository.InboundFilter) ([]model.InboundOrder, int64, error) {
		if filter.OwnerID != myOwnerID {
			t.Fatalf("expected scoped owner %d, got %d", myOwnerID, filter.OwnerID)
		}
		return []model.InboundOrder{}, 0, nil
	}
	svc := newInboundServiceForTest(inboundRepo, &fakeInventorySvc{})
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	_, _, err := svc.List(ctx, dto.InboundQuery{}, 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
}

func TestOutboundList_OwnerScopeApplied(t *testing.T) {
	myOwnerID := uint(6)
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.listFn = func(ctx context.Context, filter repository.OutboundFilter) ([]model.OutboundOrder, int64, error) {
		if filter.OwnerID != myOwnerID {
			t.Fatalf("expected scoped owner %d, got %d", myOwnerID, filter.OwnerID)
		}
		return []model.OutboundOrder{}, 0, nil
	}
	svc := newOutboundServiceForTest(outboundRepo, &fakeInventorySvc{})
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	_, _, err := svc.List(ctx, dto.OutboundQuery{}, 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
}
