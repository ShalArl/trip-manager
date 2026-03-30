"use client";

import { useEffect, useState } from "react";
import { getTrip } from "@/lib/api/trips";
import TripDetail from "@/components/trips/TripDetail";
import {TripResponse} from "@/types/trip";


export default function TripDetailPage({ params }: { params: Promise<{ id: string }> }) {
    const [trip, setTrip] = useState<TripResponse | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isEditable, setIsEditable] = useState(false);

    useEffect(() => {
        // Load user ID from localStorage once on mount
        const userJson = localStorage.getItem("user");
        let currentUserId: string | null = null;
        
        if (userJson) {
            try {
                const user = JSON.parse(userJson);
                currentUserId = user.id;
            } catch (err) {
                console.error("Failed to parse user from localStorage", err);
            }
        }

        // Extract trip ID from params and fetch trip
        params.then(async ({ id }) => {
            try {
                const data = await getTrip(id);
                setTrip(data);
                
                // Calculate isEditable after trip is loaded
                const canEdit = data.createdBy?.id === currentUserId;
                setIsEditable(canEdit);
                
                console.log("Trip Daten:", data);
                console.log("IsEditable:", canEdit, "CurrentUserId:", currentUserId, "TripOwner:", data.createdBy?.id);
            } catch (err) {
                console.error("Failed to fetch trip:", err);
                setError("Fehler beim Laden der Reise");
            }
        });
    }, [params]);

    if (error) {
        return <div>{error}</div>;
    }

    if (!trip) {
        return <div>Lädt...</div>;
    }


    return <TripDetail trip={trip} isEditable={isEditable} />;
}