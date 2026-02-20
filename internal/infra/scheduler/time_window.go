package scheduler

import (
	"fmt"
	"time"
)

// TimeWindow represents a time window for task execution
type TimeWindow struct {
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
}

// DefaultTimeWindow returns the default time window (00:00 - 23:59 for testing)
// TODO: Make this configurable via environment variable
func DefaultTimeWindow() *TimeWindow {
	return &TimeWindow{
		StartHour:   0,
		StartMinute: 0,
		EndHour:     23,
		EndMinute:   59,
	}
}

// IsActive checks if the current time is within the time window
func (tw *TimeWindow) IsActive(t time.Time) bool {
	hour, min, _ := t.Clock()

	now := hour*60 + min
	start := tw.StartHour*60 + tw.StartMinute
	end := tw.EndHour*60 + tw.EndMinute

	return now >= start && now < end
}

// IsActiveNow checks if the current time is within the time window
func (tw *TimeWindow) IsActiveNow() bool {
	return tw.IsActive(time.Now())
}

// String returns a string representation of the time window
func (tw *TimeWindow) String() string {
	return fmt.Sprintf("%02d:%02d - %02d:%02d", tw.StartHour, tw.StartMinute, tw.EndHour, tw.EndMinute)
}
