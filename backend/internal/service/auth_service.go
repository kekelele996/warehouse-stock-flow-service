package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

// AuthService 认证服务。
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserView, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Me(ctx context.Context) (*dto.UserView, error)
}

type authService struct {
	db         txProvider
	jwtSecret  string
	jwtExpire  int
	userRepo   repository.UserRepository
	ownerRepo  repository.OwnerRepository
	auditSvc   AuditService
	logger     *slog.Logger
}

// NewAuthService 构造认证服务。
func NewAuthService(db txProvider, jwtSecret string, jwtExpire int, userRepo repository.UserRepository, ownerRepo repository.OwnerRepository, auditSvc AuditService, logger *slog.Logger) AuthService {
	return &authService{db: db, jwtSecret: jwtSecret, jwtExpire: jwtExpire, userRepo: userRepo, ownerRepo: ownerRepo, auditSvc: auditSvc, logger: logger}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserView, error) {
	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("auth register: hash password: %w", err)
	}
	var created *model.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		owner := &model.Owner{
			Name:             req.OwnerName,
			ContactName:      req.ContactName,
			Phone:            req.Phone,
			SettlementMethod: constants.SettlementMonthly,
			Status:           constants.OwnerStatusActive,
		}
		if err := s.ownerRepo.WithTx(tx).Create(ctx, owner); err != nil {
			return fmt.Errorf("auth register: create owner %s: %w", req.OwnerName, err)
		}
		user := &model.User{
			Username: req.Username,
			Password: hashed,
			Name:     req.Name,
			Role:     constants.RoleOwner,
			OwnerID:  &owner.ID,
			Status:   constants.UserStatusActive,
		}
		if err := s.userRepo.WithTx(tx).Create(ctx, user); err != nil {
			return fmt.Errorf("auth register: create user %s: %w", req.Username, err)
		}
		created = user
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.InfoContext(ctx, constants.LogRegisterUser, "username", req.Username, "role", constants.RoleOwner)
	return toUserView(created), nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		s.logger.WarnContext(ctx, constants.LogLoginFailed, "username", req.Username, "reason", "user not found")
		return nil, util.WrapAppError(
			util.NewAppError(constants.CodeCredentialError, 401, "货主账号用户名或密码错误"),
			err,
		)
	}
	if !util.CheckPassword(user.Password, req.Password) {
		s.logger.WarnContext(ctx, constants.LogLoginFailed, "username", req.Username, "reason", "password mismatch")
		return nil, util.WrapAppError(
			util.NewAppError(constants.CodeCredentialError, 401, "货主账号用户名或密码错误"),
			errors.New("password mismatch"),
		)
	}
	if user.Status != constants.UserStatusActive {
		s.logger.WarnContext(ctx, constants.LogLoginFailed, "username", req.Username, "reason", "user disabled")
		return nil, util.NewAppError(constants.CodeForbidden, 403, "用户账号已被禁用，请联系管理员")
	}
	token, err := util.GenerateToken(s.jwtSecret, s.jwtExpire, user.ID, user.Username, user.Role, user.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("auth login: generate token for user %s: %w", req.Username, err)
	}
	s.logger.InfoContext(ctx, constants.LogLoginSuccess, "username", user.Username, "role", user.Role)
	return &dto.LoginResponse{Token: token, User: toUserView(user)}, nil
}

func (s *authService) Me(ctx context.Context) (*dto.UserView, error) {
	claims, ok := util.CurrentUser(ctx)
	if !ok {
		return nil, util.NewAppError(constants.CodeUnauthorized, 401, constants.MsgUnauthorized)
	}
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth me: %w", err)
	}
	return toUserView(user), nil
}

// toUserView 用户模型转视图。
func toUserView(user *model.User) *dto.UserView {
	return &dto.UserView{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
		RoleText: util.RoleText(user.Role),
		OwnerID:  user.OwnerID,
		Status:   user.Status,
	}
}
