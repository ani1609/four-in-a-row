'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { BarChart3, Users, Timer, Gamepad2 } from 'lucide-react';
import type { GameMetrics } from '@/lib/types';

interface MetricsPanelProps {
  metrics: GameMetrics | null;
  isLoading: boolean;
}

export function MetricsPanel({ metrics, isLoading }: MetricsPanelProps) {
  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}m ${secs}s`;
  };

  const stats = [
    {
      label: 'Total Games',
      value: metrics?.totalGames ?? 0,
      icon: Gamepad2,
      color: 'text-primary',
    },
    {
      label: 'Total Players',
      value: metrics?.totalPlayers ?? 0,
      icon: Users,
      color: 'text-[var(--player-2)]',
    },
    {
      label: 'Avg Duration',
      value: metrics ? formatDuration(metrics.averageDuration) : '0m 0s',
      icon: Timer,
      color: 'text-[var(--player-1)]',
    },
    {
      label: 'Games Today',
      value: metrics?.gamesToday ?? 0,
      icon: BarChart3,
      color: 'text-green-500',
    },
  ];

  return (
    <Card className="border-border/50 bg-card/50 backdrop-blur">
      <CardHeader className="pb-4">
        <CardTitle className="flex items-center gap-2 text-xl">
          <BarChart3 className="h-5 w-5 text-primary" />
          Game Statistics
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="grid grid-cols-2 gap-4">
            {[...Array(4)].map((_, i) => (
              <div
                key={i}
                className="h-20 bg-muted/50 rounded-lg animate-pulse"
              />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-4">
            {stats.map((stat) => (
              <div
                key={stat.label}
                className="p-4 rounded-lg bg-secondary/30 border border-border/30"
              >
                <div className="flex items-center gap-2 mb-2">
                  <stat.icon className={`h-4 w-4 ${stat.color}`} />
                  <span className="text-xs text-muted-foreground uppercase tracking-wide">
                    {stat.label}
                  </span>
                </div>
                <p className="text-2xl font-bold">{stat.value}</p>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
