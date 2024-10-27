package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/courage173/quiz-api/internal/models"
	"github.com/spf13/cobra"
	"net/http"
)

// Define answers for three different users
var userAnswers = map[string]models.Submission{
	"user1": {
		UserName: "user1",
		Answers: []models.Answer{
			{QuestionID: 1, OptionID: 2},
			{QuestionID: 2, OptionID: 3},
			{QuestionID: 3, OptionID: 4},
			{QuestionID: 4, OptionID: 2},
			{QuestionID: 5, OptionID: 1},
			{QuestionID: 6, OptionID: 3},
			{QuestionID: 7, OptionID: 2},
			{QuestionID: 8, OptionID: 2},
			{QuestionID: 9, OptionID: 3},
			{QuestionID: 10, OptionID: 1},
		},
	},
	"user2": {
		UserName: "user2",
		Answers: []models.Answer{
			{QuestionID: 1, OptionID: 3},
			{QuestionID: 2, OptionID: 1},
			{QuestionID: 3, OptionID: 2},
			{QuestionID: 4, OptionID: 3},
			{QuestionID: 5, OptionID: 4},
			{QuestionID: 6, OptionID: 2},
			{QuestionID: 7, OptionID: 3},
			{QuestionID: 8, OptionID: 1},
			{QuestionID: 9, OptionID: 4},
			{QuestionID: 10, OptionID: 2},
		},
	},
	"user3": {
		UserName: "user3",
		Answers: []models.Answer{
			{QuestionID: 1, OptionID: 1},
			{QuestionID: 2, OptionID: 2},
			{QuestionID: 3, OptionID: 3},
			{QuestionID: 4, OptionID: 4},
			{QuestionID: 5, OptionID: 1},
			{QuestionID: 6, OptionID: 3},
			{QuestionID: 7, OptionID: 1},
			{QuestionID: 8, OptionID: 4},
			{QuestionID: 9, OptionID: 2},
			{QuestionID: 10, OptionID: 3},
		},
	},
}

var selectedUser string

var answerCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your answers",
	Run: func(cmd *cobra.Command, args []string) {
		// Use selectedUser to get the specific user's answers
		answerRequest, exists := userAnswers[selectedUser]
		if !exists {
			fmt.Println("User not found. Please specify a valid user with --user user1 flag.")
			return
		}

		// Convert request to JSON
		data, err := json.Marshal(answerRequest)
		if err != nil {
			fmt.Println("Error encoding answers:", err)
			return
		}

		// Send POST request to the API
		resp, err := http.Post("http://localhost:4000/v1/quiz/submit", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error submitting answers:", err)
			return
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Println("Result:", result)
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
	// Add user flag to specify which user's answers to submit
	answerCmd.Flags().StringVarP(&selectedUser, "user", "u", "", "Specify the user (user1, user2, user3)")
}
