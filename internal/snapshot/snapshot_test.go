package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlens/internal/snapshot"
)

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{"APP_ENV": "staging", "PORT": "8080"}

	if err := snapshot.Save(dir, "staging-v1", ".env.staging", vars); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	path := filepath.Join(dir, "staging-v1.snap.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected snapshot file to exist at %s", path)
	}
}

func TestLoad_ReturnsCorrectVars(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	if err := snapshot.Save(dir, "prod-v2", ".env.prod", vars); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(dir, "prod-v2")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.Label != "prod-v2" {
		t.Errorf("Label = %q, want %q", snap.Label, "prod-v2")
	}
	if snap.Source != ".env.prod" {
		t.Errorf("Source = %q, want %q", snap.Source, ".env.prod")
	}
	if snap.Vars["DB_HOST"] != "localhost" {
		t.Errorf("Vars[DB_HOST] = %q, want %q", snap.Vars["DB_HOST"], "localhost")
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLoad_MissingLabel_ReturnsError(t *testing.T) {
	dir := t.TempDir()

	_, err := snapshot.Load(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}

func TestList_ReturnsLabels(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{"KEY": "value"}

	_ = snapshot.Save(dir, "snap-a", "a.env", vars)
	_ = snapshot.Save(dir, "snap-b", "b.env", vars)

	labels, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(labels) != 2 {
		t.Errorf("List() returned %d labels, want 2", len(labels))
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()

	labels, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty list, got %v", labels)
	}
}
