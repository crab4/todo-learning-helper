package cli

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/crab4/todo-learning-helper/internal/storage"
	"github.com/crab4/todo-learning-helper/internal/task"
	"github.com/fatih/color"
)

//Оставил файлик на память, да и в конце одна функция полезная

type App struct {
	storage *storage.Storage
}

func New(s *storage.Storage) *App {
	return &App{storage: s}
}
func (a *App) Run(args []string) error {
	if len(args) < 1 {
		a.usage()
		return nil
	}

	switch args[0] {
	case "add":
		return a.cmdAdd(args[1:])
	case "list":
		return a.cmdList(args[1:])
	case "done":
		return a.cmdDone(args[1:])
	case "rm":
		return a.cmdRm(args[1:])
	case "clear":
		return a.cmdClear(args[1:])
	case "help", "-h", "--help":
		a.usage()
		return nil
	default:
		return fmt.Errorf("unknown command:%s", args[0])
	}
}

func (a *App) usage() {
	fmt.Println("gotodo - terminal task manager")
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println("	gotodo add <title> [--priority low|medium|high] [--due YYYY-MM-DD]")
	fmt.Println("gotodo list [--filter all|done|pending]")
	fmt.Println("gotodo done <id>")
	fmt.Println("gotodo rm <id>")
	fmt.Println("gotodo clear")
}

func (a *App) cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	priority := fs.String("priority", "medium", "Priority: low, medium, high")
	due := fs.String("Due", "", "Due date (YYYY-MM-DD)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: gotodo add <title> [--priority ...] [--due ...]")
	}

	title := fs.Arg(0)
	t := task.New(0, title, task.Priority(*priority), *due)
	if err := t.Validate(); err != nil {
		return err
	}

	tasks, err := a.storage.Load()
	if err != nil {
		return err
	}
	maxId := 0
	for _, existing := range tasks {
		if existing.ID > maxId {
			maxId = existing.ID
		}
	}
	t.ID = maxId + 1

	tasks = append(tasks, t)
	if err := a.storage.Save(tasks); err != nil {
		return err
	}

	color.Green("Added task #%d: %s", t.ID, t.Title)
	return nil
}

func (a *App) cmdList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	filter := fs.String("filter", "all", "Filter: all, done, pending")
	if err := fs.Parse(args); err != nil {
		return err
	}

	tasks, err := a.storage.Load()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		color.Yellow("No tasks. yes. ADd one with: gotodo add <title>")
	}

	printed := 0
	for _, t := range tasks {
		switch *filter {
		case "done":
			if !t.Done {
				continue
			}
		case "pending":
			if t.Done {
				continue
			}
		case "all":
		default:
			return fmt.Errorf("invalid filter: %s", *filter)
		}
		PrintTask(t)
		printed++
	}
	if printed == 0 {
		color.Yellow("No tasks match filter: %s", *filter)
	}
	return nil
}

func (a *App) cmdDone(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: gotodo done <id>")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid id: %s", args[0])
	}

	tasks, err := a.storage.Load()
	if err != nil {
		return fmt.Errorf("Failed load storage: %w", err)
	}

	for _, t := range tasks {
		if t.ID == id {
			t.Done = true
			if err := a.storage.Save(tasks); err != nil {
				return err
			}
			color.Green("Marked task #%d as done", id)
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", id)
}

func (a *App) cmdRm(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: gotodo rm <id>")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid id: %s", args[0])
	}

	tasks, err := a.storage.Load()
	if err != nil {
		return fmt.Errorf("Failed load storage: %w", err)
	}

	filtered := make([]*task.Task, 0, len(tasks))
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			filtered = append(filtered, t)
		}
	}
	if !found {
		return fmt.Errorf("task #%d not found", id)
	}

	if err := a.storage.Save(filtered); err != nil {
		return err
	}
	color.Green("Removed task #%d", id)
	return nil
}

func (a *App) cmdClear(args []string) error {
	if err := a.storage.Save([]*task.Task{}); err != nil {
		return err
	}
	color.Green("All tasks cleared")
	return nil
}

func PrintTask(t *task.Task) {
	status := "[ ]"
	if t.Done {
		status = "[x]"
	}
	var priorityColor *color.Color
	switch t.Priority {
	case task.PriorityHigh:
		priorityColor = color.New(color.FgRed, color.Bold)
	case task.PriorityMedium:
		priorityColor = color.New(color.FgYellow)
	case task.PriorityLow:
		priorityColor = color.New(color.FgCyan)
	default:
		priorityColor = color.New(color.FgWhite)
	}
	line := fmt.Sprintf("%s #%d %s", status, t.ID, t.Title)
	if t.Done {
		color.New(color.FgGreen).Println(line)
	} else {
		fmt.Println(line)
	}
	fmt.Printf(" ")
	priorityColor.Printf("[%s]", t.Priority)

	if t.Due != "" {
		if t.IsOverdue() {
			color.Red(" due:%s (overdue!)", t.Due)
		} else {
			fmt.Printf(" due:%s", t.Due)
		}
	} else {
		fmt.Println()
	}
}
