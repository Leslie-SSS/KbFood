import { test, expect, setBarkKey } from "./fixtures";
import { NotificationModal } from "./pom/NotificationModal";
import { ConfirmModal } from "./pom/ConfirmModal";

/**
 * E2E Tests: Price Notification Flow
 *
 * Tests setting up price notifications/alerts on products.
 */

test.describe("Price Notification - Setting Up", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should open notification modal when clicking notification button", async ({
    page,
    mainPage,
  }) => {
    // Wait for products to load
    await mainPage.verifyProductsLoaded();

    // Click notification button on first product
    await mainPage.clickNotificationButton(0);

    // Notification modal should open
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Verify modal header
    await expect(page.locator("h2:has-text('价格提醒')")).toBeVisible();
  });

  test("should show product info in notification modal", async ({
    page,
    mainPage,
  }) => {
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Verify product info is shown
    await notificationModal.verifyProductInfo();
  });

  test("should show current price in modal", async ({ page, mainPage }) => {
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Current price should be visible
    const price = await notificationModal.getCurrentPrice();
    expect(price).toBeGreaterThan(0);
  });

  test("should have target price input", async ({ page, mainPage }) => {
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Target price input should be visible
    await expect(page.locator("input[type='number']")).toBeVisible();
  });

  test("should have price slider", async ({ page, mainPage }) => {
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Price slider should be visible
    await expect(page.locator("input[type='range']")).toBeVisible();
  });

  test("should have preset discount buttons", async ({ page, mainPage }) => {
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Preset buttons should be visible
    await expect(page.locator("button:has-text('9折')")).toBeVisible();
    await expect(page.locator("button:has-text('8折')")).toBeVisible();
    await expect(page.locator("button:has-text('7折')")).toBeVisible();
    await expect(page.locator("button:has-text('半价')")).toBeVisible();
  });
});

test.describe("Price Notification - Input Validation", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);
  });

  test("should disable submit button when target price is empty", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Submit should be disabled without target price
    await notificationModal.verifySubmitDisabled();
  });

  test("should disable submit button when target price equals current price", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Enter same price
    await notificationModal.enterTargetPrice(currentPrice);

    // Should show error
    expect(await notificationModal.hasError()).toBe(true);

    // Submit should be disabled
    await notificationModal.verifySubmitDisabled();
  });

  test("should disable submit button when target price is higher than current price", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Enter higher price
    await notificationModal.enterTargetPrice(currentPrice + 100);

    // Should show error
    expect(await notificationModal.hasError()).toBe(true);

    // Submit should be disabled
    await notificationModal.verifySubmitDisabled();
  });

  test("should disable submit button when target price is too low", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Enter very low price
    await notificationModal.enterTargetPrice(0.01);

    // Should show error
    expect(await notificationModal.hasError()).toBe(true);

    // Submit should be disabled
    await notificationModal.verifySubmitDisabled();
  });

  test("should enable submit button when target price is valid", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Enter valid target price (80% of current)
    await notificationModal.enterTargetPrice(currentPrice * 0.8);

    // Submit should be enabled
    await notificationModal.verifySubmitEnabled();
  });

  test("should show savings when valid price is entered", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Enter valid target price
    await notificationModal.enterTargetPrice(currentPrice * 0.8);

    // Should show savings
    expect(await notificationModal.hasSavingsDisplay()).toBe(true);
  });
});

test.describe("Price Notification - Preset Discounts", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);
  });

  test("should set target price when clicking 10% off (9折)", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Click 9折
    await notificationModal.click10PercentOff();

    // Input should show 90% of current price
    const inputValue = await page.locator("input[type='number']").inputValue();
    const expectedPrice = (currentPrice * 0.9).toFixed(2);
    expect(parseFloat(inputValue).toFixed(2)).toBe(expectedPrice);

    // Submit should be enabled
    await notificationModal.verifySubmitEnabled();
  });

  test("should set target price when clicking 20% off (8折)", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Click 8折
    await notificationModal.click20PercentOff();

    // Input should show 80% of current price
    const inputValue = await page.locator("input[type='number']").inputValue();
    const expectedPrice = (currentPrice * 0.8).toFixed(2);
    expect(parseFloat(inputValue).toFixed(2)).toBe(expectedPrice);
  });

  test("should set target price when clicking 30% off (7折)", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Click 7折
    await notificationModal.click30PercentOff();

    // Input should show 70% of current price
    const inputValue = await page.locator("input[type='number']").inputValue();
    const expectedPrice = (currentPrice * 0.7).toFixed(2);
    expect(parseFloat(inputValue).toFixed(2)).toBe(expectedPrice);
  });

  test("should set target price when clicking 50% off (半价)", async ({
    page,
  }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    const currentPrice = await notificationModal.getCurrentPrice();

    // Click 半价
    await notificationModal.click50PercentOff();

    // Input should show 50% of current price
    const inputValue = await page.locator("input[type='number']").inputValue();
    const expectedPrice = (currentPrice * 0.5).toFixed(2);
    expect(parseFloat(inputValue).toFixed(2)).toBe(expectedPrice);
  });

  test("should highlight selected preset", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Click 8折
    await notificationModal.click20PercentOff();

    // 8折 button should be highlighted
    const selectedButton = page.locator(
      "button.border-primary-400:has-text('8折')",
    );
    await expect(selectedButton).toBeVisible();
  });
});

