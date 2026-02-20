package entity

import (
	"time"
)

// Sales status constants
const (
	SalesStatusSold   = 0
	SalesStatusOnSale = 1
)

// MasterProduct represents a standard product in the master catalog (mainly for DT platform)
type MasterProduct struct {
	ID            string    `json:"id" db:"id"`
	Region        string    `json:"region" db:"region"`
	Platform      string    `json:"platform" db:"platform"`
	StandardTitle string    `json:"standardTitle" db:"standard_title"`
	Price         float64   `json:"price" db:"price"`
	Status        int       `json:"status" db:"status"`
	TrustScore    int       `json:"trustScore" db:"trust_score"`
	CreateTime    time.Time `json:"createTime" db:"create_time"`
	UpdateTime    time.Time `json:"updateTime" db:"update_time"`
}

// IsOnSale returns true if the master product is on sale
func (m *MasterProduct) IsOnSale() bool {
	return m.Status == SalesStatusOnSale
}

// IncrementTrustScore increases the trust score
func (m *MasterProduct) IncrementTrustScore() {
	m.TrustScore++
}
