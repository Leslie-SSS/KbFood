import { Page, Locator, expect } from "@playwright/test";

/**
 * Base Page Object Model with common functionality
 */
export abstract class BasePage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  /**
   * Navigate to a specific path
   */
  async goto(path: string = "/") {
    await this.page.goto(path);
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to be fully loaded
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState("networkidle");
  }

  /**
   * Take a screenshot for debugging
   */
  async takeScreenshot(name: string) {
    await this.page.screenshot({ path: `artifacts/${name}.png` });
  }

  /**
   * Wait for a toast message to appear
   */
  async waitForToast(message?: string) {
    const toast = this.page.locator("[data-testid=toast]");
    await toast.waitFor({ state: "visible", timeout: 5000 });
    if (message) {
      await expect(toast).toContainText(message);
    }
    return toast;
  }

  /**
   * Clear localStorage (useful for test isolation)
   */
  async clearLocalStorage() {
    await this.page.evaluate(() => window.localStorage.clear());
  }

  /**
   * Set localStorage value
   */
  async setLocalStorage(key: string, value: string) {
    await this.page.evaluate(
      ({ k, v }) => window.localStorage.setItem(k, v),
      { k: key, v: value },
    );
  }

  /**
   * Get localStorage value
   */
  async getLocalStorage(key: string): Promise<string | null> {
    return this.page.evaluate((k) => window.localStorage.getItem(k), key);
  }
}
