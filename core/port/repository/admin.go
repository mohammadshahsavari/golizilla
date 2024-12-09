package repository

import (
	"context"
	"golizilla/core/domain/model"
	myContext "golizilla/adapters/http/handler/context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAdminRepository interface {
	GetAllUsers(ctx, userCtx context.Context, page, pageSize int) ([]model.User, int64, error)
	GetAllQuestions(ctx, userCtx context.Context, page, pageSize int) ([]model.Question, int64, error)     // may move to question module
	GetAllQuestionnaires(ctx, userCtx context.Context, page, pageSize int) ([]model.Questionnaire, int64, error) // may move to Questionnares module
	GetAllRoles(ctx, userCtx context.Context, page, pageSize int) ([]model.Role, int64, error)
	GetAnswersByUserIDAndQuestionnaireID(ctx, userCtx context.Context, userID, questionnaireID uuid.UUID, page, pageSize int) ([]model.Answer, int64, error)
	// GivePermissionToUserByID
	// DeleteUserPermissionsByID
}

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) IAdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAllUsers(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) ([]model.User, int64, error) {

	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	var users []model.User
	var totalRecords int64

	// Count total records for pagination info
	if err := db.WithContext(ctx).Model(&model.User{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	err := db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users).Error
	return users, totalRecords, err
}

func (r *AdminRepository) GetAllQuestions(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) ([]model.Question, int64, error) {

	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	var questions []model.Question
	var totalRecords int64

	// Count total records for pagination info
	if err := db.WithContext(ctx).Model(&model.Question{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	err := db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&questions).Error
	return questions, totalRecords, err
}

func (r *AdminRepository) GetAllQuestionnaires(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) ([]model.Questionnaire, int64, error) {

	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	var Questionnaires []model.Questionnaire
	var totalRecords int64

	// Count total records for pagination info
	if err := db.WithContext(ctx).Model(&model.Questionnaire{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	err := db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&Questionnaires).Error
	return Questionnaires, totalRecords, err
}

func (r *AdminRepository) GetAllRoles(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) ([]model.Role, int64, error) {

	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	var roles []model.Role
	var totalRecords int64

	// Count total records for pagination info
	if err := db.WithContext(ctx).Model(&model.Role{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	err := db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&roles).Error
	return roles, totalRecords, err
}

func (r *AdminRepository) GetAnswersByUserIDAndQuestionnaireID(
	ctx, userCtx context.Context, userID, questionnaireID uuid.UUID, page, pageSize int,
) ([]model.Answer, int64, error) {

	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	var answers []model.Answer
	var totalRecords int64

    // Perform the query to get answers by UserID and QuestionnaireID
	// Count total records for pagination info
    err := db.WithContext(ctx).Preload("Questionnaire").
		Joins("JOIN questions ON answers.question_id = questions.id").
        Joins("JOIN questionnaires ON questions.questionnaire_id = questionnaires.id").
        Where("answers.user_id = ? AND questionnaires.id = ?", userID, questionnaireID).
        Count(&totalRecords).Error
    if err != nil {
        return nil, 0, err
    }

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	err = db.WithContext(ctx).Preload("Questionnaire").
		Joins("JOIN questions ON answers.question_id = questions.id").
        Joins("JOIN questionnaires ON questions.questionnaire_id = questionnaires.id").
        Where("answers.user_id = ? AND questionnaires.id = ?", userID, questionnaireID).
		Offset(offset).Limit(pageSize).Find(&answers).Error
	return answers, totalRecords, err
}
