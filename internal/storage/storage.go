package storage

import (
	"errors"
	"sync"

	"github.com/courage173/quiz-api/internal/models"
)

// Storage interface defines accessible methods for storage operations
type Storage interface {
	GetQuestions() []models.Question
	GetSubmissions() []models.Result
	AddUserSubmission(submission models.Result)
	GetCorrectOption(questionID int) (models.Option, error)
}

type memoryStorage struct {
	Questions   []models.Question
	Submissions []models.Result
	Mutex       sync.RWMutex
}

func NewStorage() Storage {
	return &memoryStorage{
		Questions: []models.Question{
			{
				ID:   1,
				Text: "What is the capital of France?",
				Options: []models.Option{
					{ID: 1, Text: "Berlin", IsCorrect: false},
					{ID: 2, Text: "London", IsCorrect: false},
					{ID: 3, Text: "Paris", IsCorrect: true},
					{ID: 4, Text: "Madrid", IsCorrect: false},
				},
			},
			{
				ID:   2,
				Text: "Which planet is known as the Red Planet?",
				Options: []models.Option{
					{ID: 1, Text: "Earth", IsCorrect: false},
					{ID: 2, Text: "Mars", IsCorrect: true},
					{ID: 3, Text: "Jupiter", IsCorrect: false},
					{ID: 4, Text: "Saturn", IsCorrect: false},
				},
			},
			// Add more questions as needed
		},
		Submissions: make([]models.Result, 0),
	}
}

// AddUserSubmission adds a user's submission to storage
func (s *memoryStorage) AddUserSubmission(submissionResult models.Result) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Submissions = append(s.Submissions, submissionResult)

}

// GetQuestion retrieves a question by its ID
func (s *memoryStorage) GetQuestion(id int) (models.Question, error) {
	for _, question := range s.Questions {
		if question.ID == id {
			return question, nil
		}
	}
	return models.Question{}, errors.New("question not found")
}

// GetCorrectOption retrieves the correct option for a given question
func (s *memoryStorage) GetCorrectOption(questionId int) (models.Option, error) {
	question, err := s.GetQuestion(questionId)
	if err != nil {
		return models.Option{}, err
	}
	for _, option := range question.Options {
		if option.IsCorrect {
			return option, nil
		}
	}
	return models.Option{}, errors.New("option not found")
}

func (s *memoryStorage) GetQuestions() []models.Question {
	return s.Questions
}

func (s *memoryStorage) GetSubmissions() []models.Result {
	return s.Submissions
}
