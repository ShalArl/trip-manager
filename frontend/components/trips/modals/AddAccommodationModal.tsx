"use client";

import { useState } from "react";
import { CreateAccommodationRequest } from "@/types/accommodation";
import PlaceAutocomplete, { PlaceValue } from "@/components/shared/PlaceAutocomplete";

type Props = {
    isOpen: boolean;
    onCloseAction: () => void;
    onAddAction: (accommodation: CreateAccommodationRequest) => void;
};

export default function AddAccommodationModal({ isOpen, onCloseAction, onAddAction }: Props) {
    const [name, setName] = useState("");
    const [location, setLocation] = useState<PlaceValue | null>(null);
    const [address, setAddress] = useState("");
    const [checkIn, setCheckIn] = useState("");
    const [checkOut, setCheckOut] = useState("");
    const [pricePerNight, setPricePerNight] = useState("");
    const [notes, setNotes] = useState("");
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!name.trim()) { setError("Bitte einen Namen eingeben"); return; }
        if (!location) { setError("Bitte einen Ort auswählen"); return; }
        setError(null);

        onAddAction({
            location,
            name,
            address: address || undefined,
            checkIn: checkIn ? new Date(checkIn).toISOString() : undefined,
            checkOut: checkOut ? new Date(checkOut).toISOString() : undefined,
            pricePerNight: pricePerNight ? parseFloat(pricePerNight) : undefined,
            notes: notes || undefined,
        });

        // Reset
        setName("");
        setLocation(null);
        setAddress("");
        setCheckIn("");
        setCheckOut("");
        setPricePerNight("");
        setNotes("");
        onCloseAction();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-3xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800 max-h-[90vh] overflow-y-auto">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                    Unterkunft hinzufügen
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Name */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Name <span className="text-red-500">*</span>
                        </label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="z.B. Hotel Muster"
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Location */}
                    <PlaceAutocomplete
                        label="Ort"
                        value={location}
                        onChange={setLocation}
                        placeholder="z.B. Paris, Amsterdam..."
                        required
                    />

                    {/* Address */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Adresse
                        </label>
                        <input
                            type="text"
                            value={address}
                            onChange={(e) => setAddress(e.target.value)}
                            placeholder="z.B. Musterstraße 1, 12345 Musterstadt"
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Check-in & Check-out */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Check-in</label>
                            <input
                                type="datetime-local"
                                value={checkIn}
                                onChange={(e) => setCheckIn(e.target.value)}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Check-out</label>
                            <input
                                type="datetime-local"
                                value={checkOut}
                                onChange={(e) => setCheckOut(e.target.value)}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                    </div>

                    {/* Price */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Preis pro Nacht (€)
                        </label>
                        <input
                            type="number"
                            min="0"
                            step="0.01"
                            value={pricePerNight}
                            onChange={(e) => setPricePerNight(e.target.value)}
                            placeholder="z.B. 89.99"
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Notes */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Notizen</label>
                        <textarea
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                            placeholder="z.B. Buchungsnummer AB1234"
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