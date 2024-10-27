package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var username string

// getSubmissionCmd represents the command to get a user's quiz submission
var getSubmissionCmd = &cobra.Command{
	Use:   "score",
	Short: "Retrieve a user's quiz submission details",
	RunE: func(cmd *cobra.Command, args []string) error {
		if username == "" {
			return fmt.Errorf("username is required")
		}

		url := fmt.Sprintf("http://localhost:4000/v1/quiz/submission/%s", username)
		fmt.Printf("Fetching submission for user: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching submission: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to retrieve submission, status: %s", resp.Status)
		}

		var submissionResponse struct {
			Message               string `json:"message"`
			Score                 int    `json:"score"`
			Rank                  string `json:"rank"`
			TotalQuestionAnswered int    `json:"total_question_answered"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&submissionResponse); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		fmt.Printf("User Submission for %s:\n", username)
		fmt.Printf("Message: %s\n", submissionResponse.Message)
		fmt.Printf("Score: %d\n", submissionResponse.Score)
		fmt.Printf("Rank: %s\n", submissionResponse.Rank)
		fmt.Printf("Total Questions Answered: %d\n", submissionResponse.TotalQuestionAnswered)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getSubmissionCmd)

	// Define flags
	getSubmissionCmd.Flags().StringVarP(&username, "user", "u", "", "Username for the submission (required)")
	getSubmissionCmd.MarkFlagRequired("username")
}
