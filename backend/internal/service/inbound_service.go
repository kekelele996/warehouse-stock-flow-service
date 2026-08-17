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

// InboundService 入库单服务。
type InboundService interface {
	Create(ctx context.Context, req *dto.InboundCreateRequest) (*dto.InboundOrderView, error)
	List(ctx context.Context, query dto.InboundQuery, page, pageSize int) ([]dto.InboundOrderView, int64, error)
	Get(ctx context.Context, id uint) (*dto.InboundOrderView, error)
	Receive(ctx context.Context, id uint) (*dto.InboundOrderView, error)
	QC(ctx context.Context, id uint, req *dto.InboundQCRequest) (*dto.InboundOrderView, error)
	Shelve(ctx context.Context, id uint, req *dto.InboundShelveRequest) (*dto.InboundOrderView, error)
	Complete(ctx context.Context, id uint) (*dto.InboundOrderView, error)
}

type inboundService struct {
	db          txProvider
	inboundRepo repository.InboundRepository
	ownerRepo   repository.OwnerRepository
	productRepo repository.ProductRepository
	binRepo     repository.BinLocationRepository
	userRepo    repository.UserRepository
	inventorySvc InventoryService
	auditSvc    AuditService
	logger      *slog.Logger
}

// NewInboundService 构造入库单服务。
func NewInboundService(
	db txProvider,
	inboundRepo repository.InboundRepository,
	ownerRepo repository.OwnerRepository,
	productRepo repository.ProductRepository,
	binRepo repository.BinLocationRepository,
	userRepo repository.UserRepository,
	inventorySvc InventoryService,
	auditSvc AuditService,
	logger *slog.Logger,
) InboundService {
	return &inboundService{db: db, inboundRepo: inboundRepo, ownerRepo: ownerRepo, productRepo: productRepo, binRepo: binRepo, userRepo: userRepo, inventorySvc: inventorySvc, auditSvc: auditSvc, logger: logger}
}

func (s *inboundService) Create(ctx context.Context, req *dto.InboundCreateRequest) (*dto.InboundOrderView, error) {
	owner, err := s.ownerRepo.FindByID(ctx, req.OwnerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "所属货主不存在，请检查货主ID字段"), err)
		}
		return nil, fmt.Errorf("inbound create owner: %w", err)
	}
	if owner.Status == constants.OwnerStatusSuspended {
		return nil, util.NewAppError(constants.CodeOwnerSuspended, 409, "货主已暂停合作，无法创建入库单")
	}
	if err := s.checkOwnerScope(ctx, req.OwnerID); err != nil {
		return nil, err
	}
	claims, _ := util.CurrentUser(ctx)
	order := &model.InboundOrder{
		OrderNo:             util.GenerateOrderNo("INB"),
		OwnerID:             req.OwnerID,
		SupplierName:        req.SupplierName,
		ExpectedArrivalDate: req.ExpectedArrivalDate,
		Status:              constants.InboundStatusPending,
		Remark:              req.Remark,
		CreatedBy:           claimsUserID(claims),
		Items:               make([]model.InboundItem, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "商品不存在，请检查商品ID字段"), err)
			}
			return nil, fmt.Errorf("inbound create product: %w", err)
		}
		if product.OwnerID != req.OwnerID {
			return nil, util.NewAppError(constants.CodeForbidden, 403, "商品不属于所选货主，请检查商品与货主匹配关系")
		}
		order.Items = append(order.Items, model.InboundItem{
			ProductID:   item.ProductID,
			BatchNo:     item.BatchNo,
			ExpectedQty: item.ExpectedQty,
		})
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		return s.inboundRepo.WithTx(tx).Create(ctx, order)
	})
	if err != nil {
		return nil, fmt.Errorf("inbound create: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogCreateInbound, "order_no", order.OrderNo, "owner_id", order.OwnerID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "inbound", "inbound.create", "InboundOrder", order.OrderNo, "创建入库单："+order.OrderNo)
	return s.toView(ctx, order), nil
}

func (s *inboundService) List(ctx context.Context, query dto.InboundQuery, page, pageSize int) ([]dto.InboundOrderView, int64, error) {
	if claims, ok := util.CurrentUser(ctx); ok && claims.Role == constants.RoleOwner && claims.OwnerID != nil {
		query.OwnerID = *claims.OwnerID
	}
	filter := repository.InboundFilter{Status: query.Status, OwnerID: query.OwnerID, Keyword: query.Keyword, Page: page, PageSize: pageSize}
	orders, total, err := s.inboundRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("inbound list: %w", err)
	}
	views := make([]dto.InboundOrderView, 0, len(orders))
	for i := range orders {
		views = append(views, *s.toView(ctx, &orders[i]))
	}
	return views, total, nil
}

