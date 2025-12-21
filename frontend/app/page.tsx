'use client';

import { useState, useEffect, useCallback } from 'react';
import useSWR from 'swr';
import { JoinGameForm } from '@/components/join-game-form';
import { LeaderboardPanel } from '@/components/leaderboard-panel';
import { MetricsPanel } from '@/components/metrics-panel';
import { GameHistoryPanel } from '@/components/game-history-panel';
import { ConnectionOverlay } from '@/components/connection-overlay';
import { GameBoard } from '@/components/game-board';
import { GameHeader } from '@/components/game-header';
import { GameResultModal } from '@/components/game-result-modal';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useWebSocketContext } from '@/lib/contexts';
import { fetchLeaderboard, fetchMetrics, fetchRecentGames } from '@/lib/api';
import type {
  WSMessage,
  GameState,
  Grid,
  LastMove,
  RecentGame,
} from '@/lib/types';
import { Gamepad2, Trophy, Play, History } from 'lucide-react';
import { Navbar } from '@/components/navbar';

// Create empty 6x7 grid
function createEmptyGrid(): Grid {
  return Array.from({ length: 6 }, () => Array(7).fill(0));
}

export default function HomePage() {
  const {
    connectionState,
    sendMessage,
    reconnect,
    saveSession,
    removeSession,
    subscribeToMessages,
  } = useWebSocketContext();

  const [error, setError] = useState<string | null>(null);
  const [isWaiting, setIsWaiting] = useState(false);

  // View state
  const [viewMode, setViewMode] = useState<'lobby' | 'game'>('lobby');
  const [sidebarTab, setSidebarTab] = useState<'play' | 'stats' | 'history'>(
    'play'
  );

  // Game state
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [lastMove, setLastMove] = useState<LastMove | null>(null);
  const [hoverColumn, setHoverColumn] = useState<number | null>(null);
  const [showResult, setShowResult] = useState(false);

  // Data fetching
  const { data: leaderboard, isLoading: leaderboardLoading } = useSWR(
    'leaderboard',
    fetchLeaderboard,
    {
      refreshInterval: 30000,
      fallbackData: [],
    }
  );

  const { data: metrics, isLoading: metricsLoading } = useSWR(
    'metrics',
    fetchMetrics,
    { refreshInterval: 30000 }
  );

  const { data: recentGames, isLoading: recentGamesLoading } = useSWR<
    RecentGame[]
  >('recent-games', fetchRecentGames, {
    refreshInterval: 30000,
    fallbackData: [],
  });

  // Handle incoming WebSocket messages
  const handleMessage = useCallback(
    (message: WSMessage) => {
      switch (message.type) {
        case 'GAME_START':
          saveSession({
            playerId: message.payload.you.playerId!,
            username: message.payload.you.username,
            gameId: message.payload.gameId,
          });
          setGameState({
            gameId: message.payload.gameId,
            you: message.payload.you,
            opponent: message.payload.opponent,
            grid: createEmptyGrid(),
            currentTurn: message.payload.yourTurn
              ? message.payload.you.symbol
              : message.payload.opponent.symbol,
            yourTurn: message.payload.yourTurn,
            moveNumber: 0,
            status: 'active',
          });
          setViewMode('game');
          setIsWaiting(false);
          break;

        case 'GAME_UPDATE':
          setGameState((prev: GameState | null) => {
            if (!prev) return prev;
            return {
              ...prev,
              grid: message.payload.grid,
              currentTurn: message.payload.currentTurn,
              yourTurn: prev.you.symbol === message.payload.currentTurn,
              moveNumber: message.payload.moveNumber,
            };
          });
          setLastMove(message.payload.lastMove);
          break;

        case 'GAME_OVER':
          setGameState((prev: GameState | null) => {
            if (!prev) return prev;
            return {
              ...prev,
              status: 'finished',
              winner: message.payload.winner,
            };
          });
          setShowResult(true);
          removeSession();
          break;

        case 'RECONNECT':
          saveSession({
            playerId: message.payload.you.playerId!,
            username: message.payload.you.username,
            gameId: message.payload.gameId,
          });
          setGameState({
            gameId: message.payload.gameId,
            you: message.payload.you,
            opponent: message.payload.opponent,
            grid: message.payload.grid,
            currentTurn: message.payload.currentTurn,
            yourTurn: message.payload.yourTurn,
            moveNumber: message.payload.moveNumber,
            status: 'active',
          });
          setViewMode('game');
          break;

        case 'ERROR':
          setError(message.payload);
          setIsWaiting(false);
          if (message.payload.includes('not found')) {
            removeSession();
          }
          break;
      }
    },
    [saveSession, removeSession]
  );

  const handleJoin = useCallback(
    (username: string) => {
      setError(null);
      setIsWaiting(true);
      sendMessage({
        type: 'JOIN_QUEUE',
        payload: username,
      });
    },
    [sendMessage]
  );

  const handleColumnClick = useCallback(
    (column: number) => {
      if (!gameState || !gameState.yourTurn || gameState.status !== 'active')
        return;
      sendMessage({
        type: 'MOVE',
        payload: {
          gameId: gameState.gameId,
          column,
        },
      });
    },
    [gameState, sendMessage]
  );

  const handleGameCleanup = useCallback(() => {
    setGameState(null);
    setLastMove(null);
    setViewMode('lobby');
    setShowResult(false);
  }, []);

  // Subscribe to WebSocket messages
  useEffect(() => {
    const unsubscribe = subscribeToMessages(handleMessage);
    return () => {
      unsubscribe();
    };
  }, [subscribeToMessages, handleMessage]);

  // Warn before leaving if in queue or active game
  useEffect(() => {
    const shouldWarn = isWaiting || gameState?.status === 'active';

    if (!shouldWarn) return;

    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      e.returnValue = ''; // Browser will show default confirmation dialog
      return '';
    };

    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [isWaiting, gameState]);

  // Handle visibility change (tab switch/minimize)
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.hidden && (isWaiting || gameState?.status === 'active')) {
        // Update disconnect timestamp when tab becomes hidden
        if (gameState) {
          saveSession({
            playerId: gameState.you.playerId!,
            username: gameState.you.username,
            gameId: gameState.gameId,
            disconnectTime: Date.now(),
          });
        }
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);

    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [isWaiting, gameState, saveSession]);

  return (
    <div className="min-h-screen bg-background">
      <ConnectionOverlay state={connectionState} onRetry={reconnect} />
      {/* Header */}
      <Navbar viewMode={viewMode} gameState={gameState} />

      {/* Main content */}
      <main className="container mx-auto px-4 py-6">
        <div className="grid gap-6 lg:grid-cols-12">
          {/* Left sidebar */}
          <div className="lg:col-span-4 xl:col-span-3 space-y-4">
            <Tabs
              value={sidebarTab}
              onValueChange={(v) => setSidebarTab(v as typeof sidebarTab)}
            >
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="play" className="gap-1">
                  <Play className="h-3 w-3" />
                  Play
                </TabsTrigger>
                <TabsTrigger value="stats" className="gap-1">
                  <Trophy className="h-3 w-3" />
                  Stats
                </TabsTrigger>
                <TabsTrigger value="history" className="gap-1">
                  <History className="h-3 w-3" />
                  History
                </TabsTrigger>
              </TabsList>

              <TabsContent value="play" className="mt-4 space-y-4">
                <JoinGameForm
                  onJoin={handleJoin}
                  connectionState={connectionState}
                  error={error}
                  isWaiting={isWaiting}
                />
              </TabsContent>

              <TabsContent value="stats" className="mt-4 space-y-4">
                <LeaderboardPanel
                  entries={leaderboard || []}
                  isLoading={leaderboardLoading}
                />
                <MetricsPanel
                  metrics={metrics || null}
                  isLoading={metricsLoading}
                />
              </TabsContent>

              <TabsContent value="history" className="mt-4">
                <GameHistoryPanel
                  games={recentGames || []}
                  isLoading={recentGamesLoading}
                />
              </TabsContent>
            </Tabs>
          </div>

          {/* Main game area */}
          <div className="lg:col-span-8 xl:col-span-9">
            <div className="bg-card/30 border border-border/50 rounded-xl p-6 min-h-150 flex flex-col">
              {/* Lobby view - waiting state */}
              {viewMode === 'lobby' && !isWaiting && (
                <div className="flex-1 flex flex-col items-center justify-center text-center gap-6">
                  <div className="p-6 rounded-full bg-primary/10">
                    <Gamepad2 className="h-16 w-16 text-primary" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold mb-2">Ready to Play?</h2>
                    <p className="text-muted-foreground max-w-md">
                      Enter your username in the Play tab and click Find Match
                      to start a game against another player.
                    </p>
                  </div>
                  <div className="flex items-center gap-8 text-sm text-muted-foreground">
                    <div className="flex items-center gap-2">
                      <div className="w-4 h-4 rounded-full bg-player-1" />
                      <span>Player 1</span>
                    </div>
                    <span>vs</span>
                    <div className="flex items-center gap-2">
                      <div className="w-4 h-4 rounded-full bg-player-2" />
                      <span>Player 2</span>
                    </div>
                  </div>
                  {/* Preview board */}
                  <div className="opacity-30">
                    <GameBoard
                      grid={createEmptyGrid()}
                      onColumnClick={() => {}}
                      canInteract={false}
                      lastMove={null}
                      hoverColumn={null}
                      setHoverColumn={() => {}}
                    />
                  </div>
                </div>
              )}

              {/* Lobby view - waiting for opponent */}
              {viewMode === 'lobby' && isWaiting && (
                <div className="flex-1 flex flex-col items-center justify-center text-center gap-6">
                  <div className="relative">
                    <div className="p-6 rounded-full bg-primary/10 animate-pulse">
                      <Gamepad2 className="h-16 w-16 text-primary" />
                    </div>
                    <div className="absolute inset-0 rounded-full border-4 border-primary/30 border-t-primary animate-spin" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold mb-2">
                      Finding Opponent...
                    </h2>
                    <p className="text-muted-foreground">
                      Waiting for another player to join. This usually takes a
                      few seconds.
                    </p>
                  </div>
                  <Button variant="outline" onClick={() => setIsWaiting(false)}>
                    Cancel
                  </Button>
                </div>
              )}

              {/* Live game view */}
              {viewMode === 'game' && gameState && (
                <div className="flex-1 flex flex-col items-center gap-6">
                  <GameHeader
                    you={gameState.you}
                    opponent={gameState.opponent}
                    currentTurn={gameState.currentTurn}
                    yourTurn={gameState.yourTurn}
                    moveNumber={gameState.moveNumber}
                  />
                  <GameBoard
                    grid={gameState.grid}
                    onColumnClick={handleColumnClick}
                    canInteract={
                      gameState.yourTurn && gameState.status === 'active'
                    }
                    lastMove={lastMove}
                    hoverColumn={hoverColumn}
                    setHoverColumn={setHoverColumn}
                    currentPlayerSymbol={gameState.you.symbol}
                  />
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
      {/* Result modal */}
      {gameState && (
        <GameResultModal
          isOpen={showResult}
          winner={gameState.winner || null}
          yourUsername={gameState.you.username}
          onClickOkay={handleGameCleanup}
        />
      )}
    </div>
  );
}
