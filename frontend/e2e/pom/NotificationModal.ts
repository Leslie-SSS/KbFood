import { Page, Locator, expect } from "@playwright/test";
import { BasePage } from "./BasePage";

/**
 * Page Object Model for the Notification Modal
 * Set price alerts on products
 */
export class NotificationModal extends BasePage {
  // Modal container
  readonly modal: Locator;
  readonly header: Locator;
  readonly closeButton: Locator;

  // Product info
  readonly productInfo: Locator;
  readonly currentPrice: Locator;
  readonly productTitle: Locator;

  // Form elements
  readonly targetPriceInput: Locator;
  readonly priceSlider: Locator;
  readonly presetButtons: Locator;

  // Action buttons
  readonly submitButton: Locator;
  readonly cancelButton: Locator;

  // Feedback elements
  readonly savingsDisplay: Locator;
  readonly discountPercent: Locator;
  readonly errorDisplay: Locator;

  constructor(page: Page) {
    super(page);

    // Modal container - find by unique header text
    this.modal = page.locator("h2:has-text('价格提醒'), h2:has-text('修改提醒')").locator("..").locator("..");
    this.header = page.locator("h2:has-text('价格提醒'), h2:has-text('修改提醒')");
    this.closeButton = this.header.locator("..").locator("button").first();

    // Product info
    this.productInfo = page.locator("p:has-text('商品')").locator("..");
    this.currentPrice = page.locator("p:has-text('现价')").locator("..").locator("p.text-xl");
    this.productTitle = page.locator("p:has-text('商品')").locator("..").locator("p.line-clamp-1");

    // Form elements
    this.targetPriceInput = page.locator("input[type='number']");
    this.priceSlider = page.locator("input[type='range']");
    this.presetButtons = page.locator("button:has-text('折'), button:has-text('半价')");

    // Action buttons
    this.submitButton = page.locator("button:has-text('开启提醒'), button:has-text('更新提醒')");
    this.cancelButton = page.locator("button:has-text('取消')");

    // Feedback elements
    this.savingsDisplay = page.locator("text=预计省下").locator("..");
    this.discountPercent = page.locator("span:has-text('%')");
    this.errorDisplay = page.locator("text=请输入有效价格, text=目标价格需低于当前价格, text=目标价格过低");
  }

  /**
   * Check if the notification modal is visible
   */
  async isVisible(): Promise<boolean> {
    try {
      await this.header.waitFor({ state: "visible", timeout: 3000 });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Wait for notification modal to appear
   */
  async waitForModal() {
    await this.header.waitFor({ state: "visible", timeout: 5000 });
  }

  /**
   * Get current price from modal
   */
  async getCurrentPrice(): Promise<number> {
    const priceText = await this.currentPrice.textContent();
    // Extract number from "¥XX.XX"
    const match = priceText?.match(/¥?(\d+\.?\d*)/);
    return match ? parseFloat(match[1]) : 0;
  }

  /**
   * Get product title from modal
   */
  async getProductTitle(): Promise<string> {
    return (await this.productTitle.textContent()) || "";
  }

  /**
   * Enter target price
   */
  async enterTargetPrice(price: number) {
    await this.targetPriceInput.fill(price.toString());
    // Wait for validation
    await this.page.waitForTimeout(100);
  }

  /**
   * Set price using slider
   */
  async setPriceWithSlider(percent: number) {
    await this.priceSlider.fill(percent.toString());
  }

  /**
   * Click a preset discount button
   */
  async clickPreset(presetLabel: string) {
    await this.page.locator(`button:has-text('${presetLabel}')`).click();
  }

  /**
   * Click 10% off (9折)
   */
  async click10PercentOff() {
    await this.clickPreset("9折");
  }

  /**
   * Click 20% off (8折)
   */
  async click20PercentOff() {
    await this.clickPreset("8折");
  }

  /**
   * Click 30% off (7折)
   */
  async click30PercentOff() {
    await this.clickPreset("7折");
  }

  /**
   * Click 50% off (半价)
   */
  async click50PercentOff() {
    await this.clickPreset("半价");
  }

  /**
   * Click submit button
   */
  async clickSubmit() {
    await this.submitButton.click();
  }

  /**
   * Click cancel button
   */
  async clickCancel() {
    await this.cancelButton.click();
  }

  /**
   * Close modal by clicking X button
   */
  async closeModal() {
    await this.closeButton.click();
  }

  /**
   * Check if submit button is enabled
   */
  async isSubmitEnabled(): Promise<boolean> {
    return await this.submitButton.isEnabled();
  }

  /**
   * Check if error is displayed
   */
  async hasError(): Promise<boolean> {
    try {
      await this.errorDisplay.waitFor({ state: "visible", timeout: 500 });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Check if savings display is visible
   */
  async hasSavingsDisplay(): Promise<boolean> {
    try {
      await this.savingsDisplay.waitFor({ state: "visible", timeout: 500 });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Set notification with target price
   */
  async setNotification(targetPrice: number) {
    await this.enterTargetPrice(targetPrice);
    await this.clickSubmit();
    // Wait for modal to close
    await this.modal.waitFor({ state: "hidden", timeout: 10000 });
  }

  /**
   * Set notification using preset discount
   */
  async setNotificationWithPreset(presetLabel: string) {
    await this.clickPreset(presetLabel);
    await this.clickSubmit();
    await this.modal.waitFor({ state: "hidden", timeout: 10000 });
  }

  /**
   * Verify modal shows product info
   */
  async verifyProductInfo() {
    await expect(this.productTitle).toBeVisible();
    await expect(this.currentPrice).toBeVisible();
  }

  /**
   * Verify submit button is disabled (invalid price)
   */
  async verifySubmitDisabled() {
    await expect(this.submitButton).toBeDisabled();
  }

  /**
   * Verify submit button is enabled (valid price)
   */
  async verifySubmitEnabled() {
    await expect(this.submitButton).toBeEnabled();
  }
}
