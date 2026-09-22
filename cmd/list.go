package cmd

import (
	"fmt"

	"github.com/crab4/todo-learning-helper/internal/cli"
	"github.com/spf13/cobra"
)

var filter string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  `List tasks with optional filtering by completion status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tasks, err := store.Load()
		if err != nil {
			return fmt.Errorf("load tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks yet. Add one with: gotodo add <title>")
			return nil
		}

		printed := 0
		for _, t := range tasks {
			switch filter {
			case "done":
				if !t.Done {
					continue
				}
			case "pending":
				if t.Done {
					continue
				}
			case "all":
				// все
			default:
				return fmt.Errorf("invalid filter: %s", filter)
			}

			cli.PrintTask(t)
			printed++
		}

		if printed == 0 {
			fmt.Printf("No tasks match filter: %s\n", filter)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&filter, "filter", "all", "Filter: all, done, pending")
}
