import axios from "axios";
import { settingsService } from "./settingsService";

// Auto-detect API URL based on environment
const getApiBaseUrl = (): string => {
  const { protocol, hostname, port } = window.location;

  // Production (no port, HTTPS) - use relative path through nginx proxy
  if (port === "" && protocol === "https:") {
    return "/api";
  }

  // Production on port 8200 (legacy)
  if (port === "8200") {
    return `${protocol}://${hostname}:8200/api`;
  }

  // Development with port 8000 (legacy)
  if (port === "8000") {
    return `${protocol}://${hostname}:9000/api`;
  }

  // Preview mode on ports 3000-3010 - connect to backend on 9000
  if (
    [
      "3000",
      "3001",
      "3002",
      "3003",
      "3004",
      "3005",
      "3006",
      "3007",
      "3008",
      "3009",
      "3010",
    ].includes(port)
  ) {
    return `${protocol}://${hostname}:9000/api`;
  }

  // Vite dev server - proxy to backend
  return "/api";
};

export const api = axios.create({
  baseURL: getApiBaseUrl(),
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor - add X-User-ID header
api.interceptors.request.use(
  (config) => {
    const barkKey = settingsService.getBarkKey();
    if (barkKey) {
      config.headers["X-User-ID"] = normalizeBarkKey(barkKey);
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error("API Error:", error);
    return Promise.reject(error);
  },
);

// Normalize Bark key - extract device key from URL if needed
function normalizeBarkKey(input: string): string {
  const trimmed = input.trim();
  if (trimmed.startsWith("http")) {
    const parts = trimmed.split("/");
    return parts[parts.length - 1] || trimmed;
  }
  return trimmed;
}
