import { getAuthHeaders } from "@/lib/api/auth";
import { LocationResponse, CreateLocationRequest, UpdateLocationRequest } from "@/types/location";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getLocations(tripId: string): Promise<LocationResponse[]> {
    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations`, {
        method: "GET",
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Locations: ${response.status}`);
    }
    const data = await response.json();
    return data.data as LocationResponse[];
}

export async function createLocation(tripId: string, req: CreateLocationRequest): Promise<LocationResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations`, {
        method: "POST",
        headers,
        body: JSON.stringify(req),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler Backend Response:`, errorData);
        throw new Error(`Fehler beim Erstellen der Location: ${response.status}`);
    }
    return response.json();
}

export async function updateLocation(tripId: string, locationId: string, req: UpdateLocationRequest): Promise<LocationResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations/${locationId}`, {
        method: "PUT",
        headers,
        body: JSON.stringify(req),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler Backend Response:`, errorData);
        throw new Error(`Fehler beim Aktualisieren der Location: ${response.status}`);
    }
    return response.json();
}

export async function deleteLocation(tripId: string, locationId: string): Promise<void> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations/${locationId}`, {
        method: "DELETE",
        headers,
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Löschen der Location: ${response.status}`);
    }
}