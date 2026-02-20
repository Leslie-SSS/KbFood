import { test as base, Page } from "@playwright/test";
import {
  MainPage,
  WelcomeModal,
  SettingsModal,
  NotificationModal,
  ConfirmModal,
  ProductCardPOM,
} from "../pom";

/**
 * Custom test fixtures for E2E tests
 * Provides pre-configured Page Object Models
 */

type Fixtures = {
  mainPage: MainPage;
  welcomeModal: WelcomeModal;
  settingsModal: SettingsModal;
  notificationModal: NotificationModal;
  confirmModal: ConfirmModal;
  authenticatedPage: Page;
};

export const test = base.extend<Fixtures>({
  // Main page fixture
  mainPage: async ({ page }, use) => {
    const mainPage = new MainPage(page);
    await use(mainPage);
  },

  // Welcome modal fixture
  welcomeModal: async ({ page }, use) => {
    const welcomeModal = new WelcomeModal(page);
    await use(welcomeModal);
  },

  // Settings modal fixture
  settingsModal: async ({ page }, use) => {
    const settingsModal = new SettingsModal(page);
    await use(settingsModal);
  },

  // Notification modal fixture
  notificationModal: async ({ page }, use) => {
    const notificationModal = new NotificationModal(page);
    await use(notificationModal);
  },

  // Confirm modal fixture
  confirmModal: async ({ page }, use) => {
    const confirmModal = new ConfirmModal(page);
    await use(confirmModal);
  },

  // Authenticated page (with Bark Key set)
  authenticatedPage: async ({ page }, use) => {
    // Set Bark Key in localStorage before navigating
    await page.addInitScript(() => {
      window.localStorage.setItem("barkKey", "test-bark-key-12345");
    });
    await use(page);
  },
});

/**
 * Helper to create a ProductCardPOM for a specific card
 */
export function createProductCard(page: Page, index: number): ProductCardPOM {
  const mainPage = new MainPage(page);
  return new ProductCardPOM(page, mainPage.productCards.nth(index));
}

/**
 * Helper to clear all user data (Bark Key, etc.)
 */
export async function clearUserData(page: Page) {
  await page.evaluate(() => {
    window.localStorage.clear();
  });
}

/**
 * Helper to set Bark Key in localStorage
 */
export async function setBarkKey(page: Page, key: string) {
  await page.evaluate((k) => {
    window.localStorage.setItem("barkKey", k);
  }, key);
}

/**
 * Helper to wait for products to load
 */
export async function waitForProducts(page: Page, minCount: number = 1) {
  const mainPage = new MainPage(page);
  await mainPage.gotoAndWaitForProducts();
  const count = await mainPage.getProductCount();
  if (count < minCount) {
    throw new Error(
      `Expected at least ${minCount} products, but found ${count}`,
    );
  }
}

export { expect } from "@playwright/test";
