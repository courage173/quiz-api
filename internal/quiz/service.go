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

	result := &models.Result{UserName: submission.UserName, Score: correctCount}
	scoreRankPercentage := s.calculateScoreRankPercentage(result.Score)
	s.storage.AddUserSubmission(*result)

	formattedValue := fmt.Sprintf("%.2f%%", scoreRankPercentage)

	message := fmt.Sprintf("You were better than %s of all quizzers", formattedValue)

	return models.SubmissionResponse{
		Message:            message,
		Score:              result.Score,
		Rank:               formattedValue,
		TotalQuestionCount: len(submission.Answers),
	}, nil
}

func (s service) calculateScoreRankPercentage(score int) float64 {
	if score == 0 {
		return 0.0
	}

	submissions := s.storage.GetSubmissions()
	totalSubmissions := len(submissions)

	if totalSubmissions == 0 {
		return 100.0
	}

	lowerScoreCount := 0
	for _, submission := range submissions {
		if submission.Score < score {
			lowerScoreCount++
		}
	}

	percentageBetter := (float64(lowerScoreCount) / float64(totalSubmissions)) * 100.0
	return percentageBetter
}

func (s service) GetQuestions() []models.Question {
	return s.storage.GetQuestions()
}
