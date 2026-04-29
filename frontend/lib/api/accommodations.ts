import { getAuthHeaders } from "@/lib/api/auth";
import { AccommodationResponse, CreateAccommodationRequest, UpdateAccommodationRequest } from "@/types/accommodation";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getAccommodations(tripId: string): Promise<AccommodationResponse[]> {
    const response = await fetch(`${API_URL}/api/trips/${tripId}/accommodations`, {
        method: "GET",
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Unterkünfte: ${response.status}`);
    }
    const data = await response.json();
    return data.data as AccommodationResponse[];
}

export async function createAccommodation(tripId: string, req: CreateAccommodationRequest): Promise<AccommodationResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/accommodations`, {
        method: "POST",
        headers,
        body: JSON.stringify(req),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler Backend Response:`, errorData);
        throw new Error(`Fehler beim Erstellen der Unterkunft: ${response.status}`);
    }
    return response.json();
}

export async function updateAccommodation(tripId: string, accommodationId: string, req: UpdateAccommodationRequest): Promise<AccommodationResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/accommodations/${accommodationId}`, {
        method: "PUT",
        headers,
        body: JSON.stringify(req),
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Aktualisieren der Unterkunft: ${response.status}`);
    }
    return response.json();
}

export async function deleteAccommodation(tripId: string, accommodationId: string): Promise<void> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/accommodations/${accommodationId}`, {
        method: "DELETE",
        headers,
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Löschen der Unterkunft: ${response.status}`);
    }
}