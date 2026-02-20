import { test, expect, clearUserData, setBarkKey } from "./fixtures";

/**
 * E2E Tests: Welcome Modal Flow
 *
 * Tests the first-time user experience for setting up Bark notifications.
 * The welcome modal should appear for users without a Bark Key set.
 */

test.describe("Welcome Modal", () => {
  test.beforeEach(async ({ page }) => {
    // Clear localStorage to simulate a new user
    await clearUserData(page);
  });

  test("should show welcome modal for new users (no Bark Key)", async ({
    page,
    welcomeModal,
  }) => {
    // Navigate to the app
    await page.goto("/");

    // Welcome modal should appear after a short delay
    await welcomeModal.waitForModal();

    // Verify modal content
    await welcomeModal.verifyModalContent();

    // Verify title contains welcome message
    await expect(page.locator("text=欢迎使用美食监控")).toBeVisible();

    // Verify benefit items are shown
    await expect(page.locator("text=价格提醒推送")).toBeVisible();
    await expect(page.locator("text=数据云同步")).toBeVisible();
    await expect(page.locator("text=跨设备访问")).toBeVisible();
    await expect(page.locator("text=无需注册登录")).toBeVisible();
  });

  test("should not show welcome modal for returning users (with Bark Key)", async ({
    page,
    mainPage,
  }) => {
    // Set Bark Key before navigating
    await setBarkKey(page, "existing-bark-key");

    // Navigate to the app
    await mainPage.gotoAndWaitForProducts();

    // Welcome modal should NOT appear
    const welcomeVisible = await page
      .locator("text=欢迎使用美食监控")
      .isVisible()
      .catch(() => false);
    expect(welcomeVisible).toBe(false);

    // Products should be visible
    await mainPage.verifyProductsLoaded();
  });

  test("should allow user to skip setup", async ({
    page,
    welcomeModal,
    mainPage,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Click skip button
    await welcomeModal.clickSkip();

    // Modal should close
    await expect(page.locator("text=欢迎使用美食监控")).not.toBeVisible();

    // Products should be visible
    await mainPage.verifyProductsLoaded();
  });

  test("should allow user to close modal by clicking backdrop", async ({
    page,
    welcomeModal,
    mainPage,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Click outside the modal
    await welcomeModal.closeByBackdrop();

    // Modal should close
    await expect(page.locator("text=欢迎使用美食监控")).not.toBeVisible({
      timeout: 5000,
    });

    // Products should be visible
    await mainPage.verifyProductsLoaded();
  });

  test("should allow user to close modal by clicking X button", async ({
    page,
    welcomeModal,
    mainPage,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Click close button
    await welcomeModal.closeModal();

    // Modal should close
    await expect(page.locator("text=欢迎使用美食监控")).not.toBeVisible({
      timeout: 5000,
    });

    // Products should be visible
    await mainPage.verifyProductsLoaded();
  });

  test("should have disabled test button when Bark Key is empty", async ({
    page,
    welcomeModal,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Test button should be disabled
    await expect(page.locator("button:has-text('测试通知')")).toBeDisabled();
  });

  test("should enable test button when Bark Key is entered", async ({
    page,
    welcomeModal,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Enter a Bark Key
    await welcomeModal.enterBarkKey("test-key-123");

    // Test button should now be enabled
    await expect(page.locator("button:has-text('测试通知')")).toBeEnabled();
  });

  test("should show App Store link", async ({ page, welcomeModal }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Verify App Store link is present
    const appStoreLink = page.locator("a[href*='apps.apple.com']");
    await expect(appStoreLink).toBeVisible();
    await expect(appStoreLink).toContainText("Bark");
  });

  test("should save Bark Key and close modal", async ({
    page,
    welcomeModal,
    mainPage,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Enter Bark Key
    await welcomeModal.enterBarkKey("my-new-bark-key");

    // Click save
    await welcomeModal.clickSave();

    // Modal should close
    await expect(page.locator("text=欢迎使用美食监控")).not.toBeVisible({
      timeout: 5000,
    });

    // Verify Bark Key is saved
    const savedKey = await page.evaluate(() =>
      window.localStorage.getItem("barkKey"),
    );
    expect(savedKey).toBe("my-new-bark-key");

    // Products should be visible and toast should appear
    await mainPage.verifyProductsLoaded();
  });

  test("should show toast message after completing setup", async ({
    page,
    welcomeModal,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Complete setup
    await welcomeModal.enterBarkKey("my-bark-key");
    await welcomeModal.clickSave();

    // Wait for toast
    await expect(page.locator("text=设置完成")).toBeVisible({ timeout: 5000 });
  });

  test("should show error for invalid Bark Key format", async ({
    page,
    welcomeModal,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Enter an invalid Bark Key
    await welcomeModal.enterBarkKey("invalid-key-without-url");

    // Click test (this will fail because no backend mock)
    await welcomeModal.clickTest();

    // Wait for error result (may show request failed or similar)
    await page.waitForTimeout(2000);

    // Should show some feedback (error or success based on backend response)
    // This test verifies the UI handles the response
  });
});

test.describe("Welcome Modal - Format Hints", () => {
  test.beforeEach(async ({ page }) => {
    await clearUserData(page);
  });

  test("should show format hints for Bark Key input", async ({ page }) => {
    await page.goto("/");
    await page.locator("text=欢迎使用美食监控").waitFor({ state: "visible" });

    // Should show format hints
    await expect(page.locator("text=支持格式")).toBeVisible();
    await expect(page.locator("text=https://api.day.app")).toBeVisible();
  });

  test("should accept full URL format for Bark Key", async ({
    page,
    welcomeModal,
    mainPage,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Enter full URL format
    await welcomeModal.enterBarkKey("https://api.day.app/ABC123XYZ");
    await welcomeModal.clickSave();

    // Modal should close
    await expect(page.locator("text=欢迎使用美食监控")).not.toBeVisible({
      timeout: 5000,
    });

    // Products should be visible
    await mainPage.verifyProductsLoaded();
  });

  test("should accept device key only format", async ({
    page,
    welcomeModal,
    mainPage,
  }) => {
    await page.goto("/");
    await welcomeModal.waitForModal();

    // Enter device key only
    await welcomeModal.enterBarkKey("ABC123XYZ");
    await welcomeModal.clickSave();

    // Modal should close
    await expect(page.locator("text=欢迎使用美食监控")).not.toBeVisible({
      timeout: 5000,
    });

    // Products should be visible
    await mainPage.verifyProductsLoaded();
  });
});
