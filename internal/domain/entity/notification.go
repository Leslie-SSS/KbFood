package entity

import (
	"time"
)

// NotificationConfig represents a price notification configuration
type NotificationConfig struct {
	ActivityID     string     `json:"activityId" db:"activity_id"`
	UserID         string     `json:"userId" db:"user_id"`
	TargetPrice    float64    `json:"targetPrice" db:"target_price"`
	LastNotifyTime *time.Time `json:"lastNotifyTime" db:"last_notify_time"`
	CreateTime     time.Time  `json:"createTime" db:"create_time"`
	UpdateTime     time.Time  `json:"updateTime" db:"update_time"`
}

// ShouldNotify checks if a notification should be sent
func (n *NotificationConfig) ShouldNotify(currentPrice float64) bool {
	// Check if price condition is met
	if currentPrice > n.TargetPrice {
		return false
	}

	// Check if already notified today (local time)
	if n.LastNotifyTime != nil && sameCalendarDayInLocation(time.Now(), *n.LastNotifyTime, time.Now().Location()) {
		return false // Already notified today (local time)
	}

	return true
}

// MarkNotified marks the notification as sent
func (n *NotificationConfig) MarkNotified() {
	now := time.Now()
	n.LastNotifyTime = &now
}

// HasNotifiedToday returns true if already notified today
func (n *NotificationConfig) HasNotifiedToday() bool {
	if n.LastNotifyTime == nil {
		return false
	}

	return sameCalendarDayInLocation(time.Now(), *n.LastNotifyTime, time.Now().Location())
}

func sameCalendarDayInLocation(left, right time.Time, loc *time.Location) bool {
	if loc == nil {
		loc = time.Local
	}

	left = left.In(loc)
	right = right.In(loc)

	leftYear, leftMonth, leftDay := left.Date()
	rightYear, rightMonth, rightDay := right.Date()

	return leftYear == rightYear && leftMonth == rightMonth && leftDay == rightDay
}
