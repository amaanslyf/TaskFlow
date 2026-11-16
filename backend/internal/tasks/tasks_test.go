package tasks

import (
	"testing"
)

func TestValidStatuses(t *testing.T) {
	valid := []string{"todo", "in_progress", "completed"}
	for _, s := range valid {
		if !ValidStatuses[s] {
			t.Errorf("status %q should be valid", s)
		}
	}

	invalid := []string{"done", "pending", "archived", "random", ""}
	for _, s := range invalid {
		if ValidStatuses[s] {
			t.Errorf("status %q should be invalid", s)
		}
	}
}

