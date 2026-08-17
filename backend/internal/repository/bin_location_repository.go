package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
)

// BinFilter 库位查询过滤条件。
type BinFilter struct {
	Area     string
	Status   string
	Page     int
	PageSize int
}

// BinLocationRepository 库位仓储接口。
type BinLocationRepository interface {
	WithTx(tx *gorm.DB) BinLocationRepository
	Create(ctx context.Context, bin *model.BinLocation) error
	Update(ctx context.Context, bin *model.BinLocation) error
	FindByID(ctx context.Context, id uint) (*model.BinLocation, error)
	FindByCode(ctx context.Context, area, rackNo string, layerNo, columnNo int) (*model.BinLocation, error)
	List(ctx context.Context, filter BinFilter) ([]model.BinLocation, int64, error)
	FindAvailable(ctx context.Context, storageRequirement string, limit int) ([]model.BinLocation, error)
	CountOccupied(ctx context.Context) (int64, error)
	CountAll(ctx context.Context) (int64, error)
	SumOccupancyRate(ctx context.Context) (float64, error)
}

type binLocationRepository struct {
	db *gorm.DB
}

// NewBinLocationRepository 构造库位仓储。
func NewBinLocationRepository(db *gorm.DB) BinLocationRepository {
	return &binLocationRepository{db: db}
}

func (r *binLocationRepository) WithTx(tx *gorm.DB) BinLocationRepository {
	return &binLocationRepository{db: tx}
}

func (r *binLocationRepository) Create(ctx context.Context, bin *model.BinLocation) error {
	if err := r.db.WithContext(ctx).Create(bin).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create bin location: %w", ErrDuplicate)
		}
		return fmt.Errorf("create bin location: %w", err)
	}
	return nil
}

func (r *binLocationRepository) Update(ctx context.Context, bin *model.BinLocation) error {
	if err := r.db.WithContext(ctx).Save(bin).Error; err != nil {
		return fmt.Errorf("update bin location: %w", err)
	}
	return nil
}

func (r *binLocationRepository) FindByID(ctx context.Context, id uint) (*model.BinLocation, error) {
	var bin model.BinLocation
	if err := r.db.WithContext(ctx).First(&bin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find bin location by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find bin location by id %d: %w", id, err)
	}
	return &bin, nil
}

func (r *binLocationRepository) FindByCode(ctx context.Context, area, rackNo string, layerNo, columnNo int) (*model.BinLocation, error) {
	var bin model.BinLocation
	err := r.db.WithContext(ctx).Where("area = ? AND rack_no = ? AND layer_no = ? AND column_no = ?", area, rackNo, layerNo, columnNo).First(&bin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find bin location by code: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find bin location by code: %w", err)
	}
	return &bin, nil
}

func (r *binLocationRepository) List(ctx context.Context, filter BinFilter) ([]model.BinLocation, int64, error) {
	var bins []model.BinLocation
	var total int64
	query := r.db.WithContext(ctx).Model(&model.BinLocation{})
	if filter.Area != "" {
		query = query.Where("area = ?", filter.Area)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count bin locations: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("area asc, rack_no asc, layer_no asc, column_no asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&bins).Error; err != nil {
		return nil, 0, fmt.Errorf("list bin locations: %w", err)
	}
	return bins, total, nil
}

// FindAvailable 查询满足存储要求且可用的库位。
func (r *binLocationRepository) FindAvailable(ctx context.Context, storageRequirement string, limit int) ([]model.BinLocation, error) {
	var bins []model.BinLocation
	query := r.db.WithContext(ctx).Where("status = ? AND occupancy_rate < 100", constants.BinStatusAvailable)
	if storageRequirement != "" {
		query = query.Where("storage_requirement = ?", storageRequirement)
	}
	if err := query.Order("occupancy_rate asc").Limit(limit).Find(&bins).Error; err != nil {
		return nil, fmt.Errorf("find available bin locations: %w", err)
	}
	return bins, nil
}

func (r *binLocationRepository) CountOccupied(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.BinLocation{}).Where("status = ?", constants.BinStatusOccupied).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count occupied bin locations: %w", err)
	}
	return total, nil
}

func (r *binLocationRepository) CountAll(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.BinLocation{}).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count bin locations: %w", err)
	}
	return total, nil
}

func (r *binLocationRepository) SumOccupancyRate(ctx context.Context) (float64, error) {
	var sum float64
	if err := r.db.WithContext(ctx).Model(&model.BinLocation{}).Select("COALESCE(SUM(occupancy_rate), 0)").Scan(&sum).Error; err != nil {
		return 0, fmt.Errorf("sum bin occupancy rate: %w", err)
	}
	return sum, nil
}
