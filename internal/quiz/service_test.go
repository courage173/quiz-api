package quiz_test

import (
	"errors"
	"testing"

	"github.com/courage173/quiz-api/internal/models"
	"github.com/courage173/quiz-api/internal/quiz"
	"github.com/courage173/quiz-api/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the storage interface for testing
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Count() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockStorage) GetCorrectOption(questionID int) (models.Option, error) {
	args := m.Called(questionID)
	return args.Get(0).(models.Option), args.Error(1)
}

func (m *MockStorage) GetQuestions() []models.Question {
	args := m.Called()
	return args.Get(0).([]models.Question)
}

func (m *MockStorage) GetSubmissions() []models.Result {
	args := m.Called()
	return args.Get(0).([]models.Result)
}

func (m *MockStorage) AddUserSubmission(result models.Result) {
	m.Called(result)
}

func (m *MockStorage) GetUserSubmission(userName string) (models.Result, error) {
	args := m.Called(userName)
	return args.Get(0).(models.Result), args.Error(1)
}

func (m *MockStorage) CalculateScoreRankPercentage(score int) float64 {
	args := m.Called(score)
	return args.Get(0).(float64)
}

// MockLogger is a mock implementation of the logger interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func TestSubmitQuiz(t *testing.T) {
	t.Run("valid submission", func(t *testing.T) {
		mockStorage := new(MockStorage)
		logger, _ := log.NewForTest()
		service := quiz.NewService(mockStorage, logger)

		mockStorage.On("Count").Return(2)
		mockStorage.On("GetCorrectOption", 1).Return(models.Option{ID: 2}, nil)
		mockStorage.On("GetCorrectOption", 2).Return(models.Option{ID: 3}, nil)
		mockStorage.On("AddUserSubmission", mock.Anything).Return()

		submission := models.Submission{
			UserName: "Charlie",
			Answers: []models.Answer{
				{QuestionID: 1, OptionID: 2},
				{QuestionID: 2, OptionID: 3},
			},
		}

		response, err := service.SubmitQuiz(submission)

		assert.NoError(t, err)
		assert.Equal(t, "You got 2 questions out of 2", response.Message)
		assert.Equal(t, 2, response.Score)
		assert.Equal(t, 2, response.TotalQuestionAnswered)
	})

	t.Run("invalid submission with fewer answers", func(t *testing.T) {
		mockStorage := new(MockStorage)
		logger, _ := log.NewForTest()
		service := quiz.NewService(mockStorage, logger)

		mockStorage.On("Count").Return(3)

		submission := models.Submission{
			UserName: "Charlie",
			Answers: []models.Answer{
				{QuestionID: 1, OptionID: 2},
				{QuestionID: 2, OptionID: 3},
			},
		}

		_, err := service.SubmitQuiz(submission)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid submission: expected 3 answers, got 2")
	})

	t.Run("error fetching correct option", func(t *testing.T) {
		mockStorage := new(MockStorage)
		logger, _ := log.NewForTest()
		service := quiz.NewService(mockStorage, logger)

		mockStorage.On("Count").Return(2)
		mockStorage.On("GetCorrectOption", 1).Return(models.Option{ID: 2}, nil)
		mockStorage.On("GetCorrectOption", 2).Return(models.Option{}, errors.New("database error"))

		submission := models.Submission{
			UserName: "Charlie",
			Answers: []models.Answer{
				{QuestionID: 1, OptionID: 2},
				{QuestionID: 2, OptionID: 3},
			},
		}

		_, err := service.SubmitQuiz(submission)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestGetQuestions(t *testing.T) {
	mockStorage := new(MockStorage)
	logger, _ := log.NewForTest()
	service := quiz.NewService(mockStorage, logger)

	questions := []models.Question{
		{ID: 1, Text: "Question 1", Options: []models.Option{{ID: 1, Text: "Option 1"}}},
		{ID: 2, Text: "Question 2", Options: []models.Option{{ID: 2, Text: "Option 2"}}},
	}

	mockStorage.On("GetQuestions").Return(questions)

	result := service.GetQuestions()

	assert.Equal(t, questions, result)
	mockStorage.AssertCalled(t, "GetQuestions")
}

func TestGetUserSubmission(t *testing.T) {
	t.Run("successful retrieval of user submission", func(t *testing.T) {
		mockStorage := new(MockStorage)
		logger, _ := log.NewForTest()
		service := quiz.NewService(mockStorage, logger)

		mockStorage.On("GetUserSubmission", "Charlie").Return(models.Result{
			UserName:              "Charlie",
			Score:                 8,
			TotalQuestionAnswered: 10,
		}, nil)
		mockStorage.On("CalculateScoreRankPercentage", 8).Return(75.00)

		response, err := service.GetUserSubmission("Charlie")

		assert.NoError(t, err)
		assert.Equal(t, "You were better than 75.00% of all quizzers", response.Message)
		assert.Equal(t, 8, response.Score)
		assert.Equal(t, "75.00%", response.Rank)
		assert.Equal(t, 10, response.TotalQuestionAnswered)
	})

	t.Run("user submission not found", func(t *testing.T) {
		mockStorage := new(MockStorage)
		logger, _ := log.NewForTest()
		service := quiz.NewService(mockStorage, logger)

		mockStorage.On("GetUserSubmission", "UnknownUser").Return(models.Result{}, errors.New("user submission not found"))

		_, err := service.GetUserSubmission("UnknownUser")

		assert.Error(t, err)
		assert.EqualError(t, err, "user submission not found")
	})
}
