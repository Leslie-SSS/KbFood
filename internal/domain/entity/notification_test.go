package entity

import (
	"testing"
	"time"
)

func TestSameCalendarDayInLocation_ConvertsToLocalDay(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)

	now := time.Date(2026, 3, 28, 2, 0, 0, 0, loc)
	lastNotifyUTC := time.Date(2026, 3, 27, 18, 0, 0, 0, time.UTC)

	if !sameCalendarDayInLocation(now, lastNotifyUTC, loc) {
		t.Fatal("expected UTC timestamp to match the same local calendar day")
	}
}

func TestSameCalendarDayInLocation_DetectsDifferentLocalDays(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)

	now := time.Date(2026, 3, 28, 9, 0, 0, 0, loc)
	lastNotifyUTC := time.Date(2026, 3, 27, 15, 59, 59, 0, time.UTC)

	if sameCalendarDayInLocation(now, lastNotifyUTC, loc) {
		t.Fatal("expected different local calendar days")
	}
}
