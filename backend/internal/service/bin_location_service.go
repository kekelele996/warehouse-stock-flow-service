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
	"github.com/wmsflow/wmsflow/internal/util"
)

// BinLocationService 库位服务。
type BinLocationService interface {
	Create(ctx context.Context, req *dto.BinCreateRequest) (*dto.BinView, error)
	BatchCreate(ctx context.Context, req *dto.BinBatchCreateRequest) (int, error)
	Update(ctx context.Context, id uint, req *dto.BinUpdateRequest) (*dto.BinView, error)
	List(ctx context.Context, query dto.BinQuery, page, pageSize int) ([]dto.BinView, int64, error)
	Get(ctx context.Context, id uint) (*dto.BinView, error)
	Contents(ctx context.Context, id uint) ([]dto.BinContent, error)
	Recommend(ctx context.Context, productID uint, quantity int) (*dto.BinView, error)
}

type binLocationService struct {
	db          txProvider
	binRepo     repository.BinLocationRepository
	productRepo repository.ProductRepository
	invRepo     repository.InventoryRepository
	ownerRepo   repository.OwnerRepository
	auditSvc    AuditService
	logger      *slog.Logger
}

// NewBinLocationService 构造库位服务。
func NewBinLocationService(
	db txProvider,
	binRepo repository.BinLocationRepository,
	productRepo repository.ProductRepository,
	invRepo repository.InventoryRepository,
	ownerRepo repository.OwnerRepository,
	auditSvc AuditService,
	logger *slog.Logger,
) BinLocationService {
	return &binLocationService{db: db, binRepo: binRepo, productRepo: productRepo, invRepo: invRepo, ownerRepo: ownerRepo, auditSvc: auditSvc, logger: logger}
}

func (s *binLocationService) Create(ctx context.Context, req *dto.BinCreateRequest) (*dto.BinView, error) {
	if _, err := s.binRepo.FindByCode(ctx, req.Area, req.RackNo, req.LayerNo, req.ColumnNo); err == nil {
		return nil, util.NewAppError(constants.CodeDuplicate, 409, "库位编码已存在，请检查区域/货架/层/列字段")
	}
	bin := &model.BinLocation{
		Area:               req.Area,
		RackNo:             req.RackNo,
		LayerNo:            req.LayerNo,
		ColumnNo:           req.ColumnNo,
		Capacity:           req.Capacity,
		OccupancyRate:      0,
		StorageRequirement: req.StorageRequirement,
		Status:             constants.BinStatusAvailable,
	}
	if err := s.binRepo.Create(ctx, bin); err != nil {
		return nil, fmt.Errorf("bin create: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogCreateBin, "bin_id", bin.ID, "area", bin.Area, "rack", bin.RackNo, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "bin", "bin.create", "BinLocation", fmt.Sprintf("%d", bin.ID), "创建库位："+dto.ToBinView(bin).Code)
	return dto.ToBinView(bin), nil
}

func (s *binLocationService) BatchCreate(ctx context.Context, req *dto.BinBatchCreateRequest) (int, error) {
	created := 0
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.binRepo.WithTx(tx)
		for layer := req.LayerStart; layer <= req.LayerEnd; layer++ {
			for column := req.ColumnStart; column <= req.ColumnEnd; column++ {
				bin := &model.BinLocation{
					Area:               req.Area,
					RackNo:             req.RackNo,
					LayerNo:            layer,
					ColumnNo:           column,
					Capacity:           req.Capacity,
					OccupancyRate:      0,
					StorageRequirement: req.StorageRequirement,
					Status:             constants.BinStatusAvailable,
				}
				if err := repo.Create(ctx, bin); err != nil {
					if errors.Is(err, repository.ErrDuplicate) {
						continue
					}
					return fmt.Errorf("bin batch create layer=%d column=%d: %w", layer, column, err)
				}
				created++
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("bin batch create: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogBatchCreateBin, "count", created, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "bin", "bin.batch_create", "BinLocation", "", fmt.Sprintf("批量创建库位 %d 个", created))
	return created, nil
}

func (s *binLocationService) Update(ctx context.Context, id uint, req *dto.BinUpdateRequest) (*dto.BinView, error) {
	bin, err := s.binRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "库位不存在"), err)
		}
		return nil, fmt.Errorf("bin update find: %w", err)
	}
	bin.Capacity = req.Capacity
	bin.Status = req.Status
	if err := s.binRepo.Update(ctx, bin); err != nil {
		return nil, fmt.Errorf("bin update: %w", err)
	}
	s.logger.InfoContext(ctx, constants.LogUpdateBin, "bin_id", bin.ID, "operator", currentUsername(ctx), "role", currentRole(ctx))
	_ = s.auditSvc.Record(ctx, "bin", "bin.update", "BinLocation", fmt.Sprintf("%d", bin.ID), "更新库位："+dto.ToBinView(bin).Code)
	return dto.ToBinView(bin), nil
}

func (s *binLocationService) List(ctx context.Context, query dto.BinQuery, page, pageSize int) ([]dto.BinView, int64, error) {
	bins, total, err := s.binRepo.List(ctx, repository.BinFilter{Area: query.Area, Status: query.Status, Page: page, PageSize: pageSize})
	if err != nil {
		return nil, 0, fmt.Errorf("bin list: %w", err)
	}
	views := make([]dto.BinView, 0, len(bins))
	for i := range bins {
		views = append(views, *dto.ToBinView(&bins[i]))
	}
	return views, total, nil
}

func (s *binLocationService) Get(ctx context.Context, id uint) (*dto.BinView, error) {
	bin, err := s.binRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "库位不存在"), err)
		}
		return nil, fmt.Errorf("bin get: %w", err)
	}
	return dto.ToBinView(bin), nil
}

