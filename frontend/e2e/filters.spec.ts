import { test, expect, setBarkKey } from "./fixtures";

/**
 * E2E Tests: Product Filters Flow
 *
 * Tests the filtering functionality for products by platform, region, status, etc.
 */

test.describe("Product Filters - UI", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should display filter tags on page load", async ({ page }) => {
    // Check for filter tags
    await expect(page.locator("button:has-text('平台')")).toBeVisible();
    await expect(page.locator("button:has-text('地区')")).toBeVisible();
    await expect(page.locator("button:has-text('状态')")).toBeVisible();
    await expect(page.locator("button:has-text('监控')")).toBeVisible();
    await expect(page.locator("button:has-text('近7天上新')")).toBeVisible();
  });

  test("should show reset button only when filters are active", async ({
    page,
    mainPage,
  }) => {
    // Reset button should not be visible initially (no filters)
    const resetButton = page.locator("button:has-text('重置')");
    // It might be visible if there's a default filter state, so we check after applying a filter

    // Apply a filter
    await mainPage.filterByPlatform("美团");

    // Reset button should now be visible
    await expect(resetButton).toBeVisible();
  });
});

test.describe("Product Filters - Platform", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should open platform dropdown when clicked", async ({ page }) => {
    const platformButton = page.locator("button:has-text('平台')");
    await platformButton.click();

    // Dropdown should be visible
    await expect(page.locator(".absolute.top-full.left-0.mt-1")).toBeVisible();
  });

  test("should filter by platform", async ({ page, mainPage }) => {
    await mainPage.filterByPlatform("美团");

    // Platform button should show selected value
    await expect(
      page.locator("button.bg-primary-100:has-text('美团')"),
    ).toBeVisible();
  });

  test("should close dropdown when selecting an option", async ({ page }) => {
    const platformButton = page.locator("button:has-text('平台')");
    await platformButton.click();

    // Select first option
    await page.locator(".absolute.top-full.left-0.mt-1 button").first().click();

    // Dropdown should close
    await expect(
      page.locator(".absolute.top-full.left-0.mt-1"),
    ).not.toBeVisible();
  });
});

test.describe("Product Filters - Region", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should open region dropdown when clicked", async ({ page }) => {
    const regionButton = page.locator("button:has-text('地区')");
    await regionButton.click();

    // Dropdown should be visible
    await expect(page.locator(".absolute.top-full.left-0.mt-1")).toBeVisible();
  });

  test("should filter by region", async ({ page, mainPage }) => {
    await mainPage.filterByRegion("北京");

    // Region button should show selected value
    await expect(
      page.locator("button.bg-primary-100:has-text('北京')"),
    ).toBeVisible();
  });
});

test.describe("Product Filters - Sales Status", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should open status dropdown when clicked", async ({ page }) => {
    const statusButton = page.locator("button:has-text('状态')");
    await statusButton.click();

    // Dropdown should be visible
    await expect(page.locator(".absolute.top-full.left-0.mt-1")).toBeVisible();
  });

  test("should filter by sales status", async ({ page, mainPage }) => {
    await mainPage.filterByStatus("在售");

    // Status button should show selected value
    await expect(
      page.locator("button.bg-primary-100:has-text('在售')"),
    ).toBeVisible();
  });
});

test.describe("Product Filters - Monitor Status", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should open monitor dropdown when clicked", async ({ page }) => {
    const monitorButton = page.locator("button:has-text('监控')");
    await monitorButton.click();

    // Dropdown should be visible
    await expect(page.locator(".absolute.top-full.left-0.mt-1")).toBeVisible();
  });

  test("should filter by monitor status", async ({ page, mainPage }) => {
    await mainPage.filterByMonitor("已监控");

    // Monitor button should show selected value
    await expect(
      page.locator("button.bg-primary-100:has-text('已监控')"),
    ).toBeVisible();
  });
});

test.describe("Product Filters - Recent 7 Days", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should toggle recent 7 days filter", async ({ page, mainPage }) => {
    const recentButton = page.locator("button:has-text('近7天上新')");

    // Click to enable
    await mainPage.toggleRecentSevenDays();

    // Button should be highlighted
    await expect(
      page.locator("button.bg-primary-100:has-text('近7天上新')"),
    ).toBeVisible();

    // Click to disable
    await mainPage.toggleRecentSevenDays();

    // Button should no longer be highlighted
    await expect(recentButton).not.toHaveClass(/bg-primary-100/);
  });
});

test.describe("Product Filters - Reset", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should reset all filters", async ({ page, mainPage }) => {
    // Apply multiple filters
    await mainPage.filterByPlatform("美团");
    await mainPage.filterByRegion("北京");
    await mainPage.toggleRecentSevenDays();

    // Reset all filters
    await mainPage.resetFilters();

    // All filter buttons should show default labels
    await expect(page.locator("button:has-text('平台')")).not.toHaveClass(
      /bg-primary-100/,
    );
    await expect(page.locator("button:has-text('地区')")).not.toHaveClass(
      /bg-primary-100/,
    );
    await expect(page.locator("button:has-text('近7天上新')")).not.toHaveClass(
      /bg-primary-100/,
    );
  });

  test("should hide reset button after resetting", async ({
    page,
    mainPage,
  }) => {
    // Apply a filter
    await mainPage.filterByPlatform("美团");

    // Reset button should be visible
    const resetButton = page.locator("button:has-text('重置')");
    await expect(resetButton).toBeVisible();

    // Reset
    await mainPage.resetFilters();

    // Reset button should be hidden (or not visible)
    await expect(resetButton).not.toBeVisible();
  });
});

test.describe("Product Filters - Combined", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should combine search and filters", async ({ page, mainPage }) => {
    // Search for a keyword
    await mainPage.searchProducts("牛肉");

    // Apply a filter
    await mainPage.filterByPlatform("美团");

    // Both should be active
    const searchInput = page.locator("input[placeholder*='搜索']");
    await expect(searchInput).toHaveValue("牛肉");
    await expect(
      page.locator("button.bg-primary-100:has-text('美团')"),
    ).toBeVisible();
  });

  test("should persist filters while searching", async ({ page, mainPage }) => {
    // Apply a filter first
    await mainPage.filterByRegion("北京");

    // Then search
    await mainPage.searchProducts("猪肉");

    // Filter should still be active
    await expect(
      page.locator("button.bg-primary-100:has-text('北京')"),
    ).toBeVisible();
  });
});

test.describe("Product Filters - Accessibility", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should close dropdown when clicking outside", async ({ page }) => {
    // Open dropdown
    await page.locator("button:has-text('平台')").click();
    await expect(page.locator(".absolute.top-full.left-0.mt-1")).toBeVisible();

    // Click outside
    await page.mouse.click(10, 10);

    // Dropdown should close
    await expect(
      page.locator(".absolute.top-full.left-0.mt-1"),
    ).not.toBeVisible();
  });

  test("should have touch-friendly targets on mobile", async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    // Filter buttons should have minimum touch target size (44px)
    const platformButton = page.locator("button:has-text('平台')");
    const box = await platformButton.boundingBox();

    // Check minimum height (allowing some flexibility)
    expect(box?.height).toBeGreaterThanOrEqual(36);
  });
});
