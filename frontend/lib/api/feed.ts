import { components } from "@/generated/types";
import { getAuthHeaders } from "./auth";

type FeedResponse = components["schemas"]["FeedResponse"];
type FeedTrip = components["schemas"]["FeedTrip"];

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getFeed(
  limit: number,
  offset: number
): Promise<{ data: FeedTrip[]; total: number }> {
  const response = await fetch(
    `${API_URL}/api/feed?limit=${limit}&offset=${offset}`,
    { method: "GET", headers: await getAuthHeaders() }
  );
  if (!response.ok) {
    throw new Error(`Fehler beim Laden des Feeds: ${response.status}`);
  }
  const data: FeedResponse = await response.json();
  return { data: data.data, total: data.total };
}

export async function getPersonalFeed(
  limit: number,
  offset: number
): Promise<{ data: FeedTrip[]; total: number }> {
  const response = await fetch(
    `${API_URL}/api/feed/personal?limit=${limit}&offset=${offset}`,
    {
      method: "GET",
      headers: await getAuthHeaders(),
    }
  );
  if (!response.ok) {
    throw new Error(`Fehler beim Laden des persönlichen Feeds: ${response.status}`);
  }
  const data: FeedResponse = await response.json();
  return { data: data.data, total: data.total };
}