import { Page, Locator, expect } from "@playwright/test";
import { BasePage } from "./BasePage";

/**
 * Page Object Model for the main application page
 * Contains product listing, search, filters, and navigation
 */
export class MainPage extends BasePage {
  // Locators for main components
  readonly navbar: Locator;
  readonly searchBar: Locator;
  readonly searchInput: Locator;
  readonly searchButton: Locator;
  readonly productGrid: Locator;
  readonly productCards: Locator;
  readonly productCount: Locator;
  readonly filterContainer: Locator;
  readonly resetFiltersButton: Locator;

  // Filter tags
  readonly platformFilter: Locator;
  readonly regionFilter: Locator;
  readonly statusFilter: Locator;
  readonly monitorFilter: Locator;
  readonly recentSevenDaysButton: Locator;

  // Navbar buttons
  readonly settingsButton: Locator;
  readonly blockedButton: Locator;

  constructor(page: Page) {
    super(page);

    // Navbar
    this.navbar = page.locator("nav");
    this.settingsButton = page.locator("button[title=]");
    this.blockedButton = page.locator("button:has-text('屏蔽')").first();

    // Search
    this.searchBar = page.locator("input[placeholder*='搜索']");
    this.searchInput = page.locator("input[placeholder*='搜索']");
    this.searchButton = page.locator("button:has-text('搜索')");

    // Product grid
    this.productGrid = page.locator(
      ".grid.grid-cols-1.sm\\:grid-cols-2.lg\\:grid-cols-3",
    );
    this.productCards = page.locator(
      ".bg-white.rounded-xl.p-3.sm\\:p-4.border.border-slate-200",
    );
    this.productCount = page.locator("text=/共 \\d+ 个产品/");

    // Filters
    this.filterContainer = page.locator(
      ".py-2.sm\\:py-2\\.5.bg-slate-50.border-b",
    );
    this.resetFiltersButton = page.locator("button:has-text('重置')");

    // Filter dropdowns (these open dropdowns when clicked)
    this.platformFilter = page.locator("button:has-text('平台')");
    this.regionFilter = page.locator("button:has-text('地区')");
    this.statusFilter = page.locator("button:has-text('状态')");
    this.monitorFilter = page.locator("button:has-text('监控')");
    this.recentSevenDaysButton = page.locator("button:has-text('近7天上新')");

    // Navbar buttons (using title attribute)
    this.settingsButton = page.locator("button[title='设置']");
    this.blockedButton = page.locator("button[title='屏蔽列表']");
  }

  /**
   * Navigate to the main page and wait for products to load
   */
  async gotoAndWaitForProducts() {
    await this.goto("/");
    // Wait for products to load (either skeleton or actual products)
    await this.page.waitForSelector(
      ".grid.grid-cols-1.sm\\:grid-cols-2.lg\\:grid-cols-3",
      { timeout: 15000 },
    );
  }

  /**
   * Search for products by keyword
   */
  async searchProducts(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.searchButton.click();
    // Wait for API response
    await this.page.waitForResponse(
      (resp) =>
        resp.url().includes("/api/products") && resp.status() === 200,
      { timeout: 10000 },
    );
    await this.waitForPageLoad();
  }

  /**
   * Search using Enter key (triggers debounced search)
   */
  async searchWithEnter(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.searchInput.press("Enter");
    // Wait for debounce and API response
    await this.page.waitForTimeout(500);
    await this.page.waitForResponse(
      (resp) =>
        resp.url().includes("/api/products") && resp.status() === 200,
      { timeout: 10000 },
    ).catch(() => {
      // Response might already have happened, continue
    });
    await this.waitForPageLoad();
  }

  /**
   * Get the number of visible product cards
   */
  async getProductCount(): Promise<number> {
    return await this.productCards.count();
  }

  /**
   * Click on a product card by index
   */
  async clickProductCard(index: number) {
    await this.productCards.nth(index).click();
  }

  /**
   * Get the first product card
   */
  getFirstProductCard(): Locator {
    return this.productCards.first();
  }

  /**
   * Select a filter option from a dropdown
   */
  async selectFilterOption(
    filterButton: Locator,
    optionText: string,
  ) {
    await filterButton.click();
    // Wait for dropdown to open
    await this.page.waitForSelector(".absolute.top-full.left-0.mt-1", {
      state: "visible",
    });
    // Click the option
    await this.page
      .locator(`.absolute.top-full.left-0.mt-1 button:has-text("${optionText}")`)
      .click();
    // Wait for API response
    await this.page.waitForResponse(
      (resp) =>
        resp.url().includes("/api/products") && resp.status() === 200,
      { timeout: 10000 },
    ).catch(() => {});
    await this.waitForPageLoad();
  }

  /**
   * Filter by platform
   */
  async filterByPlatform(platform: string) {
    await this.selectFilterOption(this.platformFilter, platform);
  }

  /**
   * Filter by region
   */
  async filterByRegion(region: string) {
    await this.selectFilterOption(this.regionFilter, region);
  }

  /**
   * Filter by sales status
   */
  async filterByStatus(status: string) {
    await this.selectFilterOption(this.statusFilter, status);
  }

  /**
   * Filter by monitor status
   */
  async filterByMonitor(status: string) {
    await this.selectFilterOption(this.monitorFilter, status);
  }

  /**
   * Toggle "Recent 7 Days" filter
   */
  async toggleRecentSevenDays() {
    await this.recentSevenDaysButton.click();
    await this.page.waitForResponse(
      (resp) =>
        resp.url().includes("/api/products") && resp.status() === 200,
      { timeout: 10000 },
    ).catch(() => {});
    await this.waitForPageLoad();
  }

  /**
   * Reset all filters
   */
  async resetFilters() {
    await this.resetFiltersButton.click();
    await this.page.waitForResponse(
      (resp) =>
        resp.url().includes("/api/products") && resp.status() === 200,
      { timeout: 10000 },
    ).catch(() => {});
    await this.waitForPageLoad();
  }

  /**
   * Open settings modal
   */
  async openSettings() {
    await this.settingsButton.click();
    await this.page.waitForSelector("text=通知设置", { state: "visible" });
  }

  /**
   * Open blocked products modal
   */
  async openBlockedProducts() {
    await this.blockedButton.click();
    await this.page.waitForSelector("text=屏蔽列表", { state: "visible" });
  }

  /**
   * Verify products are displayed
   */
  async verifyProductsLoaded() {
    await expect(this.productCards.first()).toBeVisible({ timeout: 15000 });
  }

  /**
   * Verify no products found state
   */
  async verifyNoProductsFound() {
    await expect(
      this.page.locator("text=未找到相关产品"),
    ).toBeVisible();
  }

  /**
   * Get a product card notification button
   */
  getProductNotificationButton(cardIndex: number): Locator {
    return this.productCards
      .nth(cardIndex)
      .locator("button:has-text('提醒')");
  }

  /**
   * Get a product card block button
   */
  getProductBlockButton(cardIndex: number): Locator {
    return this.productCards.nth(cardIndex).locator("button:has-text('屏蔽')");
  }

  /**
   * Click notification button on a product card
   */
  async clickNotificationButton(cardIndex: number) {
    await this.productCards
      .nth(cardIndex)
      .locator("button:has-text('提醒')")
      .click();
  }

  /**
   * Click block button on a product card
   */
  async clickBlockButton(cardIndex: number) {
    await this.productCards
      .nth(cardIndex)
      .locator("button:has-text('屏蔽')")
      .click();
  }
}
