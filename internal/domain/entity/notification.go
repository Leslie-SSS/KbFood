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
	if n.LastNotifyTime != nil {
		now := time.Now()
		lastYear, lastMonth, lastDay := n.LastNotifyTime.Date()
		nowYear, nowMonth, nowDay := now.Date()

		if lastYear == nowYear && lastMonth == nowMonth && lastDay == nowDay {
			return false // Already notified today (local time)
		}
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

	now := time.Now()
	lastYear, lastMonth, lastDay := n.LastNotifyTime.Date()
	nowYear, nowMonth, nowDay := now.Date()

	return lastYear == nowYear && lastMonth == nowMonth && lastDay == nowDay
}
