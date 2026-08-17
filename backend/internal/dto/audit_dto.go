package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// AuditQuery 操作日志查询参数。
type AuditQuery struct {
	Module string `form:"module"`
	Action string `form:"action"`
	UserID uint   `form:"user_id"`
}

// AuditView 操作日志视图。
type AuditView struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	RoleText   string    `json:"role_text"`
	Module     string    `json:"module"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Detail     string    `json:"detail"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToAuditView 模型转视图。
func ToAuditView(l *model.OperationLog) *AuditView {
	return &AuditView{
		ID:         l.ID,
		UserID:     l.UserID,
		Username:   l.Username,
		Role:       l.Role,
		RoleText:   roleText(l.Role),
		Module:     l.Module,
		Action:     l.Action,
		EntityType: l.EntityType,
		EntityID:   l.EntityID,
		Detail:     l.Detail,
		IP:         l.IP,
		CreatedAt:  l.CreatedAt,
	}
}
