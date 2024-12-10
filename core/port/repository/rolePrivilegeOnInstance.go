package repository

import (
	"context"
	"fmt"
	myContext "golizilla/adapters/http/handler/context"
	"golizilla/adapters/persistence/logger"
	"golizilla/core/domain/model"
	"golizilla/internal/logmessages"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRolePrivilegeOnInstanceRepository interface {
	Add(ctx context.Context, userCtx context.Context, rolePrivelegeOnInsance *model.RolePrivilegeOnInstance) error
	Delete(ctx context.Context, userCtx context.Context, roleId uuid.UUID, privilegeId string, questionnaireId uuid.UUID) error
	GetRolePrivilegesOnInstance(ctx context.Context, userCtx context.Context, roleId uuid.UUID) ([]model.RolePrivilegeOnInstance, error)
	HasPrivilegesOnInsance(ctx context.Context, userCtx context.Context, roleId uuid.UUID, questionnariId uuid.UUID, privileges ...string) (bool, error)
}

type rolePrivilegeOnInstanceRepository struct {
	db *gorm.DB
}

func NewRolePrivilegeOnInstanceRepository(db *gorm.DB) IRolePrivilegeOnInstanceRepository {
	return &rolePrivilegeOnInstanceRepository{db: db}
}

func (r *rolePrivilegeOnInstanceRepository) Add(ctx context.Context, userCtx context.Context, rolePrivelegeOnInsance *model.RolePrivilegeOnInstance) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(rolePrivelegeOnInsance).Error
	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilegeOnInstanceRepository) Delete(ctx context.Context, userCtx context.Context, roleId uuid.UUID, privilegeId string, questionnaireId uuid.UUID) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Where("role_id = ? AND privilege_id = ? AND questionnaire_id =?", roleId, privilegeId, questionnaireId).Delete(&model.RolePrivilegeOnInstance{
		RoleId:          roleId,
		PrivilegeId:     privilegeId,
		QuestionnaireId: questionnaireId,
	}).Error

	if err != nil {
		//log
	}
	return err
}

func (r *rolePrivilegeOnInstanceRepository) GetRolePrivilegesOnInstance(ctx context.Context, userCtx context.Context, roleId uuid.UUID) ([]model.RolePrivilegeOnInstance, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var rolePrivilegeOnInstance []model.RolePrivilegeOnInstance
	err := db.WithContext(ctx).Where("role_id = ?", roleId).Find(&rolePrivilegeOnInstance).Error
	if err != nil {
		//log
	}

	return rolePrivilegeOnInstance, err
}

func (r *rolePrivilegeOnInstanceRepository) HasPrivilegesOnInsance(
	ctx context.Context,
	userCtx context.Context,
	roleId uuid.UUID,
	questionnariId uuid.UUID,
	privileges ...string) (bool, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var rolePrivilegeOnInstance []model.RolePrivilegeOnInstance
	// Query matching both role_id and privilege_id
	result := db.WithContext(ctx).Where("role_id = ? AND privilege_id IN (?) AND questionnaire_id = ?", roleId, privileges, questionnariId).Find(&rolePrivilegeOnInstance)

	if result.Error != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogRolePrivilegeOnInstance,
			Message: "role privilege on instance not found",
		})
		return false, result.Error
	}
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogRolePrivilegeOnInstance,
		Message: fmt.Sprintf("rows=%d role_id=%s privileges=%s q_id=%s", result.RowsAffected, roleId, privileges[0], questionnariId),
	})

	return result.RowsAffected > 0, result.Error
}
