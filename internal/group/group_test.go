package group_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/group"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "mydb",
		"REDIS_HOST":  "127.0.0.1",
		"REDIS_PORT":  "6379",
		"PORT":        "8080",
		"DEBUG":       "true",
	}
}

func TestRun_DefaultOptions_GroupsByFirstSegment(t *testing.T) {
	result := group.Run(baseEnv(), group.DefaultOptions())

	if _, ok := result.Groups["DB"]; !ok {
		t.Error("expected group DB")
	}
	if _, ok := result.Groups["REDIS"]; !ok {
		t.Error("expected group REDIS")
	}
}

func TestRun_UngroupedKeys_PlacedInUngroupedBucket(t *testing.T) {
	result := group.Run(baseEnv(), group.DefaultOptions())

	ug, ok := result.Groups["(ungrouped)"]
	if !ok {
		t.Fatal("expected ungrouped bucket")
	}
	if _, has := ug["PORT"]; !has {
		t.Error("expected PORT in ungrouped")
	}
	if _, has := ug["DEBUG"]; !has {
		t.Error("expected DEBUG in ungrouped")
	}
}

func TestRun_OrderIsSorted(t *testing.T) {
	result := group.Run(baseEnv(), group.DefaultOptions())

	for i := 1; i < len(result.Order); i++ {
		if result.Order[i] < result.Order[i-1] {
			t.Errorf("order not sorted: %v", result.Order)
		}
	}
}

func TestRun_CustomDelimiter_GroupsCorrectly(t *testing.T) {
	env := map[string]string{
		"app.host": "localhost",
		"app.port": "8080",
		"db.host":  "db.example.com",
	}
	opts := group.DefaultOptions()
	opts.Delimiter = "."

	result := group.Run(env, opts)

	if _, ok := result.Groups["app"]; !ok {
		t.Error("expected group app")
	}
	if _, ok := result.Groups["db"]; !ok {
		t.Error("expected group db")
	}
}

func TestRun_MaxDepth2_UsesFirstTwoSegments(t *testing.T) {
	env := map[string]string{
		"AWS_S3_BUCKET": "my-bucket",
		"AWS_S3_REGION": "us-east-1",
		"AWS_EC2_AMI":   "ami-12345",
	}
	opts := group.DefaultOptions()
	opts.MaxDepth = 2

	result := group.Run(env, opts)

	if _, ok := result.Groups["AWS_S3"]; !ok {
		t.Error("expected group AWS_S3")
	}
	if _, ok := result.Groups["AWS_EC2"]; !ok {
		t.Error("expected group AWS_EC2")
	}
}

func TestRun_EmptyEnv_ReturnsEmptyResult(t *testing.T) {
	result := group.Run(map[string]string{}, group.DefaultOptions())

	if len(result.Groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(result.Groups))
	}
	if len(result.Order) != 0 {
		t.Errorf("expected empty order, got %d", len(result.Order))
	}
}

func TestRun_CustomUngroupedLabel(t *testing.T) {
	env := map[string]string{
		"STANDALONE": "yes",
	}
	opts := group.DefaultOptions()
	opts.Ungrouped = "misc"

	result := group.Run(env, opts)

	if _, ok := result.Groups["misc"]; !ok {
		t.Error("expected group misc for ungrouped keys")
	}
}
