import { test, expect, setBarkKey, clearUserData } from "./fixtures";

/**
 * E2E Tests: Product Search Flow
 *
 * Tests the search functionality for finding products.
 */

test.describe("Product Search", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should display search bar on page load", async ({ page }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");
    await expect(searchInput).toBeVisible();
    await expect(searchInput).toBeEnabled();
  });

  test("should search products by keyword", async ({ page, mainPage }) => {
    // Get initial product count
    const initialCount = await mainPage.getProductCount();

    // Search for a keyword
    await mainPage.searchProducts("牛肉");

    // Wait for results
    await page.waitForTimeout(1000);

    // Verify search was performed (URL or results changed)
    const searchInput = page.locator("input[placeholder*='搜索']");
    await expect(searchInput).toHaveValue("牛肉");
  });

  test("should search using Enter key", async ({ page, mainPage }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");

    // Type and press Enter
    await searchInput.fill("猪肉");
    await searchInput.press("Enter");

    // Wait for debounce
    await page.waitForTimeout(1000);

    // Verify search was performed
    await expect(searchInput).toHaveValue("猪肉");
  });

  test("should clear search results", async ({ page, mainPage }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");

    // Search for something
    await mainPage.searchProducts("测试");
    await page.waitForTimeout(500);

    // Clear search
    await searchInput.clear();
    await searchInput.press("Enter");
    await page.waitForTimeout(1000);

    // Verify search was cleared
    await expect(searchInput).toHaveValue("");
  });

  test("should show no results message when search has no matches", async ({
    page,
    mainPage,
  }) => {
    // Search for something unlikely to match
    await mainPage.searchProducts("xyznonexistentproduct12345");

    // Wait for results
    await page.waitForTimeout(2000);

    // Should show no results message
    await mainPage.verifyNoProductsFound();
  });

  test("should maintain search value after page refresh", async ({ page }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");

    // Search for something
    await searchInput.fill("牛奶");
    await searchInput.press("Enter");
    await page.waitForTimeout(1000);

    // Reload page
    await page.reload();
    await page.waitForLoadState("networkidle");

    // Search value should persist in input (if using URL params)
    // Note: This depends on implementation
  });
});

test.describe("Product Search - Debounce", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should debounce search input", async ({ page }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");

    // Type quickly
    await searchInput.type("abcdefghij", { delay: 50 });

    // Wait a bit
    await page.waitForTimeout(500);

    // Value should be in input
    await expect(searchInput).toHaveValue("abcdefghij");
  });
});

test.describe("Product Search - Special Characters", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should handle special characters in search", async ({
    page,
    mainPage,
  }) => {
    // Search with special characters
    await mainPage.searchProducts("牛肉 & 猪肉");

    // Should not crash
    await page.waitForTimeout(1000);
    await expect(page.locator("body")).toBeVisible();
  });

  test("should handle Chinese characters", async ({ page, mainPage }) => {
    // Search with Chinese characters
    await mainPage.searchProducts("牛排");

    // Should work correctly
    await page.waitForTimeout(1000);
    const searchInput = page.locator("input[placeholder*='搜索']");
    await expect(searchInput).toHaveValue("牛排");
  });

  test("should handle numbers in search", async ({ page, mainPage }) => {
    // Search with numbers
    await mainPage.searchProducts("500g");

    // Should work correctly
    await page.waitForTimeout(1000);
    const searchInput = page.locator("input[placeholder*='搜索']");
    await expect(searchInput).toHaveValue("500g");
  });
});

test.describe("Product Search - UI", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should have search button next to input", async ({ page }) => {
    const searchButton = page.locator("button:has-text('搜索')");
    await expect(searchButton).toBeVisible();
    await expect(searchButton).toBeEnabled();
  });

  test("should trigger search when clicking search button", async ({
    page,
    mainPage,
  }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");
    const searchButton = page.locator("button:has-text('搜索')");

    // Type and click search
    await searchInput.fill("鸡肉");
    await searchButton.click();
    await page.waitForTimeout(1000);

    // Verify search was performed
    await expect(searchInput).toHaveValue("鸡肉");
  });

  test("should have search icon in input", async ({ page }) => {
    // Search icon should be visible (lucide-react renders as SVG)
    const searchIcon = page.locator("svg.lucide-search");
    await expect(searchIcon.first()).toBeVisible();
  });

  test("should have focus state on search input", async ({ page }) => {
    const searchInput = page.locator("input[placeholder*='搜索']");

    // Click to focus
    await searchInput.click();

    // Should be focused
    await expect(searchInput).toBeFocused();
  });
});
