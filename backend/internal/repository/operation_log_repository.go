package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
)

// OperationLogFilter 操作日志查询过滤条件。
type OperationLogFilter struct {
	Module   string
	Action   string
	UserID   uint
	Page     int
	PageSize int
}

// OperationLogRepository 操作日志仓储接口。
type OperationLogRepository interface {
	WithTx(tx *gorm.DB) OperationLogRepository
	Create(ctx context.Context, log *model.OperationLog) error
	List(ctx context.Context, filter OperationLogFilter) ([]model.OperationLog, int64, error)
}

type operationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository 构造操作日志仓储。
func NewOperationLogRepository(db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) WithTx(tx *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db: tx}
}

func (r *operationLogRepository) Create(ctx context.Context, log *model.OperationLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create operation log: %w", err)
	}
	return nil
}

func (r *operationLogRepository) List(ctx context.Context, filter OperationLogFilter) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64
	query := r.db.WithContext(ctx).Model(&model.OperationLog{})
	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count operation logs: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("list operation logs: %w", err)
	}
	return logs, total, nil
}
