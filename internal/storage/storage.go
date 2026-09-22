package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/crab4/todo-learning-helper/internal/task"
)

type Storage struct {
	path string
}

func New(path string) *Storage {
	return &Storage{path: path}
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	return filepath.Join(home, ".gotodo", "tasks.join"), nil
}

func (s *Storage) Load() ([]*task.Task, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*task.Task{}, nil
		}
		return nil, fmt.Errorf("readFile: %w", err)
	}

	if len(data) == 0 {
		return []*task.Task{}, nil
	}

	var tasks []*task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal tasks: %w", err)
	}
	return tasks, nil
}

func (s *Storage) Save(tasks []*task.Task) error {
	if err := s.ensureDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (s *Storage) ensureDir() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return nil
}
