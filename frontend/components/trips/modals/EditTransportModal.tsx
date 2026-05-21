"use client";

import { useState } from "react";
import { TransportResponse, UpdateTransportRequest } from "@/types/transport";
import PlaceAutocomplete, { PlaceValue } from "@/components/shared/PlaceAutocomplete";

type Props = {
    isOpen: boolean;
    transport: TransportResponse;
    onCloseAction: () => void;
    onSaveAction: (req: UpdateTransportRequest) => void;
    onDeleteAction: () => void;
};

export default function EditTransportModal({ isOpen, transport, onCloseAction, onSaveAction, onDeleteAction }: Props) {
    const [type, setType] = useState<"flight" | "train" | "car" | "bus">(
        (transport.type as "flight" | "train" | "car" | "bus") ?? "flight"
    );
    const [from, setFrom] = useState<PlaceValue | null>(
        transport.from ?? null
    );
    const [to, setTo] = useState<PlaceValue | null>(
        transport.to ?? null
    );
    const [departureTime, setDepartureTime] = useState(
        transport.departureTime
            ? new Date(transport.departureTime).toISOString().slice(0, 16)
            : ""
    );
    const [arrivalTime, setArrivalTime] = useState(
        transport.arrivalTime
            ? new Date(transport.arrivalTime).toISOString().slice(0, 16)
            : ""
    );
    const [notes, setNotes] = useState(transport.notes ?? "");
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!from) { setError("Bitte einen Startort auswählen"); return; }
        if (!to) { setError("Bitte einen Zielort auswählen"); return; }
        setError(null);

        onSaveAction({
            from,
            to,
            departureTime: departureTime ? new Date(departureTime).toISOString() : undefined,
            arrivalTime: arrivalTime ? new Date(arrivalTime).toISOString() : undefined,
            type,
            notes: notes || undefined,
        });
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-3xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800 max-h-[90vh] overflow-y-auto">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                    Transport bearbeiten
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Type */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Typ <span className="text-red-500">*</span>
                        </label>
                        <select
                            value={type}
                            onChange={(e) => setType(e.target.value as typeof type)}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                        >
                            <option value="flight">✈️ Flug</option>
                            <option value="train">🚂 Zug</option>
                            <option value="car">🚗 Auto</option>
                            <option value="bus">🚌 Bus</option>
                        </select>
                    </div>

                    {/* From */}
                    <PlaceAutocomplete
                        label="Von"
                        value={from}
                        onChange={setFrom}
                        placeholder="z.B. Flughafen Frankfurt, Berlin..."
                        required
                    />

                    {/* To */}
                    <PlaceAutocomplete
                        label="Nach"
                        value={to}
                        onChange={setTo}
                        placeholder="z.B. Paris, Rom, Wien..."
                        required
                    />

                    {/* Departure & Arrival */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Abfahrt</label>
                            <input
                                type="datetime-local"
                                value={departureTime}
                                onChange={(e) => setDepartureTime(e.target.value)}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Ankunft</label>
                            <input
                                type="datetime-local"
                                value={arrivalTime}
                                onChange={(e) => setArrivalTime(e.target.value)}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                    </div>

                    {/* Notes */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Notizen</label>
                        <textarea
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                            rows={3}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none"
                        />
                    </div>

                    {/* Error */}
                    {error && (
                        <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
                    )}

                    {/* Buttons */}
                    <div className="flex gap-3 pt-4">
                        <button
                            type="button"
                            onClick={onDeleteAction}
                            className="px-4 py-2 text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 hover:bg-red-100 dark:hover:bg-red-950/50 rounded-lg font-medium transition-colors"
                        >
                            Löschen
                        </button>
                        <button
                            type="button"
                            onClick={onCloseAction}
                            className="flex-1 px-4 py-2 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-lg font-medium transition-colors"
                        >
                            Abbrechen
                        </button>
                        <button
                            type="submit"
                            className="flex-1 px-4 py-2 text-white bg-sky-600 hover:bg-sky-700 rounded-lg font-medium transition-colors"
                        >
                            Speichern
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}