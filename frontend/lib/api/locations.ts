import { getAuthHeaders } from "@/lib/api/auth";
import { uploadToPresignedUrl, getPresignedUrl } from "@/lib/api/uploads";
import { LocationResponse, CreateLocationRequest, UpdateLocationRequest } from "@/types/location";
import { components } from "@/generated/types";

type LocationImageResponse = components["schemas"]["LocationImageResponse"];

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getLocations(tripId: string): Promise<LocationResponse[]> {
    const response = await fetch(`${API_URL}/api/locations/${tripId}`, {
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
    const response = await fetch(`${API_URL}/api/locations/${tripId}`, {
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
    const response = await fetch(`${API_URL}/api/locations/${tripId}/${locationId}`, {
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
    const response = await fetch(`${API_URL}/api/locations/${tripId}/${locationId}`, {
        method: "DELETE",
        headers,
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Löschen der Location: ${response.status}`);
    }
}

export async function addLocationImage(
    tripId: string,
    locationId: string,
    file: File,
    sequence?: number
): Promise<LocationImageResponse> {
    // 1. Presigned URL holen
    const ticket = await getPresignedUrl({
        fileName: file.name,
        mediaType: "location",
    });

    console.log("[addLocationImage] key:", ticket.key); 

    // 2. Datei direkt zu S3/MinIO hochladen
    await uploadToPresignedUrl({ url: ticket.presignedUrl, file });

    // 3. Key ans Backend schicken
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/locations/${tripId}/${locationId}/images`, {
        method: "POST",
        headers,
        body: JSON.stringify({ imageKey: ticket.key, sequence }),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Hinzufügen des Bildes:`, errorData);
        throw new Error(`Fehler beim Hinzufügen des Bildes: ${response.status}`);
    }

    return response.json();
}

export async function deleteLocationImage(
    tripId: string,
    locationId: string,
    imageId: string
): Promise<void> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/locations/${tripId}/${locationId}/images/${imageId}`, {
        method: "DELETE",
        headers,
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Löschen des Bildes: ${response.status}`);
    }
}