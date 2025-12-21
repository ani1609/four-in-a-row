'use client';

import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  useCallback,
  ReactNode,
} from 'react';
import {
  createWebSocket,
  sendMessage as sendWSMessage,
  parseMessage,
  getStoredSession,
  storeSession,
  clearSession,
  type GameSession,
} from '../websocket';
import type { WSMessage, ConnectionState } from '../types';

interface WebSocketContextValue {
  connectionState: ConnectionState;
  isConnected: boolean;
  sendMessage: (message: { type: string; payload: unknown }) => void;
  reconnect: () => void;
  saveSession: (session: GameSession) => void;
  removeSession: () => void;
  getSession: () => GameSession | null;
  subscribeToMessages: (callback: (message: WSMessage) => void) => () => void;
}

const WebSocketContext = createContext<WebSocketContextValue | null>(null);

interface WebSocketProviderProps {
  children: ReactNode;
}

export function WebSocketProvider({ children }: WebSocketProviderProps) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const messageListenersRef = useRef<Set<(message: WSMessage) => void>>(
    new Set()
  );
  const [connectionState, setConnectionState] =
    useState<ConnectionState>('connecting');

  // Clear reconnection timeout
  const clearReconnectTimeout = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  }, []);

  // Notify all message listeners
  const notifyListeners = useCallback((message: WSMessage) => {
    messageListenersRef.current.forEach((callback) => {
      try {
        callback(message);
      } catch (error) {
        console.error('Error in message listener:', error);
      }
    });
  }, []);

  // Handle incoming messages
  const handleMessage = useCallback(
    (event: MessageEvent) => {
      const message = parseMessage(event.data);
      if (!message) {
        console.warn('Failed to parse WebSocket message:', event.data);
        return;
      }
      notifyListeners(message);
    },
    [notifyListeners]
  );

  // Handle WebSocket errors
  const handleError = useCallback(() => {
    console.error('WebSocket connection error');
  }, []);

  // Connect to WebSocket (defined early to be used by other callbacks)
  const connect = useCallback(() => {
    // Don't create a new connection if already connected
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected');
      return;
    }

    // Close existing connection if any
    if (wsRef.current) {
      wsRef.current.onopen = null;
      wsRef.current.onmessage = null;
      wsRef.current.onerror = null;
      wsRef.current.onclose = null;
      wsRef.current.close();
      wsRef.current = null;
    }

    try {
      console.log('Creating WebSocket connection');
      setConnectionState('connecting');

      const ws = createWebSocket();
      wsRef.current = ws;

      // Handle WebSocket open
      ws.onopen = () => {
        console.log('WebSocket connection established');
        setConnectionState('connected');
        clearReconnectTimeout();

        // Attempt to reconnect to existing session
        const session = getStoredSession();
        if (session && wsRef.current) {
          // Check if session is still valid (within 5 min timeout)
          const FIVE_MINUTES = 5 * 60 * 1000;
          const disconnectTime = session.disconnectTime || Date.now();
          const elapsed = Date.now() - disconnectTime;

          if (elapsed < FIVE_MINUTES) {
            console.log('Attempting to reconnect to existing session');
            sendWSMessage(wsRef.current, {
              type: 'RECONNECT',
              payload: session.playerId,
            });
          } else {
            console.log('Session expired (>5 min), clearing session');
            clearSession();
          }
        }
      };

      ws.onmessage = handleMessage;
      ws.onerror = handleError;

      // Handle WebSocket close
      ws.onclose = () => {
        console.log('WebSocket connection closed');
        setConnectionState('disconnected');

        // Simple reconnection after 2 seconds
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log('Attempting to reconnect...');
          setConnectionState('reconnecting');
          connect();
        }, 2000);
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      setConnectionState('disconnected');
    }
  }, [handleMessage, handleError, clearReconnectTimeout]);

  // Send a message
  const sendMessage = useCallback(
    (message: { type: string; payload: unknown }) => {
      if (!wsRef.current) {
        console.warn('Cannot send message: WebSocket not connected');
        return;
      }

      try {
        sendWSMessage(wsRef.current, message);
      } catch (error) {
        console.error('Failed to send message:', error);
      }
    },
    []
  );

  // Session management
  const saveSession = useCallback((session: GameSession) => {
    storeSession(session);
  }, []);

  const removeSession = useCallback(() => {
    clearSession();
  }, []);

  const getSession = useCallback(() => {
    return getStoredSession();
  }, []);

  // Subscribe to messages
  const subscribeToMessages = useCallback(
    (callback: (message: WSMessage) => void) => {
      messageListenersRef.current.add(callback);
      return () => {
        messageListenersRef.current.delete(callback);
      };
    },
    []
  );

  // Manual reconnect
  const reconnect = useCallback(() => {
    clearReconnectTimeout();
    setConnectionState('connecting');
    connect();
  }, [connect, clearReconnectTimeout]);

  // Initialize connection on mount
  useEffect(() => {
    connect();

    // Cleanup on unmount
    return () => {
      clearReconnectTimeout();
      if (wsRef.current) {
        wsRef.current.onopen = null;
        wsRef.current.onmessage = null;
        wsRef.current.onerror = null;
        wsRef.current.onclose = null;
        wsRef.current.close();
      }
    };
  }, [connect, clearReconnectTimeout]);

  const value: WebSocketContextValue = {
    connectionState,
    isConnected: connectionState === 'connected',
    sendMessage,
    reconnect,
    saveSession,
    removeSession,
    getSession,
    subscribeToMessages,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocketContext() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error(
      'useWebSocketContext must be used within WebSocketProvider'
    );
  }
  return context;
}
