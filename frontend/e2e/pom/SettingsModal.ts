import { Page, Locator, expect } from "@playwright/test";
import { BasePage } from "./BasePage";

/**
 * Page Object Model for the Settings Modal
 * Configure Bark Key for notifications
 */
export class SettingsModal extends BasePage {
  // Modal container
  readonly modal: Locator;
  readonly header: Locator;
  readonly closeButton: Locator;

  // Form elements
  readonly barkKeyInput: Locator;
  readonly clearInputButton: Locator;
  readonly testButton: Locator;
  readonly saveButton: Locator;
  readonly cancelButton: Locator;

  // Feedback elements
  readonly testResult: Locator;
  readonly successMessage: Locator;
  readonly errorMessage: Locator;

  // Help link
  readonly appStoreLink: Locator;

  constructor(page: Page) {
    super(page);

    // Modal container - find by unique header text
    this.modal = page.locator("text=通知设置").locator("..").locator("..");
    this.header = page.locator("h2:has-text('通知设置')");
    this.closeButton = this.header.locator("..").locator("button").first();

    // Form elements
    this.barkKeyInput = page.locator("input#barkKey").or(
      page.locator("input[placeholder*='URL']").or(
        page.locator("input[placeholder*='Key']")
      )
    );
    this.clearInputButton = this.barkKeyInput.locator("..").locator("button");
    this.testButton = page.locator("button:has-text('测试')");
    this.saveButton = page.locator("button:has-text('保存设置')");
    this.cancelButton = page.locator("button:has-text('取消')");

    // Feedback elements
    this.testResult = page.locator("div.bg-gradient-to-r.from-green-50, div.bg-gradient-to-r.from-red-50");
    this.successMessage = page.locator("text=发送成功");
    this.errorMessage = page.locator("text=发送失败");

    // Help link
    this.appStoreLink = page.locator("a[href*='apps.apple.com']");
  }

  /**
   * Check if the settings modal is visible
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
   * Wait for settings modal to appear
   */
  async waitForModal() {
    await this.header.waitFor({ state: "visible", timeout: 5000 });
  }

  /**
   * Enter Bark Key
   */
  async enterBarkKey(key: string) {
    await this.barkKeyInput.fill(key);
  }

  /**
   * Clear Bark Key input
   */
  async clearBarkKey() {
    await this.barkKeyInput.clear();
  }

  /**
   * Get current Bark Key value
   */
  async getBarkKeyValue(): Promise<string> {
    return await this.barkKeyInput.inputValue();
  }

  /**
   * Click test notification button
   */
  async clickTest() {
    await this.testButton.click();
  }

  /**
   * Click save button
   */
  async clickSave() {
    await this.saveButton.click();
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
   * Close modal by clicking backdrop
   */
  async closeByBackdrop() {
    // Click outside the modal (on backdrop)
    await this.page.mouse.click(10, 10);
  }

  /**
   * Wait for test result
   */
  async waitForTestResult(expectedSuccess: boolean) {
    const text = expectedSuccess ? "发送成功" : "发送失败";
    await this.page.locator(`text=${text}`).waitFor({ state: "visible", timeout: 15000 });
  }

  /**
   * Save Bark Key
   */
  async saveBarkKey(key: string) {
    await this.enterBarkKey(key);
    await this.clickSave();
    // Wait for modal to close
    await this.modal.waitFor({ state: "hidden", timeout: 5000 });
  }

  /**
   * Test notification with Bark Key
   */
  async testNotification(key: string): Promise<boolean> {
    await this.enterBarkKey(key);
    await this.clickTest();
    try {
      await this.waitForTestResult(true);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Clear Bark Key and save (removes key)
   */
  async clearBarkKeyAndSave() {
    await this.clearBarkKey();
    await this.clickSave();
    await this.modal.waitFor({ state: "hidden", timeout: 5000 });
  }

  /**
   * Verify modal content
   */
  async verifyModalContent() {
    await expect(this.header).toBeVisible();
    await expect(this.barkKeyInput).toBeVisible();
    await expect(this.testButton).toBeVisible();
    await expect(this.saveButton).toBeVisible();
    await expect(this.cancelButton).toBeVisible();
  }
}
