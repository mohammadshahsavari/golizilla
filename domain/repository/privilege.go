package repository

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	Add(privilege *model.Privilege) error
	GetById(id uuid.UUID) (*model.Privilege, error)
}

type privilgeRepository struct {
	db *gorm.DB
}

func NewPrivilegeRepository(db *gorm.DB) IPrivilegeRepository {
	return &privilgeRepository{db: db}
}

func (r *privilgeRepository) Add(privilege *model.Privilege) error {
	err := r.db.Create(privilege).Error
	if err != nil {
		//log
	}
	return err
}

func (r *privilgeRepository) GetById(id uuid.UUID) (*model.Privilege, error) {
	var privilege model.Privilege
	err := r.db.Find(&privilege, id).Error
	if err != nil {
		//log
	}

	return &privilege, err
}
