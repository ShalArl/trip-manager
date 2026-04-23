"use client";

import { useState } from "react";
import { CreateLocationRequest } from "@/types/location";

type Props = {
    isOpen: boolean;
    onCloseAction: () => void;
    onAddAction: (location: CreateLocationRequest) => void;
};

export default function AddLocationModal({ isOpen, onCloseAction, onAddAction }: Props) {
    const [formData, setFormData] = useState({
        name: "",
        city: "",
        country: "",
        shortDescription: "",
        dateFrom: "",
        dateTo: "",
        latitude: "",
        longitude: "",
        notes: "",
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.name || !formData.city || !formData.country || !formData.shortDescription || !formData.dateFrom || !formData.dateTo) {
            alert("Bitte alle erforderlichen Felder ausfüllen");
            return;
        }

        onAddAction({
            name: formData.name,
            city: formData.city,
            country: formData.country,
            shortDescription: formData.shortDescription,
            dateFrom: formData.dateFrom,
            dateTo: formData.dateTo,
            latitude: formData.latitude ? parseFloat(formData.latitude) : undefined,
            longitude: formData.longitude ? parseFloat(formData.longitude) : undefined,
            notes: formData.notes || undefined,
        });

        setFormData({ name: "", city: "", country: "", shortDescription: "", dateFrom: "", dateTo: "", latitude: "", longitude: "", notes: "" });
        onCloseAction();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-3xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                    Neuen Ort hinzufügen
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Name */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Ortsname *</label>
                        <input type="text" value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} placeholder="z.B. Paris" className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500" />
                    </div>

                    {/* City & Country */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Stadt *</label>
                            <input type="text" value={formData.city} onChange={(e) => setFormData({ ...formData, city: e.target.value })} placeholder="z.B. Paris" className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500" />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Land *</label>
                            <input type="text" value={formData.country} onChange={(e) => setFormData({ ...formData, country: e.target.value })} placeholder="z.B. Frankreich" className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500" />
                        </div>
                    </div>

                    {/* Short Description */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Kurzbeschreibung *</label>
                        <input type="text" value={formData.shortDescription} onChange={(e) => setFormData({ ...formData, shortDescription: e.target.value })} placeholder="z.B. Hauptstadt Frankreichs" className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500" />
                    </div>

                    {/* Date From & To */}
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

                    {/* Latitude & Longitude */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Breite</label>
                            <input type="number" step="any" value={formData.latitude} onChange={(e) => setFormData({ ...formData, latitude: e.target.value })} placeholder="z.B. 48.85" className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500" />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Länge</label>
                            <input type="number" step="any" value={formData.longitude} onChange={(e) => setFormData({ ...formData, longitude: e.target.value })} placeholder="z.B. 2.29" className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500" />
                        </div>
                    </div>

                    {/* Notes */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Notizen</label>
                        <textarea value={formData.notes} onChange={(e) => setFormData({ ...formData, notes: e.target.value })} placeholder="z.B. Hotel im 8. Arrondissement" rows={3} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none" />
                    </div>

                    {/* Buttons */}
                    <div className="flex gap-3 pt-4">
                        <button type="button" onClick={onCloseAction} className="flex-1 px-4 py-2 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-lg font-medium transition-colors">Abbrechen</button>
                        <button type="submit" className="flex-1 px-4 py-2 text-white bg-sky-600 hover:bg-sky-700 rounded-lg font-medium transition-colors">Hinzufügen</button>
                    </div>
                </form>
            </div>
        </div>
    );
}