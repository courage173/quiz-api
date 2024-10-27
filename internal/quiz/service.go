package quiz

import (
	"fmt"

	"github.com/courage173/quiz-api/internal/storage"

	"github.com/courage173/quiz-api/internal/models"

	"github.com/courage173/quiz-api/pkg/log"
)

type Service interface {
	SubmitQuiz(submission models.Submission) (models.SubmissionResponse, error)
	GetQuestions() []models.Question
	GetUserSubmission(userName string) (models.GetSubmissionResponse, error)
}

type service struct {
	storage storage.Storage
	logger  log.Logger
}

func NewService(storage storage.Storage, logger log.Logger) Service {
	return service{
		storage,
		logger,
	}
}

func (s service) SubmitQuiz(submission models.Submission) (models.SubmissionResponse, error) {
	correctCount := 0

	totalQuestions := s.storage.Count()
	fmt.Println("totalQuestions: ", totalQuestions)

	if totalQuestions != len(submission.Answers) {
		return models.SubmissionResponse{}, fmt.Errorf("invalid submission: expected %d answers, got %d", totalQuestions, len(submission.Answers))
	}

	for _, question := range submission.Answers {
		correctOption, err := s.storage.GetCorrectOption(question.QuestionID)
		if err != nil {
			s.logger.Error("Error getting correct option: ", err)
			return models.SubmissionResponse{}, err
		}
		if correctOption.ID == question.OptionID {
			correctCount++
		}
	}

	result := &models.Result{UserName: submission.UserName, Score: correctCount, TotalQuestionAnswered: totalQuestions}

	s.storage.AddUserSubmission(*result)

	message := fmt.Sprintf("You got %d questions out of %d", result.Score, totalQuestions)

	return models.SubmissionResponse{
		Message:               message,
		Score:                 result.Score,
		TotalQuestionAnswered: totalQuestions,
	}, nil
}

func (s service) GetUserSubmission(userName string) (models.GetSubmissionResponse, error) {
	submission, err := s.storage.GetUserSubmission(userName)

	if err != nil {
		return models.GetSubmissionResponse{}, err
	}

	scoreRankPercentage := s.storage.CalculateScoreRankPercentage(submission.Score)
	formattedValue := fmt.Sprintf("%.2f%%", scoreRankPercentage)

	message := fmt.Sprintf("You were better than %s of all quizzers", formattedValue)

	return models.GetSubmissionResponse{
		Message:               message,
		Score:                 submission.Score,
		Rank:                  formattedValue,
		TotalQuestionAnswered: submission.TotalQuestionAnswered,
	}, nil
}

func (s service) GetQuestions() []models.Question {
	return s.storage.GetQuestions()
}