test.describe("Price Notification - Price Slider", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);
  });

  test("should update target price when using slider", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Move slider to 80%
    await notificationModal.setPriceWithSlider(80);

    // Input should reflect slider value
    const inputValue = await page.locator("input[type='number']").inputValue();
    expect(parseFloat(inputValue)).toBeGreaterThan(0);
  });

  test("should show percentage display next to slider", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Move slider
    await notificationModal.setPriceWithSlider(70);

    // Should show percentage
    await expect(page.locator("text=70%")).toBeVisible();
  });
});

test.describe("Price Notification - Submission", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
    await mainPage.verifyProductsLoaded();
  });

  test("should close modal and show toast after successful submission", async ({
    page,
    mainPage,
  }) => {
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Set valid target price
    await notificationModal.click20PercentOff();

    // Submit
    await notificationModal.clickSubmit();

    // Modal should close
    await expect(page.locator("h2:has-text('价格提醒')")).not.toBeVisible({
      timeout: 10000,
    });

    // Toast should appear
    await expect(page.locator("text=监控已设置")).toBeVisible({
      timeout: 5000,
    });
  });

  test("should show loading state during submission", async ({
    page,
    mainPage,
  }) => {
    await mainPage.clickNotificationButton(0);

    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    // Set valid target price
    await notificationModal.click20PercentOff();

    // Submit
    await notificationModal.clickSubmit();

    // Should show loading state
    await expect(page.locator("text=保存中")).toBeVisible();
  });
});

test.describe("Price Notification - Cancel", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
    await mainPage.verifyProductsLoaded();
    await mainPage.clickNotificationButton(0);
  });

  test("should close modal when clicking cancel button", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    await notificationModal.clickCancel();

    // Modal should close
    await expect(page.locator("h2:has-text('价格提醒')")).not.toBeVisible({
      timeout: 5000,
    });
  });

  test("should close modal when clicking X button", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    await notificationModal.closeModal();

    // Modal should close
    await expect(page.locator("h2:has-text('价格提醒')")).not.toBeVisible({
      timeout: 5000,
    });
  });

  test("should close modal when clicking backdrop", async ({ page }) => {
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();

    await notificationModal.closeByBackdrop();

    // Modal should close
    await expect(page.locator("h2:has-text('价格提醒')")).not.toBeVisible({
      timeout: 5000,
    });
  });
});

test.describe("Price Notification - Existing Notification", () => {
  test.beforeEach(async ({ page, mainPage }) => {
    await setBarkKey(page, "test-key");
    await mainPage.gotoAndWaitForProducts();
  });

  test("should show cancel notification button when product has notification", async ({
    page,
    mainPage,
  }) => {
    await mainPage.verifyProductsLoaded();

    // Set up notification first
    await mainPage.clickNotificationButton(0);
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();
    await notificationModal.click20PercentOff();
    await notificationModal.clickSubmit();

    // Wait for toast
    await page.waitForTimeout(2000);

    // Button should now show "取消监控"
    await expect(
      page.locator("button:has-text('取消监控')").first(),
    ).toBeVisible();
  });

  test("should show confirm modal when canceling notification", async ({
    page,
    mainPage,
  }) => {
    await mainPage.verifyProductsLoaded();

    // Set up notification first
    await mainPage.clickNotificationButton(0);
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();
    await notificationModal.click20PercentOff();
    await notificationModal.clickSubmit();

    // Wait for toast
    await page.waitForTimeout(2000);

    // Click cancel notification
    await page.locator("button:has-text('取消监控')").first().click();

    // Confirm modal should appear
    const confirmModal = new ConfirmModal(page);
    await confirmModal.waitForModal();

    // Verify title
    const title = await confirmModal.getTitle();
    expect(title).toContain("取消监控");
  });

  test("should cancel notification when confirmed", async ({
    page,
    mainPage,
  }) => {
    await mainPage.verifyProductsLoaded();

    // Set up notification first
    await mainPage.clickNotificationButton(0);
    const notificationModal = new NotificationModal(page);
    await notificationModal.waitForModal();
    await notificationModal.click20PercentOff();
    await notificationModal.clickSubmit();

    // Wait for toast
    await page.waitForTimeout(2000);

    // Click cancel notification
    await page.locator("button:has-text('取消监控')").first().click();

    // Confirm
    const confirmModal = new ConfirmModal(page);
    await confirmModal.waitForModal();
    await confirmModal.clickConfirm();

    // Toast should show success
    await expect(page.locator("text=已取消监控")).toBeVisible({
      timeout: 5000,
    });
  });
});
