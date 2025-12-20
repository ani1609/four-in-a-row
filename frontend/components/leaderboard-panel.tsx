"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Trophy, Medal } from "lucide-react";
import type { LeaderboardEntry } from "@/lib/types";

interface LeaderboardPanelProps {
  entries: LeaderboardEntry[];
  isLoading: boolean;
}

export function LeaderboardPanel({
  entries,
  isLoading,
}: LeaderboardPanelProps) {
  return (
    <Card className="border-border/50 bg-card/50 backdrop-blur">
      <CardHeader className="pb-4">
        <CardTitle className="flex items-center gap-2 text-xl">
          <Trophy className="h-5 w-5 text-player-2" />
          Top Players
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-10 bg-muted/50 rounded animate-pulse" />
            ))}
          </div>
        ) : entries.length === 0 ? (
          <p className="text-muted-foreground text-sm text-center py-4">
            No players yet. Be the first!
          </p>
        ) : (
          <div className="space-y-2">
            {entries.map((entry, index) => (
              <div
                key={entry.name}
                className="flex items-center justify-between p-3 rounded-lg bg-secondary/30 hover:bg-secondary/50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <span
                    className={`text-sm font-mono w-6 text-center ${
                      index === 0
                        ? "text-player-2"
                        : index === 1
                        ? "text-gray-400"
                        : index === 2
                        ? "text-orange-400"
                        : "text-muted-foreground"
                    }`}
                  >
                    {index < 3 ? (
                      <Medal className="h-4 w-4 inline" />
                    ) : (
                      `#${index + 1}`
                    )}
                  </span>
                  <span className="font-medium">{entry.name}</span>
                </div>
                <span className="text-sm text-muted-foreground">
                  {entry.wins} wins
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
