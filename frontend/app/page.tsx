"use client";

import { useState, useEffect } from "react";
import { getTrips } from "@/lib/api/trips";
import { useUserContext } from "@/lib/context/UserContext";
import { logout } from "@/lib/api/auth";
import Navbar from "@/components/global/Navbar";
import Hero from "@/components/home/Hero";
import FeatureGrid from "@/components/home/FeatureGrid";
import TripList from "@/components/trips/TripList";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {TripResponse} from "@/types/trip";
import {LoadingSpinner} from "@/components/global/LoadingSpinner";


export default function Home() {
  const { user, isLoading, updateUser } = useUserContext();
  const [trips, setTrips] = useState<TripResponse[]>([]);
  const router = useRouter();

  useEffect(() => {
    if (user) {
      getTrips().then(setTrips).catch(console.error);
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
    return <LoadingSpinner />;
  }

  return (
    <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
      <Hero />
      <div className="mx-auto max-w-4xl px-6 pt-6 flex justify-end">
        <Link
          href="/trips/new"
          className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
        >
          + Neue Reise
        </Link>
      </div>
      <TripList trips={trips} />
      <FeatureGrid />
    </div>
  );
}