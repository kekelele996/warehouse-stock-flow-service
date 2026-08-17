package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

// OwnerService 货主服务。
type OwnerService interface {
	Create(ctx context.Context, req *dto.OwnerCreateRequest) (*dto.OwnerView, error)
	List(ctx context.Context, query dto.OwnerQuery, page, pageSize int) ([]dto.OwnerView, int64, error)
	Get(ctx context.Context, id uint) (*dto.OwnerView, error)
	GetStats(ctx context.Context, id uint) (*dto.OwnerStatsView, error)
	Update(ctx context.Context, id uint, req *dto.OwnerUpdateRequest) (*dto.OwnerView, error)
	AdjustCredit(ctx context.Context, id uint, req *dto.OwnerCreditRequest) (*dto.OwnerView, error)
	Suspend(ctx context.Context, id uint) (*dto.OwnerView, error)
	Activate(ctx context.Context, id uint) (*dto.OwnerView, error)
}

type ownerService struct {
	ownerRepo     repository.OwnerRepository
	productRepo   repository.ProductRepository
	inboundRepo   repository.InboundRepository
	outboundRepo  repository.OutboundRepository
	inventoryRepo repository.InventoryRepository
	auditSvc      AuditService
	logger        *slog.Logger
}

// NewOwnerService 构造货主服务。
func NewOwnerService(
	ownerRepo repository.OwnerRepository,
	productRepo repository.ProductRepository,
	inboundRepo repository.InboundRepository,
	outboundRepo repository.OutboundRepository,
	inventoryRepo repository.InventoryRepository,
	auditSvc AuditService,
	logger *slog.Logger,
) OwnerService {
	return &ownerService{
		ownerRepo:     ownerRepo,
		productRepo:   productRepo,
		inboundRepo:   inboundRepo,
		outboundRepo:  outboundRepo,
		inventoryRepo: inventoryRepo,
		auditSvc:      auditSvc,
		logger:        logger,
	}
}

func (s *ownerService) Create(ctx context.Context, req *dto.OwnerCreateRequest) (*dto.OwnerView, error) {
	owner := &model.Owner{
		Name:             req.Name,
		ContactName:      req.ContactName,
		Phone:            req.Phone,
		Email:            req.Email,
		Address:          req.Address,
		SettlementMethod: req.SettlementMethod,
		CreditLimit:      req.CreditLimit,
		CurrentDebt:      0,
		Status:           constants.OwnerStatusActive,
	}
	if err := s.ownerRepo.Create(ctx, owner); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.WrapAppError(
				util.NewAppError(constants.CodeDuplicate, 409, "货主名称已存在，请检查货主名称字段"),
				err,
			)
		}
		return nil, fmt.Errorf("owner create: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogCreateOwner, "owner_id", owner.ID, "name", owner.Name, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "owner", "owner.create", "Owner", fmt.Sprintf("%d", owner.ID), "创建货主："+owner.Name)
	return dto.ToOwnerView(owner), nil
}

func (s *ownerService) List(ctx context.Context, query dto.OwnerQuery, page, pageSize int) ([]dto.OwnerView, int64, error) {
	filter := repository.OwnerFilter{Keyword: query.Keyword, Status: query.Status, Page: page, PageSize: pageSize}
	if claims, ok := util.CurrentUser(ctx); ok && claims.Role == constants.RoleOwner {
		if claims.OwnerID == nil {
			return []dto.OwnerView{}, 0, nil
		}
		owner, err := s.ownerRepo.FindByID(ctx, *claims.OwnerID)
		if err != nil {
			return nil, 0, fmt.Errorf("owner list scope: %w", err)
		}
		return []dto.OwnerView{*dto.ToOwnerView(owner)}, 1, nil
	}
	owners, total, err := s.ownerRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("owner list: %w", err)
	}
	views := make([]dto.OwnerView, 0, len(owners))
	for i := range owners {
		views = append(views, *dto.ToOwnerView(&owners[i]))
	}
	return views, total, nil
}

func (s *ownerService) Get(ctx context.Context, id uint) (*dto.OwnerView, error) {
	if err := s.checkScope(ctx, id); err != nil {
		return nil, err
	}
	owner, err := s.ownerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "货主不存在"), err)
		}
		return nil, fmt.Errorf("owner get: %w", err)
	}
	return dto.ToOwnerView(owner), nil
}

