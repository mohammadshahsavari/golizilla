package repository

import (
	"context"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	Add(ctx context.Context, privilege *model.Privilege) error
	GetById(ctx context.Context, id uuid.UUID) (*model.Privilege, error)
}

type privilgeRepository struct {
	db *gorm.DB
}

func NewPrivilegeRepository(db *gorm.DB) IPrivilegeRepository {
	return &privilgeRepository{db: db}
}

func (r *privilgeRepository) Add(ctx context.Context, privilege *model.Privilege) error {
	err := r.db.WithContext(ctx).Create(privilege).Error
	if err != nil {
		//log
	}
	return err
}

func (r *privilgeRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Privilege, error) {
	var privilege model.Privilege
	err := r.db.WithContext(ctx).Find(&privilege, id).Error
	if err != nil {
		//log
	}

	return &privilege, err
}
