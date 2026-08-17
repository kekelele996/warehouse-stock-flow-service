package repository

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/model"
)

func newTestOwner(name string) *model.Owner {
	now := time.Now()
	return &model.Owner{
		Name: name, ContactName: "测试联系人", Phone: "13800000000",
		Email: "test@example.com", Address: "测试地址",
		SettlementMethod: constants.SettlementMonthly, CreditLimit: 10000, CurrentDebt: 0,
		Status: constants.OwnerStatusActive, CreatedAt: now, UpdatedAt: now,
	}
}

func newTestProduct(ownerID uint, sku string) *model.Product {
	now := time.Now()
	return &model.Product{
		OwnerID: ownerID, Name: "测试商品", SKU: sku, Barcode: "9000" + sku,
		Category: "测试类", Spec: "1件", Unit: "件",
		StorageRequirement: constants.StorageNormal, Volume: 1, Weight: 1, Price: 10,
		CreatedAt: now, UpdatedAt: now,
	}
}

func newTestBin(area, rack string, layer, column int) *model.BinLocation {
	now := time.Now()
	return &model.BinLocation{
		Area: area, RackNo: rack, LayerNo: layer, ColumnNo: column,
		Capacity: 8, OccupancyRate: 0, StorageRequirement: constants.StorageNormal,
		Status: constants.BinStatusAvailable, CreatedAt: now, UpdatedAt: now,
	}
}
