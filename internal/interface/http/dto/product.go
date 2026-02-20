package dto

import (
	"time"

	"kbfood/internal/domain/entity"
)

// ProductDTO represents a product response
type ProductDTO struct {
	ID                 int64     `json:"id"`
	ActivityID         string    `json:"activityId"`
	Platform           string    `json:"platform"`
	Region             string    `json:"region"`
	Title              string    `json:"title"`
	ShopName           string    `json:"shopName"`
	OriginalPrice      float64   `json:"originalPrice"`
	CurrentPrice       float64   `json:"currentPrice"`
	SalesStatus        int       `json:"salesStatus"`
	SalesStatusText    string    `json:"salesStatusText"`
	ActivityCreateTime time.Time `json:"activityCreateTime"`
	CreateTime         time.Time `json:"createTime"`
	UpdateTime         time.Time `json:"updateTime"`
	Discount           float64   `json:"discount,omitempty"`
	DropRate           float64   `json:"dropRate,omitempty"`
	HasNotification    bool      `json:"hasNotification,omitempty"`
	TargetPrice        *float64  `json:"targetPrice,omitempty"`
}

// PriceTrendDTO represents a price trend point
type PriceTrendDTO struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

// NotificationDTO represents a notification config response
type NotificationDTO struct {
	ActivityID     string  `json:"activityId"`
	TargetPrice    float64 `json:"targetPrice"`
	LastNotifyTime *string `json:"lastNotifyTime,omitempty"`
}

// FromEntity converts a Product entity to DTO
func FromEntity(p *entity.Product) ProductDTO {
	if p == nil {
		return ProductDTO{}
	}
	return ProductDTO{
		ID:                 p.ID,
		ActivityID:         p.ActivityID,
		Platform:           p.Platform,
		Region:             p.Region,
		Title:              p.Title,
		ShopName:           p.ShopName,
		OriginalPrice:      p.OriginalPrice,
		CurrentPrice:       p.CurrentPrice,
		SalesStatus:        p.SalesStatus,
		SalesStatusText:    p.SalesStatusText(),
		ActivityCreateTime: p.ActivityCreateTime,
		CreateTime:         p.CreateTime,
		UpdateTime:         p.UpdateTime,
		Discount:           p.Discount(),
		DropRate:           p.DropRate(),
	}
}

// FromMasterEntity converts a MasterProduct entity to DTO
func FromMasterEntity(m *entity.MasterProduct) ProductDTO {
	if m == nil {
		return ProductDTO{}
	}
	statusText := "未知"
	switch m.Status {
	case 1:
		statusText = "在售"
	case 0:
		statusText = "售罄"
	case -1:
		statusText = "下架"
	}

	// Use entity's platform if set, otherwise default to "探探糖"
	platform := m.Platform
	if platform == "" {
		platform = "探探糖"
	}

	return ProductDTO{
		ActivityID:      m.ID,
		Platform:        platform,
		Region:          m.Region,
		Title:           m.StandardTitle,
		ShopName:        platform + "精选",
		CurrentPrice:    m.Price,
		OriginalPrice:   m.Price,
		SalesStatus:     m.Status,
		SalesStatusText: statusText,
		CreateTime:      m.CreateTime,
		UpdateTime:      m.UpdateTime,
	}
}

// FromMasterEntities converts multiple MasterProduct entities to DTOs
func FromMasterEntities(products []*entity.MasterProduct) []ProductDTO {
	result := make([]ProductDTO, 0, len(products))
	for _, p := range products {
		if p != nil {
			result = append(result, FromMasterEntity(p))
		}
	}
	return result
}

// FromEntityWithNotification converts a Product entity to DTO with notification info
func FromEntityWithNotification(p *entity.Product, noti *entity.NotificationConfig) ProductDTO {
	dto := FromEntity(p)
	if noti != nil {
		dto.HasNotification = true
		dto.TargetPrice = &noti.TargetPrice
	}
	return dto
}

// FromEntities converts multiple Product entities to DTOs
func FromEntities(products []*entity.Product) []ProductDTO {
	result := make([]ProductDTO, 0, len(products))
	for _, p := range products {
		if p != nil {
			result = append(result, FromEntity(p))
		}
	}
	return result
}

// FromTrendEntity converts a PriceTrend entity to DTO
func FromTrendEntity(t *entity.PriceTrend) PriceTrendDTO {
	if t == nil {
		return PriceTrendDTO{}
	}
	return PriceTrendDTO{
		Date:  t.RecordDate.Format("2006-01-02"),
		Price: t.Price,
	}
}

// FromTrendEntities converts multiple PriceTrend entities to DTOs
func FromTrendEntities(trends []*entity.PriceTrend) []PriceTrendDTO {
	result := make([]PriceTrendDTO, 0, len(trends))
	for _, t := range trends {
		if t != nil {
			result = append(result, FromTrendEntity(t))
		}
	}
	return result
}

// FromNotificationEntity converts a NotificationConfig entity to DTO
func FromNotificationEntity(n *entity.NotificationConfig) NotificationDTO {
	if n == nil {
		return NotificationDTO{}
	}
	var lastNotifyTime *string
	if n.LastNotifyTime != nil {
		t := n.LastNotifyTime.Format("2006-01-02 15:04:05")
		lastNotifyTime = &t
	}

	return NotificationDTO{
		ActivityID:     n.ActivityID,
		TargetPrice:    n.TargetPrice,
		LastNotifyTime: lastNotifyTime,
	}
}