func (s *inboundService) Get(ctx context.Context, id uint) (*dto.InboundOrderView, error) {
	order, err := s.inboundRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "入库单不存在"), err)
		}
		return nil, fmt.Errorf("inbound get: %w", err)
	}
	if err := s.checkOwnerScope(ctx, order.OwnerID); err != nil {
		return nil, err
	}
	return s.toView(ctx, order), nil
}

// Receive 收货：Pending → Received。
func (s *inboundService) Receive(ctx context.Context, id uint) (*dto.InboundOrderView, error) {
	var result *dto.InboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.inboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitInbound(ctx, repo, order, constants.InboundStatusReceived); err != nil {
			return err
		}
		now := time.Now()
		order.ActualArrivalDate = &now
		claims, _ := util.CurrentUser(ctx)
		uid := claimsUserID(claims)
		order.KeeperID = &uid
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "inbound", "inbound.receive", "InboundOrder", order.OrderNo, "收货完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "收货")
	}
	s.logger.InfoContext(ctx, constants.LogReceiveInbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// QC 收货质检：Received → QCInProgress。
func (s *inboundService) QC(ctx context.Context, id uint, req *dto.InboundQCRequest) (*dto.InboundOrderView, error) {
	var result *dto.InboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.inboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitInbound(ctx, repo, order, constants.InboundStatusQCInProgress); err != nil {
			return err
		}
		itemMap := make(map[uint]*model.InboundItem, len(order.Items))
		for i := range order.Items {
			itemMap[order.Items[i].ID] = &order.Items[i]
		}
		for _, q := range req.Items {
			item, ok := itemMap[q.ItemID]
			if !ok {
				return util.NewAppError(constants.CodeInvalidParams, 400, fmt.Sprintf("质检明细ID %d 不属于该入库单", q.ItemID))
			}
			item.ActualQty = q.ActualQty
			item.QCResult = q.QCResult
			if err := repo.UpdateItem(ctx, item); err != nil {
				return err
			}
		}
		claims, _ := util.CurrentUser(ctx)
		uid := claimsUserID(claims)
		order.QCInspectorID = &uid
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "inbound", "inbound.qc", "InboundOrder", order.OrderNo, "质检完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "质检")
	}
	s.logger.InfoContext(ctx, constants.LogQCInbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// Shelve 上架：QCInProgress → Shelved，写入库存并指派库位。
func (s *inboundService) Shelve(ctx context.Context, id uint, req *dto.InboundShelveRequest) (*dto.InboundOrderView, error) {
	var result *dto.InboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.inboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitInbound(ctx, repo, order, constants.InboundStatusShelved); err != nil {
			return err
		}
		itemMap := make(map[uint]*model.InboundItem, len(order.Items))
		for i := range order.Items {
			itemMap[order.Items[i].ID] = &order.Items[i]
		}
		for _, reqItem := range req.Items {
			item, ok := itemMap[reqItem.ItemID]
			if !ok {
				return util.NewAppError(constants.CodeInvalidParams, 400, fmt.Sprintf("上架明细ID %d 不属于该入库单", reqItem.ItemID))
			}
			if item.QCResult == constants.QCResultFail {
				return util.NewAppError(constants.CodeConflict, 409, fmt.Sprintf("质检不合格的商品不能上架（明细ID %d）", reqItem.ItemID))
			}
			bin, err := s.binRepo.WithTx(tx).FindByID(ctx, reqItem.BinLocationID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return util.NewAppError(constants.CodeNotFound, 404, fmt.Sprintf("库位ID %d 不存在", reqItem.BinLocationID))
				}
				return err
			}
			product, err := s.productRepo.WithTx(tx).FindByID(ctx, item.ProductID)
			if err != nil {
				return err
			}
			if bin.StorageRequirement != product.StorageRequirement {
				return util.NewAppError(constants.CodeConflict, 409, fmt.Sprintf("库位存储要求 %s 与商品 %s 的存储要求 %s 不匹配", bin.StorageRequirement, product.Name, product.StorageRequirement))
			}
			item.BinLocationID = &reqItem.BinLocationID
			if err := repo.UpdateItem(ctx, item); err != nil {
				return err
			}
			qty := item.ActualQty
			if err := s.inventorySvc.AddStock(ctx, tx, item.ProductID, order.OwnerID, reqItem.BinLocationID, item.BatchNo, qty); err != nil {
				return err
			}
		}
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "inbound", "inbound.shelve", "InboundOrder", order.OrderNo, "上架完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "上架")
	}
	s.logger.InfoContext(ctx, constants.LogShelveInbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// Complete 完成入库：Shelved → Completed。
func (s *inboundService) Complete(ctx context.Context, id uint) (*dto.InboundOrderView, error) {
	var result *dto.InboundOrderView
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.inboundRepo.WithTx(tx)
		order, err := repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.transitInbound(ctx, repo, order, constants.InboundStatusCompleted); err != nil {
			return err
		}
		if err := repo.Update(ctx, order); err != nil {
			return err
		}
		_ = s.auditSvc.Record(ctx, "inbound", "inbound.complete", "InboundOrder", order.OrderNo, "入库单完成："+order.OrderNo)
		result = s.toView(ctx, order)
		return nil
	})
	if err != nil {
		return nil, s.wrapTransitionError(err, "完成")
	}
	s.logger.InfoContext(ctx, constants.LogCompleteInbound, "order_no", result.OrderNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return result, nil
}

// transitInbound 校验并执行入库单状态流转。
func (s *inboundService) transitInbound(ctx context.Context, repo repository.InboundRepository, order *model.InboundOrder, to string) error {
	if !constants.CanTransitInbound(order.Status, to) {
		return util.NewAppError(constants.CodeBadStateTransition, 409,
			fmt.Sprintf("入库单 %s 状态 %s 不能流转为 %s", order.OrderNo, order.Status, to))
	}
	order.Status = to
	return nil
}

// wrapTransitionError 将业务错误包装为入库操作的统一错误。
func (s *inboundService) wrapTransitionError(err error, action string) error {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		return err
	}
	if errors.Is(err, repository.ErrNotFound) {
		return util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "入库单不存在，无法"+action), err)
	}
	return fmt.Errorf("inbound %s: %w", action, err)
}

