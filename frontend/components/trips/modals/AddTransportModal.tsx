"use client";

import { useState } from "react";
import { CreateTransportRequest } from "@/types/transport";

type Location = {
    id: string;
    name: string;
};

type Props = {
    isOpen: boolean;
    locations: Location[];
    onCloseAction: () => void;
    onAddAction: (transport: CreateTransportRequest) => void;
};

export default function AddTransportModal({ isOpen, locations, onCloseAction, onAddAction }: Props) {
    const [formData, setFormData] = useState({
        fromLocationId: "",
        toLocationId: "",
        departureTime: "",
        arrivalTime: "",
        type: "flight" as "flight" | "train" | "car" | "bus",
        notes: "",
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.fromLocationId || !formData.toLocationId) {
            alert("Bitte Von und Nach auswählen");
            return;
        }

        onAddAction({
            fromLocationId: formData.fromLocationId,
            toLocationId: formData.toLocationId,
            departureTime: formData.departureTime ? `${formData.departureTime}:00Z` : undefined,
            arrivalTime: formData.arrivalTime ? `${formData.arrivalTime}:00Z` : undefined,
            type: formData.type,
            notes: formData.notes || undefined,
        });

        setFormData({ fromLocationId: "", toLocationId: "", departureTime: "", arrivalTime: "", type: "flight", notes: "" });
        onCloseAction();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-3xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                    Transport hinzufügen
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Type */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Typ *
                        </label>
                        <select
                            value={formData.type}
                            onChange={(e) => setFormData({ ...formData, type: e.target.value as "flight" | "train" | "car" | "bus" })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                        >
                            <option value="flight">✈️ Flug</option>
                            <option value="train">🚂 Zug</option>
                            <option value="car">🚗 Auto</option>
                            <option value="bus">🚌 Bus</option>
                        </select>
                    </div>

                    {/* From & To als Dropdown */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Von *
                            </label>
                            <select
                                value={formData.fromLocationId}
                                onChange={(e) => setFormData({ ...formData, fromLocationId: e.target.value })}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            >
                                <option value="">Auswählen</option>
                                {locations.map((loc) => (
                                    <option key={loc.id} value={loc.id}>{loc.name}</option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Nach *
                            </label>
                            <select
                                value={formData.toLocationId}
                                onChange={(e) => setFormData({ ...formData, toLocationId: e.target.value })}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            >
                                <option value="">Auswählen</option>
                                {locations.map((loc) => (
                                    <option key={loc.id} value={loc.id}>{loc.name}</option>
                                ))}
                            </select>
                        </div>
                    </div>

                    {/* Departure & Arrival */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Abfahrt
                            </label>
                            <input
                                type="datetime-local"
                                value={formData.departureTime}
                                onChange={(e) => setFormData({ ...formData, departureTime: e.target.value })}
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Ankunft
                            </label>
                            <input
                                type="datetime-local"
                                value={formData.arrivalTime}
                                onChange={(e) => setFormData({ ...formData, arrivalTime: e.target.value })}
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
                            value={formData.notes}
                            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                            placeholder="z.B. Buchungsnummer AB1234"
                            rows={3}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none"
                        />
                    </div>

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