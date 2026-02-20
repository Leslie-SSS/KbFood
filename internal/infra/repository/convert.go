package repository

import (
	"database/sql"
	"time"
)

// Helper functions for null conversions used across repositories

func sqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func sqlNullInt32(i int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

func sqlNullInt64(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: true}
}

func sqlNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: true}
}

func sqlNullTime(t interface{}) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	if time, ok := t.(time.Time); ok {
		return sql.NullTime{Time: time, Valid: true}
	}
	return sql.NullTime{}
}

func sqlNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func sqlNullTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func stringFromNull(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func float64FromNull(nf sql.NullFloat64) float64 {
	if !nf.Valid {
		return 0
	}
	return nf.Float64
}

func intFromNull(ni sql.NullInt32) int {
	if !ni.Valid {
		return 0
	}
	return int(ni.Int32)
}

func int64FromNull(ni sql.NullInt64) int64 {
	if !ni.Valid {
		return 0
	}
	return ni.Int64
}

func timeFromNull(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

// SQLite date/time conversion helpers
// SQLite stores dates as strings (RFC3339 format)

func timeToSQLite(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// dateToSQLite formats a time as date only (YYYY-MM-DD)
// Use this for date fields where time component should not affect uniqueness
func dateToSQLite(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// parseSQLiteDate parses a date string (YYYY-MM-DD or RFC3339) to time.Time
// Handles both formats for backward compatibility
func parseSQLiteDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	// Try date-only format first
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		return t
	}
	// Fall back to RFC3339 format
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func sqlNullInt64FromInt(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: true}
}

func sqlNullFloat64FromFloat(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

func sqlNullStringFromTimePtr(t *time.Time) sql.NullString {
	if t == nil || t.IsZero() {
		return sql.NullString{}
	}
	return sql.NullString{String: timeToSQLite(*t), Valid: true}
}

func sqlNullStringFromTime(t time.Time) sql.NullString {
	if t.IsZero() {
		return sql.NullString{}
	}
	return sql.NullString{String: timeToSQLite(t), Valid: true}
}
