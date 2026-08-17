package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
)

// InboundFilter 入库单查询过滤条件。
type InboundFilter struct {
	Status   string
	OwnerID  uint
	Keyword  string
	Page     int
	PageSize int
}

// InboundRepository 入库单仓储接口。
type InboundRepository interface {
	WithTx(tx *gorm.DB) InboundRepository
	Create(ctx context.Context, order *model.InboundOrder) error
	Update(ctx context.Context, order *model.InboundOrder) error
	UpdateItem(ctx context.Context, item *model.InboundItem) error
	FindByID(ctx context.Context, id uint) (*model.InboundOrder, error)
	List(ctx context.Context, filter InboundFilter) ([]model.InboundOrder, int64, error)
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)
	CountToday(ctx context.Context, ownerID uint) (int64, error)
	CountByStatus(ctx context.Context, status string, ownerID uint) (int64, error)
	ListIncomplete(ctx context.Context, ownerID uint, limit int) ([]model.InboundOrder, error)
}

type inboundRepository struct {
	db *gorm.DB
}

// NewInboundRepository 构造入库单仓储。
func NewInboundRepository(db *gorm.DB) InboundRepository {
	return &inboundRepository{db: db}
}

func (r *inboundRepository) WithTx(tx *gorm.DB) InboundRepository {
	return &inboundRepository{db: tx}
}

func (r *inboundRepository) Create(ctx context.Context, order *model.InboundOrder) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create inbound order: %w", ErrDuplicate)
		}
		return fmt.Errorf("create inbound order: %w", err)
	}
	return nil
}

func (r *inboundRepository) Update(ctx context.Context, order *model.InboundOrder) error {
	if err := r.db.WithContext(ctx).Omit("Items").Save(order).Error; err != nil {
		return fmt.Errorf("update inbound order: %w", err)
	}
	return nil
}

func (r *inboundRepository) UpdateItem(ctx context.Context, item *model.InboundItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("update inbound item: %w", err)
	}
	return nil
}

func (r *inboundRepository) FindByID(ctx context.Context, id uint) (*model.InboundOrder, error) {
	var order model.InboundOrder
	err := r.db.WithContext(ctx).Preload("Items").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find inbound order by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find inbound order by id %d: %w", id, err)
	}
	return &order, nil
}

func (r *inboundRepository) List(ctx context.Context, filter InboundFilter) ([]model.InboundOrder, int64, error) {
	var orders []model.InboundOrder
	var total int64
	query := r.db.WithContext(ctx).Model(&model.InboundOrder{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OwnerID > 0 {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}
	if filter.Keyword != "" {
		query = query.Where("order_no ILIKE ? OR supplier_name ILIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count inbound orders: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("list inbound orders: %w", err)
	}
	return orders, total, nil
}

func (r *inboundRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.InboundOrder{}).Where("owner_id = ?", ownerID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count inbound orders by owner %d: %w", ownerID, err)
	}
	return total, nil
}

func (r *inboundRepository) CountToday(ctx context.Context, ownerID uint) (int64, error) {
	var total int64
	start := time.Now().Truncate(24 * time.Hour)
	query := r.db.WithContext(ctx).Model(&model.InboundOrder{}).Where("created_at >= ?", start)
	if ownerID > 0 {
		query = query.Where("owner_id = ?", ownerID)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count today inbound orders: %w", err)
	}
	return total, nil
}

func (r *inboundRepository) ListIncomplete(ctx context.Context, ownerID uint, limit int) ([]model.InboundOrder, error) {
	var orders []model.InboundOrder
	query := r.db.WithContext(ctx).Where("status <> ?", "Completed")
	if ownerID > 0 {
		query = query.Where("owner_id = ?", ownerID)
	}
	if err := query.Order("id desc").Limit(limit).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("list incomplete inbound orders: %w", err)
	}
	return orders, nil
}

func (r *inboundRepository) CountByStatus(ctx context.Context, status string, ownerID uint) (int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&model.InboundOrder{}).Where("status = ?", status)
	if ownerID > 0 {
		query = query.Where("owner_id = ?", ownerID)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count inbound orders by status %s: %w", status, err)
	}
	return total, nil
}
