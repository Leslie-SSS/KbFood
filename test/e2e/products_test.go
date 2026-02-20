// +build e2e

// Package e2e provides end-to-end tests for the Food-Go API
package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"kbfood/internal/domain/entity"
	"kbfood/internal/interface/http/dto"
	"kbfood/test"

	"github.com/labstack/echo/v4"
)

// ProductQueryResponse represents the API response for product queries
type ProductQueryResponse struct {
	Code    int             `json:"code"`
	Data    []dto.ProductDTO `json:"data"`
	Message string          `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TestQueryProducts_EmptyResult tests querying products when no products exist
func TestQueryProducts_EmptyResult(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	e := setupTestRouter(t)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute
	resp, err := http.Get(srv.URL + "/api/products")
	if err != nil {
		t.Fatalf("Failed to request products: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result ProductQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Code != http.StatusOK {
		t.Errorf("Expected code 200, got %d", result.Code)
	}

	if len(result.Data) != 0 {
		t.Errorf("Expected empty data array, got %d items", len(result.Data))
	}
}

// TestQueryProducts_WithKeyword tests querying products with a keyword filter
func TestQueryProducts_WithKeyword(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	// Create test products
	factory := test.NewProductFactory()
	p1 := factory.WithActivityID("act_001").Create()
	p1.Title = "巧克力草莓蛋糕"

	p2 := factory.WithActivityID("act_002").Create()
	p2.Title = "提拉米苏"

	// Insert products into database
	// ... (repository calls)

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute - search for "巧克力"
	u, _ := url.Parse(srv.URL + "/api/products")
	q := u.Query()
	q.Set("keyword", "巧克力")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatalf("Failed to request products: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	var result ProductQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 product, got %d", len(result.Data))
	}

	if !strings.Contains(result.Data[0].Title, "巧克力") {
		t.Errorf("Expected title to contain '巧克力', got %s", result.Data[0].Title)
	}
}

// TestQueryProducts_WithPlatform tests querying products with platform filter
func TestQueryProducts_WithPlatform(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	// Create test products
	factory := test.NewProductFactory()
	p1 := factory.WithActivityID("act_001").WithPlatform("DT").Create()
	p2 := factory.WithActivityID("act_002").WithPlatform("TTT").Create()

	// Insert products
	// ... (repository calls)

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute - filter by platform "DT"
	u, _ := url.Parse(srv.URL + "/api/products")
	q := u.Query()
	q.Set("platform", "DT")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatalf("Failed to request products: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	var result ProductQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Data) == 0 {
		t.Error("Expected at least one DT product")
	}

	for _, p := range result.Data {
		if p.Platform != "DT" {
			t.Errorf("Expected platform DT, got %s", p.Platform)
		}
	}
}

// TestQueryProducts_WithSalesStatus tests querying products with sales status filter
func TestQueryProducts_WithSalesStatus(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	// Create test products
	factory := test.NewProductFactory()
	p1 := factory.WithActivityID("act_001").WithSalesStatus(entity.SalesStatusOnSale).Create()
	p2 := factory.WithActivityID("act_002").WithSalesStatus(entity.SalesStatusSold).Create()

	// Insert products
	// ... (repository calls)

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute - filter by sales status "1" (on sale)
	u, _ := url.Parse(srv.URL + "/api/products")
	q := u.Query()
	q.Set("salesStatus", "1")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatalf("Failed to request products: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	var result ProductQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	for _, p := range result.Data {
		if p.SalesStatus != entity.SalesStatusOnSale {
			t.Errorf("Expected sales status %d, got %d", entity.SalesStatusOnSale, p.SalesStatus)
		}
	}
}

// TestGetPriceTrend tests retrieving price trend for a product
func TestGetPriceTrend(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	// Create test trends
	activityID := "test_activity_001"
	factory := test.NewTrendFactory()
	factory.ActivityID = activityID
	trends := factory.CreateList(7) // 7 days of trend data

	// Insert trends
	// ... (repository calls)

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute
	resp, err := http.Get(srv.URL + "/api/products/" + activityID + "/trend")
	if err != nil {
		t.Fatalf("Failed to request price trend: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		Code int                `json:"code"`
		Data []dto.PriceTrendDTO `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Data) != 7 {
		t.Errorf("Expected 7 trend points, got %d", len(result.Data))
	}
}

