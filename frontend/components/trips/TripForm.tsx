"use client"

import React, { useState } from "react";

import { components } from "@/generated/types";
type CreateTripRequest = components["schemas"]["CreateTripRequest"];

type Props = {
    onCreateTripAction: (trip: CreateTripRequest) => void;
};

export default function TripForm({ onCreateTripAction }: Props) {
    const [title, setTitle] = useState("");
    const [startDate, setStartDate] = useState("");
    const [endDate, setEndDate] = useState("");
    const [shortDescription, setShortDescription] = useState("");
    const [longDescription, setLongDescription] = useState("");
    const [error, setError] = useState("");


    function handleSubmit(e: React.SubmitEvent) {
        e.preventDefault();
        setError("");

        if (!title.trim() || !startDate.trim() || !shortDescription.trim() || !longDescription.trim()) {
            setError("Bitte alle Felder ausfüllen.");
            return;
        }

        onCreateTripAction({ title, startDate, endDate, shortDescription, description: longDescription });
    }

    return (
        <div className="max-w-2xl mx-auto px-6 py-12">
            <h1 className="text-3xl font-bold text-zinc-900 dark:text-white mb-2">
                Neue Reise erstellen
            </h1>
            <p className="text-zinc-500 dark:text-zinc-400 text-sm mb-8">
                Füll die Felder aus um deine Reise anzulegen.
            </p>

            <form onSubmit={handleSubmit} className="space-y-5">

                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        Titel
                    </label>
                    <input
                        type="text"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        placeholder="Familienreise nach Norwegen"
                        className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        Startdatum
                    </label>
                    <input
                        type="date"
                        value={startDate}
                        onChange={(e) => setStartDate(e.target.value)}
                        className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        Enddatum
                    </label>
                    <input
                        type="date"
                        value={endDate}
                        onChange={(e) => setEndDate(e.target.value)}
                        className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        Kurzbeschreibung
                        <span className="ml-2 text-zinc-400 font-normal">({shortDescription.length}/80)</span>
                    </label>
                    <input
                        type="text"
                        value={shortDescription}
                        onChange={(e) => setShortDescription(e.target.value)}
                        placeholder="Erkunde die Fjorde Südnorwegens."
                        maxLength={80}
                        className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        Detaillierte Beschreibung
                    </label>
                    <textarea
                        value={longDescription}
                        onChange={(e) => setLongDescription(e.target.value)}
                        placeholder="Beschreibe deine Reise im Detail..."
                        rows={5}
                        className="w-full px-4 py-3 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm resize-none"
                    />
                </div>

                {error && (
                    <p className="text-sm text-red-500 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-800/50 rounded-xl px-4 py-3">
                        {error}
                    </p>
                )}

                <button
                    type="submit"
                    className="w-full h-12 rounded-xl bg-sky-600 hover:bg-sky-700 active:scale-[0.98] text-white font-semibold text-sm transition-all shadow-md shadow-sky-500/20"
                >
                    Reise erstellen
                </button>

            </form>
        </div>
    )
}