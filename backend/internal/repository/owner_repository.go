package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
)

// OwnerFilter 货主查询过滤条件。
type OwnerFilter struct {
	Keyword string
	Status  string
	Page    int
	PageSize int
}

// OwnerRepository 货主仓储接口。
type OwnerRepository interface {
	WithTx(tx *gorm.DB) OwnerRepository
	Create(ctx context.Context, owner *model.Owner) error
	Update(ctx context.Context, owner *model.Owner) error
	FindByID(ctx context.Context, id uint) (*model.Owner, error)
	FindByName(ctx context.Context, name string) (*model.Owner, error)
	List(ctx context.Context, filter OwnerFilter) ([]model.Owner, int64, error)
	Count(ctx context.Context) (int64, error)
}

type ownerRepository struct {
	db *gorm.DB
}

// NewOwnerRepository 构造货主仓储。
func NewOwnerRepository(db *gorm.DB) OwnerRepository {
	return &ownerRepository{db: db}
}

func (r *ownerRepository) WithTx(tx *gorm.DB) OwnerRepository {
	return &ownerRepository{db: tx}
}

func (r *ownerRepository) Create(ctx context.Context, owner *model.Owner) error {
	if err := r.db.WithContext(ctx).Create(owner).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create owner: %w", ErrDuplicate)
		}
		return fmt.Errorf("create owner: %w", err)
	}
	return nil
}

func (r *ownerRepository) Update(ctx context.Context, owner *model.Owner) error {
	if err := r.db.WithContext(ctx).Save(owner).Error; err != nil {
		return fmt.Errorf("update owner: %w", err)
	}
	return nil
}

func (r *ownerRepository) FindByID(ctx context.Context, id uint) (*model.Owner, error) {
	var owner model.Owner
	if err := r.db.WithContext(ctx).First(&owner, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find owner by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find owner by id %d: %w", id, err)
	}
	return &owner, nil
}

func (r *ownerRepository) FindByName(ctx context.Context, name string) (*model.Owner, error) {
	var owner model.Owner
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&owner).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find owner by name %s: %w", name, ErrNotFound)
		}
		return nil, fmt.Errorf("find owner by name %s: %w", name, err)
	}
	return &owner, nil
}

func (r *ownerRepository) List(ctx context.Context, filter OwnerFilter) ([]model.Owner, int64, error) {
	var owners []model.Owner
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Owner{})
	if filter.Keyword != "" {
		kw := "%" + filter.Keyword + "%"
		query = query.Where("name ILIKE ? OR contact_name ILIKE ? OR phone ILIKE ?", kw, kw, kw)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count owners: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&owners).Error; err != nil {
		return nil, 0, fmt.Errorf("list owners: %w", err)
	}
	return owners, total, nil
}

func (r *ownerRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Owner{}).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count owners: %w", err)
	}
	return total, nil
}
