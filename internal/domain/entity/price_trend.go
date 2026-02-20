package entity

import (
	"errors"
	"fmt"
	"time"
)

// PriceTrend represents a price trend record
type PriceTrend struct {
	ID         int64     `json:"id" db:"id"`
	ActivityID string    `json:"activityId" db:"activity_id"`
	Price      float64   `json:"price" db:"price"`
	RecordDate time.Time `json:"recordDate" db:"record_date"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
}

// BlockedProduct represents a blocked (hidden) product
type BlockedProduct struct {
	ActivityID string    `json:"activityId" db:"activity_id"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
}

// PlatformProductDTO represents a product from a platform API
type PlatformProductDTO struct {
	ActivityID         string    `json:"activityId"`
	Platform           string    `json:"platform"`
	Region             string    `json:"region"`
	Title              string    `json:"title"`
	ShopName           string    `json:"shopName"`
	OriginalPrice      float64   `json:"originalPrice"`
	CurrentPrice       float64   `json:"currentPrice"`
	SalesStatus        int       `json:"salesStatus"`
	ActivityCreateTime time.Time `json:"activityCreateTime"`
}

// DTInputDTO represents an input from DT platform
type DTInputDTO struct {
	Title     string
	Price     float64
	Status    int
	CrawlTime int64
	Region    string
}

// NewPriceTrend creates a new price trend record
func NewPriceTrend(activityID string, price float64, recordDate time.Time) (*PriceTrend, error) {
	if activityID == "" {
		return nil, errors.New("activityID cannot be empty")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative: %f", price)
	}
	if recordDate.IsZero() {
		return nil, errors.New("recordDate cannot be zero")
	}

	return &PriceTrend{
		ActivityID: activityID,
		Price:      price,
		RecordDate: recordDate,
		CreateTime: time.Now(),
	}, nil
}

// IsLowerThan checks if this price is lower than another
func (p *PriceTrend) IsLowerThan(other float64) bool {
	return p.Price < other
}
