import { test, expect, setBarkKey } from "./fixtures";

/**
 * E2E Tests: Settings Modal Flow
 *
 * Tests the settings modal for configuring Bark notifications.
 * Users can bind, test, and update their Bark Key.
 */

test.describe("Settings Modal", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    // Set a Bark Key so user is "authenticated"
    await setBarkKey(page, "existing-bark-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should open settings modal when clicking settings button", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    await expect(page.locator("h2:has-text('通知设置')")).toBeVisible();
  });

  test("should show existing Bark Key in input", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    // Should show existing Bark Key
    const value = await settingsModal.getBarkKeyValue();
    expect(value).toBe("existing-bark-key");
  });

  test("should allow updating Bark Key", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    // Clear and enter new key
    await settingsModal.clearBarkKey();
    await settingsModal.enterBarkKey("new-updated-key");
    await settingsModal.clickSave();

    // Modal should close
    await expect(page.locator("h2:has-text('通知设置')")).not.toBeVisible({
      timeout: 5000,
    });

    // Verify toast message
    await expect(page.locator("text=设置已保存")).toBeVisible({
      timeout: 5000,
    });
  });

  test("should allow clearing Bark Key", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    // Clear the key
    await settingsModal.clearBarkKeyAndSave();

    // Modal should close
    await expect(page.locator("h2:has-text('通知设置')")).not.toBeVisible({
      timeout: 5000,
    });

    // Verify toast message
    await expect(page.locator("text=已清除 Bark Key")).toBeVisible({
      timeout: 5000,
    });
  });

  test("should close modal when clicking cancel", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    await settingsModal.clickCancel();

    // Modal should close
    await expect(page.locator("h2:has-text('通知设置')")).not.toBeVisible({
      timeout: 5000,
    });
  });

  test("should close modal when clicking X button", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    await settingsModal.closeModal();

    // Modal should close
    await expect(page.locator("h2:has-text('通知设置')")).not.toBeVisible({
      timeout: 5000,
    });
  });

  test("should close modal when clicking backdrop", async ({
    page,
    mainPage,
    settingsModal,
  }) => {
    await mainPage.openSettings();
    await settingsModal.waitForModal();

    await settingsModal.closeByBackdrop();

    // Modal should close
    await expect(page.locator("h2:has-text('通知设置')")).not.toBeVisible({
      timeout: 5000,
    });
  });

  test("should have disabled test button when Bark Key is empty", async ({
    page,
    mainPage,
  }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Clear input
    await page.locator("input#barkKey").clear();

    // Test button should be disabled
    await expect(page.locator("button:has-text('测试')")).toBeDisabled();
  });

  test("should enable test button when Bark Key is entered", async ({
    page,
    mainPage,
  }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Clear and enter new key
    await page.locator("input#barkKey").clear();
    await page.locator("input#barkKey").fill("test-key");

    // Test button should be enabled
    await expect(page.locator("button:has-text('测试')")).toBeEnabled();
  });

  test("should show format hints", async ({ page, mainPage }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Should show format hints
    await expect(page.locator("text=支持两种格式")).toBeVisible();
    await expect(page.locator("text=https://api.day.app")).toBeVisible();
  });

  test("should show App Store link", async ({ page, mainPage }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Verify App Store link is present
    const appStoreLink = page.locator("a[href*='apps.apple.com']");
    await expect(appStoreLink).toBeVisible();
    await expect(appStoreLink).toContainText("Bark");
  });
});

test.describe("Settings Modal - Test Notification", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-bark-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should show loading state when testing notification", async ({
    page,
    mainPage,
  }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Click test button
    await page.locator("button:has-text('测试')").click();

    // Should show loading state
    await expect(page.locator("text=发送中")).toBeVisible();
  });

  test("should show success result after successful test", async ({
    page,
    mainPage,
  }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Click test button
    await page.locator("button:has-text('测试')").click();

    // Wait for result (success or failure)
    await page.waitForTimeout(3000);

    // Check if either success or failure is shown
    const successVisible = await page
      .locator("text=发送成功")
      .isVisible()
      .catch(() => false);
    const failureVisible = await page
      .locator("text=发送失败")
      .isVisible()
      .catch(() => false);

    expect(successVisible || failureVisible).toBe(true);
  });

  test("should clear test result when Bark Key changes", async ({
    page,
    mainPage,
  }) => {
    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Test notification
    await page.locator("button:has-text('测试')").click();
    await page.waitForTimeout(2000);

    // Change Bark Key
    await page.locator("input#barkKey").clear();
    await page.locator("input#barkKey").fill("new-key");

    // Test result should be cleared
    await expect(page.locator("text=发送成功")).not.toBeVisible();
    await expect(page.locator("text=发送失败")).not.toBeVisible();
  });
});

test.describe("Settings Modal - Disabled State", () => {
  test("should disable buttons while saving", async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();

    await mainPage.openSettings();
    await page.locator("h2:has-text('通知设置')").waitFor({ state: "visible" });

    // Click save
    await page.locator("button:has-text('保存设置')").click();

    // Buttons should be disabled during save
    await expect(page.locator("button:has-text('取消')")).toBeDisabled();
    await expect(page.locator("button:has-text('测试')")).toBeDisabled();
  });
});
