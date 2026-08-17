package service

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
)

type fakeBinRepo struct {
	repository.BinLocationRepository
	findByIDFn   func(ctx context.Context, id uint) (*model.BinLocation, error)
	findByCodeFn func(ctx context.Context, area, rackNo string, layerNo, columnNo int) (*model.BinLocation, error)
	createFn     func(ctx context.Context, bin *model.BinLocation) error
	updateFn     func(ctx context.Context, bin *model.BinLocation) error
	listFn       func(ctx context.Context, filter repository.BinFilter) ([]model.BinLocation, int64, error)
	findAvailFn  func(ctx context.Context, req string, limit int) ([]model.BinLocation, error)
	countAllFn   func(ctx context.Context) (int64, error)
	countOccFn   func(ctx context.Context) (int64, error)
	sumRateFn    func(ctx context.Context) (float64, error)
}

func (f *fakeBinRepo) WithTx(tx *gorm.DB) repository.BinLocationRepository { return f }
func (f *fakeBinRepo) FindByID(ctx context.Context, id uint) (*model.BinLocation, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return f.BinLocationRepository.FindByID(ctx, id)
}
func (f *fakeBinRepo) FindByCode(ctx context.Context, area, rackNo string, layerNo, columnNo int) (*model.BinLocation, error) {
	if f.findByCodeFn != nil {
		return f.findByCodeFn(ctx, area, rackNo, layerNo, columnNo)
	}
	return f.BinLocationRepository.FindByCode(ctx, area, rackNo, layerNo, columnNo)
}
func (f *fakeBinRepo) Create(ctx context.Context, bin *model.BinLocation) error {
	if f.createFn != nil {
		return f.createFn(ctx, bin)
	}
	return f.BinLocationRepository.Create(ctx, bin)
}
func (f *fakeBinRepo) Update(ctx context.Context, bin *model.BinLocation) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, bin)
	}
	return f.BinLocationRepository.Update(ctx, bin)
}
func (f *fakeBinRepo) List(ctx context.Context, filter repository.BinFilter) ([]model.BinLocation, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.BinLocationRepository.List(ctx, filter)
}
func (f *fakeBinRepo) FindAvailable(ctx context.Context, req string, limit int) ([]model.BinLocation, error) {
	if f.findAvailFn != nil {
		return f.findAvailFn(ctx, req, limit)
	}
	return f.BinLocationRepository.FindAvailable(ctx, req, limit)
}
func (f *fakeBinRepo) CountAll(ctx context.Context) (int64, error) {
	if f.countAllFn != nil {
		return f.countAllFn(ctx)
	}
	return f.BinLocationRepository.CountAll(ctx)
}
func (f *fakeBinRepo) CountOccupied(ctx context.Context) (int64, error) {
	if f.countOccFn != nil {
		return f.countOccFn(ctx)
	}
	return f.BinLocationRepository.CountOccupied(ctx)
}
func (f *fakeBinRepo) SumOccupancyRate(ctx context.Context) (float64, error) {
	if f.sumRateFn != nil {
		return f.sumRateFn(ctx)
	}
	return f.BinLocationRepository.SumOccupancyRate(ctx)
}

type fakeInvRepo struct {
	repository.InventoryRepository
	listFn       func(ctx context.Context, filter repository.InventoryFilter) ([]model.Inventory, int64, error)
	sumOwnerFn   func(ctx context.Context) ([]repository.InventorySummaryRow, error)
	sumOwnerIDFn func(ctx context.Context, ownerID uint) (float64, error)
}

func (f *fakeInvRepo) WithTx(tx *gorm.DB) repository.InventoryRepository { return f }
func (f *fakeInvRepo) List(ctx context.Context, filter repository.InventoryFilter) ([]model.Inventory, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.InventoryRepository.List(ctx, filter)
}
func (f *fakeInvRepo) SumValueByOwner(ctx context.Context) ([]repository.InventorySummaryRow, error) {
	if f.sumOwnerFn != nil {
		return f.sumOwnerFn(ctx)
	}
	return f.InventoryRepository.SumValueByOwner(ctx)
}
func (f *fakeInvRepo) SumValueByOwnerID(ctx context.Context, ownerID uint) (float64, error) {
	if f.sumOwnerIDFn != nil {
		return f.sumOwnerIDFn(ctx, ownerID)
	}
	return f.InventoryRepository.SumValueByOwnerID(ctx, ownerID)
}

