import type { LeaderboardEntry, GameMetrics, RecentGame } from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function fetchLeaderboard(): Promise<LeaderboardEntry[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/leaderboard`);
    if (!response.ok) {
      throw new Error(`Failed to fetch leaderboard: ${response.statusText}`);
    }
    const data = await response.json();
    return data || [];
  } catch (error) {
    console.error("Error fetching leaderboard:", error);
    return [];
  }
}

export async function fetchMetrics(): Promise<GameMetrics | null> {
  try {
    const response = await fetch(`${API_BASE_URL}/metrics`);
    if (!response.ok) {
      throw new Error(`Failed to fetch metrics: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching metrics:", error);
    return null;
  }
}

export async function fetchRecentGames(): Promise<RecentGame[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/recent-games`);
    if (!response.ok) {
      throw new Error(`Failed to fetch recent games: ${response.statusText}`);
    }
    const data = await response.json();
    return data || [];
  } catch (error) {
    console.error("Error fetching recent games:", error);
    return [];
  }
}
