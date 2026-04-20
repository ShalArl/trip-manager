import { LocationResponse, CreateLocationRequest } from "@/types/location";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getLocations(tripId: string): Promise<LocationResponse[]> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Locations: ${response.status}`);
    }

    const data = await response.json();
    return data.data as LocationResponse[];
}

export async function createLocation(tripId: string, req: CreateLocationRequest): Promise<LocationResponse> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify(req),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler Backend Response:`, errorData);
        throw new Error(`Fehler beim Erstellen der Location: ${response.status}`);
    }

    const data = await response.json();
    return data as LocationResponse;
}

export async function deleteLocation(tripId: string, locationId: string): Promise<void> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/locations/${locationId}`, {
        method: "DELETE",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Löschen der Location: ${response.status}`);
    }
}