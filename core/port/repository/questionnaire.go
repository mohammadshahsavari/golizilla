package repository

import (
	"context"
	myContext "golizilla/adapters/http/handler/context"
	"golizilla/adapters/persistence/logger"
	"golizilla/core/domain/model"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuestionnaireRepository interface {
	Add(ctx context.Context, userCtx context.Context, questionnaire *model.Questionnaire) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	Update(ctx context.Context, userCtx context.Context, id uuid.UUID, questionnaire map[string]interface{}) error
	GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Questionnaire, error)
	GetByOwnerId(ctx context.Context, userCtx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error)
	IsOwner(ctx context.Context, userCtx context.Context, userId uuid.UUID, questionnariId uuid.UUID) (bool, error)
}

type questionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) IQuestionnaireRepository {
	return &questionnaireRepository{
		db: db,
	}
}

func (r *questionnaireRepository) Add(ctx context.Context, userCtx context.Context, questionnaire *model.Questionnaire) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(questionnaire).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Delete(&model.Questionnaire{}, id).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) Update(ctx context.Context, userCtx context.Context, id uuid.UUID, questionnaire map[string]interface{}) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Model(&model.Questionnaire{}).Where("id = ?", id).Updates(questionnaire).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *questionnaireRepository) GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Questionnaire, error) {
	// Retrieve the database instance from context
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	// Fetch the questionnaire
	var questionnaire model.Questionnaire
	result := db.WithContext(ctx).First(&questionnaire, "id = ?", id)

	// Check if the query was successful
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return a custom error if no record is found
			return nil, apperrors.ErrQuestionnaireNotFound
		}
		// Log and return other errors
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireRepository,
			Message: apperrors.ErrQuestionnaireNotFound.Error(),
		})
		return nil, result.Error
	}

	return &questionnaire, nil
}

func (r *questionnaireRepository) GetByOwnerId(ctx context.Context, userCtx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error) {
	// Retrieve the database instance from context
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	// Fetch questionnaires by owner ID
	var questionnaires []model.Questionnaire
	result := db.WithContext(ctx).Where("owner_id = ?", ownerId).Find(&questionnaires)

	// Check for errors during query execution
	if result.Error != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireRepository,
			Message: "error fetching questionnaires by owner ID",
		})
		return nil, result.Error
	}

	// Check if no records were found
	if len(questionnaires) == 0 {
		return nil, apperrors.ErrQuestionnaireNotFound
	}

	return questionnaires, nil
}

func (r *questionnaireRepository) IsOwner(ctx context.Context, userCtx context.Context, userId uuid.UUID, questionnariId uuid.UUID) (bool, error) {
	// Retrieve the database instance from context
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	// Check ownership
	var questionnaire model.Questionnaire
	result := db.WithContext(ctx).Where("owner_id = ? AND id = ?", userId, questionnariId).First(&questionnaire)

	// Handle errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, apperrors.ErrQuestionnaireNotFound
		}
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireRepository,
			Message: "error checking questionnaire ownership",
		})
		return false, result.Error
	}

	return true, nil
}
