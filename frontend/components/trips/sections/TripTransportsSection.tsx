"use client";

import { useState } from "react";
import { createTransport, updateTransport, deleteTransport } from "@/lib/api/transports";
import { TransportResponse, CreateTransportRequest, UpdateTransportRequest } from "@/types/transport";
import AddTransportModal from "@/components/trips/modals/AddTransportModal";
import EditTransportModal from "@/components/trips/modals/EditTransportModal";

// ─── Types ────────────────────────────────────────────────────────────────────

interface Props {
    tripId: string;
    isEditable: boolean;
    transports: TransportResponse[];
    onTransportsChange: (transports: TransportResponse[]) => void;
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function TripTransportsSection({ tripId, isEditable, transports, onTransportsChange }: Props) {
    const [detailTransport, setDetailTransport] = useState<TransportResponse | null>(null);
    const [showAddModal, setShowAddModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // ── Handlers ─────────────────────────────────────────────────────────────

    const handleAdd = async (req: CreateTransportRequest) => {
        setError(null);
        try {
            const created = await createTransport(tripId, req);
            onTransportsChange([...transports, created]);
            setShowAddModal(false);
        } catch (err) {
            setError("Fehler beim Erstellen des Transports");
            console.error("[TripTransportsSection] createTransport:", err);
        }
    };

    const handleUpdate = async (req: UpdateTransportRequest) => {
        if (!detailTransport) return;
        setError(null);
        try {
            const updated = await updateTransport(tripId, detailTransport.id!, req);
            onTransportsChange(transports.map((t) => t.id === updated.id ? updated : t));
            setShowEditModal(false);
            setDetailTransport(null);
        } catch (err) {
            setError("Fehler beim Aktualisieren des Transports");
            console.error("[TripTransportsSection] updateTransport:", err);
        }
    };

    const handleDelete = async () => {
        if (!detailTransport) return;
        setError(null);
        try {
            await deleteTransport(tripId, detailTransport.id!);
            onTransportsChange(transports.filter((t) => t.id !== detailTransport.id));
            setShowEditModal(false);
            setDetailTransport(null);
        } catch (err) {
            setError("Fehler beim Löschen des Transports");
            console.error("[TripTransportsSection] deleteTransport:", err);
        }
    };

    // ── Render ────────────────────────────────────────────────────────────────

    return (
        <>
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center gap-2">
                        <span className="text-xl">🚀</span>
                        <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                            Transporte
                            <span className="ml-2 text-sm font-normal text-zinc-400 dark:text-zinc-500">
                                ({transports.length})
                            </span>
                        </h2>
                    </div>
                    {isEditable && (
                        <button
                            onClick={() => setShowAddModal(true)}
                            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                        >
                            <span>+</span> Transport
                        </button>
                    )}
                </div>

                {/* Error */}
                {error && <p className="mb-4 text-sm text-red-600 dark:text-red-400">{error}</p>}

                {/* Empty state */}
                {transports.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500">
                        <span className="text-4xl mb-3 opacity-30">🚌</span>
                        <p className="text-sm">Kein Transport hinzugefügt</p>
                    </div>
                ) : (
                    <div className="space-y-3">
                        {transports.map((t) => {
                            const typeEmoji = { flight: "✈️", train: "🚂", car: "🚗", bus: "🚌" }[t.type ?? "flight"] ?? "🚗";
                            return (
                                <div key={t.id} className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700 hover:border-sky-200 dark:hover:border-sky-800 transition-colors">
                                    <div className="flex items-center justify-between gap-4">
                                        <div className="flex items-center gap-3 min-w-0">
                                            <div className="shrink-0 w-10 h-10 rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 flex items-center justify-center text-xl shadow-sm">
                                                {typeEmoji}
                                            </div>
                                            <div className="min-w-0">
                                                <div className="flex items-center gap-2 flex-wrap">
                                                    <span className="font-semibold text-zinc-900 dark:text-white truncate">{t.from?.name || "–"}</span>
                                                    <span className="text-zinc-400 dark:text-zinc-500 shrink-0">→</span>
                                                    <span className="font-semibold text-zinc-900 dark:text-white truncate">{t.to?.name || "–"}</span>
                                                </div>
                                                <div className="flex items-center gap-2 mt-0.5 flex-wrap">
                                                    <span className="text-xs text-zinc-400 dark:text-zinc-500">
                                                        {t.from?.city && t.from?.country ? `${t.from.city}, ${t.from.country}` : ""}
                                                    </span>
                                                    {t.from?.city && t.to?.city && (
                                                        <span className="text-xs text-zinc-300 dark:text-zinc-600">·</span>
                                                    )}
                                                    <span className="text-xs text-zinc-400 dark:text-zinc-500">
                                                        {t.to?.city && t.to?.country ? `${t.to.city}, ${t.to.country}` : ""}
                                                    </span>
                                                </div>
                                                {(t.departureTime || t.arrivalTime) && (
                                                    <div className="flex items-center gap-3 mt-1 flex-wrap">
                                                        {t.departureTime && (
                                                            <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                                🕐 {new Date(t.departureTime).toLocaleString("de-DE", { day: "2-digit", month: "2-digit", year: "numeric", hour: "2-digit", minute: "2-digit" })}
                                                            </span>
                                                        )}
                                                        {t.departureTime && t.arrivalTime && (
                                                            <span className="text-xs text-zinc-300 dark:text-zinc-600">→</span>
                                                        )}
                                                        {t.arrivalTime && (
                                                            <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                                {new Date(t.arrivalTime).toLocaleString("de-DE", { day: "2-digit", month: "2-digit", year: "numeric", hour: "2-digit", minute: "2-digit" })}
                                                            </span>
                                                        )}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                        {isEditable && (
                                            <button
                                                onClick={() => { setDetailTransport(t); setShowEditModal(true); }}
                                                className="shrink-0 p-2 text-zinc-400 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                                aria-label="Transport bearbeiten"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                                </svg>
                                            </button>
                                        )}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>

            {/* Modals */}
            <AddTransportModal
                isOpen={showAddModal}
                onCloseAction={() => setShowAddModal(false)}
                onAddAction={handleAdd}
            />
            {detailTransport && (
                <EditTransportModal
                    isOpen={showEditModal}
                    transport={detailTransport}
                    onCloseAction={() => { setShowEditModal(false); setDetailTransport(null); }}
                    onSaveAction={handleUpdate}
                    onDeleteAction={handleDelete}
                />
            )}
        </>
    );
}