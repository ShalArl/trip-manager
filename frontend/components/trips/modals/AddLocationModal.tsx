"use client";

import { useState } from "react";
import { CreateLocationRequest } from "@/types/location";
import PlaceAutocomplete, { PlaceValue } from "@/components/shared/PlaceAutocomplete";

type Props = {
    isOpen: boolean;
    onCloseAction: () => void;
    onAddAction: (location: CreateLocationRequest) => void;
};

export default function AddLocationModal({ isOpen, onCloseAction, onAddAction }: Props) {
    const [place, setPlace] = useState<PlaceValue | null>(null);
    const [shortDescription, setShortDescription] = useState("");
    const [dateFrom, setDateFrom] = useState("");
    const [dateTo, setDateTo] = useState("");
    const [notes, setNotes] = useState("");
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!place) { setError("Bitte einen Ort auswählen"); return; }
        if (!shortDescription.trim()) { setError("Bitte eine Kurzbeschreibung eingeben"); return; }
        if (!dateFrom || !dateTo) { setError("Bitte Von- und Bis-Datum angeben"); return; }
        setError(null);

        onAddAction({
            name: place.name,
            city: place.city,
            country: place.country,
            latitude: place.lat ?? undefined,
            longitude: place.lng ?? undefined,
            shortDescription,
            dateFrom,
            dateTo,
            notes: notes || undefined,
        });

        // Reset
        setPlace(null);
        setShortDescription("");
        setDateFrom("");
        setDateTo("");
        setNotes("");
        onCloseAction();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-3xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800 max-h-[90vh] overflow-y-auto">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                    Neuen Ort hinzufügen
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Place Autocomplete */}
                    <PlaceAutocomplete
                        label="Ort"
                        value={place}
                        onChange={setPlace}
                        placeholder="z.B. Paris, Tokio, New York..."
                        required
                    />

                    {/* Short Description */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Kurzbeschreibung <span className="text-red-500">*</span>
                        </label>
                        <input
                            type="text"
                            value={shortDescription}
                            onChange={(e) => setShortDescription(e.target.value)}
                            placeholder="z.B. Hauptstadt Frankreichs"
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Date From & To */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Von <span className="text-red-500">*</span>
                            </label>
                            <input
                                type="date"
                                value={dateFrom}
                                onChange={(e) => setDateFrom(e.target.value)}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Bis <span className="text-red-500">*</span>
                            </label>
                            <input
                                type="date"
                                value={dateTo}
                                onChange={(e) => setDateTo(e.target.value)}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                    </div>

                    {/* Notes */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Notizen
                        </label>
                        <textarea
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                            placeholder="z.B. Hotel im 8. Arrondissement"
                            rows={3}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none"
                        />
                    </div>

                    {/* Error */}
                    {error && <p className="text-sm text-red-600 dark:text-red-400">{error}</p>}

                    {/* Buttons */}
                    <div className="flex gap-3 pt-4">
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
                            Hinzufügen
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}