type fakeInboundRepo struct {
	repository.InboundRepository
	createFn       func(ctx context.Context, order *model.InboundOrder) error
	findByIDFn     func(ctx context.Context, id uint) (*model.InboundOrder, error)
	updateFn       func(ctx context.Context, order *model.InboundOrder) error
	updateItemFn   func(ctx context.Context, item *model.InboundItem) error
	listFn         func(ctx context.Context, filter repository.InboundFilter) ([]model.InboundOrder, int64, error)
	countByOwnerFn func(ctx context.Context, ownerID uint) (int64, error)
	countTodayFn   func(ctx context.Context, ownerID uint) (int64, error)
	countStatusFn  func(ctx context.Context, status string, ownerID uint) (int64, error)
	listIncFn      func(ctx context.Context, ownerID uint, limit int) ([]model.InboundOrder, error)
}

func (f *fakeInboundRepo) WithTx(tx *gorm.DB) repository.InboundRepository { return f }
func (f *fakeInboundRepo) Create(ctx context.Context, order *model.InboundOrder) error {
	if f.createFn != nil {
		return f.createFn(ctx, order)
	}
	return f.InboundRepository.Create(ctx, order)
}
func (f *fakeInboundRepo) FindByID(ctx context.Context, id uint) (*model.InboundOrder, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return f.InboundRepository.FindByID(ctx, id)
}
func (f *fakeInboundRepo) Update(ctx context.Context, order *model.InboundOrder) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, order)
	}
	return f.InboundRepository.Update(ctx, order)
}
func (f *fakeInboundRepo) UpdateItem(ctx context.Context, item *model.InboundItem) error {
	if f.updateItemFn != nil {
		return f.updateItemFn(ctx, item)
	}
	return f.InboundRepository.UpdateItem(ctx, item)
}
func (f *fakeInboundRepo) List(ctx context.Context, filter repository.InboundFilter) ([]model.InboundOrder, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.InboundRepository.List(ctx, filter)
}
func (f *fakeInboundRepo) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	if f.countByOwnerFn != nil {
		return f.countByOwnerFn(ctx, ownerID)
	}
	return f.InboundRepository.CountByOwner(ctx, ownerID)
}
func (f *fakeInboundRepo) CountToday(ctx context.Context, ownerID uint) (int64, error) {
	if f.countTodayFn != nil {
		return f.countTodayFn(ctx, ownerID)
	}
	return f.InboundRepository.CountToday(ctx, ownerID)
}
func (f *fakeInboundRepo) CountByStatus(ctx context.Context, status string, ownerID uint) (int64, error) {
	if f.countStatusFn != nil {
		return f.countStatusFn(ctx, status, ownerID)
	}
	return f.InboundRepository.CountByStatus(ctx, status, ownerID)
}
func (f *fakeInboundRepo) ListIncomplete(ctx context.Context, ownerID uint, limit int) ([]model.InboundOrder, error) {
	if f.listIncFn != nil {
		return f.listIncFn(ctx, ownerID, limit)
	}
	return f.InboundRepository.ListIncomplete(ctx, ownerID, limit)
}

type fakeOutboundRepo struct {
	repository.OutboundRepository
	createFn       func(ctx context.Context, order *model.OutboundOrder) error
	findByIDFn     func(ctx context.Context, id uint) (*model.OutboundOrder, error)
	updateFn       func(ctx context.Context, order *model.OutboundOrder) error
	updateItemFn   func(ctx context.Context, item *model.OutboundItem) error
	listFn         func(ctx context.Context, filter repository.OutboundFilter) ([]model.OutboundOrder, int64, error)
	countByOwnerFn func(ctx context.Context, ownerID uint) (int64, error)
	countTodayFn   func(ctx context.Context, ownerID uint) (int64, error)
	countStatusFn  func(ctx context.Context, status string, ownerID uint) (int64, error)
	listIncFn      func(ctx context.Context, ownerID uint, limit int) ([]model.OutboundOrder, error)
}

