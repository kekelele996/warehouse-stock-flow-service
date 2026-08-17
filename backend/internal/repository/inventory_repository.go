package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
)

// InventoryFilter 库存查询过滤条件。
type InventoryFilter struct {
	OwnerID       uint
	ProductID     uint
	BinLocationID uint
	Keyword       string
	Page          int
	PageSize      int
}

// InventoryRepository 库存仓储接口。
type InventoryRepository interface {
	WithTx(tx *gorm.DB) InventoryRepository
	Create(ctx context.Context, inv *model.Inventory) error
	Update(ctx context.Context, inv *model.Inventory) error
	FindByID(ctx context.Context, id uint) (*model.Inventory, error)
	FindByProductAndBin(ctx context.Context, productID, binLocationID uint, batchNo string) (*model.Inventory, error)
	List(ctx context.Context, filter InventoryFilter) ([]model.Inventory, int64, error)
	ListForUpdateByProduct(ctx context.Context, productID uint) ([]model.Inventory, error)
	DecrementQuantity(ctx context.Context, id uint, qty int) error
	SumValueByOwner(ctx context.Context) ([]InventorySummaryRow, error)
	SumValueByOwnerID(ctx context.Context, ownerID uint) (float64, error)
}

// InventorySummaryRow 货主库存金额汇总行。
type InventorySummaryRow struct {
	OwnerID    uint
	OwnerName  string
	TotalQty   int
	TotalValue float64
}

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 构造库存仓储。
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) WithTx(tx *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: tx}
}

func (r *inventoryRepository) Create(ctx context.Context, inv *model.Inventory) error {
	if err := r.db.WithContext(ctx).Create(inv).Error; err != nil {
		return fmt.Errorf("create inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepository) Update(ctx context.Context, inv *model.Inventory) error {
	if err := r.db.WithContext(ctx).Save(inv).Error; err != nil {
		return fmt.Errorf("update inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepository) FindByID(ctx context.Context, id uint) (*model.Inventory, error) {
	var inv model.Inventory
	if err := r.db.WithContext(ctx).First(&inv, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find inventory by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find inventory by id %d: %w", id, err)
	}
	return &inv, nil
}

func (r *inventoryRepository) FindByProductAndBin(ctx context.Context, productID, binLocationID uint, batchNo string) (*model.Inventory, error) {
	var inv model.Inventory
	query := r.db.WithContext(ctx).Where("product_id = ? AND bin_location_id = ?", productID, binLocationID)
	if batchNo != "" {
		query = query.Where("batch_no = ?", batchNo)
	}
	if err := query.First(&inv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find inventory by product %d bin %d: %w", productID, binLocationID, ErrNotFound)
		}
		return nil, fmt.Errorf("find inventory by product %d bin %d: %w", productID, binLocationID, err)
	}
	return &inv, nil
}

func (r *inventoryRepository) List(ctx context.Context, filter InventoryFilter) ([]model.Inventory, int64, error) {
	var list []model.Inventory
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Inventory{})
	if filter.OwnerID > 0 {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}
	if filter.ProductID > 0 {
		query = query.Where("product_id = ?", filter.ProductID)
	}
	if filter.BinLocationID > 0 {
		query = query.Where("bin_location_id = ?", filter.BinLocationID)
	}
	if filter.Keyword != "" {
		query = query.Where("batch_no ILIKE ?", "%"+filter.Keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count inventory: %w", err)
	}
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	if err := query.Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list inventory: %w", err)
	}
	return list, total, nil
}

// ListForUpdateByProduct 锁定某商品全部库存行（SELECT ... FOR UPDATE）。
func (r *inventoryRepository) ListForUpdateByProduct(ctx context.Context, productID uint) ([]model.Inventory, error) {
	var list []model.Inventory
	if err := r.db.WithContext(ctx).Clauses(LockClause()).Where("product_id = ?", productID).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list inventory for update product %d: %w", productID, err)
	}
	return list, nil
}

func (r *inventoryRepository) DecrementQuantity(ctx context.Context, id uint, qty int) error {
	res := r.db.WithContext(ctx).Model(&model.Inventory{}).
		Where("id = ? AND quantity >= ?", id, qty).
		UpdateColumn("quantity", gorm.Expr("quantity - ?", qty))
	if res.Error != nil {
		return fmt.Errorf("decrement inventory quantity id %d: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("decrement inventory quantity id %d: %w", id, ErrInsufficientStock)
	}
	return nil
}

func (r *inventoryRepository) SumValueByOwner(ctx context.Context) ([]InventorySummaryRow, error) {
	var rows []InventorySummaryRow
	err := r.db.WithContext(ctx).Model(&model.Inventory{}).
		Select("inventories.owner_id, owners.name AS owner_name, SUM(inventories.quantity * products.price) AS total_value, SUM(inventories.quantity) AS total_qty").
		Joins("JOIN owners ON owners.id = inventories.owner_id").
		Joins("JOIN products ON products.id = inventories.product_id").
		Group("inventories.owner_id, owners.name").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sum inventory value by owner: %w", err)
	}
	return rows, nil
}

func (r *inventoryRepository) SumValueByOwnerID(ctx context.Context, ownerID uint) (float64, error) {
	var value float64
	err := r.db.WithContext(ctx).Model(&model.Inventory{}).
		Select("COALESCE(SUM(inventories.quantity * products.price), 0)").
		Joins("JOIN products ON products.id = inventories.product_id").
		Where("inventories.owner_id = ?", ownerID).
		Scan(&value).Error
	if err != nil {
		return 0, fmt.Errorf("sum inventory value by owner id %d: %w", ownerID, err)
	}
	return value, nil
}

// ErrInsufficientStock 库存不足哨兵错误。
var ErrInsufficientStock = errors.New("insufficient stock")