// checkOwnerScope Owner 角色只能操作自己的入库单。
func (s *inboundService) checkOwnerScope(ctx context.Context, ownerID uint) error {
	claims, ok := util.CurrentUser(ctx)
	if !ok {
		return util.NewAppError(constants.CodeUnauthorized, 401, constants.MsgUnauthorized)
	}
	if claims.Role == constants.RoleOwner {
		if claims.OwnerID == nil || *claims.OwnerID != ownerID {
			return util.NewAppError(constants.CodeRoleForbidden, 403, "货主角色无权操作其他货主的入库单")
		}
	}
	return nil
}

// toView 组装入库单视图（含货主/质检员/库管员名称与明细商品名称、库位编码）。
func (s *inboundService) toView(ctx context.Context, order *model.InboundOrder) *dto.InboundOrderView {
	ownerName, _ := s.ownerName(ctx, order.OwnerID)
	qcName, keeperName := "", ""
	if order.QCInspectorID != nil {
		if u, err := s.userRepo.FindByID(ctx, *order.QCInspectorID); err == nil {
			qcName = u.Name
		}
	}
	if order.KeeperID != nil {
		if u, err := s.userRepo.FindByID(ctx, *order.KeeperID); err == nil {
			keeperName = u.Name
		}
	}
	view := dto.ToInboundOrderView(order, ownerName, qcName, keeperName)
	view.Items = make([]dto.InboundItemView, 0, len(order.Items))
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
		view.Items = append(view.Items, dto.ToInboundItemView(item, productName, sku, binCode))
	}
	return view
}

// ownerName 查询货主名称。
func (s *inboundService) ownerName(ctx context.Context, ownerID uint) (string, error) {
	owner, err := s.ownerRepo.FindByID(ctx, ownerID)
	if err != nil {
		return "", err
	}
	return owner.Name, nil
}

// claimsUserID 获取当前用户ID。
func claimsUserID(claims *util.Claims) uint {
	if claims == nil {
		return 0
	}
	return claims.UserID
}
