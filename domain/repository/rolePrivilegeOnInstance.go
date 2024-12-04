package repository

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRolePrivilegeOnInstanceRepository interface {
	Add(rolePrivelegeOnInsance *model.RolePrivilegeOnInstance) error
	Delete(roleId uuid.UUID, privilegeId string, questionnariId uuid.UUID) error
	GetRolePrivilegesOnInstance(roleId uuid.UUID) ([]model.RolePrivilegeOnInstance, error)
}

type rolePrivilegeOnInstanceRepository struct {
	db *gorm.DB
}

func NewRolePrivilegeOnInstanceRepository(db *gorm.DB) IRolePrivilegeOnInstanceRepository {
	return &rolePrivilegeOnInstanceRepository{db: db}
}

func (r *rolePrivilegeOnInstanceRepository) Add(rolePrivelegeOnInsance *model.RolePrivilegeOnInstance) error {
	err := r.db.Create(rolePrivelegeOnInsance).Error
	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilegeOnInstanceRepository) Delete(roleId uuid.UUID, privilegeId string, questionnariId uuid.UUID) error {
	err := r.db.Delete(&model.RolePrivilegeOnInstance{
		RoleId:          roleId,
		PrivilegeId:     privilegeId,
		QuestionnaireId: questionnariId,
	}).Error

	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilegeOnInstanceRepository) GetRolePrivilegesOnInstance(roleId uuid.UUID) ([]model.RolePrivilegeOnInstance, error) {
	var rolePrivilegeOnInstance []model.RolePrivilegeOnInstance
	err := r.db.Where("role_id = ?", roleId).Find(&rolePrivilegeOnInstance).Error
	if err != nil {
		//log
	}

	return rolePrivilegeOnInstance, err
}
