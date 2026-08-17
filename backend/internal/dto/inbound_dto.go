package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// InboundItemRequest 入库明细请求。
type InboundItemRequest struct {
	ProductID   uint   `json:"product_id" binding:"required,gt=0"`
	BatchNo     string `json:"batch_no" binding:"required,min=1,max=64"`
	ExpectedQty int    `json:"expected_qty" binding:"required,gt=0"`
}

// InboundCreateRequest 创建入库单请求。
type InboundCreateRequest struct {
	OwnerID             uint                 `json:"owner_id" binding:"required,gt=0"`
	SupplierName        string               `json:"supplier_name" binding:"max=128"`
	ExpectedArrivalDate *time.Time           `json:"expected_arrival_date"`
	Remark              string               `json:"remark" binding:"max=500"`
	Items               []InboundItemRequest `json:"items" binding:"required,min=1,dive"`
}

// InboundQCItemRequest 质检明细请求。
type InboundQCItemRequest struct {
	ItemID     uint   `json:"item_id" binding:"required,gt=0"`
	ActualQty  int    `json:"actual_qty" binding:"required,min=0"`
	QCResult   string `json:"qc_result" binding:"required,oneof=Pass Fail Partial"`
}

// InboundQCRequest 收货质检请求（逐项录入实收数量与质检结果）。
type InboundQCRequest struct {
	Items []InboundQCItemRequest `json:"items" binding:"required,min=1,dive"`
}

// InboundShelveItemRequest 上架明细请求。
type InboundShelveItemRequest struct {
	ItemID        uint `json:"item_id" binding:"required,gt=0"`
	BinLocationID uint `json:"bin_location_id" binding:"required,gt=0"`
}

// InboundShelveRequest 上架指派库位请求。
type InboundShelveRequest struct {
	Items []InboundShelveItemRequest `json:"items" binding:"required,min=1,dive"`
}

// InboundQuery 入库单查询参数。
type InboundQuery struct {
	Status  string `form:"status"`
	OwnerID uint   `form:"owner_id"`
	Keyword string `form:"keyword"`
}

// InboundItemView 入库明细视图。
type InboundItemView struct {
	ID             uint    `json:"id"`
	ProductID      uint    `json:"product_id"`
	ProductName    string  `json:"product_name"`
	SKU            string  `json:"sku"`
	BatchNo        string  `json:"batch_no"`
	ExpectedQty    int     `json:"expected_qty"`
	ActualQty      int     `json:"actual_qty"`
	QCResult       string  `json:"qc_result"`
	QCResultText   string  `json:"qc_result_text"`
	BinLocationID  *uint   `json:"bin_location_id"`
	BinCode        string  `json:"bin_code"`
}

// InboundOrderView 入库单视图。
type InboundOrderView struct {
	ID                  uint              `json:"id"`
	OrderNo             string            `json:"order_no"`
	OwnerID             uint              `json:"owner_id"`
	OwnerName           string            `json:"owner_name"`
	SupplierName        string            `json:"supplier_name"`
	ExpectedArrivalDate *time.Time        `json:"expected_arrival_date"`
	ActualArrivalDate   *time.Time        `json:"actual_arrival_date"`
	Status              string            `json:"status"`
	StatusText          string            `json:"status_text"`
	QCInspectorName     string            `json:"qc_inspector_name"`
	KeeperName          string            `json:"keeper_name"`
	Remark              string            `json:"remark"`
	Items               []InboundItemView `json:"items"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// ToInboundItemView 明细模型转视图。
func ToInboundItemView(item *model.InboundItem, productName, sku, binCode string) InboundItemView {
	return InboundItemView{
		ID:            item.ID,
		ProductID:     item.ProductID,
		ProductName:   productName,
		SKU:           sku,
		BatchNo:       item.BatchNo,
		ExpectedQty:   item.ExpectedQty,
		ActualQty:     item.ExpectedQty,
		QCResult:      item.QCResult,
		QCResultText:  qcResultText(item.QCResult),
		BinLocationID: item.BinLocationID,
		BinCode:       binCode,
	}
}

// ToInboundOrderView 入库单模型转视图。
func ToInboundOrderView(o *model.InboundOrder, ownerName, qcInspectorName, keeperName string) *InboundOrderView {
	v := &InboundOrderView{
		ID:                  o.ID,
		OrderNo:             o.OrderNo,
		OwnerID:             o.OwnerID,
		OwnerName:           ownerName,
		SupplierName:        o.SupplierName,
		ExpectedArrivalDate: o.ExpectedArrivalDate,
		ActualArrivalDate:   o.ActualArrivalDate,
		Status:              o.Status,
		StatusText:          inboundStatusText(o.Status),
		QCInspectorName:     qcInspectorName,
		KeeperName:          keeperName,
		Remark:              o.Remark,
		CreatedAt:           o.CreatedAt,
		UpdatedAt:           o.UpdatedAt,
		Items:               make([]InboundItemView, 0, len(o.Items)),
	}
	for i := range o.Items {
		v.Items = append(v.Items, ToInboundItemView(&o.Items[i], "", "", ""))
	}
	return v
}
