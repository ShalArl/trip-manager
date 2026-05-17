import { getAuthHeaders } from "@/lib/api/auth";
import { TransportResponse, CreateTransportRequest, UpdateTransportRequest } from "@/types/transport";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getTransports(tripId: string): Promise<TransportResponse[]> {
    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports`, {
        method: "GET",
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Transporte: ${response.status}`);
    }
    const data = await response.json();
    return data.data as TransportResponse[];
}

export async function createTransport(tripId: string, req: CreateTransportRequest): Promise<TransportResponse> {
    const headers = await getAuthHeaders();
    const body = {
        ...req,
        departureTime: req.departureTime ? new Date(req.departureTime as string).toISOString() : undefined,
        arrivalTime: req.arrivalTime ? new Date(req.arrivalTime as string).toISOString() : undefined,
    };

    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports`, {
        method: "POST",
        headers,
        body: JSON.stringify(body),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler Backend Response:`, errorData);
        throw new Error(`Fehler beim Erstellen des Transports: ${response.status}`);
    }
    return response.json();
}

export async function updateTransport(tripId: string, transportId: string, req: UpdateTransportRequest): Promise<TransportResponse> {
    const headers = await getAuthHeaders();
    const body = {
        ...req,
        departureTime: req.departureTime ? new Date(req.departureTime as string).toISOString() : undefined,
        arrivalTime: req.arrivalTime ? new Date(req.arrivalTime as string).toISOString() : undefined,
    };

    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports/${transportId}`, {
        method: "PUT",
        headers,
        body: JSON.stringify(body),
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Aktualisieren des Transports: ${response.status}`);
    }
    return response.json();
}

export async function deleteTransport(tripId: string, transportId: string): Promise<void> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports/${transportId}`, {
        method: "DELETE",
        headers,
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Löschen des Transports: ${response.status}`);
    }
}