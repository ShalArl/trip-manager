import {TransportResponse, CreateTransportRequest, UpdateTransportRequest} from "@/types/transport";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getTransports(tripId: string): Promise<TransportResponse[]> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Transporte: ${response.status}`);
    }

    const data = await response.json();
    return data.data as TransportResponse[];
}

export async function createTransport(tripId: string, req: CreateTransportRequest): Promise<TransportResponse> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify(req),
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Erstellen des Transports: ${response.status}`);
    }

    const data = await response.json();
    return data as TransportResponse;
}

export async function updateTransport(tripId: string, transportId: string, req: UpdateTransportRequest): Promise<TransportResponse> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports/${transportId}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify(req),
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Aktualisieren des Transports: ${response.status}`);
    }

    const data = await response.json();
    return data as TransportResponse;
}

export async function deleteTransport(tripId: string, transportId: string): Promise<void> {
    const token = localStorage.getItem("token");

    const response = await fetch(`${API_URL}/api/trips/${tripId}/transports/${transportId}`, {
        method: "DELETE",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        throw new Error(`Fehler beim Löschen des Transports: ${response.status}`);
    }
}