export type PlayerSymbol = 1 | 2;
export type PlayerType = 'human' | 'bot';
export type CellValue = 0 | 1 | 2;
export type Grid = CellValue[][];

// ============= Player Info =============

export interface PlayerInfo {
  playerId?: string;
  username: string;
  symbol: PlayerSymbol;
  type: PlayerType;
  isOnline: boolean;
}

// ============= Game State =============

export interface GameState {
  gameId: string;
  you: PlayerInfo;
  opponent: PlayerInfo;
  grid: Grid;
  currentTurn: PlayerSymbol;
  yourTurn: boolean;
  moveNumber: number;
  status: 'waiting' | 'active' | 'finished';
  winner?: string;
}

// ============= Move Data =============

export interface MoveData {
  moveNumber: number;
  player: PlayerSymbol;
  column: number;
  row: number;
  timestamp: number;
}

export interface LastMove {
  player: PlayerSymbol;
  column: number;
  row: number;
}

// ============= WebSocket Messages =============

export type WSMessage =
  | GameStartMessage
  | GameUpdateMessage
  | GameOverMessage
  | ReconnectMessage
  | PlayerStatusMessage
  | ErrorMessage;

export interface GameStartMessage {
  type: 'GAME_START';
  payload: {
    gameId: string;
    you: PlayerInfo;
    opponent: PlayerInfo;
    yourTurn: boolean;
  };
}

export interface GameUpdateMessage {
  type: 'GAME_UPDATE';
  payload: {
    grid: Grid;
    currentTurn: PlayerSymbol;
    lastMove: LastMove | null;
    moveNumber: number;
  };
}

export interface GameOverMessage {
  type: 'GAME_OVER';
  payload: {
    winner: string;
  };
}

export interface ReconnectMessage {
  type: 'RECONNECT';
  payload: {
    gameId: string;
    you: PlayerInfo;
    opponent: PlayerInfo;
    grid: Grid;
    currentTurn: PlayerSymbol;
    yourTurn: boolean;
    moveNumber: number;
  };
}

export interface ErrorMessage {
  type: 'ERROR';
  payload: string;
}

export interface PlayerStatusMessage {
  type: 'PLAYER_STATUS';
  payload: {
    playerSymbol: PlayerSymbol;
    isOnline: boolean;
    timeLeft: number;
  };
}

// ============= API Responses =============

export interface LeaderboardEntry {
  name: string;
  wins: number;
}

export interface GameMetrics {
  totalGames: number;
  totalPlayers: number;
  averageDuration: number;
  gamesToday: number;
  recentActivity: HourlyActivity[];
}

export interface HourlyActivity {
  hour: number;
  gamesPlayed: number;
  averageDuration: number;
}

export interface RecentGame {
  gameId: string;
  player1: string;
  player2: string;
  winner: string;
  duration: number;
  totalMoves: number;
  playedAt: string;
}

// ============= Connection State =============

export type ConnectionState =
  | 'connecting'
  | 'connected'
  | 'reconnecting'
  | 'disconnected';
