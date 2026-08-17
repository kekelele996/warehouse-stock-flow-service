package model

import "time"

// BinLocation 库位。
type BinLocation struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Area              string    `gorm:"size:4;index;not null" json:"area"`
	RackNo            string    `gorm:"size:16;not null" json:"rack_no"`
	LayerNo           int       `gorm:"not null" json:"layer_no"`
	ColumnNo          int       `gorm:"not null" json:"column_no"`
	Capacity          float64   `gorm:"type:numeric(10,2);not null" json:"capacity"`
	OccupancyRate     float64   `gorm:"type:numeric(6,2);default:0" json:"occupancy_rate"`
	StorageRequirement string   `gorm:"size:16;not null;default:Normal" json:"storage_requirement"`
	Status            string    `gorm:"size:16;not null;default:Available" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
