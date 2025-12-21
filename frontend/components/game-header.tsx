'use client';

import type { PlayerInfo, PlayerSymbol } from '@/lib/types';
import { Badge } from '@/components/ui/badge';
import { User, Bot } from 'lucide-react';
import { cn } from '@/lib/utils';

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
              <p className="font-medium truncate">{you.username}</p>
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
              <p className="font-medium truncate">{opponent.username}</p>
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
        </div>
      </div>

      <div className="text-center text-sm text-muted-foreground">
        Move #{moveNumber}
      </div>
    </div>
  );
}
