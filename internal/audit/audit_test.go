package audit

import (
	"testing"

	"github.com/envlens/internal/diff"
)

func makeResults(entries []diff.Result) []diff.Result {
	return entries
}

func TestRun_NoFindings_WhenAllMatch(t *testing.T) {
	results := makeResults([]diff.Result{
		{Key: "APP_ENV", Status: diff.StatusMatch, ValueA: "prod", ValueB: "prod"},
	})
	report := Run(results)
	if len(report.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(report.Findings))
	}
}

func TestRun_MissingInB_ProducesWarning(t *testing.T) {
	results := makeResults([]diff.Result{
		{Key: "LOG_LEVEL", Status: diff.StatusMissingInB, ValueA: "debug", ValueB: ""},
	})
	report := Run(results)
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Severity != SeverityWarning {
		t.Errorf("expected WARNING, got %s", report.Findings[0].Severity)
	}
}

func TestRun_SensitiveKeyMissingInB_ProducesCritical(t *testing.T) {
	results := makeResults([]diff.Result{
		{Key: "DB_PASSWORD", Status: diff.StatusMissingInB, ValueA: "s3cr3t", ValueB: ""},
	})
	report := Run(results)
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Severity != SeverityCritical {
		t.Errorf("expected CRITICAL for sensitive key, got %s", report.Findings[0].Severity)
	}
}

func TestRun_MissingInA_ProducesInfo(t *testing.T) {
	results := makeResults([]diff.Result{
		{Key: "NEW_FLAG", Status: diff.StatusMissingInA, ValueA: "", ValueB: "true"},
	})
	report := Run(results)
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Severity != SeverityInfo {
		t.Errorf("expected INFO, got %s", report.Findings[0].Severity)
	}
}

func TestRun_SensitiveChanged_ProducesWarning(t *testing.T) {
	results := makeResults([]diff.Result{
		{Key: "API_KEY", Status: diff.StatusChanged, ValueA: "old", ValueB: "new"},
	})
	report := Run(results)
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Severity != SeverityWarning {
		t.Errorf("expected WARNING for sensitive changed key, got %s", report.Findings[0].Severity)
	}
}

func TestRun_GeneratedAt_IsSet(t *testing.T) {
	report := Run(nil)
	if report.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}
