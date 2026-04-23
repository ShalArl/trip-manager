import { getAuthHeaders } from "./auth";
import {CreateTripRequest, TripResponse} from "@/types/trip";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function createTrip(createTripRequest: CreateTripRequest): Promise<TripResponse> {
  const response = await fetch(`${API_URL}/api/trips`, {
    method: "POST",
    headers: await getAuthHeaders(),
    body: JSON.stringify(createTripRequest),
  });

  if (!response.ok) {
    const errorData = await response.text();
    console.error("Fehler Details:", errorData);
    throw new Error("Fehler beim Erstellen der Reise");
  }

  return response.json();
}

export async function getTrips(): Promise<TripResponse[]> {
  const response = await fetch(`${API_URL}/api/trips`, {
    method: "GET",
    headers: await getAuthHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.text();
    console.error(`Fehler beim Laden der Reisen (${response.status}):`, errorData);
    throw new Error(`Fehler beim Laden der Reisen: ${response.status} ${response.statusText}`);
  }

  const data = await response.json();
  return data.data as TripResponse[];
}

export async function getTrip(tripId: string): Promise<TripResponse> {
  const response = await fetch(`${API_URL}/api/trips/${tripId}`, {
    method: "GET",
    headers: await getAuthHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.text();
    console.error(`Fehler beim Laden der Reise (${response.status}):`, errorData);
    throw new Error(`Fehler beim Laden der Reise: ${response.status} ${response.statusText}`);
  }

  return response.json();
}

export async function getRecentPublicTrips(limit: number): Promise<TripResponse[]> {
  const response = await fetch(`${API_URL}/api/trips/recent?limit=${limit}`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error(`Fehler: ${response.status}`);
  }

  const data = await response.json();
  return data.data as TripResponse[];
}

export async function searchTrips(query: string): Promise<TripResponse[]> {
  const response = await fetch(
      `${API_URL}/api/trips/search?q=${encodeURIComponent(query)}`,
      { method: "GET" },
  );

  if (!response.ok) {
    const errorData = await response.text();
    console.error(`Fehler bei Suche (${response.status}):`, errorData);
    throw new Error(`Fehler bei Suche: ${response.status}`);
  }

  const data = await response.json();
  return data.data as TripResponse[];
}