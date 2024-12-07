package repository

import (
	"context"
	myContext "golizilla/adapters/http/handler/context"
	"golizilla/core/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	Add(ctx context.Context, userCtx context.Context, role *model.Role) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	Update(ctx context.Context, userCtx context.Context, role *model.Role) error
	GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Add(ctx context.Context, userCtx context.Context, role *model.Role) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(role).Error
	if err != nil {
		//log
	}
	return err
}

func (r *roleRepository) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Delete(&model.Role{}, id).Error
	if err != nil {
		//log
	}
	return err
}

func (r *roleRepository) Update(ctx context.Context, userCtx context.Context, role *model.Role) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Save(role).Error
	if err != nil {
		//log
	}
	return err
}

func (r *roleRepository) GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Role, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var role model.Role
	err := db.WithContext(ctx).Find(&role, id).Error
	if err != nil {
		//log
	}

	return &role, err
}
