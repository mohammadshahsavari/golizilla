package rbac

import (
	"golizilla/domain/model"

	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	DB *gorm.DB
}

func (r *RoleRepositoryImpl) CreateRole(role *model.Role) error {
	return r.DB.Create(role).Error
}

func (r *RoleRepositoryImpl) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := r.DB.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepositoryImpl) GetAllRoles() ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepositoryImpl) DeleteRole(id uint) error {
	if err := r.DB.Delete(&model.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}
