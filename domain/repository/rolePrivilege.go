package repository

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRolePrivilegeRepository interface {
	Add(rolePrivelege *model.RolePrivilege) error
	Delete(roleId uuid.UUID, privilegeId string) error
	GetRolePrivileges(roleId uuid.UUID) ([]model.RolePrivilege, error)
	GetRolePrivilegesByPrivileges(roleId uuid.UUID, privileges ...string) ([]model.RolePrivilege, error)
}

type rolePrivilege struct {
	db *gorm.DB
}

func NewRolePrivilegeRepository(db *gorm.DB) IRolePrivilegeRepository {
	return &rolePrivilege{db: db}
}

func (r *rolePrivilege) Add(rolePrivelege *model.RolePrivilege) error {
	err := r.db.Create(rolePrivelege).Error
	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilege) Delete(roleId uuid.UUID, privilegeId string) error {
	err := r.db.Delete(&model.RolePrivilege{
		RoleId:      roleId,
		PrivilegeId: privilegeId,
	}).Error

	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilege) GetRolePrivileges(roleId uuid.UUID) ([]model.RolePrivilege, error) {
	var rolePrivileges []model.RolePrivilege
	err := r.db.Where("role_id = ?", roleId).Find(&rolePrivileges).Error
	if err != nil {
		//log
	}

	return rolePrivileges, err
}

func (r *rolePrivilege) GetRolePrivilegesByPrivileges(roleId uuid.UUID, privileges ...string) ([]model.RolePrivilege, error) {
	var rolePrivileges []model.RolePrivilege
	err := r.db.Where("role_id = ? And privilege_id IN privileges", roleId).Find(&rolePrivileges).Error
	if err != nil {
		//log
	}

	return rolePrivileges, err
}
