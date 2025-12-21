import { Gamepad2 } from 'lucide-react';
import { Badge } from './ui/badge';
import { GameState } from '@/lib/types';
import { ThemeToggle } from './theme-toggle';

interface NavbarProps {
  viewMode: 'lobby' | 'game';
  gameState: GameState | null;
}

export function Navbar({ viewMode, gameState }: NavbarProps) {
  return (
    <header className="border-b border-border/50 bg-card/30 backdrop-blur-sm sticky top-0 z-40">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-primary/10">
              <Gamepad2 className="h-8 w-8 text-primary" />
            </div>
            <div>
              <h1 className="text-xl font-bold">4-in-a-Row</h1>
              <p className="text-sm text-muted-foreground">
                Real-time multiplayer
              </p>
            </div>
          </div>
          <div className="flex items-center gap-x-4">
            {/* Status badges */}
            <div className="flex items-center gap-2">
              {viewMode === 'game' && gameState && (
                <Badge
                  variant={gameState.yourTurn ? 'default' : 'secondary'}
                  className="animate-pulse"
                >
                  {gameState.yourTurn ? 'Your Turn' : "Opponent's Turn"}
                </Badge>
              )}
            </div>
            {/* Theme toggle button  */}
            <ThemeToggle />
          </div>
        </div>
      </div>
    </header>
  );
}
