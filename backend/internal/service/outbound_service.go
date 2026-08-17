package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

// OutboundService 出库单服务。
type OutboundService interface {
	Create(ctx context.Context, req *dto.OutboundCreateRequest) (*dto.OutboundOrderView, error)
	List(ctx context.Context, query dto.OutboundQuery, page, pageSize int) ([]dto.OutboundOrderView, int64, error)
	Get(ctx context.Context, id uint) (*dto.OutboundOrderView, error)
	Picking(ctx context.Context, id uint, req *dto.OutboundPickingRequest) (*dto.OutboundOrderView, error)
	Checking(ctx context.Context, id uint) (*dto.OutboundOrderView, error)
	Packing(ctx context.Context, id uint) (*dto.OutboundOrderView, error)
	Ship(ctx context.Context, id uint, req *dto.OutboundShipRequest) (*dto.OutboundOrderView, error)
	Complete(ctx context.Context, id uint) (*dto.OutboundOrderView, error)
}

type outboundService struct {
	db           txProvider
	outboundRepo repository.OutboundRepository
	ownerRepo    repository.OwnerRepository
	productRepo  repository.ProductRepository
	binRepo      repository.BinLocationRepository
	userRepo     repository.UserRepository
	inventorySvc InventoryService
	auditSvc     AuditService
	logger       *slog.Logger
}

// NewOutboundService 构造出库单服务。
func NewOutboundService(
	db txProvider,
	outboundRepo repository.OutboundRepository,
	ownerRepo repository.OwnerRepository,
	productRepo repository.ProductRepository,
	binRepo repository.BinLocationRepository,
	userRepo repository.UserRepository,
	inventorySvc InventoryService,
	auditSvc AuditService,
	logger *slog.Logger,
) OutboundService {
	return &outboundService{db: db, outboundRepo: outboundRepo, ownerRepo: ownerRepo, productRepo: productRepo, binRepo: binRepo, userRepo: userRepo, inventorySvc: inventorySvc, auditSvc: auditSvc, logger: logger}
}

func (s *outboundService) Create(ctx context.Context, req *dto.OutboundCreateRequest) (*dto.OutboundOrderView, error) {
	owner, err := s.ownerRepo.FindByID(ctx, req.OwnerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "所属货主不存在，请检查货主ID字段"), err)
		}
		return nil, fmt.Errorf("outbound create owner: %w", err)
	}
	if owner.Status == constants.OwnerStatusSuspended {
		return nil, util.NewAppError(constants.CodeOwnerSuspended, 409, "货主已暂停合作，无法创建出库单")
	}
	if err := s.checkOwnerScope(ctx, req.OwnerID); err != nil {
		return nil, err
	}
	claims, _ := util.CurrentUser(ctx)
	order := &model.OutboundOrder{
		OrderNo:          util.GenerateOrderNo("OUT"),
		OwnerID:          req.OwnerID,
		ReceiverName:     req.ReceiverName,
		ReceiverAddress:  req.ReceiverAddress,
		RequiredShipDate: req.RequiredShipDate,
		Status:           constants.OutboundStatusPending,
		CreatedBy:        claimsUserID(claims),
		Items:            make([]model.OutboundItem, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "商品不存在，请检查商品ID字段"), err)
			}
			return nil, fmt.Errorf("outbound create product: %w", err)
		}
		if product.OwnerID != req.OwnerID {
			return nil, util.NewAppError(constants.CodeForbidden, 403, "商品不属于所选货主，请检查商品与货主匹配关系")
		}
		if _, err := s.binRepo.FindByID(ctx, item.BinLocationID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "库位不存在，请检查库位ID字段"), err)
			}
			return nil, fmt.Errorf("outbound create bin: %w", err)
		}
		order.Items = append(order.Items, model.OutboundItem{
			ProductID:     item.ProductID,
			BinLocationID: &item.BinLocationID,
			ExpectedQty:   item.ExpectedQty,
		})
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		return s.outboundRepo.WithTx(tx).Create(ctx, order)
	})
	if err != nil {
		return nil, fmt.Errorf("outbound create: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogCreateOutbound, "order_no", order.OrderNo, "owner_id", order.OwnerID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "outbound", "outbound.create", "OutboundOrder", order.OrderNo, "创建出库单："+order.OrderNo)
	return s.toView(ctx, order), nil
}

func (s *outboundService) List(ctx context.Context, query dto.OutboundQuery, page, pageSize int) ([]dto.OutboundOrderView, int64, error) {
	if claims, ok := util.CurrentUser(ctx); ok && claims.Role == constants.RoleOwner && claims.OwnerID != nil {
		query.OwnerID = *claims.OwnerID
	}
	filter := repository.OutboundFilter{Status: query.Status, OwnerID: query.OwnerID, Keyword: query.Keyword, Page: page, PageSize: pageSize}
	orders, total, err := s.outboundRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("outbound list: %w", err)
	}
	views := make([]dto.OutboundOrderView, 0, len(orders))
	for i := range orders {
		views = append(views, *s.toView(ctx, &orders[i]))
	}
	return views, total, nil
}

