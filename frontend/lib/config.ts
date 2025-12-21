/**
 * Centralized configuration for the frontend application
 * Single source of truth for environment variables
 */

interface AppConfig {
  apiUrl: string;
  wsUrl: string;
}

/**
 * Get environment variable with fallback
 */
function getEnv(key: string, defaultValue: string): string {
  if (typeof window === 'undefined') {
    // Server-side: use process.env
    return process.env[key] || defaultValue;
  }
  // Client-side: use process.env (injected at build time by Next.js)
  return process.env[key] || defaultValue;
}

/**
 * Load and validate configuration
 */
function loadConfig(): AppConfig {
  const apiUrl = getEnv('NEXT_PUBLIC_API_URL', 'http://localhost:8080');
  const wsUrl = getEnv('NEXT_PUBLIC_WS_URL', 'ws://localhost:8080/ws');

  return {
    apiUrl,
    wsUrl,
  };
}

// Export the global config instance
export const config = loadConfig();
