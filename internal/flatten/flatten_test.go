package flatten_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/flatten"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db.host":    "localhost",
		"db.port":    "5432",
		"app.debug":  "true",
		"APP_SECRET": "s3cr3t",
	}
}

func TestRun_DefaultOptions_UppercasesAndReplaceDots(t *testing.T) {
	out, err := flatten.Run(baseEnv(), flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for k := range out {
		if k != strings.ToUpper(k) {
			t.Errorf("key %q is not uppercase", k)
		}
		if strings.Contains(k, ".") {
			t.Errorf("key %q still contains a dot", k)
		}
	}

	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestRun_CustomDelimiter_UsedInOutput(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Delimiter = "__"

	out, err := flatten.Run(map[string]string{"db.host": "localhost"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := out["DB__HOST"]; !ok {
		t.Errorf("expected key DB__HOST, got %v", out)
	}
}

func TestRun_Prefix_PrependedToKeys(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Prefix = "APP_"

	out, err := flatten.Run(map[string]string{"db.host": "localhost"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := out["APP_DB_HOST"]; !ok {
		t.Errorf("expected APP_DB_HOST, got %v", out)
	}
}

func TestRun_DuplicateOutputKeys_ReturnsError(t *testing.T) {
	env := map[string]string{
		"db.host": "localhost",
		"db_host": "remotehost",
	}

	_, err := flatten.Run(env, flatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for duplicate output keys, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("error message should mention 'duplicate', got: %v", err)
	}
}

func TestRun_SlashSeparator_Normalised(t *testing.T) {
	env := map[string]string{"db/host": "localhost"}
	opts := flatten.DefaultOptions()

	out, err := flatten.Run(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}

func TestRun_EmptyMap_ReturnsEmpty(t *testing.T) {
	out, err := flatten.Run(map[string]string{}, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
