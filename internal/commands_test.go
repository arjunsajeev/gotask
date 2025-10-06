package internal

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper function to capture stdout
func captureOutput(fn func()) string {
	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()

	// Save original stdout
	originalStdout := os.Stdout
	os.Stdout = w

	// Channel to capture the output
	output := make(chan string)

	// Goroutine to read from pipe
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output <- buf.String()
	}()

	// Run the function
	fn()

	// Restore stdout and close writer
	os.Stdout = originalStdout
	w.Close()

	// Get the captured output
	return <-output
}

func TestAddTask_Success(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	output := captureOutput(func() {
		err := AddTask(store, []string{"Test", "task", "description"})
		if err != nil {
			t.Errorf("AddTask failed: %v", err)
		}
	})

	expectedOutput := "✓ Task added: Test task description\n"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}

	// Verify task was actually added
	tasks := store.GetTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestAddTask_EmptyArgs(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	err := AddTask(store, []string{})
	if err == nil {
		t.Error("Expected error for empty args, got nil")
	}

	expectedError := "please provide a task description"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestListTasks_Empty(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	output := captureOutput(func() {
		err := ListTasks(store, []string{})
		if err != nil {
			t.Errorf("ListTasks failed: %v", err)
		}
	})

	expectedOutput := "No tasks found. Add one with: gotask add <description>\n"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}

func TestListTasks_WithTasks(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Add some tasks
	store.AddTask("Task 1")
	store.AddTask("Task 2")
	store.MarkDone(1) // Mark first task as done

	output := captureOutput(func() {
		err := ListTasks(store, []string{})
		if err != nil {
			t.Errorf("ListTasks failed: %v", err)
		}
	})

	// Check that output contains expected elements
	if !strings.Contains(output, "Your tasks:") {
		t.Error("Output should contain 'Your tasks:'")
	}
	if !strings.Contains(output, "[✓] 1: Task 1") {
		t.Error("Output should show completed task 1")
	}
	if !strings.Contains(output, "[ ] 2: Task 2") {
		t.Error("Output should show uncompleted task 2")
	}
}

func TestMarkDone_Success(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Add a task first
	store.AddTask("Test task")

	output := captureOutput(func() {
		err := MarkDone(store, []string{"1"})
		if err != nil {
			t.Errorf("MarkDone failed: %v", err)
		}
	})

	expectedOutput := "✓ Task 1 marked as done\n"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}

	// Verify task is actually marked as done
	tasks := store.GetTasks()
	if !tasks[0].Done {
		t.Error("Task should be marked as done")
	}
}

func TestMarkDone_InvalidID(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Test with non-numeric ID
	err := MarkDone(store, []string{"abc"})
	if err == nil {
		t.Error("Expected error for invalid ID, got nil")
	}
	if !strings.Contains(err.Error(), "invalid task ID") {
		t.Errorf("Expected 'invalid task ID' error, got '%s'", err.Error())
	}
}

func TestMarkDone_EmptyArgs(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	err := MarkDone(store, []string{})
	if err == nil {
		t.Error("Expected error for empty args, got nil")
	}

	expectedError := "please provide a task ID"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDeleteTask_Success(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// Add a task first
	store.AddTask("Test task")

	output := captureOutput(func() {
		err := DeleteTask(store, []string{"1"})
		if err != nil {
			t.Errorf("DeleteTask failed: %v", err)
		}
	})

	expectedOutput := "✓ Task 1 deleted\n"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}

	// Verify task is actually deleted
	tasks := store.GetTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(tasks))
	}
}
