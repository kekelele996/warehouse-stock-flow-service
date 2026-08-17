package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// OwnerCreateRequest 创建货主请求。
type OwnerCreateRequest struct {
	Name             string  `json:"name" binding:"required,min=1,max=128"`
	ContactName      string  `json:"contact_name" binding:"required,min=1,max=64"`
	Phone            string  `json:"phone" binding:"required,max=32"`
	Email            string  `json:"email" binding:"omitempty,email,max=128"`
	Address          string  `json:"address" binding:"max=255"`
	SettlementMethod string  `json:"settlement_method" binding:"required,oneof=Monthly PerOrder Prepaid"`
	CreditLimit      float64 `json:"credit_limit" binding:"min=0"`
}

// OwnerUpdateRequest 更新货主请求。
type OwnerUpdateRequest struct {
	ContactName string `json:"contact_name" binding:"required,min=1,max=64"`
	Phone       string `json:"phone" binding:"required,max=32"`
	Email       string `json:"email" binding:"omitempty,email,max=128"`
	Address     string `json:"address" binding:"max=255"`
}

// OwnerCreditRequest 信用额度调整请求。
type OwnerCreditRequest struct {
	CreditLimit float64 `json:"credit_limit" binding:"required,min=0"`
}

// OwnerQuery 货主查询参数。
type OwnerQuery struct {
	Keyword string `form:"keyword"`
	Status  string `form:"status"`
}

// OwnerView 货主视图。
type OwnerView struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	ContactName      string    `json:"contact_name"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	Address          string    `json:"address"`
	SettlementMethod string    `json:"settlement_method"`
	SettlementText   string    `json:"settlement_text"`
	CreditLimit      float64   `json:"credit_limit"`
	CurrentDebt      float64   `json:"current_debt"`
	Status           string    `json:"status"`
	StatusText       string    `json:"status_text"`
	CreatedAt        time.Time `json:"created_at"`
}

// OwnerStatsView 货主详情统计。
type OwnerStatsView struct {
	Owner          *OwnerView `json:"owner"`
	ProductCount   int64      `json:"product_count"`
	InboundCount   int64      `json:"inbound_count"`
	OutboundCount  int64      `json:"outbound_count"`
	InventoryValue float64    `json:"inventory_value"`
}

// ToOwnerView 模型转视图。
func ToOwnerView(o *model.Owner) *OwnerView {
	return &OwnerView{
		ID:               o.ID,
		Name:             o.Name,
		ContactName:      o.ContactName,
		Phone:            o.Phone,
		Email:            o.Email,
		Address:          o.Address,
		SettlementMethod: o.SettlementMethod,
		SettlementText:   settlementText(o.SettlementMethod),
		CreditLimit:      o.CreditLimit,
		CurrentDebt:      o.CurrentDebt,
		Status:           o.Status,
		StatusText:       ownerStatusText(o.Status),
		CreatedAt:        o.CreatedAt,
	}
}
