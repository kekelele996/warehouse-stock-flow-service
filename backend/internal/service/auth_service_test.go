package service

import (
	"context"
	"database/sql"
	"testing"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/dto"
	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/util"
)

type fakeTx struct{}

func (fakeTx) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error { return fc(nil) }

func newAuthServiceForTest(userRepo repository.UserRepository, ownerRepo repository.OwnerRepository, logRepo repository.OperationLogRepository) AuthService {
	audit := NewAuditService(logRepo, testLogger())
	return NewAuthService(fakeTx{}, "test-secret-0123456789", 72, userRepo, ownerRepo, audit, testLogger())
}

func TestAuthService_Login(t *testing.T) {
	ownerID := uint(7)
	hashed, _ := util.HashPassword("secret123")
	user := &model.User{ID: 1, Username: "owner01", Password: hashed, Name: "张三", Role: constants.RoleOwner, OwnerID: &ownerID, Status: constants.UserStatusActive}
	tests := []struct {
		name     string
		username string
		password string
		setup    func(u *fakeUserRepo)
		wantErr  bool
		wantCode int
	}{
		{
			name: "正确密码登录成功", username: "owner01", password: "secret123",
			setup: func(u *fakeUserRepo) {
				u.findByUsernameFn = func(ctx context.Context, name string) (*model.User, error) { return user, nil }
			},
			wantErr: false,
		},
		{
			name: "错误密码返回401", username: "owner01", password: "wrong",
			setup: func(u *fakeUserRepo) {
				u.findByUsernameFn = func(ctx context.Context, name string) (*model.User, error) { return user, nil }
			},
			wantErr: true, wantCode: constants.CodeCredentialError,
		},
		{
			name: "用户不存在返回401", username: "nobody", password: "secret123",
			setup: func(u *fakeUserRepo) {
				u.findByUsernameFn = func(ctx context.Context, name string) (*model.User, error) { return nil, repository.ErrNotFound }
			},
			wantErr: true, wantCode: constants.CodeCredentialError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &fakeUserRepo{}
			tt.setup(userRepo)
			svc := newAuthServiceForTest(userRepo, &fakeOwnerRepo{}, &fakeLogRepo{})
			resp, err := svc.Login(context.Background(), &dto.LoginRequest{Username: tt.username, Password: tt.password})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				appErr := util.ToAppError(err)
				if appErr.Code != tt.wantCode {
					t.Fatalf("expected code %d, got %d (%v)", tt.wantCode, appErr.Code, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil || resp.Token == "" || resp.User == nil {
				t.Fatalf("expected token and user in response")
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	userRepo := &fakeUserRepo{}
	created := &model.User{}
	userRepo.createFn = func(ctx context.Context, u *model.User) error {
		*created = *u
		u.ID = 99
		return nil
	}
	ownerRepo := &fakeOwnerRepo{}
	ownerRepo.createFn = func(ctx context.Context, o *model.Owner) error {
		o.ID = 42
		return nil
	}
	svc := newAuthServiceForTest(userRepo, ownerRepo, &fakeLogRepo{})
	view, err := svc.Register(context.Background(), &dto.RegisterRequest{
		Username: "newowner", Password: "abc12345", Name: "新货主", OwnerName: "新货主公司", ContactName: "赵六", Phone: "13900000000",
	})
	if err != nil {
		t.Fatalf("register error: %v", err)
	}
	if view.Role != constants.RoleOwner {
		t.Fatalf("expected role Owner, got %s", view.Role)
	}
	if created.OwnerID == nil || *created.OwnerID != 42 {
		t.Fatalf("expected user linked to owner 42, got %v", created.OwnerID)
	}
	if created.Password == "abc12345" {
		t.Fatalf("password must be hashed")
	}
}

func TestAuthService_Me(t *testing.T) {
	userRepo := &fakeUserRepo{}
	userRepo.findByIDFn = func(ctx context.Context, id uint) (*model.User, error) {
		return &model.User{ID: 1, Username: "admin", Name: "系统管理员", Role: constants.RoleAdmin, Status: constants.UserStatusActive}, nil
	}
	svc := newAuthServiceForTest(userRepo, &fakeOwnerRepo{}, &fakeLogRepo{})
	ctx := util.WithUser(context.Background(), &util.Claims{UserID: 1, Username: "admin", Role: constants.RoleAdmin})
	view, err := svc.Me(ctx)
	if err != nil {
		t.Fatalf("me error: %v", err)
	}
	if view.Username != "admin" || view.Role != constants.RoleAdmin {
		t.Fatalf("unexpected me result: %+v", view)
	}
}
