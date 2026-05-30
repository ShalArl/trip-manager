"use client";

import { useState } from "react";
import { AccommodationResponse, UpdateAccommodationRequest } from "@/types/accommodation";

type Location = {
    id: string;
    name: string;
};

type Props = {
    isOpen: boolean;
    accommodation: AccommodationResponse;
    locations: Location[];
    onCloseAction: () => void;
    onSaveAction: (req: UpdateAccommodationRequest) => void;
    onDeleteAction: () => void;
};

export default function EditAccommodationModal({ isOpen, accommodation, locations, onCloseAction, onSaveAction, onDeleteAction }: Props) {
    const [formData, setFormData] = useState({
        locationId: accommodation.locationId ?? "",
        name: accommodation.name ?? "",
        address: accommodation.address ?? "",
        checkIn: accommodation.checkIn ?? "",
        checkOut: accommodation.checkOut ?? "",
        pricePerNight: accommodation.pricePerNight?.toString() ?? "",
        notes: accommodation.notes ?? "",
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.locationId) {
            alert("Bitte einen Ort auswählen");
            return;
        }
        if (!formData.name.trim()) {
            alert("Bitte einen Namen eingeben");
            return;
        }
        onSaveAction({
            locationId: formData.locationId,
            name: formData.name,
            address: formData.address || undefined,
            checkIn: formData.checkIn ? formData.checkIn + ":00Z" : undefined,
            checkOut: formData.checkOut ? formData.checkOut + ":00Z" : undefined,
            pricePerNight: formData.pricePerNight ? parseFloat(formData.pricePerNight) : undefined,
            notes: formData.notes || undefined,
        });
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-3xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                    Unterkunft bearbeiten
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Name */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Name *</label>
                        <input
                            type="text"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Location */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Ort *</label>
                        <select
                            value={formData.locationId}
                            onChange={(e) => setFormData({ ...formData, locationId: e.target.value })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                        >
                            <option value="">Auswählen</option>
                            {locations.map((loc) => (
                                <option key={loc.id} value={loc.id}>{loc.name}</option>
                            ))}
                        </select>
                    </div>

                    {/* Address */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Adresse</label>
                        <input
                            type="text"
                            value={formData.address}
                            onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Check-in & Check-out */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Check-in</label>
                            <input
                                type="datetime-local"
                                value={formData.checkIn}
                                onChange={(e) => setFormData({ ...formData, checkIn: e.target.value })}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Check-out</label>
                            <input
                                type="datetime-local"
                                value={formData.checkOut}
                                onChange={(e) => setFormData({ ...formData, checkOut: e.target.value })}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                    </div>

                    {/* Price per night */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Preis pro Nacht (€)</label>
                        <input
                            type="number"
                            min="0"
                            step="0.01"
                            value={formData.pricePerNight}
                            onChange={(e) => setFormData({ ...formData, pricePerNight: e.target.value })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
                    </div>

                    {/* Notes */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Notizen</label>
                        <textarea
                            value={formData.notes}
                            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                            rows={3}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none"
                        />
                    </div>

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