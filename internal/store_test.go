package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjunsajeev/gotask/models"
)

// Helper function to create a temporary store for testing
func createTestStore(t *testing.T) (*Store, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gotask_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create store with temp file
	tempFile := filepath.Join(tempDir, "test_tasks.json")
	store := &Store{
		filePath: tempFile,
		tasks:    []models.Task{},
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return store, cleanup
}

func TestStore_AddTask(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Test adding a task
	err := store.AddTask("Test task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// Verify task was added
	tasks := store.GetTasks()
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
	if task.Title != "Test task" {
		t.Errorf("Expected title 'Test task', got '%s'", task.Title)
	}
	if task.Done {
		t.Errorf("Expected task to be not done")
	}
}

func TestStore_AddMultipleTasks(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Add multiple tasks
	tasks := []string{"Task 1", "Task 2", "Task 3"}
	for _, title := range tasks {
		err := store.AddTask(title)
		if err != nil {
			t.Fatalf("AddTask failed for '%s': %v", title, err)
		}
	}

	// Verify all tasks were added with correct IDs
	storedTasks := store.GetTasks()
	if len(storedTasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(storedTasks))
	}

	for i, task := range storedTasks {
		expectedID := i + 1
		expectedTitle := tasks[i]

		if task.ID != expectedID {
			t.Errorf("Task %d: expected ID %d, got %d", i, expectedID, task.ID)
		}
		if task.Title != expectedTitle {
			t.Errorf("Task %d: expected title '%s', got '%s'", i, expectedTitle, task.Title)
		}
	}
}

func TestStore_MarkDone(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Add a task
	err := store.AddTask("Test task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// Mark it as done
	err = store.MarkDone(1)
	if err != nil {
		t.Fatalf("MarkDone failed: %v", err)
	}

	// Verify it's marked as done
	tasks := store.GetTasks()
	if !tasks[0].Done {
		t.Errorf("Task should be marked as done")
	}
}

func TestStore_MarkDone_NonExistent(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Try to mark non-existent task as done
	err := store.MarkDone(999)
	if err == nil {
		t.Errorf("Expected error when marking non-existent task as done")
	}
}

func TestStore_DeleteTask(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Add multiple tasks
	store.AddTask("Task 1")
	store.AddTask("Task 2")
	store.AddTask("Task 3")

	// Delete middle task
	err := store.DeleteTask(2)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	// Verify task was deleted
	tasks := store.GetTasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks after deletion, got %d", len(tasks))
	}

	// Verify remaining tasks are correct
	if tasks[0].ID != 1 || tasks[0].Title != "Task 1" {
		t.Errorf("First task is incorrect: ID=%d, Title=%s", tasks[0].ID, tasks[0].Title)
	}
	if tasks[1].ID != 3 || tasks[1].Title != "Task 3" {
		t.Errorf("Second task is incorrect: ID=%d, Title=%s", tasks[1].ID, tasks[1].Title)
	}
}

func TestStore_Persistence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gotask_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "test_tasks.json")

	// Create first store and add tasks
	store1 := &Store{
		filePath: tempFile,
		tasks:    []models.Task{},
	}

	store1.AddTask("Persistent task 1")
	store1.AddTask("Persistent task 2")

	// Create second store with same file - should load existing tasks
	store2 := &Store{
		filePath: tempFile,
		tasks:    []models.Task{},
	}

	err = store2.load()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Verify tasks were loaded
	tasks := store2.GetTasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks after loading, got %d", len(tasks))
	}

	if tasks[0].Title != "Persistent task 1" {
		t.Errorf("First task title incorrect: got '%s'", tasks[0].Title)
	}
	if tasks[1].Title != "Persistent task 2" {
		t.Errorf("Second task title incorrect: got '%s'", tasks[1].Title)
	}
}
