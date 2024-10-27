package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get quiz questions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:4000/v1/quiz")
		if err != nil {
			fmt.Println("Error fetching questions:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println("Questions:", string(body))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
