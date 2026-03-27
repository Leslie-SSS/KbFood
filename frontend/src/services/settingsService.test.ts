import { describe, expect, it } from "vitest";

import { settingsService } from "./settingsService";

describe("settingsService", () => {
  it("creates and reuses a stable client ID", () => {
    settingsService.clearClientId();

    const firstId = settingsService.getOrCreateClientId();
    const secondId = settingsService.getOrCreateClientId();

    expect(firstId).toBeTruthy();
    expect(secondId).toBe(firstId);
    expect(settingsService.getClientId()).toBe(firstId);
  });

  it("normalizes Bark URLs before storing them", () => {
    settingsService.clearBarkKey();
    settingsService.setBarkKey("https://api.day.app/ABC123/?isArchive=1#foo");

    expect(settingsService.getBarkKey()).toBe("ABC123");
    expect(settingsService.getLegacyUserId()).toBe("ABC123");
  });
});