func (s *binLocationService) Contents(ctx context.Context, id uint) ([]dto.BinContent, error) {
	if _, err := s.binRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "库位不存在"), err)
		}
		return nil, fmt.Errorf("bin contents find: %w", err)
	}
	invList, _, err := s.invRepo.List(ctx, repository.InventoryFilter{BinLocationID: id, Page: 1, PageSize: 100})
	if err != nil {
		return nil, fmt.Errorf("bin contents inventory: %w", err)
	}
	contents := make([]dto.BinContent, 0, len(invList))
	for i := range invList {
		product, err := s.productRepo.FindByID(ctx, invList[i].ProductID)
		if err != nil {
			return nil, fmt.Errorf("bin contents product: %w", err)
		}
		ownerName := ""
		if owner, err := s.ownerRepo.FindByID(ctx, invList[i].OwnerID); err == nil {
			ownerName = owner.Name
		}
		contents = append(contents, dto.BinContent{
			ProductID:   product.ID,
			ProductName: product.Name,
			SKU:         product.SKU,
			BatchNo:     invList[i].BatchNo,
			Quantity:    invList[i].Quantity,
			OwnerName:   ownerName,
		})
	}
	return contents, nil
}

// Recommend 按存储要求推荐占用率最低的可用库位。
func (s *binLocationService) Recommend(ctx context.Context, productID uint, quantity int) (*dto.BinView, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.WrapAppError(util.NewAppError(constants.CodeNotFound, 404, "商品不存在，无法推荐库位"), err)
		}
		return nil, fmt.Errorf("bin recommend product: %w", err)
	}
	bins, err := s.binRepo.FindAvailable(ctx, product.StorageRequirement, 10)
	if err != nil {
		return nil, fmt.Errorf("bin recommend available: %w", err)
	}
	if len(bins) == 0 {
		return nil, util.NewAppError(constants.CodeConflict, 409, "无可用库位，请联系仓库经理分配库位")
	}
	// FindAvailable 按 occupancy_rate asc 返回，首元素占用率最低。
	bin := &bins[0]
	s.logger.InfoContext(ctx, constants.LogRecommendBin, "product_id", productID, "bin_id", bin.ID, "quantity", quantity, "operator", currentUsername(ctx), "role", currentRole(ctx))
	return dto.ToBinView(bin), nil
}
