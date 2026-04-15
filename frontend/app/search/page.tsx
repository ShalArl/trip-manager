"use client";

import { useState, useEffect } from "react";
import { useUserContext } from "@/lib/context/UserContext";
import { searchTrips, getPublicTrips } from "@/lib/api/trips";
import { components } from "@/generated/types";
import Navbar from "@/components/global/Navbar";
import Link from "next/link";
import { useRouter } from "next/navigation";

type TripResponse = components["schemas"]["TripResponse"];

export default function SearchPage() {
    const { user, updateUser } = useUserContext();
    const [query, setQuery] = useState("");
    const [trips, setTrips] = useState<TripResponse[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const router = useRouter();

    const handleLogout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("userId");
        localStorage.removeItem("user");
        updateUser(null);
    };

    useEffect(() => {
        const fetchTrips = async () => {
            setIsLoading(true);
            try {
                const results = query
                    ? await searchTrips(query)
                    : await getPublicTrips();
                setTrips(results);
            } catch (error) {
                console.error(error);
            } finally {
                setIsLoading(false);
            }
        };

        const debounce = setTimeout(fetchTrips, 300);
        return () => clearTimeout(debounce);
    }, [query]);

    return (
        <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
            <Navbar user={user} onLogout={handleLogout} />

            <div className="mx-auto max-w-4xl px-6 py-12">
                <h1 className="text-3xl font-bold text-zinc-900 dark:text-white mb-2">
                    Reisen entdecken
                </h1>
                <p className="text-zinc-500 dark:text-zinc-400 mb-8">
                    Entdecke Reisen von anderen Reisenden
                </p>

                {/* Suchleiste */}
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="Nach Reisen suchen..."
                    className="w-full px-4 py-3 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 mb-8"
                />

                {/* Ergebnisse */}
                {isLoading ? (
                    <p className="text-zinc-500 text-center">Suche...</p>
                ) : trips.length === 0 ? (
                    <p className="text-zinc-500 dark:text-zinc-400 text-center">
                        {query ? "Keine Reisen gefunden" : "Gib etwas ein um zu suchen"}
                    </p>
                ) : (
                    <div className="flex flex-col gap-4">
                        {trips.map((trip) => (
                            <Link
                                key={trip.id}
                                href={`/trips/${encodeURIComponent(trip.id)}`}
                                className="group bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl px-6 py-5 flex items-center justify-between hover:border-sky-400 dark:hover:border-sky-600 hover:shadow-md transition-all"
                            >
                                <div className="flex items-center gap-4">
                                    <div className="w-10 h-10 rounded-xl bg-sky-50 dark:bg-sky-950/50 border border-sky-100 dark:border-sky-900 flex items-center justify-center text-xl">
                                        ✈️
                                    </div>
                                    <div>
                                        <p className="font-semibold text-zinc-900 dark:text-white group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors">
                                            {trip.title}
                                        </p>
                                        <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                            {trip.startDate} · {trip.endDate}
                                        </p>
                                    </div>
                                </div>
                                <span className="text-zinc-400 dark:text-zinc-600 group-hover:text-sky-500 transition-colors text-lg">
                                    →
                                </span>
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}