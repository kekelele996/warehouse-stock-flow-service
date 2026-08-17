package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
)

// OutboundFilter 出库单查询过滤条件。
type OutboundFilter struct {
	Status   string
	OwnerID  uint
	Keyword  string
	Page     int
	PageSize int
}

// OutboundRepository 出库单仓储接口。
type OutboundRepository interface {
	WithTx(tx *gorm.DB) OutboundRepository
	Create(ctx context.Context, order *model.OutboundOrder) error
	Update(ctx context.Context, order *model.OutboundOrder) error
	UpdateItem(ctx context.Context, item *model.OutboundItem) error
	FindByID(ctx context.Context, id uint) (*model.OutboundOrder, error)
	List(ctx context.Context, filter OutboundFilter) ([]model.OutboundOrder, int64, error)
	CountByOwner(ctx context.Context, ownerID uint) (int64, error)
	CountToday(ctx context.Context, ownerID uint) (int64, error)
	CountByStatus(ctx context.Context, status string, ownerID uint) (int64, error)
	ListIncomplete(ctx context.Context, ownerID uint, limit int) ([]model.OutboundOrder, error)
}

type outboundRepository struct {
	db *gorm.DB
}

// NewOutboundRepository 构造出库单仓储。
func NewOutboundRepository(db *gorm.DB) OutboundRepository {
	return &outboundRepository{db: db}
}

func (r *outboundRepository) WithTx(tx *gorm.DB) OutboundRepository {
	return &outboundRepository{db: tx}
}

func (r *outboundRepository) Create(ctx context.Context, order *model.OutboundOrder) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create outbound order: %w", ErrDuplicate)
		}
		return fmt.Errorf("create outbound order: %w", err)
	}
	return nil
}

func (r *outboundRepository) Update(ctx context.Context, order *model.OutboundOrder) error {
	if err := r.db.WithContext(ctx).Omit("Items").Save(order).Error; err != nil {
		return fmt.Errorf("update outbound order: %w", err)
	}
	return nil
}

func (r *outboundRepository) UpdateItem(ctx context.Context, item *model.OutboundItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("update outbound item: %w", err)
	}
	return nil
}

func (r *outboundRepository) FindByID(ctx context.Context, id uint) (*model.OutboundOrder, error) {
	var order model.OutboundOrder
	err := r.db.WithContext(ctx).Preload("Items").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find outbound order by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find outbound order by id %d: %w", id, err)
	}
	return &order, nil
}

func (r *outboundRepository) List(ctx context.Context, filter OutboundFilter) ([]model.OutboundOrder, int64, error) {
	var orders []model.OutboundOrder
	var total int64
	query := r.db.WithContext(ctx).Model(&model.OutboundOrder{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OwnerID > 0 {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}
	if filter.Keyword != "" {
		query = query.Where("order_no ILIKE ? OR receiver_name ILIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count outbound orders: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("list outbound orders: %w", err)
	}
	return orders, total, nil
}

func (r *outboundRepository) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.OutboundOrder{}).Where("owner_id = ?", ownerID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count outbound orders by owner %d: %w", ownerID, err)
	}
	return total, nil
}

func (r *outboundRepository) CountToday(ctx context.Context, ownerID uint) (int64, error) {
	var total int64
	start := time.Now().Truncate(24 * time.Hour)
	query := r.db.WithContext(ctx).Model(&model.OutboundOrder{}).Where("created_at >= ?", start)
	if ownerID > 0 {
		query = query.Where("owner_id = ?", ownerID)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count today outbound orders: %w", err)
	}
	return total, nil
}

func (r *outboundRepository) ListIncomplete(ctx context.Context, ownerID uint, limit int) ([]model.OutboundOrder, error) {
	var orders []model.OutboundOrder
	query := r.db.WithContext(ctx).Where("status <> ?", "Completed")
	if ownerID > 0 {
		query = query.Where("owner_id = ?", ownerID)
	}
	if err := query.Order("id desc").Limit(limit).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("list incomplete outbound orders: %w", err)
	}
	return orders, nil
}

func (r *outboundRepository) CountByStatus(ctx context.Context, status string, ownerID uint) (int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&model.OutboundOrder{}).Where("status = ?", status)
	if ownerID > 0 {
		query = query.Where("owner_id = ?", ownerID)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count outbound orders by status %s: %w", status, err)
	}
	return total, nil
}
