package entity

import "time"

// SyncStatus represents the status of a sync job
type SyncStatus struct {
	ID           int64     `json:"id" db:"id"`
	JobName      string    `json:"jobName" db:"job_name"`
	LastRunTime  time.Time `json:"lastRunTime" db:"last_run_time"`
	Status       string    `json:"status" db:"status"` // success, failed, running
	ProductCount int       `json:"productCount" db:"product_count"`
	ErrorMessage string    `json:"errorMessage" db:"error_message"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// IsHealthy returns true if the sync job ran recently (within 30 minutes)
func (s *SyncStatus) IsHealthy() bool {
	return time.Since(s.LastRunTime) < 30*time.Minute
}

// Status constants
const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusRunning = "running"
)
