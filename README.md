# gotask

A simple, fast task manager for the command line.

## Installation

```bash
git clone https://github.com/yourusername/gotask.git
cd gotask
go build -o gotask
```

## Usage

```bash
# Add a new task
gotask add "Buy groceries"

# List all tasks
gotask list

# Mark task as done
gotask done 1

# Delete a task
gotask delete 2

# Open tasks file in editor
gotask open

# Show help
gotask help
```

## Features

- ✅ Add, list, complete, and delete tasks
- 💾 Persistent storage in `~/.gotask.json`
- 🚀 Fast and lightweight
- 🔧 Cross-platform (macOS, Linux, Windows)

## Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `add <description>` | | Add a new task |
| `list` | `ls` | List all tasks |
| `done <id>` | `complete` | Mark task as done |
| `delete <id>` | `del`, `rm` | Delete a task |
| `open` | `edit` | Open task file in default editor |
| `help` | `-h`, `--help` | Show help |

## Requirements

- Go 1.19+

## License

MIT
