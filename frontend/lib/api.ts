import type { LeaderboardEntry, GameMetrics, RecentGame } from './types';
import { config } from './config';

export async function fetchLeaderboard(): Promise<LeaderboardEntry[]> {
  try {
    const response = await fetch(`${config.apiUrl}/leaderboard`);
    if (!response.ok) {
      throw new Error(`Failed to fetch leaderboard: ${response.statusText}`);
    }
    const data = await response.json();
    return data || [];
  } catch (error) {
    console.error('Error fetching leaderboard:', error);
    return [];
  }
}

export async function fetchMetrics(): Promise<GameMetrics | null> {
  try {
    const response = await fetch(`${config.apiUrl}/metrics`);
    if (!response.ok) {
      throw new Error(`Failed to fetch metrics: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching metrics:', error);
    return null;
  }
}

export async function fetchRecentGames(): Promise<RecentGame[]> {
  try {
    const response = await fetch(`${config.apiUrl}/recent-games`);
    if (!response.ok) {
      throw new Error(`Failed to fetch recent games: ${response.statusText}`);
    }
    const data = await response.json();
    return data || [];
  } catch (error) {
    console.error('Error fetching recent games:', error);
    return [];
  }
}
