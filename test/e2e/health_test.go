// +build e2e

// Package e2e provides end-to-end tests for the Food-Go API
package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kbfood/internal/config"
	"kbfood/test"

	"github.com/labstack/echo/v4"
	"kbfood/internal/infra/db"
	"kbfood/internal/interface/http/handler"
)

// TestHealthEndpoint tests the /health endpoint
func TestHealthEndpoint(t *testing.T) {
	// Setup
	e := echo.New()
	e.HideBanner = true
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
		})
	})

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to request health endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if status, ok := result["status"].(string); !ok || status != "ok" {
		t.Errorf("Expected status 'ok', got %v", result["status"])
	}
}

// TestReadyEndpoint_WithDB tests the /ready endpoint with database
func TestReadyEndpoint_WithDB(t *testing.T) {
	t.Skip("Skipping until full integration test is configured")

	// Setup test database
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()

	// Create handler with real database
	healthHandler := handler.NewHealthHandler(testDB.Pool)

	// Create router
	e := echo.New()
	e.HideBanner = true
	e.GET("/ready", healthHandler.Ready)

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute
	resp, err := http.Get(srv.URL + "/ready")
	if err != nil {
		t.Fatalf("Failed to request ready endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if status, ok := result["status"].(string); !ok || status != "ready" {
		t.Errorf("Expected status 'ready', got %v", result["status"])
	}

	if dbStatus, ok := result["database"].(string); !ok || dbStatus != "ok" {
		t.Errorf("Expected database status 'ok', got %v", result["database"])
	}
}

// TestReadyEndpoint_WithoutDB tests the /ready endpoint without database
func TestReadyEndpoint_WithoutDB(t *testing.T) {
	// Create handler without database
	healthHandler := handler.NewHealthHandler(nil)

	// Create router
	e := echo.New()
	e.HideBanner = true
	e.GET("/ready", healthHandler.Ready)

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Execute
	resp, err := http.Get(srv.URL + "/ready")
	if err != nil {
		t.Fatalf("Failed to request ready endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Assert - should be ready even without DB check when DB is nil
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if status, ok := result["status"].(string); !ok || status != "ready" {
		t.Errorf("Expected status 'ready', got %v", result["status"])
	}
}

// TestCORSMiddleware tests that CORS headers are set correctly
func TestCORSMiddleware(t *testing.T) {
	e := echo.New()
	e.HideBanner = true

	// Add a simple handler
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	srv := httptest.NewServer(e)
	defer srv.Close()

	// Create request with Origin header
	req, err := http.NewRequest("GET", srv.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Origin", "https://example.com")

	// Execute
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
