"use client";

import { useEffect, useState } from "react";
import { getTrip } from "@/lib/api/trips";
import TripDetail from "@/components/trips/TripDetail";
import { TripResponse } from "@/types/trip";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";
import { useUserContext } from "@/lib/context/UserContext";


export default function TripDetailPage({ params }: { params: Promise<{ id: string }> }) {
    const { user, isLoading: userLoading } = useUserContext();
    const [trip, setTrip] = useState<TripResponse | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        params.then(async ({ id }) => {
            try {
                const data = await getTrip(id);
                setTrip(data);
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
        return <LoadingSpinner />;
    }

    const isEditable = !!user && trip.createdBy?.id === user.id;

    return <TripDetail trip={trip} isEditable={isEditable} onTripUpdate={setTrip} currentUser={user} />;
}