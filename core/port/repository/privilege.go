package repository

import (
	"context"
	myContext "golizilla/adapters/http/handler/context"
	"golizilla/core/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	Add(ctx context.Context, userCtx context.Context, privilege *model.Privilege) error
	GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Privilege, error)
}

type privilgeRepository struct {
	db *gorm.DB
}

func NewPrivilegeRepository(db *gorm.DB) IPrivilegeRepository {
	return &privilgeRepository{db: db}
}

func (r *privilgeRepository) Add(ctx context.Context, userCtx context.Context, privilege *model.Privilege) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(privilege).Error
	if err != nil {
		//log
	}
	return err
}

func (r *privilgeRepository) GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Privilege, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var privilege model.Privilege
	err := db.WithContext(ctx).Find(&privilege, id).Error
	if err != nil {
		//log
	}

	return &privilege, err
}
