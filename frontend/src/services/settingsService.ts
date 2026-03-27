const BARK_KEY_STORAGE_KEY = "barkKey";
const CLIENT_ID_STORAGE_KEY = "clientId";

function normalizeBarkKey(input: string): string {
  const trimmed = input.trim();
  if (!trimmed) {
    return "";
  }

  if (trimmed.toLowerCase().startsWith("http")) {
    try {
      const parsed = new URL(trimmed);
      const path = parsed.pathname.replace(/^\/+|\/+$/g, "");
      if (!path) return "";

      const parts = path.split("/");
      for (let i = parts.length - 1; i >= 0; i -= 1) {
        const segment = parts[i].trim();
        if (segment) {
          return segment;
        }
      }
      return "";
    } catch {
      const fallback = trimmed.replace(/^\/+|\/+$/g, "");
      const parts = fallback.split("/");
      for (let i = parts.length - 1; i >= 0; i -= 1) {
        const segment = parts[i].trim();
        if (segment) {
          return segment;
        }
      }
      return "";
    }
  }

  return trimmed;
}

function createClientId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }

  return `client_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`;
}

export const settingsService = {
  // Get Bark key from localStorage
  getBarkKey: (): string => {
    return localStorage.getItem(BARK_KEY_STORAGE_KEY) || "";
  },

  // Save Bark key to localStorage
  setBarkKey: (key: string): void => {
    const normalizedKey = normalizeBarkKey(key);

    if (normalizedKey) {
      localStorage.setItem(BARK_KEY_STORAGE_KEY, normalizedKey);
    } else {
      localStorage.removeItem(BARK_KEY_STORAGE_KEY);
    }
  },

  // Clear Bark key from localStorage
  clearBarkKey: (): void => {
    localStorage.removeItem(BARK_KEY_STORAGE_KEY);
  },

  // Get the stable client ID from localStorage
  getClientId: (): string => {
    return localStorage.getItem(CLIENT_ID_STORAGE_KEY) || "";
  },

  // Get or create the stable client ID used as the API user identifier
  getOrCreateClientId: (): string => {
    const existingId = localStorage.getItem(CLIENT_ID_STORAGE_KEY);
    if (existingId) {
      return existingId;
    }

    const clientId = createClientId();
    localStorage.setItem(CLIENT_ID_STORAGE_KEY, clientId);
    return clientId;
  },

  // Clear the stable client ID, primarily for tests
  clearClientId: (): void => {
    localStorage.removeItem(CLIENT_ID_STORAGE_KEY);
  },

  // Get the legacy Bark-key-based user ID for one-time backend migration
  getLegacyUserId: (): string => {
    const barkKey = localStorage.getItem(BARK_KEY_STORAGE_KEY) || "";
    return normalizeBarkKey(barkKey);
  },

  // Normalize Bark key for reuse across services
  normalizeBarkKey: (input: string): string => {
    return normalizeBarkKey(input);
  },
};
