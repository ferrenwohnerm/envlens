// Package snapshot provides functionality for saving and loading environment
// variable snapshots to disk, enabling point-in-time comparisons over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a saved state of an environment variable file.
type Snapshot struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Vars      map[string]string `json:"vars"`
}

// Save writes a snapshot to the given directory using the label as filename.
// The file is written as JSON with a .snap.json extension.
func Save(dir, label, source string, vars map[string]string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: create directory: %w", err)
	}

	snap := Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Vars:      vars,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	path := filepath.Join(dir, label+".snap.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write file: %w", err)
	}

	return nil
}

// Load reads a snapshot from disk by label within the given directory.
func Load(dir, label string) (*Snapshot, error) {
	path := filepath.Join(dir, label+".snap.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: %q not found in %s", label, dir)
		}
		return nil, fmt.Errorf("snapshot: read file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return &snap, nil
}

// List returns the labels of all snapshots stored in the given directory.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot: list directory: %w", err)
	}

	var labels []string
	for _, e := range entries {
		name := e.Name()
		if filepath.Ext(name) == ".json" && len(name) > 10 {
			labels = append(labels, name[:len(name)-len(".snap.json")])
		}
	}
	return labels, nil
}
