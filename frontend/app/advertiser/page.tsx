"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";
import { Megaphone, MapPin, TrendingUp, Heart, MessageCircle, Calendar } from "lucide-react";
import { getAdvertiserMe } from "@/lib/api/advertiser";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function getAuthHeaders(): Promise<Record<string, string>> {
    const { firebaseAuth } = await import("@/lib/api/firebase");
    const user = firebaseAuth.currentUser;
    if (!user) return { "Content-Type": "application/json" };
    const token = await user.getIdToken();
    return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

type TopDestination = {
    country: string;
    countryCode: string;
    tripCount: number;
    avgLikes: number;
};

type EngagementStats = {
    totalLikes: number;
    totalComments: number;
    avgLikesPerTrip: number;
};

type SeasonalPattern = {
    peakMonth: string;
    avgPlanningLeadDays: number;
};

type TenantInsights = {
    tenantId: string;
    topDestinations: TopDestination[];
    engagement: EngagementStats;
    seasonalPattern: SeasonalPattern;
    generatedAt: string;
};

export default function AdvertiserDashboard() {
    const router = useRouter();
    const { user, isLoading } = useUserContext();
    const [insights, setInsights] = useState<TenantInsights[]>([]);
    const [advertiser, setAdvertiser] = useState<{ id: string; name: string; tenants: string[] } | null>(null);
    const [loading, setLoading] = useState(true);
    const [contactModal, setContactModal] = useState<{ tenantId: string; tenantName: string } | null>(null);
    const [contactMessage, setContactMessage] = useState("");
    const [sending, setSending] = useState(false);
    const [sentSuccess, setSentSuccess] = useState(false);

    const handleContact = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!contactModal || !advertiser) return;
        setSending(true);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(
                `${API_URL}/api/users/advertisers/${advertiser.id}/contact/${contactModal.tenantId}`,
                {
                    method: "POST",
                    headers,
                    body: JSON.stringify({ message: contactMessage }),
                }
            );
            if (res.ok) {
                setSentSuccess(true);
                setTimeout(() => {
                    setContactModal(null);
                    setSentSuccess(false);
                    setContactMessage("");
                }, 2000);
            }
        } finally {
            setSending(false);
        }
    };

    useEffect(() => {
        if (!isLoading && !user) {
            router.push("/auth?redirect=/advertiser");
        }
    }, [user, isLoading, router]);

    useEffect(() => {
        if (!user) return;
        getAuthHeaders().then(async (headers) => {
            const adv = await getAdvertiserMe();
            if (!adv) {
                router.push("/");
                return;
            }
            setAdvertiser(adv);

            const res = await fetch(`${API_URL}/api/newsletter/insights`, { headers });
            if (res.ok) {
                const data = await res.json();
                setInsights(data);
            }
            setLoading(false);
        });
    }, [user, router]);

    if (isLoading || loading) return <LoadingSpinner />;
    if (!user) return null;

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
                <div className="mb-8 flex items-center gap-3">
                    <Megaphone className="h-8 w-8 text-amber-500" />
                    <div>
                        <h1 className="text-3xl font-bold text-zinc-900 dark:text-white">
                            Travel Insights
                        </h1>
                        <p className="text-zinc-500 mt-1">
                            Wöchentliche Reisedaten für {advertiser?.name}
                        </p>
                    </div>
                </div>

                {insights.length === 0 ? (
                    <div className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-12 text-center">
                        <Megaphone className="h-12 w-12 text-zinc-300 mx-auto mb-4" />
                        <p className="text-zinc-500 font-medium mb-2">Noch keine Insights verfügbar</p>
                        <p className="text-zinc-400 text-sm">
                            Insights werden wöchentlich generiert. Bitte komm später wieder.
                        </p>
                    </div>
                ) : insights.map((insight, i) => (
                    <div key={i} className="mb-8 space-y-6">
                        <div className="flex items-center justify-between">
                            <h2 className="text-lg font-semibold text-zinc-900 dark:text-white">
                                Reisebüro: <span className="font-mono text-sm text-zinc-500">{insight.tenantId}</span>
                            </h2>
                            <div className="flex items-center gap-3">
                                <span className="text-xs text-zinc-400">
                                    Generiert: {new Date(insight.generatedAt).toLocaleDateString("de-DE")}
                                </span>
                                <button
                                    onClick={() => setContactModal({ tenantId: insight.tenantId, tenantName: insight.tenantId })}
                                    className="flex items-center gap-1.5 px-3 py-1.5 text-sm border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors"
                                >
                                    <Megaphone className="h-3.5 w-3.5" />
                                    Kontakt aufnehmen
                                </button>
                            </div>
                        </div>

                        {/* Engagement Stats */}
                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                <div className="flex items-center gap-2 mb-2">
                                    <Heart className="h-4 w-4 text-red-500" />
                                    <p className="text-xs text-zinc-500 uppercase tracking-wide">Likes gesamt</p>
                                </div>
                                <p className="text-3xl font-bold text-zinc-900 dark:text-white">
                                    {insight.engagement?.totalLikes ?? 0}
                                </p>
                            </div>
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                <div className="flex items-center gap-2 mb-2">
                                    <MessageCircle className="h-4 w-4 text-blue-500" />
                                    <p className="text-xs text-zinc-500 uppercase tracking-wide">Kommentare</p>
                                </div>
                                <p className="text-3xl font-bold text-zinc-900 dark:text-white">
                                    {insight.engagement?.totalComments ?? 0}
                                </p>
                            </div>
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                <div className="flex items-center gap-2 mb-2">
                                    <Calendar className="h-4 w-4 text-green-500" />
                                    <p className="text-xs text-zinc-500 uppercase tracking-wide">Ø Vorlaufzeit</p>
                                </div>
                                <p className="text-3xl font-bold text-zinc-900 dark:text-white">
                                    {insight.seasonalPattern?.avgPlanningLeadDays ?? 0}
                                    <span className="text-sm font-normal text-zinc-500 ml-1">Tage</span>
                                </p>
                            </div>
                        </div>

                        {/* Top Destinationen */}
                        <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-hidden">
                            <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800 flex items-center gap-2">
                                <MapPin className="h-4 w-4 text-zinc-500" />
                                <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">
                                    Top Destinationen
                                </h3>
                            </div>
                            <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                                {(insight.topDestinations ?? []).length === 0 ? (
                                    <p className="px-6 py-6 text-sm text-zinc-500 text-center">Noch keine Daten</p>
                                ) : insight.topDestinations.map((dest, j) => (
                                    <div key={j} className="flex items-center justify-between px-6 py-4">
                                        <div className="flex items-center gap-3">
                                            <span className="text-2xl">
                                                {dest.countryCode
                                                    ? String.fromCodePoint(...dest.countryCode.toUpperCase().split('').map(c => 0x1F1E6 + c.charCodeAt(0) - 65))
                                                    : '🌍'}
                                            </span>
                                            <div>
                                                <p className="text-sm font-medium text-zinc-900 dark:text-white">{dest.country}</p>
                                                <p className="text-xs text-zinc-500">Ø {dest.avgLikes.toFixed(1)} Likes/Trip</p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-sm font-bold text-zinc-900 dark:text-white">{dest.tripCount}</p>
                                            <p className="text-xs text-zinc-500">Trips</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Saisonales Muster */}
                        {insight.seasonalPattern?.peakMonth && (
                            <div className="bg-amber-50 dark:bg-amber-950/30 rounded-xl border border-amber-200 dark:border-amber-900 p-5 flex items-center gap-4">
                                <TrendingUp className="h-8 w-8 text-amber-500 shrink-0" />
                                <div>
                                    <p className="text-sm font-semibold text-zinc-900 dark:text-white mb-1">
                                        Peak-Monat: {insight.seasonalPattern.peakMonth}
                                    </p>
                                    <p className="text-xs text-zinc-500">
                                        Reisen werden durchschnittlich {insight.seasonalPattern.avgPlanningLeadDays} Tage im Voraus geplant.
                                    </p>
                                </div>
                            </div>
                        )}
                    </div>
                ))}
            </div>

            {/* Kontakt Modal */}
            {contactModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-white dark:bg-zinc-900 rounded-xl p-6 w-full max-w-md shadow-xl">
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-bold text-zinc-900 dark:text-white">Kontakt aufnehmen</h3>
                            <button onClick={() => setContactModal(null)} className="text-zinc-400 hover:text-zinc-600">✕</button>
                        </div>
                        {sentSuccess ? (
                            <div className="text-center py-6">
                                <p className="text-green-600 font-medium">✓ Nachricht wurde gesendet!</p>
                            </div>
                        ) : (
                            <form onSubmit={handleContact} className="space-y-4">
                                <p className="text-sm text-zinc-500">
                                    Sende eine Kontaktanfrage an den Inhaber von <strong>{contactModal.tenantName}</strong>.
                                </p>
                                <textarea
                                    value={contactMessage}
                                    onChange={(e) => setContactMessage(e.target.value)}
                                    placeholder="Deine Nachricht (optional)..."
                                    rows={4}
                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)] resize-none"
                                />
                                <div className="flex gap-3">
                                    <button
                                        type="submit"
                                        disabled={sending}
                                        className="flex-1 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                    >
                                        {sending ? "Wird gesendet..." : "Kontaktanfrage senden"}
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setContactModal(null)}
                                        className="px-4 py-2 text-sm border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800"
                                    >
                                        Abbrechen
                                    </button>
                                </div>
                            </form>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}