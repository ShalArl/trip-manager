"use client";

import ProfileSettings from "@/components/settings/ProfileSettings";
import PasswordSettings from "@/components/settings/PasswordSettings";
import { useEffect, useState } from "react";
import { UserResponse } from "@/types/user";
import { getMe } from "@/lib/api/auth";
import { User, Lock } from "lucide-react";

export default function SettingsPage() {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState("profile");

  useEffect(() => {
    fetchUser();
  }, []);

  const fetchUser = async () => {
    try {
      const userData = await getMe();
      setUser(userData);
    } catch (err) {
      setError("Fehler beim Laden der Benutzerinformationen");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-zinc-50 to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
        <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="animate-pulse space-y-4">
            <div className="h-10 w-48 bg-zinc-300 dark:bg-zinc-700 rounded-lg" />
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              <div className="h-40 bg-zinc-300 dark:bg-zinc-700 rounded-xl" />
              <div className="md:col-span-3 h-96 bg-zinc-300 dark:bg-zinc-700 rounded-xl" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-zinc-50 to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
        <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="rounded-xl bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-900 p-6">
            <p className="text-sm font-medium text-red-800 dark:text-red-200">
              {error || "Fehler beim Laden der Einstellungen"}
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-zinc-50 to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8 sm:mb-12">
          <h1 className="text-4xl sm:text-5xl font-bold bg-gradient-to-r from-zinc-900 to-zinc-700 dark:from-white dark:to-zinc-300 bg-clip-text text-transparent">
            Einstellungen
          </h1>
          <p className="mt-3 text-base text-zinc-600 dark:text-zinc-400">
            Verwalte dein Profil, Passwort und Kontodaten
          </p>
        </div>

        {/* Sidebar + Content Layout */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="md:col-span-1">
            <div className="bg-white dark:bg-zinc-950 rounded-xl border border-zinc-200 dark:border-zinc-800 shadow-lg p-2 sticky top-20">
              <div className="space-y-2">
                <button
                  onClick={() => setActiveTab("profile")}
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left font-medium transition-all duration-200 ${
                    activeTab === "profile"
                      ? "bg-blue-100 dark:bg-blue-900/50 text-blue-900 dark:text-blue-100"
                      : "text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-900"
                  }`}
                >
                  <User className="h-5 w-5" />
                  <span>Profil</span>
                </button>
                <button
                  onClick={() => setActiveTab("password")}
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left font-medium transition-all duration-200 ${
                    activeTab === "password"
                      ? "bg-blue-100 dark:bg-blue-900/50 text-blue-900 dark:text-blue-100"
                      : "text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-900"
                  }`}
                >
                  <Lock className="h-5 w-5" />
                  <span>Passwort</span>
                </button>
              </div>
            </div>
          </div>

          {/* Content */}
          <div className="md:col-span-3">
            {activeTab === "profile" && user && <ProfileSettings user={user} />}
            {activeTab === "password" && <PasswordSettings />}
          </div>
        </div>
      </div>
    </div>
  );
}

