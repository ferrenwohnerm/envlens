package clone

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := baseEnv()
	dst := map[string]string{}
	opts := DefaultOptions()

	result, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", result["APP_HOST"])
	}
	if result["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080, got %q", result["APP_PORT"])
	}
}

func TestClone_PrefixReplace_RenamesKeys(t *testing.T) {
	src := map[string]string{"STAGING_DB": "stagedb", "STAGING_HOST": "stagehost"}
	dst := map[string]string{}
	opts := DefaultOptions()
	opts.PrefixReplace = map[string]string{"STAGING_": "PROD_"}

	result, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["PROD_DB"]; !ok {
		t.Error("expected PROD_DB to exist after prefix replacement")
	}
	if _, ok := result["STAGING_DB"]; ok {
		t.Error("expected STAGING_DB to not exist in result")
	}
}

func TestClone_NoOverwrite_SkipsExistingKeys(t *testing.T) {
	src := map[string]string{"APP_HOST": "newhost"}
	dst := map[string]string{"APP_HOST": "original"}
	opts := DefaultOptions()

	result, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_HOST"] != "original" {
		t.Errorf("expected APP_HOST to remain 'original', got %q", result["APP_HOST"])
	}
}

func TestClone_Overwrite_ReplacesExistingKeys(t *testing.T) {
	src := map[string]string{"APP_HOST": "newhost"}
	dst := map[string]string{"APP_HOST": "original"}
	opts := DefaultOptions()
	opts.Overwrite = true

	result, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_HOST"] != "newhost" {
		t.Errorf("expected APP_HOST=newhost, got %q", result["APP_HOST"])
	}
}

func TestClone_DryRun_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"NEW_KEY": "val"}
	dst := map[string]string{"OLD_KEY": "old"}
	opts := DefaultOptions()
	opts.DryRun = true

	_, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["NEW_KEY"]; ok {
		t.Error("DryRun should not mutate dst")
	}
}

func TestClone_DryRun_ReturnsPreviewWithNewKeys(t *testing.T) {
	src := map[string]string{"NEW_KEY": "val"}
	dst := map[string]string{"OLD_KEY": "old"}
	opts := DefaultOptions()
	opts.DryRun = true

	result, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY in dry-run result, got %q", result["NEW_KEY"])
	}
	if result["OLD_KEY"] != "old" {
		t.Errorf("expected OLD_KEY preserved in dry-run result, got %q", result["OLD_KEY"])
	}
}
