package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/arjunsajeev/gotask/models"
)

type Store struct {
	filePath string
	tasks    []models.Task
}

// NewStore creates a new store instance and loads existing tasks
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, ".gotask.json")
	store := &Store{
		filePath: filePath,
		tasks:    []models.Task{},
	}

	// Load existing tasks
	if err := store.load(); err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	return store, nil
}

func (s *Store) AddTask(description string) error {
	maxID := 0
	for _, task := range s.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	newTask := models.Task{
		ID:    maxID + 1,
		Title: description,
		Done:  false,
	}

	s.tasks = append(s.tasks, newTask)
	return s.save()
}

func (s *Store) GetTasks() []models.Task {
	return s.tasks
}

// GetFilePath returns the file path where tasks are stored
func (s *Store) GetFilePath() string {
	return s.filePath
}

// MarkDone marks a task as completed
func (s *Store) MarkDone(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = true
			return s.save()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// DeleteTask removes a task from the store
func (s *Store) DeleteTask(id int) error {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// load reads tasks from the JSON file
func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return err
	}

	return json.Unmarshal(data, &s.tasks)
}

// save writes tasks to the JSON file
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
