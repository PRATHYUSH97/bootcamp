package cmd

import (
	"fmt"
	"os"

	"task/db"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Lists all tasks you completed today",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.Completedtoday()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have not completed any task today")
			return
		}
		fmt.Println("You have completed the following tasks today")
		for _, task := range tasks {
			fmt.Printf("%s\n", task)
		}
	},
}

func init() {
	RootCmd.AddCommand(completeCmd)
}