func (f *fakeOutboundRepo) WithTx(tx *gorm.DB) repository.OutboundRepository { return f }
func (f *fakeOutboundRepo) Create(ctx context.Context, order *model.OutboundOrder) error {
	if f.createFn != nil {
		return f.createFn(ctx, order)
	}
	return f.OutboundRepository.Create(ctx, order)
}
func (f *fakeOutboundRepo) FindByID(ctx context.Context, id uint) (*model.OutboundOrder, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return f.OutboundRepository.FindByID(ctx, id)
}
func (f *fakeOutboundRepo) Update(ctx context.Context, order *model.OutboundOrder) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, order)
	}
	return f.OutboundRepository.Update(ctx, order)
}
func (f *fakeOutboundRepo) UpdateItem(ctx context.Context, item *model.OutboundItem) error {
	if f.updateItemFn != nil {
		return f.updateItemFn(ctx, item)
	}
	return f.OutboundRepository.UpdateItem(ctx, item)
}
func (f *fakeOutboundRepo) List(ctx context.Context, filter repository.OutboundFilter) ([]model.OutboundOrder, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.OutboundRepository.List(ctx, filter)
}
func (f *fakeOutboundRepo) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	if f.countByOwnerFn != nil {
		return f.countByOwnerFn(ctx, ownerID)
	}
	return f.OutboundRepository.CountByOwner(ctx, ownerID)
}
func (f *fakeOutboundRepo) CountToday(ctx context.Context, ownerID uint) (int64, error) {
	if f.countTodayFn != nil {
		return f.countTodayFn(ctx, ownerID)
	}
	return f.OutboundRepository.CountToday(ctx, ownerID)
}
func (f *fakeOutboundRepo) CountByStatus(ctx context.Context, status string, ownerID uint) (int64, error) {
	if f.countStatusFn != nil {
		return f.countStatusFn(ctx, status, ownerID)
	}
	return f.OutboundRepository.CountByStatus(ctx, status, ownerID)
}
func (f *fakeOutboundRepo) ListIncomplete(ctx context.Context, ownerID uint, limit int) ([]model.OutboundOrder, error) {
	if f.listIncFn != nil {
		return f.listIncFn(ctx, ownerID, limit)
	}
	return f.OutboundRepository.ListIncomplete(ctx, ownerID, limit)
}

type fakeLogRepo struct {
	repository.OperationLogRepository
	createFn func(ctx context.Context, log *model.OperationLog) error
	listFn   func(ctx context.Context, filter repository.OperationLogFilter) ([]model.OperationLog, int64, error)
}

func (f *fakeLogRepo) WithTx(tx *gorm.DB) repository.OperationLogRepository { return f }
func (f *fakeLogRepo) Create(ctx context.Context, log *model.OperationLog) error {
	if f.createFn != nil {
		return f.createFn(ctx, log)
	}
	return f.OperationLogRepository.Create(ctx, log)
}
func (f *fakeLogRepo) List(ctx context.Context, filter repository.OperationLogFilter) ([]model.OperationLog, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.OperationLogRepository.List(ctx, filter)
}

type fakeInventorySvc struct {
	InventoryService
	addStockFn    func(ctx context.Context, tx *gorm.DB, productID, ownerID, binID uint, batchNo string, qty int) error
	reduceStockFn func(ctx context.Context, tx *gorm.DB, productID, binID uint, qty int) error
	summaryFn     func(ctx context.Context) ([]dto.InventorySummaryItem, error)
}

func (f *fakeInventorySvc) AddStock(ctx context.Context, tx *gorm.DB, productID, ownerID, binID uint, batchNo string, qty int) error {
	if f.addStockFn != nil {
		return f.addStockFn(ctx, tx, productID, ownerID, binID, batchNo, qty)
	}
	return nil
}
func (f *fakeInventorySvc) ReduceStock(ctx context.Context, tx *gorm.DB, productID, binID uint, qty int) error {
	if f.reduceStockFn != nil {
		return f.reduceStockFn(ctx, tx, productID, binID, qty)
	}
	return nil
}
func (f *fakeInventorySvc) Summary(ctx context.Context) ([]dto.InventorySummaryItem, error) {
	if f.summaryFn != nil {
		return f.summaryFn(ctx)
	}
	return []dto.InventorySummaryItem{}, nil
}

type fakeAuditSvc struct {
	AuditService
	recordFn func(ctx context.Context, module, action, entityType, entityID, detail string) error
}

func (f *fakeAuditSvc) Record(ctx context.Context, module, action, entityType, entityID, detail string) error {
	if f.recordFn != nil {
		return f.recordFn(ctx, module, action, entityType, entityID, detail)
	}
	return nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(discardWriter{}, nil))
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }
