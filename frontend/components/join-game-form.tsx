'use client';

import type React from 'react';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Loader2, Play, Wifi, WifiOff } from 'lucide-react';
import type { ConnectionState } from '@/lib/types';

interface JoinGameFormProps {
  onJoin: (username: string) => void;
  connectionState: ConnectionState;
  error: string | null;
  isWaiting: boolean;
  hasActiveGame?: boolean;
}

export function JoinGameForm({
  onJoin,
  connectionState,
  error,
  isWaiting,
  hasActiveGame = false,
}: JoinGameFormProps) {
  const [username, setUsername] = useState('');
  const [validationError, setValidationError] = useState<string | null>(null);

  const validateUsername = (name: string): boolean => {
    const trimmed = name.trim().toLowerCase();
    if (trimmed === 'bot' || trimmed === 'draw') {
      setValidationError("Username cannot be 'bot' or 'draw'");
      return false;
    }
    setValidationError(null);
    return true;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = username.trim();
    if (trimmed && validateUsername(trimmed)) {
      onJoin(trimmed);
    }
  };

  const isConnected = connectionState === 'connected';
  const isConnecting =
    connectionState === 'connecting' || connectionState === 'reconnecting';

  return (
    <Card className="border-border/50 bg-card/50 backdrop-blur">
      <CardHeader className="pb-4">
        <CardTitle className="flex items-center gap-2 text-xl">
          <Play className="h-5 w-5 text-primary" />
          Join Game
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Input
              placeholder="Enter username"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value);
                if (e.target.value.trim()) {
                  validateUsername(e.target.value);
                }
              }}
              disabled={!isConnected || isWaiting || hasActiveGame}
              className="bg-input/50"
            />
          </div>

          {validationError && (
            <p className="text-sm text-destructive">{validationError}</p>
          )}

          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            {isConnecting ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                <span>Connecting to server...</span>
              </>
            ) : isConnected ? (
              <>
                <Wifi className="h-4 w-4 text-green-500" />
                <span>Connected</span>
              </>
            ) : (
              <>
                <WifiOff className="h-4 w-4 text-destructive" />
                <span>Disconnected</span>
              </>
            )}
          </div>

          {error && <p className="text-sm text-destructive">{error}</p>}

          {hasActiveGame && (
            <div className="flex items-center gap-2 text-sm text-amber-600 dark:text-amber-500 bg-amber-50 dark:bg-amber-950/20 p-3 rounded-md">
              <Play className="h-4 w-4" />
              <span>You are currently in an active game</span>
            </div>
          )}

          {isWaiting && (
            <div className="flex items-center gap-2 text-sm text-primary">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>Waiting for opponent...</span>
            </div>
          )}

          <Button
            type="submit"
            className="w-full"
            disabled={
              !username.trim() || !isConnected || isWaiting || !!validationError || hasActiveGame
            }
          >
            {isWaiting ? 'Searching...' : hasActiveGame ? 'Already in Game' : 'Find Match'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
