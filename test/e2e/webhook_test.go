// +build e2e

// Package e2e provides end-to-end tests for the Food-Go API
package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kbfood/internal/interface/http/handler"
	"kbfood/test"

	"github.com/labstack/echo/v4"
)

// DTWebhookRequest represents a DT platform webhook request
type DTWebhookRequest struct {
	Items []DTWebhookItem `json:"items"`
}

// DTWebhookItem represents a single item in DT webhook
type DTWebhookItem struct {
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	Status    int     `json:"status"`
	CrawlTime int64   `json:"crawlTime"`
	Region    string  `json:"region"`
}

// WebhookResponse represents the webhook response
type WebhookResponse struct {
	Code    int    `json:"code"`
	Data    struct {
		Received int `json:"received"`
		Promoted int `json:"promoted"`
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

// TestHandleDTPush_Success tests successful DT webhook processing
func TestHandleDTPush_Success(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	e := setupWebhookRouter(t)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create webhook request
	reqBody := DTWebhookRequest{
		Items: []DTWebhookItem{
			{
				Title:     "巧克力草莓蛋糕",
				Price:     68.0,
				Status:    1,
				CrawlTime: 1706745600,
				Region:    "华北",
			},
			{
				Title:     "提拉米苏蛋糕",
				Price:     45.0,
				Status:    1,
				CrawlTime: 1706745600,
				Region:    "华东",
			},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Code != http.StatusOK {
		t.Errorf("Expected code 200, got %d", result.Code)
	}

	if result.Data.Received != 2 {
		t.Errorf("Expected received 2, got %d", result.Data.Received)
	}

	// Promoted count depends on whether items match master products
	// For new items, promoted should be 0 initially
	if result.Data.Promoted < 0 || result.Data.Promoted > 2 {
		t.Errorf("Expected promoted between 0 and 2, got %d", result.Data.Promoted)
	}
}

// TestHandleDTPush_EmptyItems tests validation error for empty items array
func TestHandleDTPush_EmptyItems(t *testing.T) {
	// Setup
	cleaningService := &mockDataCleaningService{}
	h := handler.NewExternalHandler(cleaningService)

	e := echo.New()
	e.HideBanner = true
	e.POST("/api/external/dt/push", h.HandleDTPush)

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create webhook request with empty items
	reqBody := DTWebhookRequest{
		Items: []DTWebhookItem{},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["code"].(float64) != 400 {
		t.Errorf("Expected code 400, got %v", result["code"])
	}
}

// TestHandleDTPush_TooManyItems tests validation error for too many items
func TestHandleDTPush_TooManyItems(t *testing.T) {
	// Setup
	cleaningService := &mockDataCleaningService{}
	h := handler.NewExternalHandler(cleaningService)

	e := echo.New()
	e.HideBanner = true
	e.POST("/api/external/dt/push", h.HandleDTPush)

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create webhook request with 1001 items (over the limit)
	items := make([]DTWebhookItem, 1001)
	for i := range items {
		items[i] = DTWebhookItem{
			Title:     "Test Product",
			Price:     10.0,
			Status:    1,
			CrawlTime: 1706745600,
			Region:    "华北",
		}
	}

	reqBody := DTWebhookRequest{Items: items}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

// TestHandleDTPush_InvalidJSON tests validation error for invalid JSON
func TestHandleDTPush_InvalidJSON(t *testing.T) {
	// Setup
	cleaningService := &mockDataCleaningService{}
	h := handler.NewExternalHandler(cleaningService)

	e := echo.New()
	e.HideBanner = true
	e.POST("/api/external/dt/push", h.HandleDTPush)

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute with invalid JSON
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader([]byte("invalid json")))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

// TestHandleDTPush_SingleItem tests processing a single item
func TestHandleDTPush_SingleItem(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	e := setupWebhookRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create webhook request with single item
	reqBody := DTWebhookRequest{
		Items: []DTWebhookItem{
			{
				Title:     "巧克力草莓蛋糕",
				Price:     68.0,
				Status:    1,
				CrawlTime: 1706745600,
				Region:    "华北",
			},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Data.Received != 1 {
		t.Errorf("Expected received 1, got %d", result.Data.Received)
	}
}

// TestHandleDTPush_MaxItems tests processing the maximum allowed items (1000)
func TestHandleDTPush_MaxItems(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	e := setupWebhookRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create webhook request with 1000 items (at the limit)
	items := make([]DTWebhookItem, 1000)
	for i := range items {
		items[i] = DTWebhookItem{
			Title:     "Test Product",
			Price:     10.0,
			Status:    1,
			CrawlTime: 1706745600,
			Region:    "华北",
		}
	}

	reqBody := DTWebhookRequest{Items: items}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert - should succeed
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestHandleDTPush_WithSoldOutItem tests processing items with sold status
func TestHandleDTPush_WithSoldOutItem(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	e := setupWebhookRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create webhook request with sold out item
	reqBody := DTWebhookRequest{
		Items: []DTWebhookItem{
			{
				Title:     "巧克力草莓蛋糕",
				Price:     68.0,
				Status:    0, // Sold out
				CrawlTime: 1706745600,
				Region:    "华北",
			},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/external/dt/push", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to post webhook: %v", err)
	}
	defer resp.Body.Close()

	// Assert - should still succeed
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// Helper functions

func setupWebhookRouter(t *testing.T) *echo.Echo {
	cleaningService := &mockDataCleaningService{}
	h := handler.NewExternalHandler(cleaningService)

	e := echo.New()
	e.HideBanner = true
	e.POST("/api/external/dt/push", h.HandleDTPush)
	return e
}

func setupWebhookRouterWithDB(t *testing.T, pool interface{}) *echo.Echo {
	// This would create a full router with real cleaning service
	return setupWebhookRouter(t)
}

// mockDataCleaningService is a mock implementation for testing
type mockDataCleaningService struct{}

func (m *mockDataCleaningService) ProcessIncomingItem(ctx interface{}, input interface{}, region string) (interface{}, error) {
	// Mock implementation - always returns success
	return nil, nil
}
