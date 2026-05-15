package mask_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/mask"
	"github.com/yourorg/envlens/internal/parser"
	"os"
)

// TestApply_WithParsedFile ensures mask.Apply works correctly on output
// produced by parser.ParseFile, simulating the real CLI pipeline.
func TestApply_WithParsedFile(t *testing.T) {
	content := `APP_ENV=production
DB_PASSWORD=supersecret
API_KEY=key-abc-123
PORT=8080
`
	f, err := os.CreateTemp("", "envlens-mask-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()

	env, err := parser.ParseFile(f.Name())
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	masked := mask.Apply(env, mask.Options{})

	if masked["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["API_KEY"] != "***" {
		t.Errorf("API_KEY should be masked, got %q", masked["API_KEY"])
	}
	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be preserved, got %q", masked["APP_ENV"])
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should be preserved, got %q", masked["PORT"])
	}
}
