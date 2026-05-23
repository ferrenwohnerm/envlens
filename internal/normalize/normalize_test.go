package normalize

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db_host":       "localhost",
		"  api_key  ":   "secret",
		"app-name":      "envlens",
		"PORT":          "  8080  ",
		"Already_Upper": "value",
	}
}

func TestApply_DefaultOptions_UppercasesKeys(t *testing.T) {
	res := Apply(map[string]string{"db_host": "localhost"}, DefaultOptions())
	if _, ok := res.Vars["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist after normalization")
	}
	if _, ok := res.Vars["db_host"]; ok {
		t.Error("expected original key db_host to be absent")
	}
}

func TestApply_DefaultOptions_TrimsValues(t *testing.T) {
	res := Apply(map[string]string{"PORT": "  8080  "}, DefaultOptions())
	if got := res.Vars["PORT"]; got != "8080" {
		t.Errorf("expected trimmed value '8080', got %q", got)
	}
}

func TestApply_DefaultOptions_TrimsKeys(t *testing.T) {
	res := Apply(map[string]string{"  api_key  ": "secret"}, DefaultOptions())
	if _, ok := res.Vars["API_KEY"]; !ok {
		t.Error("expected API_KEY after key trimming and uppercasing")
	}
}

func TestApply_DefaultOptions_ReplacesHyphens(t *testing.T) {
	res := Apply(map[string]string{"app-name": "envlens"}, DefaultOptions())
	if _, ok := res.Vars["APP_NAME"]; !ok {
		t.Error("expected APP_NAME after hyphen replacement and uppercasing")
	}
}

func TestApply_RecordsRenamedKeys(t *testing.T) {
	res := Apply(map[string]string{"db_host": "localhost", "PORT": "8080"}, DefaultOptions())
	if _, ok := res.Renamed["db_host"]; !ok {
		t.Error("expected db_host to appear in Renamed map")
	}
	// PORT is already uppercase — should not appear in Renamed
	if _, ok := res.Renamed["PORT"]; ok {
		t.Error("PORT was not renamed and should not appear in Renamed map")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"db_host": "localhost"}
	Apply(input, DefaultOptions())
	if _, ok := input["DB_HOST"]; ok {
		t.Error("Apply must not mutate the input map")
	}
	if _, ok := input["db_host"]; !ok {
		t.Error("original key should remain in input map")
	}
}

func TestApply_DisabledOptions_LeavesKeysUnchanged(t *testing.T) {
	opts := Options{
		UppercaseKeys:  false,
		TrimValues:     false,
		TrimKeys:       false,
		ReplaceHyphens: false,
	}
	input := map[string]string{"app-name": "  envlens  "}
	res := Apply(input, opts)
	if v, ok := res.Vars["app-name"]; !ok || v != "  envlens  " {
		t.Errorf("expected key and value to be unchanged, got key present=%v val=%q", ok, v)
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renamed keys, got %v", res.Renamed)
	}
}
