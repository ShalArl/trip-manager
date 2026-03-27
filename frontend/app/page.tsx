"use client";

import { useState } from "react";
import { User } from "@/types/user";
import { Trip } from "@/types/trip"
import { mockTrips } from "@/lib/mock-trips";

import AuthPage from "@/components/auth/AuthPage";
import Navbar from "@/components/home/Navbar";
import Hero from "@/components/home/Hero";
import FeatureGrid from "@/components/home/FeatureGrid";
import TripList from "@/components/trips/TripList";

export default function Home() {
  const [user, setUser] = useState<User | null>(() => {
    if (typeof window === "undefined") return null;
    const saved = localStorage.getItem("user");
    return saved ? JSON.parse(saved) : null;
  });

  const [trips] = useState<Trip[]>(mockTrips);

  const handleAuth = (user: User) => {
    localStorage.setItem("user", JSON.stringify(user));
    setUser(user);
  }

  const handleLogout = () => {
    localStorage.removeItem("user");
    setUser(null);
  };

  if (!user) {
    return <AuthPage onAuthAction={handleAuth} />;
  }

  return (
    <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
      <Navbar user={user} onLogout={handleLogout} />
      <Hero />
      <TripList trips={trips} />
      <FeatureGrid />
    </div>
  );
}