// TestBlockProduct tests blocking a product
func TestBlockProduct(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	activityID := "test_activity_001"

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute - block product
	resp, err := http.Post(srv.URL+"/api/products/"+activityID+"/block", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to block product: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify product is blocked by querying blocked products
	resp2, err := http.Get(srv.URL + "/api/products/blocked")
	if err != nil {
		t.Fatalf("Failed to request blocked products: %v", err)
	}
	defer resp2.Body.Close()

	var result ProductQueryResponse
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	found := false
	for _, p := range result.Data {
		if p.ActivityID == activityID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected blocked product to be in blocked list")
	}
}

// TestUnblockProduct tests unblocking a product
func TestUnblockProduct(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	activityID := "test_activity_001"

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// First block the product
	http.Post(srv.URL+"/api/products/"+activityID+"/block", "application/json", nil)

	// Execute - unblock product
	resp, err := http.Post(srv.URL+"/api/products/unblock/"+activityID, "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to unblock product: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify product is not in blocked list
	resp2, err := http.Get(srv.URL + "/api/products/blocked")
	if err != nil {
		t.Fatalf("Failed to request blocked products: %v", err)
	}
	defer resp2.Body.Close()

	var result ProductQueryResponse
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	for _, p := range result.Data {
		if p.ActivityID == activityID {
			t.Error("Expected product to be removed from blocked list")
		}
	}
}

// TestCreateNotification tests creating a price notification
func TestCreateNotification(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	e := setupTestRouterWithDB(t, testDB.Pool)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create notification request
	reqBody := map[string]interface{}{
		"activityId":  "test_activity_001",
		"targetPrice": 50.0,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/products/notifications", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify notification was created by querying with monitorStatus
	resp2, err := http.Get(srv.URL + "/api/products?monitorStatus=true")
	if err != nil {
		t.Fatalf("Failed to request products: %v", err)
	}
	defer resp2.Body.Close()

	var result ProductQueryResponse
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Find the product with notification
	for _, p := range result.Data {
		if p.ActivityID == "test_activity_001" {
			if !p.HasNotification {
				t.Error("Expected product to have notification")
			}
			if p.TargetPrice == nil || *p.TargetPrice != 50.0 {
				t.Errorf("Expected target price 50.0, got %v", p.TargetPrice)
			}
			return
		}
	}

	t.Error("Expected to find product with notification")
}

// TestCreateNotification_MissingActivityId tests validation error for missing activityId
func TestCreateNotification_MissingActivityId(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	e := setupTestRouter(t)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create notification request without activityId
	reqBody := map[string]interface{}{
		"targetPrice": 50.0,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/products/notifications", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	var result ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Message == "" {
		t.Error("Expected error message")
	}
}

// TestCreateNotification_InvalidTargetPrice tests validation error for invalid target price
func TestCreateNotification_InvalidTargetPrice(t *testing.T) {
	t.Skip("Skipping until full test infrastructure is set up")

	// Setup
	e := setupTestRouter(t)
	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create notification request with invalid target price
	reqBody := map[string]interface{}{
		"activityId":  "test_activity_001",
		"targetPrice": -10.0,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Execute
	resp, err := http.Post(srv.URL+"/api/products/notifications", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

// Helper functions

func setupTestRouter(t *testing.T) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	return e
}

func setupTestRouterWithDB(t *testing.T, pool interface{}) *echo.Echo {
	// This would create a full router with all handlers
	// For now, return a basic router
	e := echo.New()
	e.HideBanner = true
	return e
}
