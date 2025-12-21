/**
 * Centralized configuration for the frontend application
 * Single source of truth for environment variables
 *
 * IMPORTANT: Next.js only replaces NEXT_PUBLIC_* env vars at build time
 * when directly referenced (e.g., process.env.NEXT_PUBLIC_API_URL).
 * Dynamic key access (process.env[key]) will NOT work!
 */

interface AppConfig {
  apiUrl: string;
  wsUrl: string;
}

/**
 * Load and validate configuration
 * Must directly reference process.env.NEXT_PUBLIC_* for Next.js to replace at build time
 */
function loadConfig(): AppConfig {
  // Direct references are required for Next.js build-time replacement
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';

  return {
    apiUrl,
    wsUrl,
  };
}

// Export the global config instance
export const config = loadConfig();
