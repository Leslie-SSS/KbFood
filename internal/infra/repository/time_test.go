package repository

import (
	"testing"
	"time"
)

func TestParseSQLiteTime_SupportsSQLiteDatetime(t *testing.T) {
	got := parseSQLiteTime("2026-03-27 12:20:00")
	if got.IsZero() {
		t.Fatal("expected SQLite datetime to parse")
	}
	if got.Location() != time.UTC {
		t.Fatalf("expected SQLite datetime to be interpreted as UTC, got %v", got.Location())
	}

	want := time.Date(2026, 3, 27, 12, 20, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("unexpected parsed time: got %v want %v", got, want)
	}
}

func TestParseSQLiteTime_SupportsRFC3339(t *testing.T) {
	got := parseSQLiteTime("2026-03-27T20:20:00+08:00")
	if got.IsZero() {
		t.Fatal("expected RFC3339 time to parse")
	}

	want := time.Date(2026, 3, 27, 20, 20, 0, 0, time.FixedZone("CST", 8*60*60))
	if !got.Equal(want) {
		t.Fatalf("unexpected parsed time: got %v want %v", got, want)
	}
}
