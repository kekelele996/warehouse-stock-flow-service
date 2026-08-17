package model

import "time"

// InboundOrder 入库单。
type InboundOrder struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	OrderNo            string         `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	OwnerID            uint           `gorm:"index;not null" json:"owner_id"`
	SupplierName       string         `gorm:"size:128" json:"supplier_name"`
	ExpectedArrivalDate *time.Time    `json:"expected_arrival_date"`
	ActualArrivalDate  *time.Time     `json:"actual_arrival_date"`
	Status             string         `gorm:"size:16;index;not null" json:"status"`
	QCInspectorID      *uint          `gorm:"index" json:"qc_inspector_id"`
	KeeperID           *uint          `gorm:"index" json:"keeper_id"`
	Remark             string         `gorm:"size:500" json:"remark"`
	CreatedBy          uint           `json:"created_by"`
	Items              []InboundItem  `gorm:"foreignKey:InboundOrderID" json:"items"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// InboundItem 入库明细。
type InboundItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	InboundOrderID uint      `gorm:"index;not null" json:"inbound_order_id"`
	ProductID      uint      `gorm:"index;not null" json:"product_id"`
	BatchNo        string    `gorm:"size:64" json:"batch_no"`
	ExpectedQty    int       `gorm:"not null" json:"expected_qty"`
	ActualQty      int       `gorm:"default:0" json:"actual_qty"`
	QCResult       string    `gorm:"size:16" json:"qc_result"`
	BinLocationID  *uint     `gorm:"index" json:"bin_location_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
