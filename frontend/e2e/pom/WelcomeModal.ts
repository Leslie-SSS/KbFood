import { Page, Locator, expect } from "@playwright/test";
import { BasePage } from "./BasePage";

/**
 * Page Object Model for the Welcome Modal
 * First-time user experience for setting up Bark notifications
 */
export class WelcomeModal extends BasePage {
  // Modal container
  readonly modal: Locator;
  readonly header: Locator;
  readonly closeButton: Locator;

  // Form elements
  readonly barkKeyInput: Locator;
  readonly testButton: Locator;
  readonly saveButton: Locator;
  readonly skipButton: Locator;

  // Feedback elements
  readonly testResult: Locator;
  readonly successMessage: Locator;
  readonly errorMessage: Locator;

  // Help link
  readonly appStoreLink: Locator;

  constructor(page: Page) {
    super(page);

    // Modal container
    this.modal = page.locator("text=欢迎使用美食监控").locator("..").locator("..");
    this.header = page.locator("text=欢迎使用美食监控");
    this.closeButton = page.locator("text=欢迎使用美食监控")
      .locator("..")
      .locator("button").first();

    // Form elements
    this.barkKeyInput = page.locator("input[placeholder*='URL']").or(
      page.locator("input[placeholder*='Key']"),
    );
    this.testButton = page.locator("button:has-text('测试通知')");
    this.saveButton = page.locator("button:has-text('完成设置')");
    this.skipButton = page.locator("button:has-text('暂不设置')");

    // Feedback elements
    this.testResult = page.locator("text=发送成功").or(page.locator("text=发送失败")).locator("..").locator("..");
    this.successMessage = page.locator("text=通知发送成功");
    this.errorMessage = page.locator("text=发送失败");

    // Help link
    this.appStoreLink = page.locator("a[href*='apps.apple.com']");
  }

  /**
   * Check if the welcome modal is visible
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
   * Wait for welcome modal to appear
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
   * Click skip button
   */
  async clickSkip() {
    await this.skipButton.click();
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
   * Set up Bark Key and save
   */
  async setupBarkKey(key: string) {
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
   * Skip the welcome setup
   */
  async skipSetup() {
    await this.clickSkip();
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
    await expect(this.skipButton).toBeVisible();
  }
}
