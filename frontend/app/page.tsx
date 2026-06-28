"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { Building2, Megaphone, Globe2, ArrowRight, MapPin, Users, TrendingUp } from "lucide-react";
import Link from "next/link";

export default function LandingPage() {
    const { user, isLoading } = useUserContext();
    const router = useRouter();

    useEffect(() => {
        if (!isLoading && user) {
            router.push("/trips");
        }
    }, [user, isLoading, router]);

    if (isLoading) return null;

    return (
        <div className="min-h-screen bg-white dark:bg-zinc-950">
            {/* Hero */}
            <div className="relative overflow-hidden">
                <div className="mx-auto max-w-7xl px-6 py-24 sm:py-32 lg:px-8">
                    <div className="text-center">
                        <div className="flex justify-center mb-6">
                            <span className="text-6xl">🌍</span>
                        </div>
                        <h1 className="text-5xl sm:text-7xl font-bold tracking-tight text-zinc-900 dark:text-white mb-6">
                            Trip Manager
                        </h1>
                        <p className="text-xl text-zinc-500 dark:text-zinc-400 max-w-2xl mx-auto mb-12">
                            Die Plattform für Reisende, Reisebüros und Werbepartner.
                            Plane, verwalte und entdecke Reisen – alles an einem Ort.
                        </p>
                        <div className="flex flex-col sm:flex-row gap-4 justify-center">
                            <Link
                                href="/auth"
                                className="px-8 py-4 text-base font-semibold bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white rounded-xl transition-colors flex items-center justify-center gap-2"
                            >
                                Jetzt starten
                                <ArrowRight className="h-5 w-5" />
                            </Link>
                            <Link
                                href="/search"
                                className="px-8 py-4 text-base font-semibold border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-900 rounded-xl transition-colors"
                            >
                                Reisen entdecken
                            </Link>
                        </div>
                    </div>
                </div>
            </div>

            {/* Zielgruppen */}
            <div className="mx-auto max-w-7xl px-6 py-24 lg:px-8">
                <h2 className="text-3xl font-bold text-center text-zinc-900 dark:text-white mb-4">
                    Für wen ist Trip Manager?
                </h2>
                <p className="text-center text-zinc-500 dark:text-zinc-400 mb-16">
                    Wähle deinen Bereich und leg direkt los.
                </p>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    {/* Reisende */}
                    <div className="bg-zinc-50 dark:bg-zinc-900 rounded-2xl p-8 border border-zinc-200 dark:border-zinc-800 flex flex-col">
                        <div className="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex items-center justify-center mb-6">
                            <Globe2 className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                        </div>
                        <h3 className="text-xl font-bold text-zinc-900 dark:text-white mb-3">Reisende</h3>
                        <p className="text-zinc-500 dark:text-zinc-400 text-sm flex-1 mb-6">
                            Plane deine nächste Reise, verwalte Unterkünfte und Transporte,
                            teile Erlebnisse mit anderen und entdecke inspirierende Trips.
                        </p>
                        <div className="space-y-2 mb-8">
                            {["Reisepläne erstellen", "Unterkünfte & Transport", "Feed & Inspiration", "Newsletter"].map((f) => (
                                <div key={f} className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400">
                                    <span className="text-green-500">✓</span>
                                    {f}
                                </div>
                            ))}
                        </div>
                        <Link
                            href="/auth"
                            className="w-full py-3 text-center text-sm font-semibold bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white rounded-xl transition-colors"
                        >
                            Kostenlos registrieren
                        </Link>
                    </div>

                    {/* Reisebüros */}
                    <div className="bg-zinc-50 dark:bg-zinc-900 rounded-2xl p-8 border border-zinc-200 dark:border-zinc-800 flex flex-col relative">
                        <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                            <span className="bg-[var(--brand-primary)] text-white text-xs font-semibold px-3 py-1 rounded-full">
                                Beliebt
                            </span>
                        </div>
                        <div className="w-12 h-12 bg-purple-100 dark:bg-purple-900/30 rounded-xl flex items-center justify-center mb-6">
                            <Building2 className="h-6 w-6 text-purple-600 dark:text-purple-400" />
                        </div>
                        <h3 className="text-xl font-bold text-zinc-900 dark:text-white mb-3">Reisebüros</h3>
                        <p className="text-zinc-500 dark:text-zinc-400 text-sm flex-1 mb-6">
                            Verwalte dein Team, erstelle eigenes Branding und behalte
                            den Überblick über alle Reisepläne deiner Mitarbeiter.
                        </p>
                        <div className="space-y-2 mb-8">
                            {["Team-Verwaltung", "Eigenes Branding", "Mitarbeiter einladen", "Nutzungsübersicht"].map((f) => (
                                <div key={f} className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400">
                                    <span className="text-green-500">✓</span>
                                    {f}
                                </div>
                            ))}
                        </div>
                        <Link
                            href="/business"
                            className="w-full py-3 text-center text-sm font-semibold bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white rounded-xl transition-colors"
                        >
                            Reisebüro registrieren
                        </Link>
                    </div>

                    {/* Advertiser */}
                    <div className="bg-zinc-50 dark:bg-zinc-900 rounded-2xl p-8 border border-zinc-200 dark:border-zinc-800 flex flex-col">
                        <div className="w-12 h-12 bg-amber-100 dark:bg-amber-900/30 rounded-xl flex items-center justify-center mb-6">
                            <Megaphone className="h-6 w-6 text-amber-600 dark:text-amber-400" />
                        </div>
                        <h3 className="text-xl font-bold text-zinc-900 dark:text-white mb-3">Werbepartner</h3>
                        <p className="text-zinc-500 dark:text-zinc-400 text-sm flex-1 mb-6">
                            Erhalte wöchentliche Travel Insights über beliebte Reiseziele,
                            Trends und Engagement-Daten deiner zugewiesenen Reisebüros.
                        </p>
                        <div className="space-y-2 mb-8">
                            {["Weekly Travel Insights", "Top-Destinationen", "Engagement-Trends", "Saisonale Muster"].map((f) => (
                                <div key={f} className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400">
                                    <span className="text-green-500">✓</span>
                                    {f}
                                </div>
                            ))}
                        </div>
                        <Link
                            href="/auth?type=advertiser"
                            className="w-full py-3 text-center text-sm font-semibold border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl transition-colors"
                        >
                            Als Werbepartner anmelden
                        </Link>
                    </div>
                </div>
            </div>

            {/* Stats */}
            <div className="bg-zinc-50 dark:bg-zinc-900 border-y border-zinc-200 dark:border-zinc-800">
                <div className="mx-auto max-w-7xl px-6 py-16 lg:px-8">
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-8 text-center">
                        <div>
                            <div className="flex justify-center mb-2">
                                <MapPin className="h-6 w-6 text-[var(--brand-primary)]" />
                            </div>
                            <p className="text-4xl font-bold text-zinc-900 dark:text-white mb-1">200+</p>
                            <p className="text-sm text-zinc-500">Reisewarnungen überwacht</p>
                        </div>
                        <div>
                            <div className="flex justify-center mb-2">
                                <Users className="h-6 w-6 text-[var(--brand-primary)]" />
                            </div>
                            <p className="text-4xl font-bold text-zinc-900 dark:text-white mb-1">Multi-Tenant</p>
                            <p className="text-sm text-zinc-500">Isolierte Reisebüro-Umgebungen</p>
                        </div>
                        <div>
                            <div className="flex justify-center mb-2">
                                <TrendingUp className="h-6 w-6 text-[var(--brand-primary)]" />
                            </div>
                            <p className="text-4xl font-bold text-zinc-900 dark:text-white mb-1">Weekly</p>
                            <p className="text-sm text-zinc-500">Travel Insights für Partner</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Footer */}
            <div className="mx-auto max-w-7xl px-6 py-12 lg:px-8 text-center">
                <p className="text-sm text-zinc-400">
                    © 2026 Trip Manager · <a href="mailto:admin@neatnode.xyz" className="hover:text-zinc-600">Kontakt</a>
                </p>
            </div>
        </div>
    );
}