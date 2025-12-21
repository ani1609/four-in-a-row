// WebSocket client for real-time game communication

import type { WSMessage } from './types';

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';

// Session storage keys
const SESSION_KEY = 'game_session';

/**
 * Game session data stored in localStorage
 */
export interface GameSession {
  playerId: string;
  username: string;
  gameId: string;
  disconnectTime?: number; // Timestamp when user disconnected
}

// ============= Session Management =============

/**
 * Store a game session in localStorage
 * @param session - The game session to store
 */
export function storeSession(session: GameSession): void {
  if (typeof window !== 'undefined') {
    const sessionWithTime = {
      ...session,
      disconnectTime: session.disconnectTime || Date.now(),
    };
    localStorage.setItem(SESSION_KEY, JSON.stringify(sessionWithTime));
  }
}

/**
 * Retrieve the stored game session from localStorage
 * @returns The stored game session or null if not found
 */
export function getStoredSession(): GameSession | null {
  if (typeof window === 'undefined') return null;

  const stored = localStorage.getItem(SESSION_KEY);
  if (!stored) return null;

  try {
    return JSON.parse(stored);
  } catch {
    return null;
  }
}

/**
 * Clear the stored game session from localStorage
 */
export function clearSession(): void {
  if (typeof window !== 'undefined') {
    localStorage.removeItem(SESSION_KEY);
  }
}

/**
 * Check if stored session is expired (5 minute timeout)
 * @returns true if session exists and is not expired
 */
export function isSessionValid(): boolean {
  const session = getStoredSession();
  if (!session) return false;

  const FIVE_MINUTES = 5 * 60 * 1000;
  const disconnectTime = session.disconnectTime || Date.now();
  const elapsed = Date.now() - disconnectTime;

  return elapsed < FIVE_MINUTES;
}

// ============= WebSocket Connection =============

/**
 * Create a new WebSocket connection
 * @returns A new WebSocket instance
 */
export function createWebSocket(): WebSocket {
  return new WebSocket(WS_URL);
}

/**
 * Send a message through a WebSocket connection
 * @param ws - The WebSocket instance
 * @param message - The message to send
 */
export function sendMessage(
  ws: WebSocket | null,
  message: { type: string; payload: unknown }
): void {
  if (!ws) {
    console.warn('WebSocket is null. Message not sent:', message);
    return;
  }

  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(message));
  } else {
    console.warn('WebSocket not open. Message not sent:', message);
  }
}

/**
 * Parse a raw WebSocket message string into a typed message object
 * @param data - The raw message data string
 * @returns The parsed message or null if parsing fails
 */
export function parseMessage(data: string): WSMessage | null {
  try {
    return JSON.parse(data) as WSMessage;
  } catch (error) {
    console.error('Failed to parse WebSocket message:', error);
    return null;
  }
}
