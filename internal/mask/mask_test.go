package mask_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/mask"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "envlens",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"LOG_LEVEL":   "info",
		"AUTH_TOKEN":  "tok_xyz",
		"REGION":      "us-east-1",
	}
}

func TestApply_MasksSensitiveKeys(t *testing.T) {
	result := mask.Apply(baseEnv(), mask.Options{})

	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN"}
	for _, k := range sensitive {
		if result[k] != "***" {
			t.Errorf("expected key %q to be masked, got %q", k, result[k])
		}
	}
}

func TestApply_PreservesNonSensitiveKeys(t *testing.T) {
	result := mask.Apply(baseEnv(), mask.Options{})

	plain := map[string]string{"APP_NAME": "envlens", "LOG_LEVEL": "info", "REGION": "us-east-1"}
	for k, want := range plain {
		if result[k] != want {
			t.Errorf("key %q: expected %q, got %q", k, want, result[k])
		}
	}
}

func TestApply_MaskAll_MasksEverything(t *testing.T) {
	result := mask.Apply(baseEnv(), mask.Options{MaskAll: true})
	for k, v := range result {
		if v != "***" {
			t.Errorf("MaskAll: key %q should be masked, got %q", k, v)
		}
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	result := mask.Apply(baseEnv(), mask.Options{MaskAll: true, Placeholder: "<redacted>"})
	for k, v := range result {
		if v != "<redacted>" {
			t.Errorf("key %q: expected \"<redacted>\", got %q", k, v)
		}
	}
}

func TestApply_EmptyMap(t *testing.T) {
	result := mask.Apply(map[string]string{}, mask.Options{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := baseEnv()
	orig := input["DB_PASSWORD"]
	mask.Apply(input, mask.Options{})
	if input["DB_PASSWORD"] != orig {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_CustomSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"MY_SECRET":   "hunter2",
		"APP_NAME":    "envlens",
		"PRIVATE_KEY": "rsa-data",
	}
	result := mask.Apply(env, mask.Options{SensitiveKeys: []string{"MY_SECRET", "PRIVATE_KEY"}})

	if result["MY_SECRET"] != "***" {
		t.Errorf("expected MY_SECRET to be masked, got %q", result["MY_SECRET"])
	}
	if result["PRIVATE_KEY"] != "***" {
		t.Errorf("expected PRIVATE_KEY to be masked, got %q", result["PRIVATE_KEY"])
	}
	if result["APP_NAME"] != "envlens" {
		t.Errorf("expected APP_NAME to be preserved, got %q", result["APP_NAME"])
	}
}

func TestApply_ResultHasSameKeyCount(t *testing.T) {
	input := baseEnv()
	result := mask.Apply(input, mask.Options{})
	if len(result) != len(input) {
		t.Errorf("expected result to have %d keys, got %d", len(input), len(result))
	}
}
