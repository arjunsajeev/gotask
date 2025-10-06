package main

import (
	"fmt"
	"os"

	"github.com/arjunsajeev/gotask/internal"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Initialize store once at startup
	store, err := internal.NewStore()
	if err != nil {
		fmt.Printf("Error initializing task store: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Route commands to handlers
	switch command {
	case "add":
		err = internal.AddTask(store, args)
	case "list", "ls":
		err = internal.ListTasks(store, args)
	case "done", "complete":
		err = internal.MarkDone(store, args)
	case "delete", "del", "rm":
		err = internal.DeleteTask(store, args)
	case "open", "edit":
		err = internal.OpenFile(store, args)
	case "help", "-h", "--help":
		printUsage()
		return
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	// Handle any errors from command execution
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`gotask - A simple task manager

Usage:
  gotask add <description>    Add a new task
  gotask list                 List all tasks
  gotask done <id>           Mark task as done
  gotask delete <id>         Delete a task
  gotask open                Open the task file in default editor
  gotask help                Show this help

Examples:
  gotask add "Buy groceries"
  gotask list
  gotask done 1
  gotask delete 2
  gotask open`)
}
