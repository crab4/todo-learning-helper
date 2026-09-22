package cmd

import (
	"fmt"
	"strconv"

	"github.com/crab4/todo-learning-helper/internal/task"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <id>",
	Short: "Remove a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}

		tasks, err := store.Load()
		if err != nil {
			return fmt.Errorf("load tasks: %w", err)
		}

		filtered := make([]*task.Task, 0, len(tasks))
		found := false
		for _, t := range tasks {
			if t.ID == id {
				found = true
				continue
			}
			filtered = append(filtered, t)
		}
		if !found {
			return fmt.Errorf("task #%d not found", id)
		}

		if err := store.Save(filtered); err != nil {
			return fmt.Errorf("save tasks: %w", err)
		}
		fmt.Printf("Removed task #%d\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
