package cmd

import (
	"fmt"

	"github.com/crab4/todo-learning-helper/internal/task"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.Save([]*task.Task{}); err != nil {
			return fmt.Errorf("clear tasks: %w", err)
		}
		fmt.Println("All tasks cleared")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
