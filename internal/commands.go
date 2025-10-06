package internal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// AddTask handles the "add" command
func AddTask(store *Store, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a task description")
	}

	description := strings.Join(args, " ")
	if err := store.AddTask(description); err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}

	fmt.Printf("✓ Task added: %s\n", description)
	return nil
}

// ListTasks handles the "list" command
func ListTasks(store *Store, args []string) error {
	tasks := store.GetTasks()

	if len(tasks) == 0 {
		fmt.Println("No tasks found. Add one with: gotask add <description>")
		return nil
	}

	fmt.Println("Your tasks:")
	for _, task := range tasks {
		status := "[ ]"
		if task.Done {
			status = "[✓]"
		}
		fmt.Printf("  %s %d: %s\n", status, task.ID, task.Title)
	}

	return nil
}

// MarkDone handles the "done" command
func MarkDone(store *Store, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a task ID")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	if err := store.MarkDone(id); err != nil {
		return err
	}

	fmt.Printf("✓ Task %d marked as done\n", id)
	return nil
}

// DeleteTask handles the "delete" command
func DeleteTask(store *Store, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a task ID")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	if err := store.DeleteTask(id); err != nil {
		return err
	}

	fmt.Printf("✓ Task %d deleted\n", id)
	return nil
}

// OpenFile handles the "open" command - opens the JSON file in the default editor
func OpenFile(store *Store, args []string) error {
	filePath := store.GetFilePath()

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	fmt.Printf("✓ Opened %s in default editor\n", filePath)
	return nil
}
