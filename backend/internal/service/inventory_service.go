package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
)

// InventoryService 库存服务：出入库共享库存增减与库位占用率重算。
type InventoryService interface {
	List(ctx context.Context, query dto.InventoryQuery, page, pageSize int) ([]dto.InventoryView, int64, error)
	Summary(ctx context.Context) ([]dto.InventorySummaryItem, error)
	AddStock(ctx context.Context, tx *gorm.DB, productID, ownerID, binLocationID uint, batchNo string, quantity int) error
	ReduceStock(ctx context.Context, tx *gorm.DB, productID, binLocationID uint, quantity int) error
}

type inventoryService struct {
	invRepo repository.InventoryRepository
	binRepo repository.BinLocationRepository
	prodRepo repository.ProductRepository
	ownerRepo repository.OwnerRepository
	logger  *slog.Logger
}

// NewInventoryService 构造库存服务。
func NewInventoryService(
	invRepo repository.InventoryRepository,
	binRepo repository.BinLocationRepository,
	prodRepo repository.ProductRepository,
	ownerRepo repository.OwnerRepository,
	logger *slog.Logger,
) InventoryService {
	return &inventoryService{invRepo: invRepo, binRepo: binRepo, prodRepo: prodRepo, ownerRepo: ownerRepo, logger: logger}
}

// AddStock 增加库存（上架/入库时调用），复用方：入库上架、库存初始化。
func (s *inventoryService) AddStock(ctx context.Context, tx *gorm.DB, productID, ownerID, binLocationID uint, batchNo string, quantity int) error {
	invRepo := s.invRepo.WithTx(tx)
	existing, err := invRepo.FindByProductAndBin(ctx, productID, binLocationID, batchNo)
	if err == nil {
		existing.Quantity += quantity
		if err := invRepo.Update(ctx, existing); err != nil {
			return fmt.Errorf("inventory add stock update: %w", err)
		}
	} else if errors.Is(err, repository.ErrNotFound) {
		inv := &model.Inventory{
			ProductID:     productID,
			OwnerID:       ownerID,
			BinLocationID: binLocationID,
			BatchNo:       batchNo,
			Quantity:      quantity,
		}
		if err := invRepo.Create(ctx, inv); err != nil {
			return fmt.Errorf("inventory add stock create: %w", err)
		}
	} else {
		return fmt.Errorf("inventory add stock find: %w", err)
	}
	if err := s.recalcBinOccupancy(ctx, tx, binLocationID); err != nil {
		return fmt.Errorf("inventory add stock recalc bin: %w", err)
	}
	return nil
}

// ReduceStock 扣减库存（拣货时调用），复用方：出库拣货、库存盘点调整。
func (s *inventoryService) ReduceStock(ctx context.Context, tx *gorm.DB, productID, binLocationID uint, quantity int) error {
	invRepo := s.invRepo.WithTx(tx)
	rows, err := invRepo.ListForUpdateByProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("inventory reduce stock lock: %w", err)
	}
	var total int
	for i := range rows {
		if rows[i].BinLocationID == binLocationID {
			total += rows[i].Quantity
		}
	}
	if total < quantity {
		return fmt.Errorf("inventory reduce stock product=%d bin=%d: %w", productID, binLocationID, repository.ErrInsufficientStock)
	}
	remaining := quantity
	for i := range rows {
		if remaining <= 0 {
			break
		}
		row := &rows[i]
		if row.BinLocationID != binLocationID {
			continue
		}
		deduct := row.Quantity
		if deduct > remaining {
			deduct = remaining
		}
		if err := invRepo.DecrementQuantity(ctx, row.ID, deduct); err != nil {
			return fmt.Errorf("inventory reduce stock decrement: %w", err)
		}
		remaining -= deduct
	}
	if err := s.recalcBinOccupancy(ctx, tx, binLocationID); err != nil {
		return fmt.Errorf("inventory reduce stock recalc bin: %w", err)
	}
	return nil
}

// recalcBinOccupancy 重算库位占用率与状态（按商品体积×数量/容量）。
func (s *inventoryService) recalcBinOccupancy(ctx context.Context, tx *gorm.DB, binLocationID uint) error {
	binRepo := s.binRepo.WithTx(tx)
	bin, err := binRepo.FindByID(ctx, binLocationID)
	if err != nil {
		return fmt.Errorf("inventory recalc bin find: %w", err)
	}
	invList, _, err := s.invRepo.WithTx(tx).List(ctx, repository.InventoryFilter{BinLocationID: binLocationID, Page: 1, PageSize: 100})
	if err != nil {
		return fmt.Errorf("inventory recalc bin list: %w", err)
	}
	var used float64
	var totalQty int
	for i := range invList {
		totalQty += invList[i].Quantity
		product, err := s.prodRepo.WithTx(tx).FindByID(ctx, invList[i].ProductID)
		if err != nil {
			return fmt.Errorf("inventory recalc bin product: %w", err)
		}
		used += product.Volume * float64(invList[i].Quantity)
	}
	rate := 0.0
	if bin.Capacity > 0 {
		rate = used / bin.Capacity * 100
		if rate > 100 {
			rate = 100
		}
	}
	bin.OccupancyRate = rate
	if totalQty > 0 {
		bin.Status = constants.BinStatusOccupied
	} else if bin.Status == constants.BinStatusOccupied {
		bin.Status = constants.BinStatusAvailable
	}
	if err := binRepo.Update(ctx, bin); err != nil {
		return fmt.Errorf("inventory recalc bin update: %w", err)
	}
	return nil
}

func (s *inventoryService) List(ctx context.Context, query dto.InventoryQuery, page, pageSize int) ([]dto.InventoryView, int64, error) {
	filter := repository.InventoryFilter{
		OwnerID:       query.OwnerID,
		ProductID:     query.ProductID,
		BinLocationID: query.BinLocationID,
		Keyword:       query.Keyword,
		Page:          page,
		PageSize:      pageSize,
	}
	list, total, err := s.invRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("inventory list: %w", err)
	}
	views := make([]dto.InventoryView, 0, len(list))
	for i := range list {
		product, err := s.prodRepo.FindByID(ctx, list[i].ProductID)
		if err != nil {
			return nil, 0, fmt.Errorf("inventory list product: %w", err)
		}
		bin, err := s.binRepo.FindByID(ctx, list[i].BinLocationID)
		if err != nil {
			return nil, 0, fmt.Errorf("inventory list bin: %w", err)
		}
		ownerName := ""
		if owner, err := s.ownerRepo.FindByID(ctx, list[i].OwnerID); err == nil {
			ownerName = owner.Name
		}
		views = append(views, *dto.ToInventoryView(&list[i], product.Name, product.SKU, ownerName, dto.ToBinView(bin).Code, product.Price))
	}
	return views, total, nil
}

func (s *inventoryService) Summary(ctx context.Context) ([]dto.InventorySummaryItem, error) {
	rows, err := s.invRepo.SumValueByOwner(ctx)
	if err != nil {
		return nil, fmt.Errorf("inventory summary: %w", err)
	}
	items := make([]dto.InventorySummaryItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, dto.InventorySummaryItem{
			OwnerID:    row.OwnerID,
			OwnerName:  row.OwnerName,
			TotalQty:   row.TotalQty,
			TotalValue: row.TotalValue,
		})
	}
	return items, nil
}
