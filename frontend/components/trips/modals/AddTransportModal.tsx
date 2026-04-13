"use client";

import { useState } from "react";

type Props = {
    isOpen: boolean;
    onCloseAction: () => void;
    onAddAction: (transport: {
        from: string;
        to: string;
        date: string;
        type: string;
        notes?: string;
    }) => void;
};

export default function AddTransportModal({ isOpen, onCloseAction, onAddAction }: Props) {
    const [formData, setFormData] = useState({
        from: "",
        to: "",
        date: "",
        type: "flight",
        notes: "",
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.from || !formData.to || !formData.date) {
            alert("Bitte alle erforderlichen Felder ausfüllen");
            return;
        }

        onAddAction({
            from: formData.from,
            to: formData.to,
            date: formData.date,
            type: formData.type,
            notes: formData.notes || undefined,
        });

        setFormData({ from: "", to: "", date: "", type: "flight", notes: "" });
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
                            onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                        >
                            <option value="flight">✈️ Flug</option>
                            <option value="train">🚂 Zug</option>
                            <option value="car">🚗 Auto</option>
                            <option value="bus">🚌 Bus</option>
                        </select>
                    </div>

                    {/* From & To */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Von *
                            </label>
                            <input
                                type="text"
                                value={formData.from}
                                onChange={(e) => setFormData({ ...formData, from: e.target.value })}
                                placeholder="z.B. Zürich"
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                Nach *
                            </label>
                            <input
                                type="text"
                                value={formData.to}
                                onChange={(e) => setFormData({ ...formData, to: e.target.value })}
                                placeholder="z.B. Oslo"
                                className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500"
                            />
                        </div>
                    </div>

                    {/* Date */}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                            Datum *
                        </label>
                        <input
                            type="date"
                            value={formData.date}
                            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                            className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500"
                        />
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