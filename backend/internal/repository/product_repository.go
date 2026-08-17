package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
)

// ProductFilter 商品查询过滤条件。
type ProductFilter struct {
	OwnerID            uint
	Keyword            string
	Category           string
	StorageRequirement string
	Page               int
	PageSize           int
}

// ProductRepository 商品仓储接口。
type ProductRepository interface {
	WithTx(tx *gorm.DB) ProductRepository
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id uint) (*model.Product, error)
	FindBySKU(ctx context.Context, sku string) (*model.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]model.Product, int64, error)
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 构造商品仓储。
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create product: %w", ErrDuplicate)
		}
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find product by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find product by id %d: %w", id, err)
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find product by sku %s: %w", sku, ErrNotFound)
		}
		return nil, fmt.Errorf("find product by sku %s: %w", sku, err)
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, filter ProductFilter) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Product{})
	if filter.OwnerID > 0 {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}
	if filter.Keyword != "" {
		kw := "%" + filter.Keyword + "%"
		query = query.Where("name ILIKE ? OR sku ILIKE ? OR barcode ILIKE ?", kw, kw, kw)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.StorageRequirement != "" {
		query = query.Where("storage_requirement = ?", filter.StorageRequirement)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return products, total, nil
}

func (r *productRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("owner_id = ?", ownerID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count products by owner %d: %w", ownerID, err)
	}
	return total, nil
}
