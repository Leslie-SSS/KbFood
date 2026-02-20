// Package test provides test data factories for E2E and integration tests
package test

import (
	"time"

	"kbfood/internal/domain/entity"
)

// ProductFactory creates test Product entities
type ProductFactory struct {
	ID         int64
	ActivityID string
	Platform   string
	Region     string
	Title      string
	ShopName   string
}

// NewProductFactory creates a new product factory with defaults
func NewProductFactory() *ProductFactory {
	return &ProductFactory{
		ID:         1,
		ActivityID: "test_activity_001",
		Platform:   "DT",
		Region:     "华北",
		Title:      "巧克力草莓蛋糕",
		ShopName:   "测试烘焙店",
	}
}

// Create creates a new Product with factory defaults
func (f *ProductFactory) Create() *entity.Product {
	return &entity.Product{
		ID:                 f.ID,
		ActivityID:         f.ActivityID,
		Platform:           f.Platform,
		Region:             f.Region,
		Title:              f.Title,
		ShopName:           f.ShopName,
		OriginalPrice:      100.0,
		CurrentPrice:       70.0,
		SalesStatus:        entity.SalesStatusOnSale,
		ActivityCreateTime: time.Now().Add(-24 * time.Hour),
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
	}
}

// WithActivityID sets a custom activity ID
func (f *ProductFactory) WithActivityID(id string) *entity.Product {
	p := f.Create()
	p.ActivityID = id
	return p
}

// WithPlatform sets a custom platform
func (f *ProductFactory) WithPlatform(platform string) *entity.Product {
	p := f.Create()
	p.Platform = platform
	return p
}

// WithPrices sets custom original and current prices
func (f *ProductFactory) WithPrices(original, current float64) *entity.Product {
	p := f.Create()
	p.OriginalPrice = original
	p.CurrentPrice = current
	return p
}

// WithSalesStatus sets a custom sales status
func (f *ProductFactory) WithSalesStatus(status int) *entity.Product {
	p := f.Create()
	p.SalesStatus = status
	return p
}

// MasterProductFactory creates test MasterProduct entities
type MasterProductFactory struct {
	ID            string
	Region        string
	StandardTitle string
	Price         float64
	Status        int
	TrustScore    int
}

// NewMasterProductFactory creates a new master product factory
func NewMasterProductFactory() *MasterProductFactory {
	return &MasterProductFactory{
		ID:            "master_001",
		Region:        "广州",
		StandardTitle: "巧克力草莓蛋糕",
		Price:         70.0,
		Status:        entity.SalesStatusOnSale,
		TrustScore:    5,
	}
}

// Create creates a new MasterProduct
func (f *MasterProductFactory) Create() *entity.MasterProduct {
	return &entity.MasterProduct{
		ID:            f.ID,
		Region:        f.Region,
		StandardTitle: f.StandardTitle,
		Price:         f.Price,
		Status:        f.Status,
		TrustScore:    f.TrustScore,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
}

// NotificationFactory creates test NotificationConfig entities
type NotificationFactory struct {
	ActivityID  string
	TargetPrice float64
}

// NewNotificationFactory creates a new notification factory
func NewNotificationFactory() *NotificationFactory {
	return &NotificationFactory{
		ActivityID:  "test_activity_001",
		TargetPrice: 50.0,
	}
}

// Create creates a new NotificationConfig
func (f *NotificationFactory) Create() *entity.NotificationConfig {
	return &entity.NotificationConfig{
		ActivityID:  f.ActivityID,
		TargetPrice: f.TargetPrice,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
}

// WithActivityID sets a custom activity ID
func (f *NotificationFactory) WithActivityID(id string) *entity.NotificationConfig {
	n := f.Create()
	n.ActivityID = id
	return n
}

// TrendFactory creates test PriceTrend entities
type TrendFactory struct {
	ActivityID string
	Date       time.Time
	Price      float64
}

// NewTrendFactory creates a new trend factory
func NewTrendFactory() *TrendFactory {
	return &TrendFactory{
		ActivityID: "test_activity_001",
		Date:       time.Now(),
		Price:      70.0,
	}
}

// Create creates a new PriceTrend
func (f *TrendFactory) Create() *entity.PriceTrend {
	return &entity.PriceTrend{
		ActivityID: f.ActivityID,
		RecordDate: f.Date,
		Price:      f.Price,
		CreateTime: time.Now(),
	}
}

// CreateList creates a list of PriceTrend entities with dates in the past
func (f *TrendFactory) CreateList(count int) []*entity.PriceTrend {
	trends := make([]*entity.PriceTrend, count)
	for i := 0; i < count; i++ {
		date := time.Now().AddDate(0, 0, -(count - i - 1))
		price := 100.0 - float64(i)*5
		trends[i] = &entity.PriceTrend{
			ActivityID: f.ActivityID,
			RecordDate: date,
			Price:      price,
			CreateTime: time.Now(),
		}
	}
	return trends
}
