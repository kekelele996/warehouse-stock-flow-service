package model

import "time"

// Product 商品。
type Product struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	OwnerID           uint      `gorm:"index;not null" json:"owner_id"`
	Name              string    `gorm:"size:128;not null" json:"name"`
	SKU               string    `gorm:"size:64;uniqueIndex;not null" json:"sku"`
	Barcode           string    `gorm:"size:64;index" json:"barcode"`
	Category          string    `gorm:"size:64" json:"category"`
	Spec              string    `gorm:"size:128" json:"spec"`
	Unit              string    `gorm:"size:16;not null" json:"unit"`
	ShelfLifeDays     *int      `json:"shelf_life_days"`
	StorageRequirement string   `gorm:"size:16;not null;default:Normal" json:"storage_requirement"`
	Volume            float64   `gorm:"type:numeric(12,3);default:0" json:"volume"`
	Weight            float64   `gorm:"type:numeric(12,3);default:0" json:"weight"`
	Price             float64   `gorm:"type:numeric(14,2);default:0" json:"price"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
