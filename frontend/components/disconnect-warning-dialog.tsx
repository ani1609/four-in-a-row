'use client';

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';

interface DisconnectWarningDialogProps {
  isOpen: boolean;
  onCancel: () => void;
  onConfirm: () => void;
  isInQueue: boolean;
}

export function DisconnectWarningDialog({
  isOpen,
  onCancel,
  onConfirm,
  isInQueue,
}: DisconnectWarningDialogProps) {
  return (
    <AlertDialog open={isOpen}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            {isInQueue ? 'Leave Queue?' : 'Leave Game?'}
          </AlertDialogTitle>
          <AlertDialogDescription>
            {isInQueue
              ? "You're currently waiting for a match. If you leave now, you might miss your opponent. Are you sure you want to leave?"
              : "You're in an active game! Leaving will disconnect you. You'll have 5 minutes to reconnect, or you'll forfeit the match."}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onCancel}>Stay</AlertDialogCancel>
          <AlertDialogAction onClick={onConfirm}>Leave</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
