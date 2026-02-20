package external

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"kbfood/internal/domain/entity"
)

// BarkClient handles Bark push notifications
type BarkClient struct {
	client *resty.Client
	baseURL string
	deviceKey string
}

// NewBarkClient creates a new Bark client
func NewBarkClient(deviceKey string) *BarkClient {
	return &BarkClient{
		client: resty.New().
			SetTimeout(10 * time.Second),
		baseURL:   "https://api.day.app",
		deviceKey: deviceKey,
	}
}

// NewBarkClientWithURL creates a new Bark client with custom URL
func NewBarkClientWithURL(baseURL, deviceKey string) *BarkClient {
	return &BarkClient{
		client: resty.New().
			SetTimeout(10 * time.Second),
		baseURL:   baseURL,
		deviceKey: deviceKey,
	}
}

// Send sends a push notification
func (b *BarkClient) Send(ctx context.Context, title, body string, level string) error {
	url := fmt.Sprintf("%s/%s/%s/%s", b.baseURL, b.deviceKey,
		url.QueryEscape(title), url.QueryEscape(body))

	if level != "" {
		url += fmt.Sprintf("?level=%s", level)
	}

	resp, err := b.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return fmt.Errorf("send bark notification: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("bark notification failed: status %d", resp.StatusCode())
	}

	return nil
}

// SendPriceAlert sends a price alert notification
func (b *BarkClient) SendPriceAlert(ctx context.Context, product *entity.Product, targetPrice float64) error {
	title := "价格提醒"
	body := fmt.Sprintf("【%s %s ¥%.2f】%s\n目标价格: ¥%.2f",
		product.Platform,
		product.Region,
		product.CurrentPrice,
		product.Title,
		targetPrice,
	)

	return b.Send(ctx, title, body, "critical")
}

// SendSimple sends a simple notification
func (b *BarkClient) SendSimple(ctx context.Context, message string) error {
	return b.Send(ctx, "kbFood", message, "")
}