func (s *ownerService) GetStats(ctx context.Context, id uint) (*dto.OwnerStatsView, error) {
	if err := s.checkScope(ctx, id); err != nil {
		return nil, err
	}
	owner, err := s.ownerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "货主不存在"), err)
		}
		return nil, fmt.Errorf("owner stats: %w", err)
	}
	productCount, err := s.productRepo.CountByOwner(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("owner stats product count: %w", err)
	}
	inboundCount, err := s.inboundRepo.CountByOwner(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("owner stats inbound count: %w", err)
	}
	outboundCount, err := s.outboundRepo.CountByOwner(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("owner stats outbound count: %w", err)
	}
	value, err := s.inventoryRepo.SumValueByOwnerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("owner stats inventory value: %w", err)
	}
	return &dto.OwnerStatsView{
		Owner:          dto.ToOwnerView(owner),
		ProductCount:   productCount,
		InboundCount:   inboundCount,
		OutboundCount:  outboundCount,
		InventoryValue: value,
	}, nil
}

func (s *ownerService) Update(ctx context.Context, id uint, req *dto.OwnerUpdateRequest) (*dto.OwnerView, error) {
	if err := s.checkScope(ctx, id); err != nil {
		return nil, err
	}
	owner, err := s.ownerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "货主不存在"), err)
		}
		return nil, fmt.Errorf("owner update find: %w", err)
	}
	owner.ContactName = req.ContactName
	owner.Phone = req.Phone
	owner.Email = req.Email
	owner.Address = req.Address
	if err := s.ownerRepo.Update(ctx, owner); err != nil {
		return nil, fmt.Errorf("owner update: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogUpdateOwner, "owner_id", owner.ID, "name", owner.Name, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "owner", "owner.update", "Owner", fmt.Sprintf("%d", owner.ID), "更新货主："+owner.Name)
	return dto.ToOwnerView(owner), nil
}

func (s *ownerService) AdjustCredit(ctx context.Context, id uint, req *dto.OwnerCreditRequest) (*dto.OwnerView, error) {
	if err := s.checkScope(ctx, id); err != nil {
		return nil, err
	}
	owner, err := s.ownerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "货主不存在"), err)
		}
		return nil, fmt.Errorf("owner adjust credit find: %w", err)
	}
	owner.CreditLimit = req.CreditLimit
	if err := s.ownerRepo.Update(ctx, owner); err != nil {
		return nil, fmt.Errorf("owner adjust credit: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogAdjustOwnerCredit, "owner_id", owner.ID, "credit_limit", owner.CreditLimit, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "owner", "owner.credit", "Owner", fmt.Sprintf("%d", owner.ID), "调整货主信用额度为 "+util.FormatAmount(req.CreditLimit))
	return dto.ToOwnerView(owner), nil
}

func (s *ownerService) Suspend(ctx context.Context, id uint) (*dto.OwnerView, error) {
	owner, err := s.ownerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "货主不存在"), err)
		}
		return nil, fmt.Errorf("owner suspend find: %w", err)
	}
	owner.Status = constants.OwnerStatusSuspended
	if err := s.ownerRepo.Update(ctx, owner); err != nil {
		return nil, fmt.Errorf("owner suspend: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogSuspendOwner, "owner_id", owner.ID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "owner", "owner.suspend", "Owner", fmt.Sprintf("%d", owner.ID), "暂停货主合作："+owner.Name)
	return dto.ToOwnerView(owner), nil
}

func (s *ownerService) Activate(ctx context.Context, id uint) (*dto.OwnerView, error) {
	owner, err := s.ownerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "货主不存在"), err)
		}
		return nil, fmt.Errorf("owner activate find: %w", err)
	}
	owner.Status = constants.OwnerStatusActive
	if err := s.ownerRepo.Update(ctx, owner); err != nil {
		return nil, fmt.Errorf("owner activate: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogActivateOwner, "owner_id", owner.ID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "owner", "owner.activate", "Owner", fmt.Sprintf("%d", owner.ID), "恢复货主合作："+owner.Name)
	return dto.ToOwnerView(owner), nil
}

// checkScope Owner 角色只能访问自己的货主数据。
func (s *ownerService) checkScope(ctx context.Context, ownerID uint) error {
	claims, ok := util.CurrentUser(ctx)
	if !ok {
		return util.NewAppError(constants.CodeUnauthorized, 401, constants.MsgUnauthorized)
	}
	if claims.Role == constants.RoleOwner {
		if claims.OwnerID == nil || *claims.OwnerID != ownerID {
			return util.NewAppError(constants.CodeRoleForbidden, 403, "货主角色无权访问其他货主的数据")
		}
	}
	return nil
}

// currentUsername 获取当前用户名。
func currentUsername(ctx context.Context) string {
	if claims, ok := util.CurrentUser(ctx); ok {
		return claims.Username
	}
	return "anonymous"
}

// currentRole 获取当前角色。
func currentRole(ctx context.Context) string {
	if claims, ok := util.CurrentUser(ctx); ok {
		return claims.Role
	}
	return "anonymous"
}
