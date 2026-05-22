"use client";

import { useState } from "react";
import { createAccommodation, updateAccommodation, deleteAccommodation } from "@/lib/api/accommodations";
import { AccommodationResponse, CreateAccommodationRequest, UpdateAccommodationRequest } from "@/types/accommodation";
import AddAccommodationModal from "@/components/trips/modals/AddAccommodationModal";
import EditAccommodationModal from "@/components/trips/modals/EditAccommodationModal";

// ─── Types ────────────────────────────────────────────────────────────────────

interface Props {
    tripId: string;
    isEditable: boolean;
    accommodations: AccommodationResponse[];
    onAccommodationsChange: (accommodations: AccommodationResponse[]) => void;
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function TripAccommodationsSection({ tripId, isEditable, accommodations, onAccommodationsChange }: Props) {
    const [detailAccommodation, setDetailAccommodation] = useState<AccommodationResponse | null>(null);
    const [showAddModal, setShowAddModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // ── Handlers ─────────────────────────────────────────────────────────────

    const handleAdd = async (req: CreateAccommodationRequest) => {
        setError(null);
        try {
            const created = await createAccommodation(tripId, req);
            onAccommodationsChange([...accommodations, created]);
            setShowAddModal(false);
        } catch (err) {
            setError("Fehler beim Erstellen der Unterkunft");
            console.error("[TripAccommodationsSection] createAccommodation:", err);
        }
    };

    const handleUpdate = async (req: UpdateAccommodationRequest) => {
        if (!detailAccommodation) return;
        setError(null);
        try {
            const updated = await updateAccommodation(tripId, detailAccommodation.id!, req);
            onAccommodationsChange(accommodations.map((a) => a.id === updated.id ? updated : a));
            setShowEditModal(false);
            setDetailAccommodation(null);
        } catch (err) {
            setError("Fehler beim Aktualisieren der Unterkunft");
            console.error("[TripAccommodationsSection] updateAccommodation:", err);
        }
    };

    const handleDelete = async () => {
        if (!detailAccommodation) return;
        setError(null);
        try {
            await deleteAccommodation(tripId, detailAccommodation.id!);
            onAccommodationsChange(accommodations.filter((a) => a.id !== detailAccommodation.id));
            setShowEditModal(false);
            setDetailAccommodation(null);
        } catch (err) {
            setError("Fehler beim Löschen der Unterkunft");
            console.error("[TripAccommodationsSection] deleteAccommodation:", err);
        }
    };

    // ── Render ────────────────────────────────────────────────────────────────

    return (
        <>
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center gap-2">
                        <span className="text-xl">🏨</span>
                        <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                            Unterkünfte
                            <span className="ml-2 text-sm font-normal text-zinc-400 dark:text-zinc-500">
                                ({accommodations.length})
                            </span>
                        </h2>
                    </div>
                    {isEditable && (
                        <button
                            onClick={() => setShowAddModal(true)}
                            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                        >
                            <span>+</span> Unterkunft
                        </button>
                    )}
                </div>

                {/* Error */}
                {error && <p className="mb-4 text-sm text-red-600 dark:text-red-400">{error}</p>}

                {/* Empty state */}
                {accommodations.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500">
                        <span className="text-4xl mb-3 opacity-30">🏨</span>
                        <p className="text-sm">Keine Unterkunft hinzugefügt</p>
                    </div>
                ) : (
                    <div className="space-y-3">
                        {accommodations.map((a) => (
                            <div key={a.id} className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700 hover:border-sky-200 dark:hover:border-sky-800 transition-colors">
                                <div className="flex items-center justify-between gap-4">
                                    <div className="flex items-center gap-3 min-w-0">
                                        <div className="shrink-0 w-10 h-10 rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 flex items-center justify-center text-xl shadow-sm">
                                            🏨
                                        </div>
                                        <div className="min-w-0">
                                            <p className="font-semibold text-zinc-900 dark:text-white truncate">{a.name}</p>
                                            <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5 truncate">
                                                📍 {a.location?.name
                                                    ? `${a.location.name}, ${a.location.city}, ${a.location.country}`
                                                    : "Kein Ort angegeben"}
                                            </p>
                                            {(a.checkIn || a.checkOut) && (
                                                <div className="flex items-center gap-3 mt-1 flex-wrap">
                                                    {a.checkIn && (
                                                        <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                            🛬 Check-in: {new Date(a.checkIn).toLocaleDateString("de-DE", { day: "2-digit", month: "2-digit", year: "numeric" })}
                                                        </span>
                                                    )}
                                                    {a.checkIn && a.checkOut && (
                                                        <span className="text-xs text-zinc-300 dark:text-zinc-600">·</span>
                                                    )}
                                                    {a.checkOut && (
                                                        <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                            🛫 Check-out: {new Date(a.checkOut).toLocaleDateString("de-DE", { day: "2-digit", month: "2-digit", year: "numeric" })}
                                                        </span>
                                                    )}
                                                </div>
                                            )}
                                            {a.pricePerNight && (
                                                <p className="text-xs text-sky-600 dark:text-sky-400 mt-1 font-medium">
                                                    {a.pricePerNight} € / Nacht
                                                </p>
                                            )}
                                        </div>
                                    </div>
                                    {isEditable && (
                                        <button
                                            onClick={() => { setDetailAccommodation(a); setShowEditModal(true); }}
                                            className="shrink-0 p-2 text-zinc-400 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                            aria-label="Unterkunft bearbeiten"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                            </svg>
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Modals */}
            <AddAccommodationModal
                isOpen={showAddModal}
                onCloseAction={() => setShowAddModal(false)}
                onAddAction={handleAdd}
            />
            {detailAccommodation && (
                <EditAccommodationModal
                    isOpen={showEditModal}
                    accommodation={detailAccommodation}
                    onCloseAction={() => { setShowEditModal(false); setDetailAccommodation(null); }}
                    onSaveAction={handleUpdate}
                    onDeleteAction={handleDelete}
                />
            )}
        </>
    );
}