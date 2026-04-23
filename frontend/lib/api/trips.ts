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

export async function updateTrip(tripId: string, data: Partial<TripResponse>): Promise<TripResponse> {
    const response = await fetch(`${API_URL}/api/trips/${tripId}`, {
        method: "PUT",
        headers: await getAuthHeaders(),
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Aktualisieren der Reise (${response.status}):`, errorData);
        throw new Error(`Fehler beim Aktualisieren der Reise: ${response.status}`);
    }

    return await response.json() as TripResponse;
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

export async function getRecentPublicTrips(limit: number, offset: number): Promise<{ data: TripResponse[], total: number }> {
  const response = await fetch(`${API_URL}/api/trips/recent?limit=${limit}&offset=${offset}`, {
      method: "GET",
  });

  if (!response.ok) {
    throw new Error(`Fehler: ${response.status}`);
  }

  const data = await response.json();
  return { data: data.data as TripResponse[], total: data.total as number };
}

export async function searchTrips(query: string, limit: number, offset: number): Promise<{ data: TripResponse[], total: number }> {
  const response = await fetch(`${API_URL}/api/trips/search?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}`, {
      method: "GET",
  });

  if (!response.ok) {
    const errorData = await response.text();
    console.error(`Fehler beim Aktualisieren der Reise (${response.status}):`, errorData);
    throw new Error(`Fehler beim Aktualisieren der Reise: ${response.status}`);
  }

  const data = await response.json();
  return { data: data.data as TripResponse[], total: data.total as number };
}