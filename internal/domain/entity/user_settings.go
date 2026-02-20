package entity

import "time"

// UserSettings stores user-specific settings including Bark key
type UserSettings struct {
	UserID     string     `json:"userId" db:"user_id"`
	BarkKey    string     `json:"barkKey" db:"bark_key"`
	CreateTime time.Time  `json:"createTime" db:"create_time"`
	UpdateTime time.Time  `json:"updateTime" db:"update_time"`
}
