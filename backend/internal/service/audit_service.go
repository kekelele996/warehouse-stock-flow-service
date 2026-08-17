package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

// AuditService 操作日志服务。
type AuditService interface {
	Record(ctx context.Context, module, action, entityType, entityID, detail string) error
	List(ctx context.Context, filter repository.OperationLogFilter) ([]dto.AuditView, int64, error)
}

type auditService struct {
	logRepo repository.OperationLogRepository
	logger  *slog.Logger
}

// NewAuditService 构造操作日志服务。
func NewAuditService(logRepo repository.OperationLogRepository, logger *slog.Logger) AuditService {
	return &auditService{logRepo: logRepo, logger: logger}
}

// Record 写入操作日志；获取不到当前用户时仍记录匿名日志。
func (s *auditService) Record(ctx context.Context, module, action, entityType, entityID, detail string) error {
	claims, ok := util.CurrentUser(ctx)
	entry := &model.OperationLog{
		Module:     module,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Detail:     detail,
	}
	if ok {
		entry.UserID = claims.UserID
		entry.Username = claims.Username
		entry.Role = claims.Role
	}
	if err := s.logRepo.Create(ctx, entry); err != nil {
		return fmt.Errorf("audit record module=%s action=%s: %w", module, action, err)
	}
	s.logger.InfoContext(ctx, constants.LogAuditRecord,
		"module", module, "action", action, "entity_type", entityType, "entity_id", entityID,
		"operator", entry.Username, "role", entry.Role)
	return nil
}

// List 分页查询操作日志。
func (s *auditService) List(ctx context.Context, filter repository.OperationLogFilter) ([]dto.AuditView, int64, error) {
	logs, total, err := s.logRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("audit list: %w", err)
	}
	views := make([]dto.AuditView, 0, len(logs))
	for i := range logs {
		views = append(views, *dto.ToAuditView(&logs[i]))
	}
	return views, total, nil
}
