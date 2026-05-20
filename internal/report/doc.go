// Package report handles rendering of environment variable diff results
// into human-readable or machine-readable output formats.
//
// Supported formats:
//
//	- text: a tabular, human-friendly representation with status prefixes
//	        '-' for keys only in A, '+' for keys only in B,
//	        '~' for changed values, and ' ' for identical entries.
//
//	- json: a structured JSON object keyed by variable name, with
//	        status, value_a, and value_b fields per entry.
//
// The Write function accepts any io.Writer, making it straightforward to
// direct output to a file, buffer, or standard streams.
//
// Usage:
//
//	result := diff.Compare(envA, envB)
//	report.Write(os.Stdout, result, report.FormatText)
//
// To capture output as a string:
//
//	var buf bytes.Buffer
//	report.Write(&buf, result, report.FormatJSON)
//	fmt.Println(buf.String())
package report
