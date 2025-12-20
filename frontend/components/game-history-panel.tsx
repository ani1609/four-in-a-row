"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { History, Trophy } from "lucide-react";
import type { RecentGame } from "@/lib/types";

interface GameHistoryPanelProps {
  games: RecentGame[];
  isLoading?: boolean;
}

export function GameHistoryPanel({
  games,
  isLoading = false,
}: GameHistoryPanelProps) {
  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
    });
  };

  const isWinner = (game: RecentGame, player: string) => {
    return game.winner === player;
  };

  return (
    <Card className="border-border/50 bg-card/50 backdrop-blur">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-lg">
          <History className="h-4 w-4 text-primary" />
          Game History
        </CardTitle>
      </CardHeader>
      <CardContent className="pt-0">
        {isLoading ? (
          <div className="space-y-2">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-10 bg-muted/50 rounded animate-pulse" />
            ))}
          </div>
        ) : games.length === 0 ? (
          <p className="text-muted-foreground text-sm text-center py-4">
            No games played yet.
          </p>
        ) : (
          <div className="space-y-1.5 max-h-115 overflow-y-auto pr-1">
            {games.map((game) => (
              <div
                key={game.gameId}
                className="flex items-center gap-2 px-2.5 py-2 rounded bg-secondary/30 text-sm"
              >
                {/* Players - compact inline */}
                <div className="flex items-center gap-1 min-w-0 flex-1">
                  <span
                    className={`truncate ${
                      isWinner(game, game.player1)
                        ? "text-primary font-medium"
                        : "text-muted-foreground"
                    }`}
                  >
                    {game.player1}
                  </span>
                  <span className="text-muted-foreground/60 shrink-0">v</span>
                  <span
                    className={`truncate ${
                      isWinner(game, game.player2)
                        ? "text-primary font-medium"
                        : "text-muted-foreground"
                    }`}
                  >
                    {game.player2}
                  </span>
                </div>

                {/* Winner indicator */}
                <div className="flex items-center gap-1 shrink-0">
                  {game.winner !== "draw" ? (
                    <Trophy className="h-3 w-3 text-chart-2" />
                  ) : (
                    <span className="text-xs text-muted-foreground">Draw</span>
                  )}
                </div>

                {/* Date */}
                <span className="text-xs text-muted-foreground/70 shrink-0">
                  {formatDate(game.playedAt)}
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
