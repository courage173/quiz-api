package storage_test

import (
	"fmt"
	"testing"

	"github.com/courage173/quiz-api/internal/models"
	"github.com/courage173/quiz-api/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_AddUserSubmission(t *testing.T) {
	store := storage.NewStorage()

	// Test adding a new user submission
	submission := models.Result{
		UserName:              "testUser",
		Score:                 8,
		TotalQuestionAnswered: 10,
	}
	store.AddUserSubmission(submission)

	// Check if the submission was added correctly
	retrievedSubmission, err := store.GetUserSubmission("testUser")
	assert.NoError(t, err)
	assert.Equal(t, submission, retrievedSubmission)
}

func TestMemoryStorage_GetUserSubmission(t *testing.T) {
	store := storage.NewStorage()

	// Test retrieving an existing submission
	submission := models.Result{
		UserName:              "testUser",
		Score:                 8,
		TotalQuestionAnswered: 10,
	}
	store.AddUserSubmission(submission)

	retrievedSubmission, err := store.GetUserSubmission("testUser")
	assert.NoError(t, err)
	assert.Equal(t, submission, retrievedSubmission)

	// Test retrieving a non-existent submission
	_, err = store.GetUserSubmission("nonExistentUser")
	assert.Error(t, err)
	assert.EqualError(t, err, "submission not found")
}

func TestMemoryStorage_CalculateScoreRankPercentage(t *testing.T) {
	store := storage.NewStorage()

	// Add multiple submissions with varying scores
	store.AddUserSubmission(models.Result{UserName: "User1", Score: 5})
	store.AddUserSubmission(models.Result{UserName: "User2", Score: 7})
	store.AddUserSubmission(models.Result{UserName: "User3", Score: 8})
	store.AddUserSubmission(models.Result{UserName: "User4", Score: 10})

	//calculate for user 3
	percentage := store.CalculateScoreRankPercentage(8)
	formattedValue := fmt.Sprintf("%.2f", percentage)
	assert.Equal(t, "66.67", formattedValue)
}

func TestMemoryStorage_GetCorrectOption(t *testing.T) {
	store := storage.NewStorage()

	option, err := store.GetCorrectOption(1)
	assert.NoError(t, err)
	assert.Equal(t, models.Option{ID: 3, Text: "Paris", IsCorrect: true}, option)

	_, err = store.GetCorrectOption(999)
	assert.Error(t, err)
	assert.EqualError(t, err, "question not found")
}

func TestMemoryStorage_GetQuestions(t *testing.T) {
	store := storage.NewStorage()

	// Test that all questions are returned
	questions := store.GetQuestions()
	assert.Equal(t, 10, len(questions))

	assert.Equal(t, 1, questions[0].ID)
	assert.Equal(t, "What is the capital of France?", questions[0].Text)
}

func TestMemoryStorage_Count(t *testing.T) {
	store := storage.NewStorage()

	count := store.Count()
	assert.Equal(t, 10, count)
}
