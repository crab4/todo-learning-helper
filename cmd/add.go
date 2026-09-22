package cmd

import (
	"fmt"

	"github.com/crab4/todo-learning-helper/internal/task"
	"github.com/spf13/cobra"
)

var (
	priority string
	due      string
)

//Хоспадя, как же я заколебался переводить тебя на кобру, ты изначально был жутко нудной таской. а сейчас я тебя уже немного ненавижу(=

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new task",
	Long: `Add a new task with the given title.
Optional flags allow u to set priority and due date.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		t := task.New(0, title, task.Priority(priority), due)
		if err := t.Validate(); err != nil {
			return fmt.Errorf("Invalid task: %w", err)
		}

		tasks, err := store.Load()
		if err != nil {
			return fmt.Errorf("load tasks: %w", err)
		}

		maxID := 0
		for _, existing := range tasks {
			if existing.ID > maxID {
				maxID = existing.ID
			}
		}

		t.ID = maxID + 1
		tasks = append(tasks, t)

		if store.Save(tasks); err != nil {
			return fmt.Errorf("save tasks: %w", err)
		}

		fmt.Printf("Added task #%d: %s\n", t.ID, t.Title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Priority: low, medium, high")
	addCmd.Flags().StringVar(&due, "due", "", "Due date (YYYY-MM-DD)")
}
