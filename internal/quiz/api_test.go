package quiz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/courage173/quiz-api/internal/models"
	"github.com/courage173/quiz-api/internal/quiz"
	"github.com/courage173/quiz-api/pkg/log"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetQuestions() []models.Question {
	args := m.Called()
	return args.Get(0).([]models.Question)
}

func (m *MockService) SubmitQuiz(submission models.Submission) (models.SubmissionResponse, error) {
	args := m.Called(submission)
	return args.Get(0).(models.SubmissionResponse), args.Error(1)
}

func (m *MockService) GetUserSubmission(userName string) (models.GetSubmissionResponse, error) {
	args := m.Called(userName)
	return args.Get(0).(models.GetSubmissionResponse), args.Error(1)
}

func setupRouter(service *MockService, logger log.Logger) *routing.Router {
	router := routing.New()
	rg := router.Group("/v1")

	quiz.RegisterHandlers(rg.Group("/quiz"), service, logger)
	return router
}

func TestGetQuizHandler(t *testing.T) {
	mockService := new(MockService)
	logger, _ := log.NewForTest()
	router := setupRouter(mockService, logger)

	// Mock the GetQuestions method
	mockService.On("GetQuestions").Return([]models.Question{
		{ID: 1, Text: "Question 1", Options: []models.Option{{ID: 1, Text: "Option 1"}}},
	})

	req := httptest.NewRequest("GET", "/v1/quiz", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestSubmitQuizHandler(t *testing.T) {
	mockService := new(MockService)
	logger, _ := log.NewForTest()
	router := setupRouter(mockService, logger)

	validSubmission := models.Submission{
		UserName: "testUser",
		Answers: []models.Answer{
			{QuestionID: 1, OptionID: 2},
		},
	}

	mockService.On("SubmitQuiz", validSubmission).Return(models.SubmissionResponse{
		Message:               "You got 1 questions out of 1",
		Score:                 1,
		TotalQuestionAnswered: 1,
	}, nil)

	body, _ := json.Marshal(validSubmission)
	req := httptest.NewRequest("POST", "/v1/quiz/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserSubmissionHandler(t *testing.T) {
	mockService := new(MockService)
	logger, _ := log.NewForTest()
	router := setupRouter(mockService, logger)

	mockService.On("GetUserSubmission", "testUser").Return(models.GetSubmissionResponse{
		Message:               "You were better than 75.00% of all quizzers",
		Score:                 8,
		Rank:                  "75.00%",
		TotalQuestionAnswered: 10,
	}, nil)

	req := httptest.NewRequest("GET", "/v1/quiz/submission/testUser", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
