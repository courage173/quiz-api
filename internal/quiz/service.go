package quiz

import (
	"fmt"

	"github.com/courage173/quiz-api/internal/storage"

	"github.com/courage173/quiz-api/internal/models"

	"github.com/courage173/quiz-api/pkg/log"
)

type Service interface {
	SubmitQuiz(submission models.Submission) models.SubmissionResponse
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

func (s service) SubmitQuiz(submission models.Submission) models.SubmissionResponse {
	correctCount := 0

	for _, question := range submission.Answers {
		correctOption, err := s.storage.GetCorrectOption(question.QuestionID)
		if err != nil {
			s.logger.Error("Error getting correct option: ", err)
			continue
		}
		if correctOption.ID == question.OptionID {
			correctCount++
		}
	}
	result := &models.Result{UserName: submission.UserName, Score: correctCount}
	scoreRankPercentage := s.calculateScoreRankPercentage(result.Score)
	s.storage.AddUserSubmission(*result)

	message := fmt.Sprintf("You were better than %.2f%% of all quizzers", scoreRankPercentage)

	return models.SubmissionResponse{
		Message:            message,
		Score:              result.Score,
		Rank:               scoreRankPercentage,
		TotalQuestionCount: len(submission.Answers),
	}
}

func (s service) calculateScoreRankPercentage(score int) float64 {
	submissions := s.storage.GetSubmissions()
	totalSubmissions := len(submissions)

	fmt.Println(submissions)

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
