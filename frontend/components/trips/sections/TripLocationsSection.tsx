"use client";

import { useState } from "react";
import { MapPin, Plus, Pencil, ChevronDown, ChevronUp } from "lucide-react";

import { createLocation, updateLocation, deleteLocation } from "@/lib/api/locations";
import { LocationResponse, CreateLocationRequest, UpdateLocationRequest } from "@/types/location";
import AddLocationModal from "@/components/trips/modals/AddLocationModal";
import LocationDetailModal from "@/components/trips/modals/LocationDetailModal";

// ─── Types ────────────────────────────────────────────────────────────────────

interface TripLocationsSectionProps {
    tripId: string;
    isEditable: boolean;
    locations: LocationResponse[];
    selectedLocationId: string | null;
    onLocationsChange: (locations: LocationResponse[]) => void;
    onLocationSelect: (locationId: string | null) => void;
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function TripLocationsSection({
    tripId,
    isEditable,
    locations,
    selectedLocationId,
    onLocationsChange,
    onLocationSelect,
}: TripLocationsSectionProps) {
    const [showAddModal, setShowAddModal] = useState(false);
    const [showDetailModal, setShowDetailModal] = useState(false);
    const [activeLocation, setActiveLocation] = useState<LocationResponse | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // ── Handlers ────────────────────────────────────────────────────────────

    const handleAdd = async (req: CreateLocationRequest) => {
        setIsLoading(true);
        setError(null);
        try {
            const created = await createLocation(tripId, req);
            onLocationsChange([...locations, created]);
            setShowAddModal(false);
        } catch (err) {
            setError("Fehler beim Erstellen des Ortes");
            console.error("[TripLocationsSection] createLocation:", err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleUpdate = async (req: UpdateLocationRequest) => {
        if (!activeLocation) return;
        setIsLoading(true);
        setError(null);
        try {
            const updated = await updateLocation(tripId, activeLocation.id!, req);
            onLocationsChange(locations.map((l) => (l.id === updated.id ? updated : l)));
            setShowDetailModal(false);
            setActiveLocation(null);
        } catch (err) {
            setError("Fehler beim Aktualisieren des Ortes");
            console.error("[TripLocationsSection] updateLocation:", err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleDelete = async () => {
        if (!activeLocation) return;
        setIsLoading(true);
        setError(null);
        try {
            await deleteLocation(tripId, activeLocation.id!);
            onLocationsChange(locations.filter((l) => l.id !== activeLocation.id));
            setShowDetailModal(false);
            setActiveLocation(null);
            if (selectedLocationId === activeLocation.id) {
                onLocationSelect(null);
            }
        } catch (err) {
            setError("Fehler beim Löschen des Ortes");
            console.error("[TripLocationsSection] deleteLocation:", err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleEditClick = (e: React.MouseEvent, location: LocationResponse) => {
        e.stopPropagation();
        setActiveLocation(location);
        setShowDetailModal(true);
    };

    const handleLocationClick = (locationId: string) => {
        onLocationSelect(selectedLocationId === locationId ? null : locationId);
    };

    // ── Render ──────────────────────────────────────────────────────────────

    return (
        <>
            <section className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center gap-2">
                        <MapPin className="w-5 h-5 text-sky-500" />
                        <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                            Orte
                            <span className="ml-2 text-sm font-normal text-zinc-400 dark:text-zinc-500">
                                ({locations.length})
                            </span>
                        </h2>
                    </div>
                    {isEditable && (
                        <button
                            onClick={() => setShowAddModal(true)}
                            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                        >
                            <Plus className="w-4 h-4" />
                            Ort hinzufügen
                        </button>
                    )}
                </div>

                {/* Error */}
                {error && (
                    <p className="mb-4 text-sm text-red-600 dark:text-red-400">{error}</p>
                )}

                {/* Empty state */}
                {locations.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500">
                        <MapPin className="w-10 h-10 mb-3 opacity-30" />
                        <p className="text-sm">Keine Orte hinzugefügt</p>
                    </div>
                ) : (
                    <ul className="space-y-2">
                        {locations.map((location) => {
                            const isSelected = selectedLocationId === location.id;
                            return (
                                <li key={location.id}>
                                    <div
                                        role="button"
                                        tabIndex={0}
                                        onClick={() => handleLocationClick(location.id!)}
                                        onKeyDown={(e) => e.key === "Enter" && handleLocationClick(location.id!)}
                                        className={`w-full text-left p-4 rounded-xl border-2 transition-all cursor-pointer outline-none focus-visible:ring-2 focus-visible:ring-sky-400 ${isSelected
                                                ? "bg-sky-50 dark:bg-sky-950/30 border-sky-300 dark:border-sky-700"
                                                : "bg-zinc-50 dark:bg-zinc-800/50 border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700"
                                            }`}
                                    >
                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-3 min-w-0">
                                                <div className={`w-2 h-2 rounded-full shrink-0 ${isSelected ? "bg-sky-500" : "bg-zinc-300 dark:bg-zinc-600"}`} />
                                                <div className="min-w-0">
                                                    <p className="font-medium text-zinc-900 dark:text-white truncate">
                                                        {location.name}
                                                    </p>
                                                    <p className="text-sm text-zinc-500 dark:text-zinc-400 truncate">
                                                        {location.city}, {location.country}
                                                        {" · "}
                                                        {location.dateFrom} – {location.dateTo}
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-1 shrink-0 ml-2">
                                                {isEditable && (
                                                    <button
                                                        onClick={(e) => handleEditClick(e, location)}
                                                        className="p-2 text-zinc-400 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                                        aria-label="Ort bearbeiten"
                                                    >
                                                        <Pencil className="w-4 h-4" />
                                                    </button>
                                                )}
                                                {isSelected ? (
                                                    <ChevronUp className="w-4 h-4 text-sky-500" />
                                                ) : (
                                                    <ChevronDown className="w-4 h-4 text-zinc-400" />
                                                )}
                                            </div>
                                        </div>

                                        {/* Expanded: short description */}
                                        {isSelected && location.shortDescription && (
                                            <p className="mt-3 ml-5 text-sm text-zinc-600 dark:text-zinc-400 leading-relaxed">
                                                {location.shortDescription}
                                            </p>
                                        )}
                                    </div>
                                </li>
                            );
                        })}
                    </ul>
                )}
            </section>

            {/* Modals */}
            <AddLocationModal
                isOpen={showAddModal}
                onCloseAction={() => setShowAddModal(false)}
                onAddAction={handleAdd}
            />
            {activeLocation && (
                <LocationDetailModal
                    isOpen={showDetailModal}
                    location={activeLocation}
                    onCloseAction={() => {
                        setShowDetailModal(false);
                        setActiveLocation(null);
                    }}
                    onSaveAction={handleUpdate}
                    onDeleteAction={handleDelete}
                />
            )}
        </>
    );
}