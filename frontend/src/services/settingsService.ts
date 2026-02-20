const BARK_KEY_STORAGE_KEY = 'barkKey';

export const settingsService = {
  // Get Bark key from localStorage
  getBarkKey: (): string => {
    return localStorage.getItem(BARK_KEY_STORAGE_KEY) || '';
  },

  // Save Bark key to localStorage
  setBarkKey: (key: string): void => {
    if (key) {
      localStorage.setItem(BARK_KEY_STORAGE_KEY, key);
    } else {
      localStorage.removeItem(BARK_KEY_STORAGE_KEY);
    }
  },

  // Clear Bark key from localStorage
  clearBarkKey: (): void => {
    localStorage.removeItem(BARK_KEY_STORAGE_KEY);
  },
};
