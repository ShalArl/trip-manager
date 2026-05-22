"use client";

import Link from "next/link";
import { useState, useEffect } from "react";

import { updateTrip } from "@/lib/api/trips";
import { getTripLikes } from "@/lib/api/social";
import { getTransports } from "@/lib/api/transports";
import { getLocations } from "@/lib/api/locations";
import { getAccommodations } from "@/lib/api/accommodations";

import { components } from "@/generated/types";
import { TransportResponse } from "@/types/transport";
import { LocationResponse } from "@/types/location";
import { AccommodationResponse } from "@/types/accommodation";
import { UserResponse } from "@/types/user";
import { TripLikeResponse } from "@/types/social";

import TripLocationsSection from "@/components/trips/sections/TripLocationsSection";
import TripTransportsSection from "@/components/trips/sections/TripTransportsSection";
import TripAccommodationsSection from "@/components/trips/sections/TripAccommodationsSection";
import TripSocialSection from "@/components/trips/sections/TripSocialSection";
import EditTripModal from "./modals/EditTripModal";
import AddActivityModal from "./modals/AddActivityModal";

// ─── Types ────────────────────────────────────────────────────────────────────

type TripResponse = components["schemas"]["TripResponse"];

type Props = {
    trip: TripResponse;
    isEditable?: boolean;
    onTripUpdateAction: (trip: TripResponse) => void;
    currentUser?: UserResponse | null;
};

// ─── Component ────────────────────────────────────────────────────────────────

