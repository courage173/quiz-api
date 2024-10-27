package storage

import (
	"errors"

	"sort"

	"sync"

	"github.com/courage173/quiz-api/internal/models"
)

type Storage interface {
	GetQuestions() []models.Question
	GetUserSubmission(userName string) (models.Result, error)
	CalculateScoreRankPercentage(score int) float64
	AddUserSubmission(submission models.Result)
	GetCorrectOption(questionID int) (models.Option, error)
	Count() int
}

type memoryStorage struct {
	Questions    map[int]models.Question
	Submissions  map[string]models.Result
	ScoreTracker map[int]int
	Mutex        sync.RWMutex
}

func NewStorage() Storage {
	return &memoryStorage{
		Questions: map[int]models.Question{
			1: {
				ID:   1,
				Text: "What is the capital of France?",
				Options: []models.Option{
					{ID: 1, Text: "Berlin", IsCorrect: false},
					{ID: 2, Text: "London", IsCorrect: false},
					{ID: 3, Text: "Paris", IsCorrect: true},
					{ID: 4, Text: "Madrid", IsCorrect: false},
				},
			},
			2: {
				ID:   2,
				Text: "Which planet is known as the Red Planet?",
				Options: []models.Option{
					{ID: 1, Text: "Earth", IsCorrect: false},
					{ID: 2, Text: "Mars", IsCorrect: true},
					{ID: 3, Text: "Jupiter", IsCorrect: false},
					{ID: 4, Text: "Saturn", IsCorrect: false},
				},
			},
			3: {
				ID:   3,
				Text: "What is the largest ocean on Earth?",
				Options: []models.Option{
					{ID: 1, Text: "Atlantic Ocean", IsCorrect: false},
					{ID: 2, Text: "Indian Ocean", IsCorrect: false},
					{ID: 3, Text: "Arctic Ocean", IsCorrect: false},
					{ID: 4, Text: "Pacific Ocean", IsCorrect: true},
				},
			},
			4: {
				ID:   4,
				Text: "Which country is known as the Land of the Rising Sun?",
				Options: []models.Option{
					{ID: 1, Text: "China", IsCorrect: false},
					{ID: 2, Text: "Japan", IsCorrect: true},
					{ID: 3, Text: "Thailand", IsCorrect: false},
					{ID: 4, Text: "South Korea", IsCorrect: false},
				},
			},
			5: {
				ID:   5,
				Text: "Who wrote 'Romeo and Juliet'?",
				Options: []models.Option{
					{ID: 1, Text: "William Shakespeare", IsCorrect: true},
					{ID: 2, Text: "Charles Dickens", IsCorrect: false},
					{ID: 3, Text: "Leo Tolstoy", IsCorrect: false},
					{ID: 4, Text: "Jane Austen", IsCorrect: false},
				},
			},
			6: {
				ID:   6,
				Text: "What is the smallest prime number?",
				Options: []models.Option{
					{ID: 1, Text: "0", IsCorrect: false},
					{ID: 2, Text: "1", IsCorrect: false},
					{ID: 3, Text: "2", IsCorrect: true},
					{ID: 4, Text: "3", IsCorrect: false},
				},
			},
			7: {
				ID:   7,
				Text: "What is the powerhouse of the cell?",
				Options: []models.Option{
					{ID: 1, Text: "Nucleus", IsCorrect: false},
					{ID: 2, Text: "Mitochondria", IsCorrect: true},
					{ID: 3, Text: "Ribosome", IsCorrect: false},
					{ID: 4, Text: "Chloroplast", IsCorrect: false},
				},
			},
			8: {
				ID:   8,
				Text: "What is the boiling point of water at sea level in Celsius?",
				Options: []models.Option{
					{ID: 1, Text: "50°C", IsCorrect: false},
					{ID: 2, Text: "100°C", IsCorrect: true},
					{ID: 3, Text: "75°C", IsCorrect: false},
					{ID: 4, Text: "125°C", IsCorrect: false},
				},
			},
			9: {
				ID:   9,
				Text: "Which planet has the most moons?",
				Options: []models.Option{
					{ID: 1, Text: "Earth", IsCorrect: false},
					{ID: 2, Text: "Mars", IsCorrect: false},
					{ID: 3, Text: "Jupiter", IsCorrect: true},
					{ID: 4, Text: "Venus", IsCorrect: false},
				},
			},
			10: {
				ID:   10,
				Text: "What is the chemical symbol for gold?",
				Options: []models.Option{
					{ID: 1, Text: "Au", IsCorrect: true},
					{ID: 2, Text: "Ag", IsCorrect: false},
					{ID: 3, Text: "Gd", IsCorrect: false},
					{ID: 4, Text: "Go", IsCorrect: false},
				},
			},
		},
		Submissions:  make(map[string]models.Result),
		ScoreTracker: make(map[int]int),
	}
}

func (s *memoryStorage) AddUserSubmission(submissionResult models.Result) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Submissions[submissionResult.UserName] = submissionResult
	s.ScoreTracker[submissionResult.Score]++
}

func (s *memoryStorage) GetUserSubmission(userName string) (models.Result, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	submission, exists := s.Submissions[userName]
	if !exists {
		return models.Result{}, errors.New("submission not found")
	}
	return submission, nil
}

func (s *memoryStorage) CalculateScoreRankPercentage(score int) float64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	if score == 0 {
		return 0.0
	}

	// subtracting 1 to exclude the current submission
	totalScores := len(s.Submissions) - 1

	if totalScores == 0 {
		return 100.0
	}

	lowerScores := 0
	for s, count := range s.ScoreTracker {
		if s < score {
			lowerScores += count
		}
	}

	return float64(lowerScores) / float64(totalScores) * 100
}

func (s *memoryStorage) GetQuestion(id int) (models.Question, error) {
	question, exists := s.Questions[id]
	if !exists {
		return models.Question{}, errors.New("question not found")
	}
	return question, nil
}

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
	return models.Option{}, errors.New("correct option not found")
}

func (s *memoryStorage) GetQuestions() []models.Question {

	keys := make([]int, 0, len(s.Questions))
	for key := range s.Questions {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	questions := make([]models.Question, 0, len(s.Questions))
	for _, key := range keys {
		questions = append(questions, s.Questions[key])
	}
	return questions
}

func (s *memoryStorage) Count() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return len(s.Questions)
}
