import { Trip } from "@/types/trip";
import Link from "next/link";

type Props = {
    trip: Trip;
};

export default function TripDetail({ trip }: Props) {
    return (
        <div className="mx-auto max-w-3xl px-6 py-12">
            <Link
                href="/"
                className="inline-flex items-center gap-2 text-sm text-zinc-500 dark:text-zinc-400 hover:text-sky-600 dark:hover:text-sky-400 transition-colors mb-8"
            >
                ← Zurück zur Übersicht
            </Link>

            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">

                <div className="flex items-center gap-4 mb-6">
                    <div className="w-14 h-14 rounded-2xl bg-sky-50 dark:bg-sky-950/50 border border-sky-100 dark:border-sky-900 flex items-center justify-center text-3xl">
                        ✈️
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
                            {trip.title}
                        </h1>
                        <p className="text-sm text-zinc-500 dark:text-zinc-400">
                            {trip.destination} · {trip.startDate}
                        </p>
                    </div>
                </div>

                <div className="border-t border-zinc-100 dark:border-zinc-800 pt-6 space-y-6">
                    <div>
                        <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1">
                            Kurzbeschreibung
                        </p>
                        <p className="text-zinc-700 dark:text-zinc-300">
                            {trip.shortDescription}
                        </p>
                    </div>

                    <div>
                        <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1">
                            Details
                        </p>
                        <p className="text-zinc-700 dark:text-zinc-300 leading-relaxed">
                            {trip.description}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}