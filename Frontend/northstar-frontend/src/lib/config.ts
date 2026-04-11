/**
 * Centralized config to prevent 'process.env' leakage 
 * and provide type-safe access to environment variables.
 */
export const CONFIG = {
  API_BASE_URL: import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1',
  APP_NAME: import.meta.env.VITE_APP_NAME || 'NorthStar',
  IS_DEV: import.meta.env.DEV,
} as const;