func (s *outboundService) Get(ctx context.Context, id uint) (*dto.OutboundOrderView, error) {
	order, err := s.outboundRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "出库单不存在"), err)
		}
		return nil, fmt.Errorf("outbound get: %w", err)
	}
	if err := s.checkOwnerScope(ctx, order.OwnerID); err != nil {
		return nil, err
	}
	return s.toView(ctx, order), nil
}

// Picking 拣货：Pending → Picking，确认实拣数量并扣减库存（SELECT ... FOR UPDATE）。
func (s *outboundService) Picking(ctx context.Context, id uint, req *dto.OutboundPickingRequest) (*dto.OutboundOrderView, error) {
	var result *dto.OutboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.outboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitOutbound(ctx, repo, order, constants.OutboundStatusPicking); err != nil {
			return err
		}
		itemMap := make(map[uint]*model.OutboundItem, len(order.Items))
		for i := range order.Items {
			itemMap[order.Items[i].ID] = &order.Items[i]
		}
		for _, p := range req.Items {
			item, ok := itemMap[p.ItemID]
			if !ok {
				return util.NewAppError(constants.CodeInvalidParams, 400, fmt.Sprintf("拣货明细ID %d 不属于该出库单", p.ItemID))
			}
			item.ActualQty = p.ActualQty
			if err := repo.UpdateItem(ctx, item); err != nil {
				return err
			}
			if item.BinLocationID == nil {
				return util.NewAppError(constants.CodeConflict, 409, "出库明细缺少库位，无法扣减库存")
			}
			if err := s.inventorySvc.ReduceStock(ctx, tx, item.ProductID, *item.BinLocationID, item.ActualQty); err != nil {
				if errors.Is(err, repository.ErrInsufficientStock) {
					s.logger.WarnContext(ctx, constants.LogInsufficientStock,
						"product_id", item.ProductID, "bin_id", *item.BinLocationID, "required", item.ActualQty, "order_no", order.OrderNo)
					return util.WrapAppError(
						util.NewAppError(constants.CodeInsufficientStock, 409, fmt.Sprintf("商品ID %d 在指定库位库存不足，请检查库存与拣货数量", item.ProductID)),
						err,
					)
				}
				return err
			}
		}
		claims, _ := util.CurrentUser(ctx)
		uid := claimsUserID(claims)
		order.PickerID = &uid
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "outbound", "outbound.picking", "OutboundOrder", order.OrderNo, "拣货完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "拣货")
	}
	s.logger.InfoContext(ctx, constants.LogPickOutbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// Checking 复核：Picking → Checking。
func (s *outboundService) Checking(ctx context.Context, id uint) (*dto.OutboundOrderView, error) {
	var result *dto.OutboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.outboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitOutbound(ctx, repo, order, constants.OutboundStatusChecking); err != nil {
			return err
		}
		claims, _ := util.CurrentUser(ctx)
		uid := claimsUserID(claims)
		order.CheckerID = &uid
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "outbound", "outbound.checking", "OutboundOrder", order.OrderNo, "复核完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "复核")
	}
	s.logger.InfoContext(ctx, constants.LogCheckOutbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// Packing 打包：Checking → Packing。
func (s *outboundService) Packing(ctx context.Context, id uint) (*dto.OutboundOrderView, error) {
	var result *dto.OutboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.outboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitOutbound(ctx, repo, order, constants.OutboundStatusPacking); err != nil {
			return err
		}
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "outbound", "outbound.packing", "OutboundOrder", order.OrderNo, "打包完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "打包")
	}
	s.logger.InfoContext(ctx, constants.LogPackOutbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// Ship 发货：Packing → Shipped。
func (s *outboundService) Ship(ctx context.Context, id uint, req *dto.OutboundShipRequest) (*dto.OutboundOrderView, error) {
	var result *dto.OutboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.outboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitOutbound(ctx, repo, order, constants.OutboundStatusShipped); err != nil {
			return err
		}
		now := time.Now()
		order.ActualShipDate = &now
		order.TrackingNo = req.TrackingNo
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "outbound", "outbound.ship", "OutboundOrder", order.OrderNo, "发货完成，快递单号："+req.TrackingNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "发货")
	}
	s.logger.InfoContext(ctx, constants.LogShipOutbound, "order_no", result.OrderNo, "tracking_no", req.TrackingNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// Complete 完成出库：Shipped → Completed。
func (s *outboundService) Complete(ctx context.Context, id uint) (*dto.OutboundOrderView, error) {
	var result *dto.OutboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.outboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitOutbound(ctx, repo, order, constants.OutboundStatusCompleted); err != nil {
			return err
		}
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "outbound", "outbound.complete", "OutboundOrder", order.OrderNo, "出库单完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "完成")
	}
	s.logger.InfoContext(ctx, constants.LogCompleteOutbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// transitOutbound 校验并执行出库单状态流转。
func (s *outboundService) transitOutbound(ctx context.Context, repo repository.OutboundRepository, order *model.OutboundOrder, to string) error {
	if !constants.CanTransitOutbound(order.Status, to) {
		return util.NewAppError(constants.CodeBadStateTransition, 409,
			fmt.Sprintf("出库单 %s 状态 %s 不能流转为 %s", order.OrderNo, order.Status, to))
	}
	order.Status = to
	return nil
}

// wrapTransitionError 将业务错误包装为出库操作的统一错误。
func (s *outboundService) wrapTransitionError(err error, action string) error {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		return err
	}
	if errors.Is(err, repository.ErrNotFound) {
		return util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "出库单不存在，无法"+action), err)
	}
	return fmt.Errorf("outbound %s: %w", action, err)
}

// checkOwnerScope Owner 角色只能操作自己的出库单。
func (s *outboundService) checkOwnerScope(ctx context.Context, ownerID uint) error {
	claims, ok := util.CurrentUser(ctx)
	if !ok {
		return util.NewAppError(constants.CodeUnauthorized, 401, constants.MsgUnauthorized)
	}
	if claims.Role == constants.RoleOwner {
		if claims.OwnerID == nil || *claims.OwnerID != ownerID {
			return util.NewAppError(constants.CodeRoleForbidden, 403, "货主角色无权操作其他货主的出库单")
		}
	}
	return nil
}

// toView 组装出库单视图（含货主/拣货员/复核员名称与明细商品名称、库位编码）。
func (s *outboundService) toView(ctx context.Context, order *model.OutboundOrder) *dto.OutboundOrderView {
	ownerName, _ := s.ownerName(ctx, order.OwnerID)
	pickerName, checkerName := "", ""
	if order.PickerID != nil {
		if u, err := s.userRepo.FindByID(ctx, *order.PickerID); err == nil {
			pickerName = u.Name
		}
	}
	if order.CheckerID != nil {
		if u, err := s.userRepo.FindByID(ctx, *order.CheckerID); err == nil {
			checkerName = u.Name
		}
	}
	view := dto.ToOutboundOrderView(order, ownerName, pickerName, checkerName)
	view.Items = make([]dto.OutboundItemView, 0, len(order.Items))
	for i := range order.Items {
		item := &order.Items[i]
		productName, sku := "", ""
		if p, err := s.productRepo.FindByID(ctx, item.ProductID); err == nil {
			productName, sku = p.Name, p.SKU
		}
		binCode := ""
		if item.BinLocationID != nil {
			if b, err := s.binRepo.FindByID(ctx, *item.BinLocationID); err == nil {
				binCode = dto.ToBinView(b).Code
			}
		}
		view.Items = append(view.Items, dto.ToOutboundItemView(item, productName, sku, binCode))
	}
	return view
}

// ownerName 查询货主名称。
func (s *outboundService) ownerName(ctx context.Context, ownerID uint) (string, error) {
	owner, err := s.ownerRepo.FindByID(ctx, ownerID)
	if err != nil {
		return "", err
	}
	return owner.Name, nil
}
