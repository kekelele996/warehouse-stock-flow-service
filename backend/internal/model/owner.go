package model

import "time"

// Owner 货主。
type Owner struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:128;uniqueIndex;not null" json:"name"`
	ContactName      string    `gorm:"size:64;not null" json:"contact_name"`
	Phone            string    `gorm:"size:32;not null" json:"phone"`
	Email            string    `gorm:"size:128" json:"email"`
	Address          string    `gorm:"size:255" json:"address"`
	SettlementMethod string    `gorm:"size:16;not null" json:"settlement_method"`
	CreditLimit      float64   `gorm:"type:numeric(14,2);default:0" json:"credit_limit"`
	CurrentDebt      float64   `gorm:"type:numeric(14,2);default:0" json:"current_debt"`
	Status           string    `gorm:"size:16;default:Active" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
