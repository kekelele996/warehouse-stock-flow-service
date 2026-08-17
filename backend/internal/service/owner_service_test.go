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

func newOwnerServiceForTest(ownerRepo repository.OwnerRepository) OwnerService {
	return NewOwnerService(ownerRepo, &fakeProductRepo{}, &fakeInboundRepo{}, &fakeOutboundRepo{}, &fakeInvRepo{}, &fakeAuditSvc{}, testLogger())
}

func TestOwnerService_Create_Duplicate(t *testing.T) {
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.createFn = func(ctx context.Context, o *model.Owner) error { return repository.ErrDuplicate }
	svc := newOwnerServiceForTest(ownerRepo)
	_, err := svc.Create(context.Background(), &dto.OwnerCreateRequest{Name: "重复货主", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly})
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
	appErr := util.ToAppError(err)
	if appErr.Code != constants.CodeDuplicate {
		t.Fatalf("expected code %d, got %d", constants.CodeDuplicate, appErr.Code)
	}
}

func TestOwnerService_List_OwnerScope(t *testing.T) {
	myOwnerID := uint(5)
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) {
		return &model.Owner{ID: myOwnerID, Name: "我的货主", ContactName: "李四", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}, nil
	}
	svc := newOwnerServiceForTest(ownerRepo)
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "owner01", Role: constants.RoleOwner, OwnerID: &myOwnerID})
	list, total, err := svc.List(ctx, dto.OwnerQuery{}, 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != myOwnerID {
		t.Fatalf("owner scope failed: total=%d list=%+v", total, list)
	}
}

func TestOwnerService_Suspend(t *testing.T) {
	owner := &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) { return owner, nil }
	ownerRepo.updateFn = func(ctx context.Context, o *model.Owner) error { return nil }
	svc := newOwnerServiceForTest(ownerRepo)
	view, err := svc.Suspend(context.Background(), 1)
	if err != nil {
		t.Fatalf("suspend error: %v", err)
	}
	if view.Status != constants.OwnerStatusSuspended {
		t.Fatalf("expected suspended, got %s", view.Status)
	}
}

func TestOwnerService_AdjustCredit(t *testing.T) {
	owner := &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) { return owner, nil }
	ownerRepo.updateFn = func(ctx context.Context, o *model.Owner) error { return nil }
	svc := newOwnerServiceForTest(ownerRepo)
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "admin", Role: constants.RoleAdmin})
	view, err := svc.AdjustCredit(ctx, 1, &dto.OwnerCreditRequest{CreditLimit: 999999})
	if err != nil {
		t.Fatalf("adjust credit error: %v", err)
	}
	if view.CreditLimit != 999999 {
		t.Fatalf("expected credit limit 999999, got %v", view.CreditLimit)
	}
}

func TestOwnerService_GetStats(t *testing.T) {
	owner := &model.Owner{ID: 1, Name: "货主A", ContactName: "张三", Phone: "138", SettlementMethod: constants.SettlementMonthly, Status: constants.OwnerStatusActive}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.findByIDFn = func(ctx context.Context, id uint) (*model.Owner, error) { return owner, nil }
	productRepo := &fakeProductRepo{}
	productRepo.countFn = func(ctx context.Context, ownerID uint) (int64, error) { return 5, nil }
	inboundRepo := &fakeInboundRepo{}
	inboundRepo.countByOwnerFn = func(ctx context.Context, ownerID uint) (int64, error) { return 3, nil }
	outboundRepo := &fakeOutboundRepo{}
	outboundRepo.countByOwnerFn = func(ctx context.Context, ownerID uint) (int64, error) { return 7, nil }
	invRepo := &fakeInvRepo{}
	invRepo.sumOwnerIDFn = func(ctx context.Context, ownerID uint) (float64, error) { return 8888.5, nil }
	svc := NewOwnerService(ownerRepo, productRepo, inboundRepo, outboundRepo, invRepo, &fakeAuditSvc{}, testLogger())
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 2, Username: "admin", Role: constants.RoleAdmin})
	stats, err := svc.GetStats(ctx, 1)
	if err != nil {
		t.Fatalf("stats error: %v", err)
	}
	if stats.ProductCount != 5 || stats.InboundCount != 3 || stats.InventoryValue != 8888.5 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}
