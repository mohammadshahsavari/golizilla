package repository

import (
	"context"
	"fmt"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRolePrivilegeRepository interface {
	Add(ctx context.Context, rolePrivelege *model.RolePrivilege) error
	Delete(ctx context.Context, roleId uuid.UUID, privilegeId string) error
	GetRolePrivileges(ctx context.Context, roleId uuid.UUID) ([]model.RolePrivilege, error)
	GetRolePrivilegesByPrivileges(ctx context.Context, roleId uuid.UUID, privileges ...string) ([]model.RolePrivilege, error)
}

type rolePrivilege struct {
	db *gorm.DB
}

func NewRolePrivilegeRepository(db *gorm.DB) IRolePrivilegeRepository {
	return &rolePrivilege{db: db}
}

func (r *rolePrivilege) Add(ctx context.Context, rolePrivelege *model.RolePrivilege) error {
	err := r.db.WithContext(ctx).Create(rolePrivelege).Error
	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilege) Delete(ctx context.Context, roleId uuid.UUID, privilegeId string) error {
	err := r.db.WithContext(ctx).Delete(&model.RolePrivilege{
		RoleId:      roleId,
		PrivilegeId: privilegeId,
	}).Error

	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilege) GetRolePrivileges(ctx context.Context, roleId uuid.UUID) ([]model.RolePrivilege, error) {
	var rolePrivileges []model.RolePrivilege
	err := r.db.WithContext(ctx).Where("role_id = ?", roleId).Find(&rolePrivileges).Error
	if err != nil {
		//log
	}

	return rolePrivileges, err
}

func (r *rolePrivilege) GetRolePrivilegesByPrivileges(ctx context.Context, roleId uuid.UUID, privileges ...string) ([]model.RolePrivilege, error) {
	var rolePrivileges []model.RolePrivilege
	if len(privileges) == 0 {
		return nil, fmt.Errorf("privilegeIds cannot be empty")
	}
	// Query matching both role_id and privilege_id
	err := r.db.WithContext(ctx).Where("role_id = ? AND privilege_id IN ?", roleId, privileges).Find(&rolePrivileges).Error

	if err != nil {
		// Log the error if necessary
		// log.Printf("Error fetching role privileges: %v", err)
		return nil, err
	}

	return rolePrivileges, err
}
