import { Page, Locator, expect } from "@playwright/test";
import { BasePage } from "./BasePage";

/**
 * Page Object Model for the Confirm Modal
 * Used for confirming actions like block, delete notification
 */
export class ConfirmModal extends BasePage {
  // Modal container
  readonly modal: Locator;
  readonly title: Locator;
  readonly message: Locator;

  // Action buttons
  readonly confirmButton: Locator;
  readonly cancelButton: Locator;

  constructor(page: Page) {
    super(page);

    // Modal container - generic confirm modal
    this.modal = page.locator(".fixed.inset-0.z-50").filter({
      has: page.locator("button:has-text('确认')"),
    });
    this.title = this.modal.locator("h3, h2").first();
    this.message = this.modal.locator("p").first();

    // Action buttons
    this.confirmButton = this.modal.locator("button:has-text('确认')");
    this.cancelButton = this.modal.locator("button:has-text('取消')");
  }

  /**
   * Check if the confirm modal is visible
   */
  async isVisible(): Promise<boolean> {
    try {
      await this.modal.waitFor({ state: "visible", timeout: 3000 });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Wait for confirm modal to appear
   */
  async waitForModal() {
    await this.modal.waitFor({ state: "visible", timeout: 5000 });
  }

  /**
   * Get modal title
   */
  async getTitle(): Promise<string> {
    return (await this.title.textContent()) || "";
  }

  /**
   * Get modal message
   */
  async getMessage(): Promise<string> {
    return (await this.message.textContent()) || "";
  }

  /**
   * Click confirm button
   */
  async clickConfirm() {
    await this.confirmButton.click();
    // Wait for modal to close
    await this.modal.waitFor({ state: "hidden", timeout: 10000 });
  }

  /**
   * Click cancel button
   */
  async clickCancel() {
    await this.cancelButton.click();
    await this.modal.waitFor({ state: "hidden", timeout: 5000 });
  }

  /**
   * Verify modal content
   */
  async verifyModalContent(expectedTitle?: string, expectedMessage?: string) {
    await expect(this.modal).toBeVisible();
    if (expectedTitle) {
      await expect(this.title).toContainText(expectedTitle);
    }
    if (expectedMessage) {
      await expect(this.message).toContainText(expectedMessage);
    }
  }
}
