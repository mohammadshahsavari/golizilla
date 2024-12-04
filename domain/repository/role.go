package repository

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	Add(role *model.Role) error
	Delete(id uuid.UUID) error
	Update(role *model.Role) error
	GetById(id uuid.UUID) (*model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Add(role *model.Role) error {
	err := r.db.Create(role).Error
	if err != nil {
		//log
	}
	return err
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&model.Role{}, id).Error
	if err != nil {
		//log
	}
	return err
}

func (r *roleRepository) Update(role *model.Role) error {
	err := r.db.Save(role).Error
	if err != nil {
		//log
	}
	return err
}

func (r *roleRepository) GetById(id uuid.UUID) (*model.Role, error) {
	var role model.Role
	err := r.db.Find(&role, id).Error
	if err != nil {
		//log
	}

	return &role, err
}
