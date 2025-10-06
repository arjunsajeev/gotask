package models

import (
	"encoding/json"
	"testing"
)

func TestTaskJSONSerialization(t *testing.T) {
	// Test that Task can be serialized to/from JSON correctly
	original := Task{
		ID:    1,
		Title: "Test task",
		Done:  false,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	// Unmarshal back
	var restored Task
	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	// Compare
	if original.ID != restored.ID {
		t.Errorf("ID mismatch: expected %d, got %d", original.ID, restored.ID)
	}
	if original.Title != restored.Title {
		t.Errorf("Title mismatch: expected %s, got %s", original.Title, restored.Title)
	}
	if original.Done != restored.Done {
		t.Errorf("Done mismatch: expected %t, got %t", original.Done, restored.Done)
	}
}

func TestTaskValidation(t *testing.T) {
	// Test edge cases
	tests := []struct {
		name string
		task Task
		want string
	}{
		{
			name: "empty task",
			task: Task{ID: 0, Title: "", Done: false},
			want: `{"id":0,"title":"","done":false}`,
		},
		{
			name: "completed task",
			task: Task{ID: 42, Title: "Done task", Done: true},
			want: `{"id":42,"title":"Done task","done":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.task)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("JSON mismatch:\nwant: %s\ngot:  %s", tt.want, string(data))
			}
		})
	}
}
