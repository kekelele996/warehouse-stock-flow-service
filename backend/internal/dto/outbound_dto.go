package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// OutboundItemRequest 出库明细请求。
type OutboundItemRequest struct {
	ProductID      uint `json:"product_id" binding:"required,gt=0"`
	BinLocationID  uint `json:"bin_location_id" binding:"required,gt=0"`
	ExpectedQty    int  `json:"expected_qty" binding:"required,gt=0"`
}

// OutboundCreateRequest 创建出库单请求。
type OutboundCreateRequest struct {
	OwnerID          uint                 `json:"owner_id" binding:"required,gt=0"`
	ReceiverName     string               `json:"receiver_name" binding:"required,min=1,max=128"`
	ReceiverAddress  string               `json:"receiver_address" binding:"max=255"`
	RequiredShipDate *time.Time           `json:"required_ship_date"`
	Items            []OutboundItemRequest `json:"items" binding:"required,min=1,dive"`
}

// OutboundPickingItemRequest 拣货明细请求。
type OutboundPickingItemRequest struct {
	ItemID    uint `json:"item_id" binding:"required,gt=0"`
	ActualQty int  `json:"actual_qty" binding:"required,min=0"`
}

// OutboundPickingRequest 拣货任务确认请求。
type OutboundPickingRequest struct {
	Items []OutboundPickingItemRequest `json:"items" binding:"required,min=1,dive"`
}

// OutboundShipRequest 发货确认请求。
type OutboundShipRequest struct {
	TrackingNo string `json:"tracking_no" binding:"required,min=1,max=64"`
}

// OutboundQuery 出库单查询参数。
type OutboundQuery struct {
	Status  string `form:"status"`
	OwnerID uint   `form:"owner_id"`
	Keyword string `form:"keyword"`
}

// OutboundItemView 出库明细视图。
type OutboundItemView struct {
	ID            uint   `json:"id"`
	ProductID     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	SKU           string `json:"sku"`
	BinLocationID *uint  `json:"bin_location_id"`
	BinCode       string `json:"bin_code"`
	ExpectedQty   int    `json:"expected_qty"`
	ActualQty     int    `json:"actual_qty"`
}

// OutboundOrderView 出库单视图。
type OutboundOrderView struct {
	ID               uint               `json:"id"`
	OrderNo          string             `json:"order_no"`
	OwnerID          uint               `json:"owner_id"`
	OwnerName        string             `json:"owner_name"`
	ReceiverName     string             `json:"receiver_name"`
	ReceiverAddress  string             `json:"receiver_address"`
	RequiredShipDate *time.Time         `json:"required_ship_date"`
	ActualShipDate   *time.Time         `json:"actual_ship_date"`
	Status           string             `json:"status"`
	StatusText       string             `json:"status_text"`
	PickerName       string             `json:"picker_name"`
	CheckerName      string             `json:"checker_name"`
	TrackingNo       string             `json:"tracking_no"`
	Items            []OutboundItemView `json:"items"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// ToOutboundItemView 明细模型转视图。
func ToOutboundItemView(item *model.OutboundItem, productName, sku, binCode string) OutboundItemView {
	return OutboundItemView{
		ID:            item.ID,
		ProductID:     item.ProductID,
		ProductName:   productName,
		SKU:           sku,
		BinLocationID: item.BinLocationID,
		BinCode:       binCode,
		ExpectedQty:   item.ExpectedQty,
		ActualQty:     item.ExpectedQty,
	}
}

// ToOutboundOrderView 出库单模型转视图。
func ToOutboundOrderView(o *model.OutboundOrder, ownerName, pickerName, checkerName string) *OutboundOrderView {
	v := &OutboundOrderView{
		ID:               o.ID,
		OrderNo:          o.OrderNo,
		OwnerID:          o.OwnerID,
		OwnerName:        ownerName,
		ReceiverName:     o.ReceiverName,
		ReceiverAddress:  o.ReceiverAddress,
		RequiredShipDate: o.RequiredShipDate,
		ActualShipDate:   o.ActualShipDate,
		Status:           o.Status,
		StatusText:       outboundStatusText(o.Status),
		PickerName:       pickerName,
		CheckerName:      checkerName,
		TrackingNo:       o.TrackingNo,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
		Items:            make([]OutboundItemView, 0, len(o.Items)),
	}
	for i := range o.Items {
		v.Items = append(v.Items, ToOutboundItemView(&o.Items[i], "", "", ""))
	}
	return v
}
