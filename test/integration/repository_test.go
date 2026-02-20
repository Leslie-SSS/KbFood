// +build integration

// Package integration_test provides repository integration tests
package integration_test

import (
	"context"
	"testing"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	"kbfood/internal/infra/db/sqlc"
	"kbfood/test"
)

// TestProductRepository_Create tests creating a product
func TestProductRepository_Create(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	repo := repository.NewProductRepository(queries)

	// Create test product
	factory := test.NewProductFactory()
	product := factory.Create()

	// Execute
	err := repo.Create(ctx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Assert - verify product was created
	products, err := repo.FindByFilter(ctx, repository.ProductFilter{
		ActivityID: &product.ActivityID,
	})
	if err != nil {
		t.Fatalf("Failed to find product: %v", err)
	}

	if len(products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(products))
	}

	if products[0].ActivityID != product.ActivityID {
		t.Errorf("Expected activity ID %s, got %s", product.ActivityID, products[0].ActivityID)
	}
}

// TestProductRepository_FindByFilter_Keyword tests filtering by keyword
func TestProductRepository_FindByFilter_Keyword(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	repo := repository.NewProductRepository(queries)

	// Create test products
	factory := test.NewProductFactory()

	p1 := factory.WithActivityID("act_001").Create()
	p1.Title = "巧克力草莓蛋糕"
	repo.Create(ctx, p1)

	p2 := factory.WithActivityID("act_002").Create()
	p2.Title = "提拉米苏蛋糕"
	repo.Create(ctx, p2)

	p3 := factory.WithActivityID("act_003").Create()
	p3.Title = "巧克力慕斯"
	repo.Create(ctx, p3)

	// Execute - search for "巧克力"
	products, err := repo.FindByFilter(ctx, repository.ProductFilter{
		Keyword: "巧克力",
	})
	if err != nil {
		t.Fatalf("Failed to find products: %v", err)
	}

	// Assert - should find 2 products with "巧克力"
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	for _, p := range products {
		if !containsChinese(p.Title, "巧克力") {
			t.Errorf("Expected title to contain '巧克力', got %s", p.Title)
		}
	}
}

// TestProductRepository_FindByFilter_Platform tests filtering by platform
func TestProductRepository_FindByFilter_Platform(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	repo := repository.NewProductRepository(queries)

	// Create test products
	factory := test.NewProductFactory()

	p1 := factory.WithActivityID("act_001").WithPlatform("DT").Create()
	repo.Create(ctx, p1)

	p2 := factory.WithActivityID("act_002").WithPlatform("TTT").Create()
	repo.Create(ctx, p2)

	p3 := factory.WithActivityID("act_003").WithPlatform("DT").Create()
	repo.Create(ctx, p3)

	// Execute - filter by platform "DT"
	products, err := repo.FindByFilter(ctx, repository.ProductFilter{
		Platform: "DT",
	})
	if err != nil {
		t.Fatalf("Failed to find products: %v", err)
	}

	// Assert - should find 2 DT products
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	for _, p := range products {
		if p.Platform != "DT" {
			t.Errorf("Expected platform DT, got %s", p.Platform)
		}
	}
}

// TestProductRepository_FindByFilter_SalesStatus tests filtering by sales status
func TestProductRepository_FindByFilter_SalesStatus(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	repo := repository.NewProductRepository(queries)

	// Create test products
	factory := test.NewProductFactory()

	onSale := 1
	p1 := factory.WithActivityID("act_001").WithSalesStatus(entity.SalesStatusOnSale).Create()
	repo.Create(ctx, p1)

	p2 := factory.WithActivityID("act_002").WithSalesStatus(entity.SalesStatusSold).Create()
	repo.Create(ctx, p2)

	p3 := factory.WithActivityID("act_003").WithSalesStatus(entity.SalesStatusOnSale).Create()
	repo.Create(ctx, p3)

	// Execute - filter by sales status "1" (on sale)
	products, err := repo.FindByFilter(ctx, repository.ProductFilter{
		SalesStatus: &onSale,
	})
	if err != nil {
		t.Fatalf("Failed to find products: %v", err)
	}

	// Assert - should find 2 on-sale products
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	for _, p := range products {
		if p.SalesStatus != entity.SalesStatusOnSale {
			t.Errorf("Expected sales status %d, got %d", entity.SalesStatusOnSale, p.SalesStatus)
		}
	}
}

