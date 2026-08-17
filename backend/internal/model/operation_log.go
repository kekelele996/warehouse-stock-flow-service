package model

import "time"

// OperationLog 操作审计日志。
type OperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Username   string    `gorm:"size:64" json:"username"`
	Role       string    `gorm:"size:32" json:"role"`
	Module     string    `gorm:"size:32;index" json:"module"`
	Action     string    `gorm:"size:64;index" json:"action"`
	EntityType string    `gorm:"size:32" json:"entity_type"`
	EntityID   string    `gorm:"size:64" json:"entity_id"`
	Detail     string    `gorm:"size:500" json:"detail"`
	IP         string    `gorm:"size:64" json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}
