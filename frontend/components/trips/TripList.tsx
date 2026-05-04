import Link from "next/link";

import { components } from "@/generated/types";
type TripResponse = components["schemas"]["TripResponse"];

type Props = {
    trips: TripResponse[];
};

export default function TripList({ trips }: Props) {
    if (trips.length === 0) {
        return (
            <div className="mx-auto max-w-4xl px-6 py-12 text-center">
                <p className="text-zinc-500 dark:text-zinc-400">Noch keine Reisen geplant.</p>
                <Link
                    href="/trips/new"
                    className="mt-4 inline-flex items-center text-sky-600 dark:text-sky-400 font-medium hover:underline"
                >
                    Erste Reise erstellen →
                </Link>
            </div>
        );
    }

    return (
        <div className="mx-auto max-w-4xl px-6 py-12">
            <h2 className="text-2xl font-bold text-zinc-900 dark:text-white mb-6">
                Meine Reisen
            </h2>
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
        </div>
    );
}