import { getAuthHeaders } from "./auth";
import { components } from "@/generated/types";

type TripLikeResponse = components["schemas"]["TripLikeResponse"];
type TripCommentResponse = components["schemas"]["TripCommentResponse"];
type TripCommentListResponse = components["schemas"]["TripCommentListResponse"];
type CreateTripCommentRequest = components["schemas"]["CreateTripCommentRequest"];

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getTripLikes(tripId: string): Promise<TripLikeResponse> {
    const response = await fetch(`${API_URL}/api/social/${tripId}/likes`, {
        method: "GET",
        headers: await getAuthHeaders(),
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Likes: ${response.status}`);
    }
    return response.json();
}

export async function likeTrip(tripId: string): Promise<void> {
    const response = await fetch(`${API_URL}/api/social/${tripId}/likes`, {
        method: "POST",
        headers: await getAuthHeaders(),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Liken (${response.status}):`, errorData);
        throw new Error(`Fehler beim Liken: ${response.status}`);
    }
}

export async function unlikeTrip(tripId: string): Promise<void> {
    const response = await fetch(`${API_URL}/api/social/${tripId}/likes`, {
        method: "DELETE",
        headers: await getAuthHeaders(),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Unliken (${response.status}):`, errorData);
        throw new Error(`Fehler beim Unliken: ${response.status}`);
    }
}

export async function getTripComments(tripId: string): Promise<TripCommentListResponse> {
    const response = await fetch(`${API_URL}/api/social/${tripId}/comments`, {
        method: "GET",
    });
    if (!response.ok) {
        throw new Error(`Fehler beim Laden der Kommentare: ${response.status}`);
    }
    return response.json();
}

export async function createTripComment(tripId: string, text: string): Promise<TripCommentResponse> {
    const body: CreateTripCommentRequest = { text };
    const response = await fetch(`${API_URL}/api/social/${tripId}/comments`, {
        method: "POST",
        headers: await getAuthHeaders(),
        body: JSON.stringify(body),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Erstellen des Kommentars (${response.status}):`, errorData);
        throw new Error(`Fehler beim Erstellen des Kommentars: ${response.status}`);
    }
    return response.json();
}

export async function deleteTripComment(tripId: string, commentId: string): Promise<void> {
    const response = await fetch(`${API_URL}/api/social/${tripId}/comments/${commentId}`, {
        method: "DELETE",
        headers: await getAuthHeaders(),
    });
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Löschen des Kommentars (${response.status}):`, errorData);
        throw new Error(`Fehler beim Löschen des Kommentars: ${response.status}`);
    }
}