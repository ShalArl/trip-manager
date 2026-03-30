"use client";

import { useEffect, useState } from "react";
import { getTrip } from "@/lib/api/trips";
import TripDetail from "@/components/trips/TripDetail";
import { components } from "@/generated/types";

type TripResponse = components["schemas"]["TripResponse"];

export default function TripDetailPage({ params }: { params: Promise<{ id: string }> }) {
    const [trip, setTrip] = useState<TripResponse | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [tripId, setTripId] = useState<string | null>(null);

    // Extract trip ID from params
    useEffect(() => {
        params.then(({ id }) => setTripId(id));
    }, [params]);

    // Fetch trip when ID is available
    useEffect(() => {
        if (!tripId) return;

        async function fetchTrip() {
            try {
                const data = await getTrip(tripId as string);
                setTrip(data);
            } catch (err) {
                console.error("Failed to fetch trip:", err);
                setError("Fehler beim Laden der Reise");
            }
        }

        fetchTrip();
    }, [tripId]);

    if (error) {
        return <div>{error}</div>;
    }

    if (!trip) {
        return <div>Lädt...</div>;
    }

    console.log("Trip Daten:", trip);

    return <TripDetail trip={trip} />;
}