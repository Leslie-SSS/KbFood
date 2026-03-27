import { afterEach, describe, expect, it } from "vitest";

import { api } from "./api";
import { settingsService } from "./settingsService";

describe("api request headers", () => {
  afterEach(() => {
    settingsService.clearBarkKey();
    settingsService.clearClientId();
    delete api.defaults.adapter;
  });

  it("sends the stable client ID and the legacy Bark-based user ID", async () => {
    settingsService.clearClientId();
    settingsService.setBarkKey("https://api.day.app/LEGACY123/");

    let capturedHeaders: Record<string, unknown> | undefined;
    api.defaults.adapter = async (config) => {
      capturedHeaders = config.headers as Record<string, unknown>;
      return {
        data: { ok: true },
        status: 200,
        statusText: "OK",
        headers: {},
        config,
      };
    };

    await api.get("/status");

    expect(capturedHeaders?.["X-User-ID"]).toBe(settingsService.getClientId());
    expect(capturedHeaders?.["X-Legacy-User-ID"]).toBe("LEGACY123");
  });

  it("omits the legacy header when no Bark key exists", async () => {
    settingsService.clearBarkKey();
    settingsService.clearClientId();

    let capturedHeaders: Record<string, unknown> | undefined;
    api.defaults.adapter = async (config) => {
      capturedHeaders = config.headers as Record<string, unknown>;
      return {
        data: { ok: true },
        status: 200,
        statusText: "OK",
        headers: {},
        config,
      };
    };

    await api.get("/status");

    expect(capturedHeaders?.["X-User-ID"]).toBe(settingsService.getClientId());
    expect(capturedHeaders?.["X-Legacy-User-ID"]).toBeUndefined();
  });
});
