import { Page, Locator, expect } from "@playwright/test";
import { BasePage } from "./BasePage";

/**
 * Page Object Model for a Product Card
 * Individual product display with actions
 */
export class ProductCardPOM extends BasePage {
  readonly card: Locator;
  readonly title: Locator;
  readonly price: Locator;
  readonly originalPrice: Locator;
  readonly metaInfo: Locator;
  readonly statusBadge: Locator;
  readonly notificationBadge: Locator;
  readonly targetPriceInfo: Locator;
  readonly progressBar: Locator;
  readonly notificationButton: Locator;
  readonly blockButton: Locator;
  readonly chart: Locator;

  constructor(page: Page, cardLocator: Locator) {
    super(page);
    this.card = cardLocator;

    // Card elements
    this.title = cardLocator.locator("h3");
    this.price = cardLocator.locator("span.text-xl, span.text-2xl").first();
    this.originalPrice = cardLocator.locator("span.line-through");
    this.metaInfo = cardLocator.locator(".flex.gap-1\\.5.text-xs");
    this.statusBadge = cardLocator.locator("span:has-text('已售')");
    this.notificationBadge = cardLocator.locator("span:has-text('监控中'), span:has-text('已达标')");
    this.targetPriceInfo = cardLocator.locator("text=目标价").locator("..");
    this.progressBar = cardLocator.locator(".h-1\\.5.bg-slate-100");

    // Action buttons
    this.notificationButton = cardLocator.locator("button:has-text('提醒'), button:has-text('取消监控')");
    this.blockButton = cardLocator.locator("button:has-text('屏蔽')");

    // Chart
    this.chart = cardLocator.locator("canvas");
  }

  /**
   * Get product title
   */
  async getTitle(): Promise<string> {
    return (await this.title.textContent()) || "";
  }

  /**
   * Get current price
   */
  async getPrice(): Promise<number> {
    const priceText = await this.price.textContent();
    const match = priceText?.match(/¥?(\d+\.?\d*)/);
    return match ? parseFloat(match[1]) : 0;
  }

  /**
   * Check if product has notification set
   */
  async hasNotification(): Promise<boolean> {
    try {
      await this.notificationBadge.waitFor({ state: "visible", timeout: 500 });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Check if product is sold out
   */
  async isSoldOut(): Promise<boolean> {
    try {
      await this.statusBadge.waitFor({ state: "visible", timeout: 500 });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Click notification button
   */
  async clickNotification() {
    await this.notificationButton.click();
  }

  /**
   * Click block button
   */
  async clickBlock() {
    await this.blockButton.click();
  }

  /**
   * Click on target price info (to edit)
   */
  async clickTargetPriceInfo() {
    await this.targetPriceInfo.click();
  }

  /**
   * Verify card is visible
   */
  async verifyVisible() {
    await expect(this.card).toBeVisible();
    await expect(this.title).toBeVisible();
    await expect(this.price).toBeVisible();
  }

  /**
   * Verify notification badge shows "监控中"
   */
  async verifyMonitoring() {
    await expect(this.notificationBadge).toContainText("监控中");
  }

  /**
   * Verify notification badge shows "已达标"
   */
  async verifyTargetReached() {
    await expect(this.notificationBadge).toContainText("已达标");
  }

  /**
   * Get target price if set
   */
  async getTargetPrice(): Promise<number | null> {
    try {
      const text = await this.targetPriceInfo.textContent();
      const match = text?.match(/目标价:\s*¥(\d+\.?\d*)/);
      return match ? parseFloat(match[1]) : null;
    } catch {
      return null;
    }
  }

  /**
   * Wait for chart to load
   */
  async waitForChart() {
    await this.chart.waitFor({ state: "visible", timeout: 5000 });
  }
}
