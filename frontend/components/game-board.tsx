"use client";

import type { Grid, LastMove, PlayerSymbol } from "@/lib/types";
import { cn } from "@/lib/utils";

interface GameBoardProps {
  grid: Grid;
  onColumnClick: (column: number) => void;
  canInteract: boolean;
  lastMove: LastMove | null;
  winningCells?: { row: number; col: number }[];
  hoverColumn: number | null;
  setHoverColumn: (col: number | null) => void;
  currentPlayerSymbol?: PlayerSymbol;
}

export function GameBoard({
  grid,
  onColumnClick,
  canInteract,
  lastMove,
  winningCells = [],
  hoverColumn,
  setHoverColumn,
  currentPlayerSymbol,
}: GameBoardProps) {
  const isWinningCell = (row: number, col: number) =>
    winningCells.some((cell) => cell.row === row && cell.col === col);

  const isLastMove = (row: number, col: number) =>
    lastMove?.row === row && lastMove?.column === col;

  return (
    <div className="flex flex-col items-center gap-2">
      {/* Preview row showing where disc will drop */}
      <div className="flex gap-1 md:gap-2 px-2">
        {Array.from({ length: 7 }).map((_, col) => (
          <div
            key={col}
            className="w-10 h-10 md:w-12 md:h-12 flex items-center justify-center"
          >
            {canInteract && hoverColumn === col && currentPlayerSymbol && (
              <div
                className={cn(
                  "w-8 h-8 md:w-10 md:h-10 rounded-full opacity-50 animate-pulse",
                  currentPlayerSymbol === 1 ? "bg-player-1" : "bg-player-2"
                )}
              />
            )}
          </div>
        ))}
      </div>

      {/* Main board */}
      <div className="bg-board p-2 md:p-3 rounded-xl shadow-2xl">
        <div className="grid grid-rows-6 gap-1 md:gap-2">
          {grid.map((row, rowIndex) => (
            <div key={rowIndex} className="flex gap-1 md:gap-2">
              {row.map((cell, colIndex) => (
                <button
                  key={colIndex}
                  onClick={() => canInteract && onColumnClick(colIndex)}
                  onMouseEnter={() => canInteract && setHoverColumn(colIndex)}
                  onMouseLeave={() => setHoverColumn(null)}
                  disabled={!canInteract}
                  className={cn(
                    "w-10 h-10 md:w-12 md:h-12 rounded-full transition-all duration-200",
                    "flex items-center justify-center",
                    canInteract
                      ? "cursor-pointer hover:scale-105"
                      : "cursor-default",
                    cell === 0 && "bg-board-cell",
                    cell === 1 && "bg-player-1 shadow-lg",
                    cell === 2 && "bg-player-2 shadow-lg",
                    isWinningCell(rowIndex, colIndex) &&
                      "ring-4 ring-white/50 animate-pulse",
                    isLastMove(rowIndex, colIndex) &&
                      cell !== 0 &&
                      "ring-2 ring-white/30"
                  )}
                  aria-label={`Column ${colIndex + 1}, Row ${rowIndex + 1}, ${
                    cell === 0 ? "empty" : cell === 1 ? "Player 1" : "Player 2"
                  }`}
                />
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
