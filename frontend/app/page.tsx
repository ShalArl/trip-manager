"use client";

import { useState } from "react";
import { User } from "@/types/user";
import AuthPage from "@/components/auth/AuthPage";
import Navbar from "@/components/home/Navbar";
import Hero from "@/components/home/Hero";
import FeatureGrid from "@/components/home/FeatureGrid";

export default function Home() {
  const [user, setUser] = useState<User | null>(null);

  if (!user) {
    return <AuthPage onAuth={setUser} />;
  }

  return (
    <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
      <Navbar user={user} onLogout={() => setUser(null)} />
      <Hero onCTA={() => console.log("TODO: Reise planen")} />
      <FeatureGrid />
    </div>
  );
}