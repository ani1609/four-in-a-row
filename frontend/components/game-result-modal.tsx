"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Trophy, Frown, Handshake } from "lucide-react";

interface GameResultModalProps {
  isOpen: boolean;
  winner: string | null;
  yourUsername: string;
  onClickOkay: () => void;
}

export function GameResultModal({
  isOpen,
  winner,
  yourUsername,
  onClickOkay,
}: GameResultModalProps) {
  const isWinner = winner === yourUsername;
  const isDraw = winner === "draw";

  const getTitle = () => {
    if (isDraw) return "It's a Draw!";
    if (isWinner) return "You Won!";
    return "You Lost";
  };

  const getIcon = () => {
    if (isDraw)
      return <Handshake className="h-16 w-16 text-muted-foreground" />;
    if (isWinner) return <Trophy className="h-16 w-16 text-player-2" />;
    return <Frown className="h-16 w-16 text-player-1" />;
  };

  return (
    <Dialog open={isOpen}>
      <DialogContent className="sm:max-w-md" showCloseButton={false}>
        <DialogHeader className="text-center">
          <div className="flex justify-center mb-4">{getIcon()}</div>
          <DialogTitle className="text-2xl">{getTitle()}</DialogTitle>
          <DialogDescription>
            {isDraw
              ? "Great game! Neither player could claim victory."
              : isWinner
              ? "Congratulations on your victory!"
              : `${winner} won this round. Better luck next time!`}
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-3 mt-4">
          <Button
            variant="outline"
            className="w-full gap-2 bg-transparent"
            onClick={onClickOkay}
          >
            Okay
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
