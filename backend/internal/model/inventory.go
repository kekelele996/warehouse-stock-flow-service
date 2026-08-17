package model

import "time"

// Inventory 库存记录（按商品×批次×库位）。
type Inventory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProductID     uint      `gorm:"index:idx_inv_prod_bin,unique;not null" json:"product_id"`
	OwnerID       uint      `gorm:"index;not null" json:"owner_id"`
	BinLocationID uint      `gorm:"index:idx_inv_prod_bin,unique;not null" json:"bin_location_id"`
	BatchNo       string    `gorm:"size:64;index:idx_inv_prod_bin,unique" json:"batch_no"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
