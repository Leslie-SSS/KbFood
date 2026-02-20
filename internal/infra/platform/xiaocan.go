package platform

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	"kbfood/internal/config"
	"kbfood/internal/domain/entity"
)

// XiaoCanClient implements the XiaoCan platform client
type XiaoCanClient struct {
	cfg    *config.XiaoCanConfig
	client *resty.Client
}

// NewXiaoCanClient creates a new XiaoCan client
func NewXiaoCanClient(cfg *config.XiaoCanConfig) *XiaoCanClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15")

	return &XiaoCanClient{
		cfg:    cfg,
		client: client,
	}
}

// Name returns the platform name
func (c *XiaoCanClient) Name() string {
	return "小餐"
}

// ShouldFetch checks if fetching is allowed at this time
func (c *XiaoCanClient) ShouldFetch(now time.Time) bool {
	return true
}

// FetchProducts fetches products from XiaoCan
func (c *XiaoCanClient) FetchProducts(ctx context.Context, region string) ([]*entity.PlatformProductDTO, error) {
	// Implement XiaoCan API integration
	// This is a placeholder - actual implementation depends on XiaoCan's API spec
	return nil, nil
}

// setHeaders sets the required headers for XiaoCan API
func (c *XiaoCanClient) setHeaders(r *resty.Request) {
	r.SetHeader("x-vayne", c.cfg.XVayne)
	r.SetHeader("x-teemo", c.cfg.XTeemo)
	r.SetHeader("x-ashe", c.cfg.XAshe)
	r.SetHeader("x-nami", c.cfg.XNami)
	r.SetHeader("x-sivir", c.cfg.XSivir)
	r.SetHeader("user-id", c.cfg.UserID)
	r.SetHeader("silk-id", c.cfg.SilkID)
}