export default function TripDetail({ trip, isEditable = false, onTripUpdateAction, currentUser }: Props) {
    // ── Trip ────────────────────────────────────────────────────────────────
    const [currentTrip, setCurrentTrip] = useState<TripResponse>(trip);
    const [isEditingTrip, setIsEditingTrip] = useState(false);

    // ── Locations ───────────────────────────────────────────────────────────
    const [locations, setLocations] = useState<LocationResponse[]>([]);
    const [selectedLocationId, setSelectedLocationId] = useState<string | null>(null);
    const activeLocation = locations.find((l) => l.id === selectedLocationId);

    // ── Transports ──────────────────────────────────────────────────────────
    const [transports, setTransports] = useState<TransportResponse[]>([]);

    // ── Accommodations ──────────────────────────────────────────────────────
    const [accommodations, setAccommodations] = useState<AccommodationResponse[]>([]);

    // ── Social ──────────────────────────────────────────────────────────────
    const [likeInfo, setLikeInfo] = useState<TripLikeResponse>({ likeCount: 0, hasLiked: false });

    // ── Activities (placeholder) ─────────────────────────────────────────────
    const [activities, setActivities] = useState<any[]>([]);
    const [showAddActivityModal, setShowAddActivityModal] = useState(false);
    const selectedLocationActivities = activities.filter((a) => a.locationId === selectedLocationId);

    // ── Data fetching ────────────────────────────────────────────────────────

    useEffect(() => {
        getLocations(trip.id).then(setLocations).catch(console.error);
        getTransports(trip.id).then(setTransports).catch(console.error);
        getAccommodations(trip.id).then(setAccommodations).catch(console.error);
    }, [trip.id]);

    useEffect(() => {
        getTripLikes(trip.id).then(setLikeInfo).catch(console.error);
    }, [trip.id, currentUser]);

    // ── Handlers ─────────────────────────────────────────────────────────────

    const handleEditTrip = async (updatedTrip: Partial<TripResponse>) => {
        try {
            const updated = await updateTrip(trip.id, updatedTrip);
            setCurrentTrip(updated);
            onTripUpdateAction(updated);
            setIsEditingTrip(false);
        } catch (err) {
            console.error("[TripDetail] updateTrip:", err);
        }
    };

    const handleAddActivity = (newActivity: any) => {
        setActivities([...activities, {
            id: `act-${Date.now()}`,
            locationId: selectedLocationId!,
            ...newActivity,
        }]);
    };

    // ── Render ────────────────────────────────────────────────────────────────

    return (
        <div className="max-w-5xl px-6 py-12">
            <Link
                href="/"
                className="inline-flex items-center gap-2 text-sm text-zinc-500 dark:text-zinc-400 hover:text-sky-600 dark:hover:text-sky-400 transition-colors mb-8"
            >
                ← Zurück zur Übersicht
            </Link>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">

                {/* ── Left: Main Content ───────────────────────────────────── */}
                <div className="lg:col-span-2 space-y-6">

                    {/* ── Trip Header ─────────────────────────────────────── */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-start justify-between mb-6">
                            <div className="flex items-center gap-4">
                                <div className="w-14 h-14 rounded-2xl bg-sky-50 dark:bg-sky-950/50 border border-sky-100 dark:border-sky-900 flex items-center justify-center text-3xl">
                                    ✈️
                                </div>
                                <div>
                                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
                                        {currentTrip.title}
                                    </h1>
                                    <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                        {currentTrip.startDate} · {currentTrip.endDate}
                                    </p>
                                </div>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setIsEditingTrip(true)}
                                    className="px-3 py-1.5 text-sm font-medium text-sky-600 dark:text-sky-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                >
                                    Bearbeiten
                                </button>
                            )}
                        </div>
                        <div className="border-t border-zinc-100 dark:border-zinc-800 pt-6 space-y-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-2">
                                    Kurzbeschreibung
                                </p>
                                <p className="text-zinc-700 dark:text-zinc-300">{currentTrip.shortDescription}</p>
                            </div>
                            {currentTrip.description && (
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-2">
                                        Details
                                    </p>
                                    <p className="text-zinc-700 dark:text-zinc-300 leading-relaxed">{currentTrip.description}</p>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* ── Social ───────────────────────────────────────────── */}
                    <TripSocialSection
                        tripId={trip.id}
                        currentUser={currentUser}
                        initialLikeInfo={likeInfo}
                    />

                    {/* ── Locations ────────────────────────────────────────── */}
                    <TripLocationsSection
                        tripId={trip.id}
                        isEditable={isEditable}
                        locations={locations}
                        selectedLocationId={selectedLocationId}
                        onLocationsChange={setLocations}
                        onLocationSelect={setSelectedLocationId}
                    />

                    {/* ── Transports ───────────────────────────────────────── */}
                    <TripTransportsSection
                        tripId={trip.id}
                        isEditable={isEditable}
                        transports={transports}
                        onTransportsChange={setTransports}
                    />

                    {/* ── Accommodations ───────────────────────────────────── */}
                    <TripAccommodationsSection
                        tripId={trip.id}
                        isEditable={isEditable}
                        accommodations={accommodations}
                        onAccommodationsChange={setAccommodations}
                    />

                </div>{/* end lg:col-span-2 */}

                {/* ── Right: Activities ────────────────────────────────────── */}
                {activeLocation && (
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 h-fit">
                        <div className="flex items-center justify-between mb-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1">
                                    Aktivitäten in
                                </p>
                                <h3 className="text-lg font-bold text-zinc-900 dark:text-white">{activeLocation.name}</h3>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddActivityModal(true)}
                                    className="p-2 text-sky-600 dark:text-sky-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                >
                                    +
                                </button>
                            )}
                        </div>
                        {selectedLocationActivities.length === 0 ? (
                            <p className="text-zinc-500 dark:text-zinc-400 text-sm text-center py-4">Keine Aktivitäten</p>
                        ) : (
                            <div className="space-y-3">
                                {selectedLocationActivities.map((activity) => (
                                    <div key={activity.id} className="p-3 bg-zinc-50 dark:bg-zinc-800 rounded-lg">
                                        <p className="font-medium text-sm text-zinc-900 dark:text-white">{activity.name}</p>
                                        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">{activity.category}</p>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

            </div>{/* end grid */}

            {/* ── Modals ───────────────────────────────────────────────────── */}
            <EditTripModal
                isOpen={isEditingTrip}
                trip={currentTrip}
                onCloseAction={() => setIsEditingTrip(false)}
                onSaveAction={handleEditTrip}
            />
            <AddActivityModal
                isOpen={showAddActivityModal}
                locationId={selectedLocationId}
                locationName={activeLocation?.name || ""}
                tripStartDate={trip.startDate}
                onCloseAction={() => setShowAddActivityModal(false)}
                onAddAction={handleAddActivity}
            />
        </div>
    );
}