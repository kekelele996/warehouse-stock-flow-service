package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/wmsflow/wmsflow/internal/model"
	"github.com/wmsflow/wmsflow/internal/repository"
)

// 通用测试 fake：通过函数字段注入行为，未注入的方法调用会 panic（便于暴露测试盲区）。

type fakeUserRepo struct {
	repository.UserRepository
	createFn         func(ctx context.Context, user *model.User) error
	findByUsernameFn func(ctx context.Context, username string) (*model.User, error)
	findByIDFn       func(ctx context.Context, id uint) (*model.User, error)
}

func (f *fakeUserRepo) WithTx(tx *gorm.DB) repository.UserRepository { return f }
func (f *fakeUserRepo) Create(ctx context.Context, user *model.User) error {
	if f.createFn != nil {
		return f.createFn(ctx, user)
	}
	return f.UserRepository.Create(ctx, user)
}
func (f *fakeUserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	if f.findByUsernameFn != nil {
		return f.findByUsernameFn(ctx, username)
	}
	return f.UserRepository.FindByUsername(ctx, username)
}
func (f *fakeUserRepo) FindByID(ctx context.Context, id uint) (*model.User, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return f.UserRepository.FindByID(ctx, id)
}

type fakeOwnerRepo struct {
	repository.OwnerRepository
	createFn     func(ctx context.Context, owner *model.Owner) error
	findByIDFn   func(ctx context.Context, id uint) (*model.Owner, error)
	updateFn     func(ctx context.Context, owner *model.Owner) error
	listFn       func(ctx context.Context, filter repository.OwnerFilter) ([]model.Owner, int64, error)
	countFn      func(ctx context.Context) (int64, error)
}

func (f *fakeOwnerRepo) WithTx(tx *gorm.DB) repository.OwnerRepository { return f }
func (f *fakeOwnerRepo) Create(ctx context.Context, owner *model.Owner) error {
	if f.createFn != nil {
		return f.createFn(ctx, owner)
	}
	return f.OwnerRepository.Create(ctx, owner)
}
func (f *fakeOwnerRepo) FindByID(ctx context.Context, id uint) (*model.Owner, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return f.OwnerRepository.FindByID(ctx, id)
}
func (f *fakeOwnerRepo) Update(ctx context.Context, owner *model.Owner) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, owner)
	}
	return f.OwnerRepository.Update(ctx, owner)
}
func (f *fakeOwnerRepo) List(ctx context.Context, filter repository.OwnerFilter) ([]model.Owner, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.OwnerRepository.List(ctx, filter)
}
func (f *fakeOwnerRepo) Count(ctx context.Context) (int64, error) {
	if f.countFn != nil {
		return f.countFn(ctx)
	}
	return f.OwnerRepository.Count(ctx)
}

type fakeProductRepo struct {
	repository.ProductRepository
	createFn   func(ctx context.Context, product *model.Product) error
	findByIDFn func(ctx context.Context, id uint) (*model.Product, error)
	updateFn   func(ctx context.Context, product *model.Product) error
	listFn     func(ctx context.Context, filter repository.ProductFilter) ([]model.Product, int64, error)
	countFn    func(ctx context.Context, ownerID uint) (int64, error)
}

func (f *fakeProductRepo) WithTx(tx *gorm.DB) repository.ProductRepository { return f }
func (f *fakeProductRepo) Create(ctx context.Context, product *model.Product) error {
	if f.createFn != nil {
		return f.createFn(ctx, product)
	}
	return f.ProductRepository.Create(ctx, product)
}
func (f *fakeProductRepo) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return f.ProductRepository.FindByID(ctx, id)
}
func (f *fakeProductRepo) Update(ctx context.Context, product *model.Product) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, product)
	}
	return f.ProductRepository.Update(ctx, product)
}
func (f *fakeProductRepo) List(ctx context.Context, filter repository.ProductFilter) ([]model.Product, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filter)
	}
	return f.ProductRepository.List(ctx, filter)
}
func (f *fakeProductRepo) CountByOwner(ctx context.Context, ownerID uint) (int64, error) {
	if f.countFn != nil {
		return f.countFn(ctx, ownerID)
	}
	return f.ProductRepository.CountByOwner(ctx, ownerID)
}
