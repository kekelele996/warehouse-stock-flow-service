package model

import "time"

// OutboundOrder 出库单。
type OutboundOrder struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	OrderNo          string          `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	OwnerID          uint            `gorm:"index;not null" json:"owner_id"`
	ReceiverName     string          `gorm:"size:128;not null" json:"receiver_name"`
	ReceiverAddress  string          `gorm:"size:255" json:"receiver_address"`
	RequiredShipDate *time.Time      `json:"required_ship_date"`
	ActualShipDate   *time.Time      `json:"actual_ship_date"`
	Status           string          `gorm:"size:16;index;not null" json:"status"`
	PickerID         *uint           `gorm:"index" json:"picker_id"`
	CheckerID        *uint           `gorm:"index" json:"checker_id"`
	TrackingNo       string          `gorm:"size:64" json:"tracking_no"`
	CreatedBy        uint            `json:"created_by"`
	Items            []OutboundItem  `gorm:"foreignKey:OutboundOrderID" json:"items"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// OutboundItem 出库明细。
type OutboundItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OutboundOrderID uint     `gorm:"index;not null" json:"outbound_order_id"`
	ProductID      uint      `gorm:"index;not null" json:"product_id"`
	BinLocationID  *uint     `gorm:"index" json:"bin_location_id"`
	ExpectedQty    int       `gorm:"not null" json:"expected_qty"`
	ActualQty      int       `gorm:"default:0" json:"actual_qty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
