'use client';

import { useState, useEffect } from 'react';
import type {
  PlayerInfo,
  PlayerSymbol,
  PlayerStatusMessage,
} from '@/lib/types';
import { Badge } from '@/components/ui/badge';
import { User, Bot } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useWebSocketContext } from '@/lib/contexts';

interface GameHeaderProps {
  you: PlayerInfo;
  opponent: PlayerInfo;
  currentTurn: PlayerSymbol;
  yourTurn: boolean;
  moveNumber: number;
}

export function GameHeader({
  you,
  opponent,
  yourTurn,
  moveNumber,
}: GameHeaderProps) {
  const { subscribeToMessages } = useWebSocketContext();
  const [opponentTimeLeft, setOpponentTimeLeft] = useState<number>(0);

  useEffect(() => {
    const unsubscribe = subscribeToMessages((message) => {
      if (message.type === 'PLAYER_STATUS') {
        const payload = (message as PlayerStatusMessage).payload;
        // Only track countdown if opponent goes offline
        if (payload.playerSymbol === opponent.symbol && !payload.isOnline) {
          setOpponentTimeLeft(payload.timeLeft);
        } else if (
          payload.playerSymbol === opponent.symbol &&
          payload.isOnline
        ) {
          setOpponentTimeLeft(0);
        }
      }
    });

    return unsubscribe;
  }, [subscribeToMessages, opponent.symbol]);

  return (
    <div className="flex flex-col gap-4 w-full max-w-md mx-auto">
      <div className="flex items-center justify-between gap-4">
        {/* You */}
        <div
          className={cn(
            'flex-1 p-3 rounded-lg border transition-all',
            yourTurn
              ? 'border-primary bg-primary/10 shadow-lg'
              : 'border-border/50 bg-card/50'
          )}
        >
          <div className="flex items-center gap-2">
            <div
              className={cn(
                'w-6 h-6 rounded-full',
                you.symbol === 1 ? 'bg-player-1' : 'bg-player-2'
              )}
            />
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <p className="font-medium truncate">{you.username}</p>
                <div
                  className={cn(
                    'w-2 h-2 rounded-full',
                    you.isOnline ? 'bg-green-500' : 'bg-gray-400'
                  )}
                  title={you.isOnline ? 'Online' : 'Offline'}
                />
              </div>
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                {you.type === 'bot' ? (
                  <Bot className="h-3 w-3" />
                ) : (
                  <User className="h-3 w-3" />
                )}
                <span>You</span>
              </div>
            </div>
          </div>
          {yourTurn && (
            <Badge variant="secondary" className="mt-2 text-xs">
              Your turn
            </Badge>
          )}
        </div>

        {/* VS */}
        <div className="text-muted-foreground font-bold text-sm">VS</div>

        {/* Opponent */}
        <div
          className={cn(
            'flex-1 p-3 rounded-lg border transition-all',
            !yourTurn
              ? 'border-primary bg-primary/10 shadow-lg'
              : 'border-border/50 bg-card/50'
          )}
        >
          <div className="flex items-center gap-2">
            <div
              className={cn(
                'w-6 h-6 rounded-full',
                opponent.symbol === 1 ? 'bg-player-1' : 'bg-player-2'
              )}
            />
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <p className="font-medium truncate">{opponent.username}</p>
                <div
                  className={cn(
                    'w-2 h-2 rounded-full',
                    opponent.isOnline ? 'bg-green-500' : 'bg-gray-400'
                  )}
                  title={opponent.isOnline ? 'Online' : 'Offline'}
                />
              </div>
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                {opponent.type === 'bot' ? (
                  <Bot className="h-3 w-3" />
                ) : (
                  <User className="h-3 w-3" />
                )}
                <span>Opponent</span>
              </div>
            </div>
          </div>
          {!yourTurn && (
            <Badge variant="secondary" className="mt-2 text-xs">
              Their turn
            </Badge>
          )}
          {!opponent.isOnline && opponentTimeLeft > 0 && (
            <Badge variant="destructive" className="mt-2 text-xs">
              Reconnecting... {opponentTimeLeft}s
            </Badge>
          )}
        </div>
      </div>

      <div className="text-center text-sm text-muted-foreground">
        Move #{moveNumber}
      </div>
    </div>
  );
}
