package repository

import (
	"context"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRolePrivilegeOnInstanceRepository interface {
	Add(ctx context.Context, rolePrivelegeOnInsance *model.RolePrivilegeOnInstance) error
	Delete(ctx context.Context, roleId uuid.UUID, privilegeId string, questionnaireId uuid.UUID) error
	GetRolePrivilegesOnInstance(ctx context.Context, roleId uuid.UUID) ([]model.RolePrivilegeOnInstance, error)
}

type rolePrivilegeOnInstanceRepository struct {
	db *gorm.DB
}

func NewRolePrivilegeOnInstanceRepository(db *gorm.DB) IRolePrivilegeOnInstanceRepository {
	return &rolePrivilegeOnInstanceRepository{db: db}
}

func (r *rolePrivilegeOnInstanceRepository) Add(ctx context.Context, rolePrivelegeOnInsance *model.RolePrivilegeOnInstance) error {
	err := r.db.WithContext(ctx).Create(rolePrivelegeOnInsance).Error
	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilegeOnInstanceRepository) Delete(ctx context.Context, roleId uuid.UUID, privilegeId string, questionnaireId uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&model.RolePrivilegeOnInstance{
		RoleId:          roleId,
		PrivilegeId:     privilegeId,
		QuestionnaireId: questionnaireId,
	}).Error

	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilegeOnInstanceRepository) GetRolePrivilegesOnInstance(ctx context.Context, roleId uuid.UUID) ([]model.RolePrivilegeOnInstance, error) {
	var rolePrivilegeOnInstance []model.RolePrivilegeOnInstance
	err := r.db.WithContext(ctx).Where("role_id = ?", roleId).Find(&rolePrivilegeOnInstance).Error
	if err != nil {
		//log
	}

	return rolePrivilegeOnInstance, err
}
