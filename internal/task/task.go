package task

import (
	"fmt"
	"time"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Done      bool     `json:"done"`
	Priority  Priority `json:"priority"`
	Due       string   `json:"due,omitempty"`
	CreatedAt string   `json:"created_at"`
}

func New(id int, title string, priority Priority, due string) *Task {
	return &Task{
		ID:        id,
		Title:     title,
		Done:      false,
		Priority:  priority,
		Due:       due,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

func (t *Task) IsOverdue() bool {
	if t.Due == "" || t.Done {
		return false
	}

	due, err := time.Parse("2006-01-02", t.Due)
	if err != nil {
		return false
	}
	return due.Before(time.Now())
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title cannt be empty")
	}

	switch t.Priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
	default:
		return fmt.Errorf("invalid priority: %s", t.Priority)
	}
	if t.Due != "" {
		if _, err := time.Parse("2006-01-02", t.Due); err != nil {
			return fmt.Errorf("invalid due date, use yyyy-MM-DD format: %w", err)
		}
	}
	return nil
}
