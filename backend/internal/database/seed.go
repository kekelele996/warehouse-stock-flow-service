package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/util"
)

// Seed 幂等写入演示数据（用户/货主/商品/库位/示例单据）。
func Seed(ctx context.Context, db *gorm.DB, logger *slog.Logger) error {
	var userCount int64
	if err := db.WithContext(ctx).Model(&model.User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("seed count users: %w", err)
	}
	if userCount > 0 {
		return nil
	}
	users := []model.User{
		{Username: "admin", Name: "系统管理员", Role: constants.RoleAdmin, Status: constants.UserStatusActive},
		{Username: "manager", Name: "王经理", Role: constants.RoleWarehouseManager, Status: constants.UserStatusActive},
		{Username: "qc", Name: "李质检", Role: constants.RoleQCInspector, Status: constants.UserStatusActive},
		{Username: "picker", Name: "赵拣货", Role: constants.RolePicker, Status: constants.UserStatusActive},
		{Username: "checker", Name: "钱复核", Role: constants.RoleChecker, Status: constants.UserStatusActive},
	}
	hashed, err := util.HashPassword("wmsflow@123")
	if err != nil {
		return fmt.Errorf("seed hash password: %w", err)
	}
	now := time.Now()
	owner1 := &model.Owner{
		Name: "华建建材有限公司", ContactName: "张三", Phone: "13800000001",
		Email: "zhangsan@huajian.example", Address: "上海市青浦区仓储路1号",
		SettlementMethod: constants.SettlementMonthly, CreditLimit: 500000, CurrentDebt: 86000,
		Status: constants.OwnerStatusActive, CreatedAt: now, UpdatedAt: now,
	}
	owner2 := &model.Owner{
		Name: "启航五金商贸", ContactName: "李四", Phone: "13800000002",
		Email: "lisi@qihang.example", Address: "江苏省苏州市吴中区物流园2号",
		SettlementMethod: constants.SettlementPerOrder, CreditLimit: 200000, CurrentDebt: 12000,
		Status: constants.OwnerStatusActive, CreatedAt: now, UpdatedAt: now,
	}
	owner3 := &model.Owner{
		Name: "冷链食品供应链", ContactName: "王五", Phone: "13800000003",
		Email: "wangwu@coldchain.example", Address: "浙江省杭州市余杭区冷链基地",
		SettlementMethod: constants.SettlementPrepaid, CreditLimit: 800000, CurrentDebt: 0,
		Status: constants.OwnerStatusActive, CreatedAt: now, UpdatedAt: now,
	}
	if err := db.WithContext(ctx).Create([]*model.Owner{owner1, owner2, owner3}).Error; err != nil {
		return fmt.Errorf("seed owners: %w", err)
	}
	users = append(users,
		model.User{Username: "owner01", Name: "张三(华建)", Role: constants.RoleOwner, OwnerID: &owner1.ID, Status: constants.UserStatusActive},
		model.User{Username: "owner02", Name: "李四(启航)", Role: constants.RoleOwner, OwnerID: &owner2.ID, Status: constants.UserStatusActive},
	)
	for i := range users {
		users[i].Password = hashed
		users[i].CreatedAt = now
		users[i].UpdatedAt = now
	}
	if err := db.WithContext(ctx).Create(&users).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	products := []model.Product{
		{OwnerID: owner1.ID, Name: "硅酸盐水泥P.O42.5", SKU: "HJ-CM-001", Barcode: "690100000001", Category: "水泥", Spec: "50kg/袋", Unit: "袋", ShelfLifeDays: intPtr(180), StorageRequirement: constants.StorageNormal, Volume: 0.05, Weight: 50, Price: 26.5},
		{OwnerID: owner1.ID, Name: "热轧螺纹钢HRB400", SKU: "HJ-ST-002", Barcode: "690100000002", Category: "钢材", Spec: "Φ12mm×9m", Unit: "根", StorageRequirement: constants.StorageNormal, Volume: 0.12, Weight: 80, Price: 118.0},
		{OwnerID: owner1.ID, Name: "聚合物防水涂料", SKU: "HJ-WP-003", Barcode: "690100000003", Category: "防水材料", Spec: "20kg/桶", Unit: "桶", ShelfLifeDays: intPtr(365), StorageRequirement: constants.StorageDangerous, Volume: 0.02, Weight: 20, Price: 168.0},
		{OwnerID: owner1.ID, Name: "釉面瓷砖800×800", SKU: "HJ-TL-004", Barcode: "690100000004", Category: "瓷砖", Spec: "800×800mm", Unit: "箱", StorageRequirement: constants.StorageFragile, Volume: 0.06, Weight: 28, Price: 89.9},
		{OwnerID: owner2.ID, Name: "不锈钢角码", SKU: "QH-HW-001", Barcode: "690200000001", Category: "五金件", Spec: "50×50×3mm", Unit: "盒", StorageRequirement: constants.StorageNormal, Volume: 0.004, Weight: 2.5, Price: 12.8},
		{OwnerID: owner2.ID, Name: "电动螺丝刀套装", SKU: "QH-HW-002", Barcode: "690200000002", Category: "电动工具", Spec: "12V 双速", Unit: "套", StorageRequirement: constants.StorageNormal, Volume: 0.03, Weight: 1.8, Price: 199.0},
		{OwnerID: owner3.ID, Name: "冷冻虾仁", SKU: "CC-FD-001", Barcode: "690300000001", Category: "冷冻食品", Spec: "1kg/袋", Unit: "袋", ShelfLifeDays: intPtr(540), StorageRequirement: constants.StorageColdChain, Volume: 0.002, Weight: 1, Price: 45.0},
		{OwnerID: owner3.ID, Name: "冰淇淋(家庭装)", SKU: "CC-FD-002", Barcode: "690300000002", Category: "冷冻食品", Spec: "2L/桶", Unit: "桶", ShelfLifeDays: intPtr(365), StorageRequirement: constants.StorageColdChain, Volume: 0.003, Weight: 1.2, Price: 39.9},
	}
	for i := range products {
		products[i].CreatedAt = now
		products[i].UpdatedAt = now
	}
	if err := db.WithContext(ctx).Create(&products).Error; err != nil {
		return fmt.Errorf("seed products: %w", err)
	}
	bins := make([]model.BinLocation, 0, 36)
	for _, area := range []string{"A", "B", "C", "D"} {
		for layer := 1; layer <= 3; layer++ {
			for column := 1; column <= 3; column++ {
				bins = append(bins, model.BinLocation{
					Area: area, RackNo: "R01", LayerNo: layer, ColumnNo: column,
					Capacity: 8.0, OccupancyRate: 0, StorageRequirement: constants.StorageNormal,
					Status: constants.BinStatusAvailable, CreatedAt: now, UpdatedAt: now,
				})
			}
		}
	}
	// 将 D 区部分库位设置为冷链/危险品/易碎存储要求。
	for i := range bins {
		if bins[i].Area == "D" {
			switch bins[i].ColumnNo {
			case 1:
				bins[i].StorageRequirement = constants.StorageColdChain
			case 2:
				bins[i].StorageRequirement = constants.StorageDangerous
			case 3:
				bins[i].StorageRequirement = constants.StorageFragile
			}
		}
	}
	if err := db.WithContext(ctx).Create(&bins).Error; err != nil {
		return fmt.Errorf("seed bins: %w", err)
	}
	// 示例入库单（Pending）。
	inbound := &model.InboundOrder{
		OrderNo: "INB-" + time.Now().Format("20060102") + "-0001",
		OwnerID: owner1.ID, SupplierName: "华建建材供应商", Status: constants.InboundStatusPending,
		ExpectedArrivalDate: &now, CreatedBy: 2, CreatedAt: now, UpdatedAt: now,
		Items: []model.InboundItem{
			{ProductID: products[0].ID, BatchNo: "B20260801", ExpectedQty: 200},
			{ProductID: products[3].ID, BatchNo: "B20260802", ExpectedQty: 120},
		},
	}
	if err := db.WithContext(ctx).Create(inbound).Error; err != nil {
		return fmt.Errorf("seed inbound: %w", err)
	}
	outbound := &model.OutboundOrder{
		OrderNo: "OUT-" + time.Now().Format("20060102") + "-0001",
		OwnerID: owner2.ID, ReceiverName: "华东装饰工程公司", ReceiverAddress: "南京市江宁区装饰城3栋",
		Status: constants.OutboundStatusPending, RequiredShipDate: &now, CreatedBy: 2, CreatedAt: now, UpdatedAt: now,
		Items: []model.OutboundItem{
			{ProductID: products[4].ID, BinLocationID: &bins[0].ID, ExpectedQty: 50},
		},
	}
	if err := db.WithContext(ctx).Create(outbound).Error; err != nil {
		return fmt.Errorf("seed outbound: %w", err)
	}
	logger.Info(constants.LogSeedCompleted, "users", len(users), "owners", 3, "products", len(products), "bins", len(bins))
	return nil
}

func intPtr(v int) *int { return &v }
