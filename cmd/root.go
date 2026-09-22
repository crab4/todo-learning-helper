package cmd

import (
	"fmt"
	"os"

	"github.com/crab4/todo-learning-helper/internal/storage"
	"github.com/spf13/cobra"
)

var taskFile string
var store *storage.Storage

var rootCmd = &cobra.Command{
	Use:   "gotodo",
	Short: "Terminal Task Manager",
	Long: `gotodo is a CLI task manager that helps u organize u tasks

	Tasks are storen in s JSON file at ~/.gotodo/tasks.json by default.	`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path := taskFile
		if path == "" {
			var err error
			path, err = storage.DefaultPath()
			if err != nil {
				return fmt.Errorf("determine default path: %w", err)
			}
		}
		store = storage.New(path)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&taskFile, "file", "", "task file path(default: ~/.gotodo/tasks.json)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
