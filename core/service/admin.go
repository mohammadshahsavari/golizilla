package service

import (
	"context"
	"golizilla/core/domain/model"
	"golizilla/core/port/repository"

	"github.com/google/uuid"
)

type IAdminService interface {
	GetAllUsers(ctx, userCtx context.Context, page int, pageSize int) (PaginatedUsers, error)
	GetAllQuestions(ctx, userCtx context.Context, page int, pageSize int) (PaginatedQuestions, error)
	GetAllQuestionnaires(ctx, userCtx context.Context, page int, pageSize int) (PaginatedQuestionnaires, error)
	GetAllRoles(ctx, userCtx context.Context, page int, pageSize int) (PaginatedRoles, error)
	GetAnswersByUserIDAndQuestionnaireID(ctx, userCtx context.Context, userID, questionnaireID uuid.UUID, page, pageSize int) (PaginatedUserQuestionnaireAnswer, error)
}

type AdminService struct {
	adminRepo repository.IAdminRepository
}

func NewAdminService(adminRepo repository.IAdminRepository) IAdminService {
	return &AdminService{adminRepo: adminRepo}
}

type PaginatedUsers struct {
	Data  []model.User `json:"data"`
	Pages int          `json:"pages"`
	Page  int          `json:"page"`
}

func (s *AdminService) GetAllUsers(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) (PaginatedUsers, error) {

	users, totalRecords, err := s.adminRepo.GetAllUsers(ctx, userCtx, page, pageSize)
	if err != nil {
		return PaginatedUsers{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the output struct
	result := PaginatedUsers{
		Data:  users,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}

type PaginatedQuestions struct {
	Data  []model.Question `json:"data"`
	Pages int              `json:"pages"`
	Page  int              `json:"page"`
}

func (s *AdminService) GetAllQuestions(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) (PaginatedQuestions, error) {

	questions, totalRecords, err := s.adminRepo.GetAllQuestions(ctx, userCtx, page, pageSize)
	if err != nil {
		return PaginatedQuestions{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the output struct
	result := PaginatedQuestions{
		Data:  questions,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}

type PaginatedQuestionnaires struct {
	Data  []model.Questionnaire `json:"data"`
	Pages int                   `json:"pages"`
	Page  int                   `json:"page"`
}

func (s *AdminService) GetAllQuestionnaires(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) (PaginatedQuestionnaires, error) {

	questionnaires, totalRecords, err := s.adminRepo.GetAllQuestionnaires(ctx, userCtx, page, pageSize)
	if err != nil {
		return PaginatedQuestionnaires{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the output struct
	result := PaginatedQuestionnaires{
		Data:  questionnaires,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}

type PaginatedRoles struct {
	Data  []model.Role `json:"data"`
	Pages int          `json:"pages"`
	Page  int          `json:"page"`
}

func (s *AdminService) GetAllRoles(
	ctx context.Context,
	userCtx context.Context,
	page int,
	pageSize int,
) (PaginatedRoles, error) {

	roles, totalRecords, err := s.adminRepo.GetAllRoles(ctx, userCtx, page, pageSize)
	if err != nil {
		return PaginatedRoles{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the output struct
	result := PaginatedRoles{
		Data:  roles,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}

type PaginatedUserQuestionnaireAnswer struct {
	Data  []model.Answer `json:"data"`
	Pages int            `json:"pages"`
	Page  int            `json:"page"`
}

func (s *AdminService) GetAnswersByUserIDAndQuestionnaireID(
	ctx, userCtx context.Context, userID, questionnaireID uuid.UUID, page, pageSize int,
) (PaginatedUserQuestionnaireAnswer, error) {

	answers, totalRecords, err := s.adminRepo.GetAnswersByUserIDAndQuestionnaireID(
		ctx, userCtx, userID, questionnaireID, page, pageSize)
	if err != nil {
		return PaginatedUserQuestionnaireAnswer{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the output struct
	result := PaginatedUserQuestionnaireAnswer{
		Data:  answers,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}