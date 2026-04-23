"use client";

import ProfileSettings from "@/components/settings/ProfileSettings";
import PasswordSettings from "@/components/settings/PasswordSettings";
import { useState } from "react";
import { useUserContext } from "@/lib/context/UserContext";
import { User, Lock } from "lucide-react";
import { useRouter } from "next/navigation";
import {LoadingSpinner} from "@/components/global/LoadingSpinner";

export default function SettingsPage() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("profile");
  const { user, isLoading } = useUserContext();

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (!user) {
    router.push("/auth");
    return null;
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
             {activeTab === "profile" && <ProfileSettings user={user} />}
             {activeTab === "password" && <PasswordSettings />}
           </div>
        </div>
      </div>
    </div>
  );
}

