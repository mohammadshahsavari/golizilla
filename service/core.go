package service

import (
	"context"
	"fmt"
	"golizilla/domain/model"
	"golizilla/domain/repository"
	"golizilla/internal/logmessages"
	"golizilla/persistence/logger"

	"github.com/google/uuid"
)

type ICoreService interface {
	Start(ctx context.Context, userCtx context.Context, userID, questionnaireID uuid.UUID) (uuid.UUID, *model.Question, error)
	Submit(ctx context.Context, userCtx context.Context, submissionID, questionID uuid.UUID, answer *model.Answer) error
	Back(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) (*model.Question, error)
	Next(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) (*model.Question, error)
	End(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) error
}

type CoreService struct {
	questionRepo      repository.IQuestionRepository
	submissionRepo    repository.ISubmissionRepository
	questionnaireRepo repository.IQuestionnaireRepository
	answerRepo        repository.IAnswerRepository
}

func NewCoreService(
	questionRepo repository.IQuestionRepository,
	submissionRepo repository.ISubmissionRepository,
	questionnaireRepo repository.IQuestionnaireRepository,
	answerRepo repository.IAnswerRepository,
) ICoreService {
	return &CoreService{
		questionRepo:      questionRepo,
		submissionRepo:    submissionRepo,
		questionnaireRepo: questionnaireRepo,
		answerRepo:        answerRepo,
	}
}

func (c *CoreService) Start(ctx context.Context, userCtx context.Context, userID, questionnaireID uuid.UUID) (uuid.UUID, *model.Question, error) {
	qn, err := c.questionnaireRepo.GetById(ctx, userCtx, questionnaireID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireService,
			Message: fmt.Sprintf("failed to get questionnaire: %v", err.Error()),
		})
		return uuid.Nil, nil, fmt.Errorf("failed to get questionnaire: %w", err)
	}
	if qn == nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireService,
			Message: "questionnaire not found",
		})
		return uuid.Nil, nil, fmt.Errorf("questionnaire not found")
	}

	// Create new submission
	submission := &model.UserSubmission{
		ID:              uuid.New(),
		UserId:          userID,
		QuestionnaireId: questionnaireID,
		Status:          model.SubmissionsStatusInProgress,
	}
	if err := c.submissionRepo.CreateSubmission(ctx, userCtx, submission); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireService,
			Message: fmt.Sprintf("failed to create submission: %v", err.Error()),
		})
		return uuid.Nil, nil, err
	}

	questions, err := c.getQuestionsForQuestionnaire(ctx, userCtx, questionnaireID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireService,
			Message: fmt.Sprintf("failed to get questions: %v", err.Error()),
		})
		return uuid.Nil, nil, err
	}
	if len(questions) == 0 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireService,
			Message: "no questions available",
		})
		return submission.ID, nil, fmt.Errorf("no questions available")
	}

	submission.CurrentQuestionIndex = 0
	if err := c.submissionRepo.UpdateSubmission(ctx, userCtx, submission); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireService,
			Message: fmt.Sprintf("failed to update submission: %v", err.Error()),
		})
		return submission.ID, nil, err
	}

	return submission.ID, questions[0], nil
}

func (c *CoreService) Submit(ctx context.Context, userCtx context.Context, submissionID, questionID uuid.UUID, answer *model.Answer) error {
	submission, err := c.submissionRepo.GetSubmissionByID(ctx, userCtx, submissionID)
	if err != nil {
		return err
	}

	if submission.Status != model.SubmissionsStatusInProgress {
		return fmt.Errorf("submission not in progress")
	}

	questions, err := c.getQuestionsForQuestionnaire(ctx, userCtx, submission.QuestionnaireId)
	if err != nil {
		return err
	}
	if submission.CurrentQuestionIndex < 0 || submission.CurrentQuestionIndex >= len(questions) {
		return fmt.Errorf("no current question")
	}

	currentQuestion := questions[submission.CurrentQuestionIndex]
	if currentQuestion.ID != questionID {
		return fmt.Errorf("question does not match current index")
	}

	answer.QuestionID = questionID
	_, err = c.answerRepo.Create(ctx, userCtx, answer)
	if err != nil {
		return err
	}

	return c.submissionRepo.UpdateSubmission(ctx, userCtx, submission)
}

func (c *CoreService) Back(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) (*model.Question, error) {
	submission, err := c.submissionRepo.GetSubmissionByID(ctx, userCtx, submissionID)
	if err != nil {
		return nil, err
	}

	// TODO: Check questionnaire settings if back is allowed. For now, we assume allowed.
	if submission.CurrentQuestionIndex > 0 {
		submission.CurrentQuestionIndex--
		if err := c.submissionRepo.UpdateSubmission(ctx, userCtx, submission); err != nil {
			return nil, err
		}

		questions, err := c.getQuestionsForQuestionnaire(ctx, userCtx, submission.QuestionnaireId)
		if err != nil {
			return nil, err
		}
		return questions[submission.CurrentQuestionIndex], nil
	}

	return nil, fmt.Errorf("cannot go back")
}

func (c *CoreService) Next(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) (*model.Question, error) {
	submission, err := c.submissionRepo.GetSubmissionByID(ctx, userCtx, submissionID)
	if err != nil {
		return nil, err
	}

	questions, err := c.getQuestionsForQuestionnaire(ctx, userCtx, submission.QuestionnaireId)
	if err != nil {
		return nil, err
	}

	if submission.CurrentQuestionIndex+1 < len(questions) {
		submission.CurrentQuestionIndex++
		if err := c.submissionRepo.UpdateSubmission(ctx, userCtx, submission); err != nil {
			return nil, err
		}
		return questions[submission.CurrentQuestionIndex], nil
	}

	return nil, fmt.Errorf("no more questions")
}

func (c *CoreService) End(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) error {
	submission, err := c.submissionRepo.GetSubmissionByID(ctx, userCtx, submissionID)
	if err != nil {
		return err
	}

	submission.Status = model.SubmissionsStatusDone
	return c.submissionRepo.UpdateSubmission(ctx, userCtx, submission)
}

// getQuestionsForQuestionnaire retrieves all questions for the given questionnaire from the repo
func (c *CoreService) getQuestionsForQuestionnaire(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error) {
	questions, err := c.questionRepo.GetByQuestionnaireID(ctx, userCtx, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	return questions, nil
}
