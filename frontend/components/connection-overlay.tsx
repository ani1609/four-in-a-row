"use client";

import type { ConnectionState } from "@/lib/types";
import { Loader2, WifiOff, RefreshCcw } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ConnectionOverlayProps {
  state: ConnectionState;
  onRetry: () => void;
}

export function ConnectionOverlay({ state, onRetry }: ConnectionOverlayProps) {
  if (state === "connected") return null;

  return (
    <div className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center">
      <div className="text-center p-6 rounded-xl bg-card border border-border max-w-sm mx-4">
        {state === "connecting" || state === "reconnecting" ? (
          <>
            <Loader2 className="h-12 w-12 animate-spin mx-auto mb-4 text-primary" />
            <h3 className="text-lg font-semibold mb-2">
              {state === "connecting" ? "Connecting..." : "Reconnecting..."}
            </h3>
            <p className="text-sm text-muted-foreground">
              Please wait while we establish a connection to the game server.
            </p>
          </>
        ) : (
          <>
            <WifiOff className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
            <h3 className="text-lg font-semibold mb-2">Disconnected</h3>
            <p className="text-sm text-muted-foreground mb-4">
              Connection lost. Click to reconnect.
            </p>
            <Button onClick={onRetry} className="gap-2">
              <RefreshCcw className="h-4 w-4" />
              Reconnect
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
