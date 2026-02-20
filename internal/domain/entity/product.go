package entity

import (
	"time"
)

// Product represents a product entity
type Product struct {
	ID                 int64     `json:"id" db:"id"`
	ActivityID         string    `json:"activityId" db:"activity_id"`
	Platform           string    `json:"platform" db:"platform"`
	Region             string    `json:"region" db:"region"`
	Title              string    `json:"title" db:"title"`
	ShopName           string    `json:"shopName" db:"shop_name"`
	OriginalPrice      float64   `json:"originalPrice" db:"original_price"`
	CurrentPrice       float64   `json:"currentPrice" db:"current_price"`
	SalesStatus        int       `json:"salesStatus" db:"sales_status"`
	ActivityCreateTime time.Time `json:"activityCreateTime" db:"activity_create_time"`
	CreateTime         time.Time `json:"createTime" db:"create_time"`
	UpdateTime         time.Time `json:"updateTime" db:"update_time"`
}

// IsOnSale returns true if the product is on sale
func (p *Product) IsOnSale() bool {
	return p.SalesStatus == SalesStatusOnSale
}

// IsSold returns true if the product is sold out
func (p *Product) IsSold() bool {
	return p.SalesStatus == SalesStatusSold
}

// Discount calculates the discount rate (0-10)
func (p *Product) Discount() float64 {
	if p.OriginalPrice <= 0 {
		return 0
	}
	return (1 - p.CurrentPrice/p.OriginalPrice) * 10
}

// SalesStatusText returns the sales status text
func (p *Product) SalesStatusText() string {
	if p.IsOnSale() {
		return "在售"
	}
	return "已售"
}

// IsPriceDrop returns true if current price is lower than original
func (p *Product) IsPriceDrop() bool {
	return p.CurrentPrice < p.OriginalPrice
}

// DropRate returns the price drop rate as percentage
func (p *Product) DropRate() float64 {
	if p.OriginalPrice <= 0 {
		return 0
	}
	return (1 - p.CurrentPrice/p.OriginalPrice) * 100
}
