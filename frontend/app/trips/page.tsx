"use client";
import {useEffect, useState} from "react";
import {getTrips} from "@/lib/api/trips";
import {getTenantSettings, logout} from "@/lib/api/auth";
import {useUserContext} from "@/lib/context/UserContext";
import Hero from "@/components/home/Hero";
import FeatureGrid from "@/components/home/FeatureGrid";
import TripList from "@/components/trips/TripList";
import Link from "next/link";
import {useRouter} from "next/navigation";
import {TripResponse} from "@/types/trip";
import {LoadingSpinner} from "@/components/global/LoadingSpinner";

export default function Home() {
    const {user, isLoading, updateUser} = useUserContext();
    const [trips, setTrips] = useState<TripResponse[]>([]);
    const [maxActiveTrips, setMaxActiveTrips] = useState(0);
    const router = useRouter();

    useEffect(() => {
        if (user) {
            getTrips().then(setTrips).catch(console.error);
            getTenantSettings().then((s) => setMaxActiveTrips(s?.maxActiveTrips ?? 0)).catch(console.error);
        }
    }, [user]);

    useEffect(() => {
        if (!isLoading && !user) {
            router.push("/search");
        }
    }, [user, isLoading, router]);

    const handleLogout = async () => {
        await logout()
        updateUser(null);
        router.push("/search");
    };

    if (isLoading) {
        return <LoadingSpinner/>;
    }

    const activeTripsCount = trips.filter(
        (t) => t.status === "planned" || t.status === "ongoing"
    ).length;
    const limitReached = maxActiveTrips > 0 && activeTripsCount >= maxActiveTrips;

    return (
        <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
            <Hero/>
            <div className="mx-auto max-w-4xl px-6 pt-6 flex justify-end items-center gap-3">
                {limitReached && (
                    <span className="text-xs text-zinc-500 dark:text-zinc-400">
            Limit erreicht ({activeTripsCount}/{maxActiveTrips} aktive Reisen) ·{" "}
                        <Link href="/settings" className="text-[var(--brand-primary)] hover:underline">
              Upgrade
            </Link>
          </span>
                )}
                {limitReached ? (
                    <span
                        className="px-4 py-2 text-sm font-medium bg-zinc-200 dark:bg-zinc-800 text-zinc-400 dark:text-zinc-600 rounded-lg cursor-not-allowed"
                        title={`Du hast das Limit von ${maxActiveTrips} aktiven Reisen erreicht`}
                    >
            + Neue Reise
          </span>
                ) : (
                    <Link
                        href="/trips/new"
                        className="px-4 py-2 text-sm font-medium bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white rounded-lg transition-colors"
                    >
                        + Neue Reise
                    </Link>
                )}
            </div>
            <TripList trips={trips}/>
            <FeatureGrid/>
        </div>
    );
}