package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// InventoryQuery 库存查询参数。
type InventoryQuery struct {
	OwnerID       uint   `form:"owner_id"`
	ProductID     uint   `form:"product_id"`
	BinLocationID uint   `form:"bin_location_id"`
	Keyword       string `form:"keyword"`
}

// InventoryView 库存视图。
type InventoryView struct {
	ID            uint      `json:"id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"`
	SKU           string    `json:"sku"`
	OwnerID       uint      `json:"owner_id"`
	OwnerName     string    `json:"owner_name"`
	BinLocationID uint      `json:"bin_location_id"`
	BinCode       string    `json:"bin_code"`
	BatchNo       string    `json:"batch_no"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	TotalValue    float64   `json:"total_value"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// InventorySummaryItem 货主库存金额汇总。
type InventorySummaryItem struct {
	OwnerID   uint    `json:"owner_id"`
	OwnerName string  `json:"owner_name"`
	TotalQty  int     `json:"total_qty"`
	TotalValue float64 `json:"total_value"`
}

// ToInventoryView 模型转视图。
func ToInventoryView(inv *model.Inventory, productName, sku, ownerName, binCode string, unitPrice float64) *InventoryView {
	return &InventoryView{
		ID:            inv.ID,
		ProductID:     inv.ProductID,
		ProductName:   productName,
		SKU:           sku,
		OwnerID:       inv.OwnerID,
		OwnerName:     ownerName,
		BinLocationID: inv.BinLocationID,
		BinCode:       binCode,
		BatchNo:       inv.BatchNo,
		Quantity:      inv.Quantity,
		UnitPrice:     unitPrice,
		TotalValue:    unitPrice * float64(inv.Quantity),
		UpdatedAt:     inv.UpdatedAt,
	}
}