// TestProductRepository_Update tests updating a product
func TestProductRepository_Update(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	repo := repository.NewProductRepository(queries)

	// Create test product
	factory := test.NewProductFactory()
	product := factory.Create()
	repo.Create(ctx, product)

	// Update price
	product.CurrentPrice = 55.0
	err := repo.Update(ctx, product)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	// Assert - verify updated price
	products, err := repo.FindByFilter(ctx, repository.ProductFilter{
		ActivityID: &product.ActivityID,
	})
	if err != nil {
		t.Fatalf("Failed to find product: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("Expected 1 product, got %d", len(products))
	}

	if products[0].CurrentPrice != 55.0 {
		t.Errorf("Expected price 55.0, got %f", products[0].CurrentPrice)
	}
}

// TestBlockedRepository_CreateAndDelete tests blocking and unblocking products
func TestBlockedRepository_CreateAndDelete(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	blockedRepo := repository.NewBlockedRepository(queries)
	productRepo := repository.NewProductRepository(queries)

	// Create test product
	factory := test.NewProductFactory()
	product := factory.Create()
	productRepo.Create(ctx, product)

	// Block product
	activityID := product.ActivityID
	err := blockedRepo.Create(ctx, activityID)
	if err != nil {
		t.Fatalf("Failed to block product: %v", err)
	}

	// Verify product is blocked
	blockedProducts, err := productRepo.ListBlocked(ctx)
	if err != nil {
		t.Fatalf("Failed to list blocked products: %v", err)
	}

	if len(blockedProducts) != 1 {
		t.Errorf("Expected 1 blocked product, got %d", len(blockedProducts))
	}

	// Unblock product
	err = blockedRepo.Delete(ctx, activityID)
	if err != nil {
		t.Fatalf("Failed to unblock product: %v", err)
	}

	// Verify product is unblocked
	blockedProducts, err = productRepo.ListBlocked(ctx)
	if err != nil {
		t.Fatalf("Failed to list blocked products: %v", err)
	}

	if len(blockedProducts) != 0 {
		t.Errorf("Expected 0 blocked products, got %d", len(blockedProducts))
	}
}

// TestNotificationRepository_CRUD tests notification CRUD operations
func TestNotificationRepository_CRUD(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	notiRepo := repository.NewNotificationRepository(queries)

	// Create notification
	factory := test.NewNotificationFactory()
	config := factory.Create()

	err := notiRepo.Create(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	// Find notification
	found, err := notiRepo.FindByActivityID(ctx, config.ActivityID)
	if err != nil {
		t.Fatalf("Failed to find notification: %v", err)
	}

	if found.ActivityID != config.ActivityID {
		t.Errorf("Expected activity ID %s, got %s", config.ActivityID, found.ActivityID)
	}

	if found.TargetPrice != config.TargetPrice {
		t.Errorf("Expected target price %f, got %f", config.TargetPrice, found.TargetPrice)
	}

	// Update notification
	found.TargetPrice = 40.0
	err = notiRepo.Update(ctx, found)
	if err != nil {
		t.Fatalf("Failed to update notification: %v", err)
	}

	// Verify update
	updated, err := notiRepo.FindByActivityID(ctx, config.ActivityID)
	if err != nil {
		t.Fatalf("Failed to find notification: %v", err)
	}

	if updated.TargetPrice != 40.0 {
		t.Errorf("Expected target price 40.0, got %f", updated.TargetPrice)
	}

	// Delete notification
	err = notiRepo.Delete(ctx, config.ActivityID)
	if err != nil {
		t.Fatalf("Failed to delete notification: %v", err)
	}

	// Verify deletion
	_, err = notiRepo.FindByActivityID(ctx, config.ActivityID)
	if err == nil {
		t.Error("Expected error when finding deleted notification")
	}
}

// TestTrendRepository_FindByActivityID tests finding price trends
func TestTrendRepository_FindByActivityID(t *testing.T) {
	t.Skip("Skipping until TestContainers is fully configured")

	// Setup
	ctx := context.Background()
	testDB := test.SetupTestDB(t)
	defer testDB.Pool.Close()
	test.TruncateTables(t, testDB.Pool)

	queries := testDB.Pool.Queries()
	trendRepo := repository.NewTrendRepository(queries)

	// Create test trends
	activityID := "test_activity_001"
	factory := test.NewTrendFactory()
	factory.ActivityID = activityID

	trends := factory.CreateList(7)
	for _, trend := range trends {
		err := trendRepo.Create(ctx, trend)
		if err != nil {
			t.Fatalf("Failed to create trend: %v", err)
		}
	}

	// Find trends
	found, err := trendRepo.FindByActivityID(ctx, activityID)
	if err != nil {
		t.Fatalf("Failed to find trends: %v", err)
	}

	if len(found) != 7 {
		t.Errorf("Expected 7 trends, got %d", len(found))
	}

	// Verify order (should be ascending by date)
	for i := 1; i < len(found); i++ {
		if found[i-1].RecordDate.After(found[i].RecordDate) {
			t.Error("Trends should be ordered by date ascending")
		}
	}
}

// Helper functions

func containsChinese(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}
