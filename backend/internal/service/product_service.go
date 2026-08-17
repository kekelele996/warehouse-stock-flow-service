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

// ProductService 商品服务。
type ProductService interface {
	Create(ctx context.Context, req *dto.ProductCreateRequest) (*dto.ProductView, error)
	Update(ctx context.Context, id uint, req *dto.ProductUpdateRequest) (*dto.ProductView, error)
	List(ctx context.Context, query dto.ProductQuery, page, pageSize int) ([]dto.ProductView, int64, error)
	ListByOwner(ctx context.Context, ownerID uint, page, pageSize int) ([]dto.ProductView, int64, error)
	Get(ctx context.Context, id uint) (*dto.ProductView, error)
}

type productService struct {
	productRepo repository.ProductRepository
	ownerRepo   repository.OwnerRepository
	auditSvc    AuditService
	logger      *slog.Logger
}

// NewProductService 构造商品服务。
func NewProductService(productRepo repository.ProductRepository, ownerRepo repository.OwnerRepository, auditSvc AuditService, logger *slog.Logger) ProductService {
	return &productService{productRepo: productRepo, ownerRepo: ownerRepo, auditSvc: auditSvc, logger: logger}
}

func (s *productService) Create(ctx context.Context, req *dto.ProductCreateRequest) (*dto.ProductView, error) {
	if err := s.checkOwnerScope(ctx, req.OwnerID); err != nil {
		return nil, err
	}
	if _, err := s.ownerRepo.FindByID(ctx, req.OwnerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "所属货主不存在，请检查货主ID字段"), err)
		}
		return nil, fmt.Errorf("product create owner check: %w", err)
	}
	product := &model.Product{
		OwnerID:            req.OwnerID,
		Name:               req.Name,
		SKU:                req.SKU,
		Barcode:            req.Barcode,
		Category:           req.Category,
		Spec:               req.Spec,
		Unit:               req.Unit,
		ShelfLifeDays:      req.ShelfLifeDays,
		StorageRequirement: req.StorageRequirement,
		Volume:             req.Volume,
		Weight:             req.Weight,
		Price:              req.Price,
	}
	if err := s.productRepo.Create(ctx, product); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeDuplicate, 409, "商品SKU已存在，请检查SKU字段"), err)
		}
		return nil, fmt.Errorf("product create: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogCreateProduct, "product_id", product.ID, "sku", product.SKU, "owner_id", product.OwnerID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "product", "product.create", "Product", fmt.Sprintf("%d", product.ID), "创建商品："+product.Name)
	ownerName, _ := s.ownerName(ctx, product.OwnerID)
	return dto.ToProductView(product, ownerName), nil
}

func (s *productService) Update(ctx context.Context, id uint, req *dto.ProductUpdateRequest) (*dto.ProductView, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "商品不存在"), err)
		}
		return nil, fmt.Errorf("product update find: %w", err)
	}
	if err := s.checkOwnerScope(ctx, product.OwnerID); err != nil {
		return nil, err
	}
	product.Name = req.Name
	product.Barcode = req.Barcode
	product.Category = req.Category
	product.Spec = req.Spec
	product.Unit = req.Unit
	product.ShelfLifeDays = req.ShelfLifeDays
	product.StorageRequirement = req.StorageRequirement
	product.Volume = req.Volume
	product.Weight = req.Weight
	product.Price = req.Price
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("product update: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogUpdateProduct, "product_id", product.ID, "sku", product.SKU, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "product", "product.update", "Product", fmt.Sprintf("%d", product.ID), "更新商品："+product.Name)
	ownerName, _ := s.ownerName(ctx, product.OwnerID)
	return dto.ToProductView(product, ownerName), nil
}

func (s *productService) List(ctx context.Context, query dto.ProductQuery, page, pageSize int) ([]dto.ProductView, int64, error) {
	if claims, ok := util.CurrentUser(ctx); ok && claims.Role == constants.RoleOwner && claims.OwnerID != nil {
		query.OwnerID = *claims.OwnerID
	}
	filter := repository.ProductFilter{
		OwnerID:            query.OwnerID,
		Keyword:            query.Keyword,
		Category:           query.Category,
		StorageRequirement: query.StorageRequirement,
		Page:               page,
		PageSize:           pageSize,
	}
	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("product list: %w", err)
	}
	views := make([]dto.ProductView, 0, len(products))
	for i := range products {
		ownerName, _ := s.ownerName(ctx, products[i].OwnerID)
		views = append(views, *dto.ToProductView(&products[i], ownerName))
	}
	return views, total, nil
}

// ListByOwner 复用 List 的商品列表逻辑（按货主过滤）。
func (s *productService) ListByOwner(ctx context.Context, ownerID uint, page, pageSize int) ([]dto.ProductView, int64, error) {
	return s.List(ctx, dto.ProductQuery{OwnerID: ownerID}, page, pageSize)
}

func (s *productService) Get(ctx context.Context, id uint) (*dto.ProductView, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "商品不存在"), err)
		}
		return nil, fmt.Errorf("product get: %w", err)
	}
	if err := s.checkOwnerScope(ctx, product.OwnerID); err != nil {
		return nil, err
	}
	ownerName, _ := s.ownerName(ctx, product.OwnerID)
	return dto.ToProductView(product, ownerName), nil
}

// checkOwnerScope Owner 角色只能操作自己的商品。
func (s *productService) checkOwnerScope(ctx context.Context, ownerID uint) error {
	claims, ok := util.CurrentUser(ctx)
	if !ok {
		return util.NewAppError(constants.CodeUnauthorized, 401, constants.MsgUnauthorized)
	}
	if claims.Role == constants.RoleOwner {
		if claims.OwnerID == nil || *claims.OwnerID != ownerID {
			return util.NewAppError(constants.CodeRoleForbidden, 403, "货主角色无权操作其他货主的商品")
		}
	}
	return nil
}

// ownerName 查询货主名称，失败时返回空串。
func (s *productService) ownerName(ctx context.Context, ownerID uint) (string, error) {
	owner, err := s.ownerRepo.FindByID(ctx, ownerID)
	if err != nil {
		return "", err
	}
	return owner.Name, nil
}
