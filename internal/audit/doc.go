// Package audit analyses a slice of diff results and produces a structured
// audit report with severity-rated findings.
//
// # Overview
//
// After comparing two environment files with the diff package, callers can
// pass the resulting []diff.Result slice to audit.Run. The function inspects
// each result and emits a Finding for every non-matching key.
//
// # Severity Levels
//
//   - INFO     – key exists only in the target or a non-sensitive value changed.
//   - WARNING  – key is missing from the target, or a sensitive value changed.
//   - CRITICAL – a sensitive key (password, token, secret, …) is absent from
//     the target environment entirely.
//
// # Sensitive Key Detection
//
// A key is considered sensitive when its uppercased name contains one of the
// built-in patterns: SECRET, PASSWORD, TOKEN, API_KEY, PRIVATE, CREDENTIAL,
// AUTH, CERT, or KEY.
//
// # Usage
//
//	results := diff.Compare(sourceMap, targetMap)
//	report  := audit.Run(results)
//	for _, f := range report.Findings {
//	    fmt.Printf("[%s] %s\n", f.Severity, f.Message)
//	}
package audit
