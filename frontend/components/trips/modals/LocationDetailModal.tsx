"use client";

import React, { useState } from "react";
import { LocationResponse, UpdateLocationRequest } from "@/types/location";

type Props = {
    isOpen: boolean;
    location: LocationResponse;
    onCloseAction: () => void;
    onSaveAction: (req: UpdateLocationRequest) => void;
    onDeleteAction: () => void;
};

export default function LocationDetailModal({
    isOpen,
    location,
    onCloseAction,
    onSaveAction,
    onDeleteAction,
}: Props) {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
        name: location.name,
        city: location.city,
        country: location.country,
        shortDescription: location.shortDescription,
        dateFrom: location.dateFrom,
        dateTo: location.dateTo,
        notes: location.notes ?? "",
    });

    const handleSubmit = (e: React.SubmitEvent) => {
        e.preventDefault();
        onSaveAction({
            name: formData.name,
            city: formData.city,
            country: formData.country,
            shortDescription: formData.shortDescription,
            dateFrom: formData.dateFrom,
            dateTo: formData.dateTo,
            notes: formData.notes || undefined,
        });
        setIsEditing(false);
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-2xl shadow-xl max-w-md w-full border border-zinc-100 dark:border-zinc-800 overflow-hidden">
                {/* Header */}
                <div className="bg-gradient-to-r from-sky-500 to-sky-600 px-8 py-6">
                    <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                        📍 {location.name}
                    </h2>
                    <p className="text-sky-100 text-sm mt-1">
                        {location.city}, {location.country}
                    </p>
                </div>

                <div className="p-8">
                    {!isEditing ? (
                        // Details Ansicht
                        <div className="space-y-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Kurzbeschreibung</p>
                                <p className="text-zinc-700 dark:text-zinc-300">{location.shortDescription}</p>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Von</p>
                                    <p className="text-zinc-700 dark:text-zinc-300">{location.dateFrom}</p>
                                </div>
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Bis</p>
                                    <p className="text-zinc-700 dark:text-zinc-300">{location.dateTo}</p>
                                </div>
                            </div>
                            {location.notes && (
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Notizen</p>
                                    <p className="text-zinc-700 dark:text-zinc-300">{location.notes}</p>
                                </div>
                            )}

                            <div className="flex gap-3 pt-4 border-t border-zinc-100 dark:border-zinc-800">
                                <button
                                    onClick={onCloseAction}
                                    className="flex-1 px-4 py-3 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-xl font-semibold transition-all"
                                >
                                    Schließen
                                </button>
                                <button
                                    onClick={() => setIsEditing(true)}
                                    className="flex-1 px-4 py-3 text-white bg-gradient-to-r from-sky-600 to-sky-700 hover:from-sky-700 hover:to-sky-800 rounded-xl font-semibold transition-all"
                                >
                                    ✏️ Bearbeiten
                                </button>
                                <button
                                    onClick={onDeleteAction}
                                    className="px-4 py-3 text-white bg-red-500 hover:bg-red-600 rounded-xl font-semibold transition-all"
                                >
                                    🗑️
                                </button>
                            </div>
                        </div>
                    ) : (
                        // Edit Ansicht
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Name *</label>
                                    <input type="text" value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Stadt *</label>
                                    <input type="text" value={formData.city} onChange={(e) => setFormData({ ...formData, city: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Land *</label>
                                <input type="text" value={formData.country} onChange={(e) => setFormData({ ...formData, country: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Kurzbeschreibung *</label>
                                <input type="text" value={formData.shortDescription} onChange={(e) => setFormData({ ...formData, shortDescription: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Von *</label>
                                    <input type="date" value={formData.dateFrom} onChange={(e) => setFormData({ ...formData, dateFrom: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Bis *</label>
                                    <input type="date" value={formData.dateTo} onChange={(e) => setFormData({ ...formData, dateTo: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Notizen</label>
                                <textarea value={formData.notes} onChange={(e) => setFormData({ ...formData, notes: e.target.value })} rows={3} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none" />
                            </div>
                            <div className="flex gap-3 pt-4 border-t border-zinc-100 dark:border-zinc-800">
                                <button type="button" onClick={() => setIsEditing(false)} className="flex-1 px-4 py-3 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-xl font-semibold transition-all">
                                    Abbrechen
                                </button>
                                <button type="submit" className="flex-1 px-4 py-3 text-white bg-gradient-to-r from-sky-600 to-sky-700 hover:from-sky-700 hover:to-sky-800 rounded-xl font-semibold transition-all">
                                    ✓ Speichern
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}