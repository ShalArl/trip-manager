"use client";

import {useEffect, useState} from "react";
import {useRouter} from "next/navigation";
import {useUserContext} from "@/lib/context/UserContext";
import {useTenantContext} from "@/lib/context/TenantContext";
import {LoadingSpinner} from "@/components/global/LoadingSpinner";
import {BarChart3, Building2, Euro, Link, Megaphone, Plus, Trash2, Users} from "lucide-react";
import {Advertiser, assignTenant, createAdvertiser, listAdvertisers, removeTenant} from "@/lib/api/advertiser";
import {getAuthHeaders} from "@/lib/api/auth";
import {BarChart as RechartsBarChart, Bar, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis} from "recharts";

const API_URL = process.env.NEXT_PUBLIC_API_URL;


type Tenant = { id: string; name: string; tier: string; status: string };

export default function PlatformDashboard() {
    const router = useRouter();
    const {user, isLoading} = useUserContext();
    const {isPlatformAdmin} = useTenantContext();
    const [activeTab, setActiveTab] = useState("tenants");
    const [tenants, setTenants] = useState<Tenant[]>([]);
    const [advertisers, setAdvertisers] = useState<Advertiser[]>([]);
    const [newAdv, setNewAdv] = useState({email: "", name: "", firebaseUid: ""});
    const [creating, setCreating] = useState(false);
    const [assigningTo, setAssigningTo] = useState<string | null>(null);
    const [assignTenantId, setAssignTenantId] = useState("");
    const [config, setConfig] = useState<{
        free: { basePrice: number; freeApiCalls: number; pricePerCall: number };
        standard: { basePrice: number; freeApiCalls: number; pricePerCall: number };
        enterprise: { basePrice: number; freeApiCalls: number; pricePerCall: number };
    } | null>(null);
    const [savingConfig, setSavingConfig] = useState(false);
    const [configMsg, setConfigMsg] = useState<string | null>(null);
    const [analyticstenantId, setAnalyticsTenantId] = useState("");
    const [analyticsData, setAnalyticsData] = useState<{ date: string; calls: number }[]>([]);
    const [loadingAnalytics, setLoadingAnalytics] = useState(false);


    useEffect(() => {
        if (!isLoading && (!user || !isPlatformAdmin)) {
            router.push("/");
        }
    }, [user, isLoading, isPlatformAdmin, router]);

    useEffect(() => {
        if (!isPlatformAdmin) return;

        getAuthHeaders().then((headers) => {
            fetch(`${API_URL}/api/users/tenants/all`, {headers})
                .then((r) => r.ok ? r.json() : [])
                .then(setTenants)
                .catch(() => {
                });
            fetch(`${API_URL}/api/users/platform/config`, {headers})
                .then((r) => r.ok ? r.json() : null)
                .then(setConfig)
                .catch(() => {
                });
        });

        listAdvertisers().then(setAdvertisers).catch(() => {
        });
    }, [isPlatformAdmin]);

    const handleCreateAdvertiser = async (e: React.FormEvent) => {
        e.preventDefault();
        setCreating(true);
        try {
            const adv = await createAdvertiser(newAdv);
            setAdvertisers((prev) => [adv, ...prev]);
            setNewAdv({email: "", name: "", firebaseUid: ""});
        } catch {
        } finally {
            setCreating(false);
        }
    };

    const handleAssignTenant = async (advertiserID: string) => {
        if (!assignTenantId) return;
        await assignTenant(advertiserID, assignTenantId);
        setAdvertisers((prev) => prev.map((a) =>
            a.id === advertiserID
                ? {...a, tenants: [...(a.tenants ?? []), assignTenantId]}
                : a
        ));
        setAssigningTo(null);
        setAssignTenantId("");
    };

    const handleRemoveTenant = async (advertiserID: string, tenantId: string) => {
        await removeTenant(advertiserID, tenantId);
        setAdvertisers((prev) => prev.map((a) =>
            a.id === advertiserID
                ? {...a, tenants: a.tenants.filter((t) => t !== tenantId)}
                : a
        ));
    };

    const handleSaveConfig = async (e: React.FormEvent) => {
        e.preventDefault();
        setSavingConfig(true);
        setConfigMsg(null);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/users/platform/config`, {
                method: "PUT",
                headers,
                body: JSON.stringify(config),
            });
            if (res.ok) setConfigMsg("Gespeichert ✓");
            else setConfigMsg("Fehler beim Speichern");
        } finally {
            setSavingConfig(false);
        }
    };

    const loadAnalytics = async (tid: string) => {
        if (!tid) return;
        setLoadingAnalytics(true);
        const headers = await getAuthHeaders();
        const res = await fetch(
            `${API_URL}/api/tenants/me/usage/timeseries?tenantId=${tid}`,
            {headers}
        );
        if (res.ok) setAnalyticsData(await res.json());
        setLoadingAnalytics(false);
    };

    if (isLoading) return <LoadingSpinner/>;
    if (!user || !isPlatformAdmin) return null;

    const tabs = [
        {id: "tenants", label: "Tenants", icon: Building2},
        {id: "advertisers", label: "Advertiser", icon: Megaphone},
        {id: "pricing", label: "Preise", icon: Euro},
        {id: "analytics", label: "Analytics", icon: BarChart3}
    ];

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-zinc-900 dark:text-white flex items-center gap-3">
                        <Users className="h-8 w-8 text-[var(--brand-primary)]"/>
                        Platform Dashboard
                    </h1>
                    <p className="text-zinc-500 mt-1">Verwalte alle Tenants und Advertiser der Plattform</p>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-2 gap-4 mb-8">
                    <div
                        className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                        <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Tenants</p>
                        <p className="text-3xl font-bold text-zinc-900 dark:text-white">{tenants.length}</p>
                    </div>
                    <div
                        className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                        <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Advertiser</p>
                        <p className="text-3xl font-bold text-zinc-900 dark:text-white">{advertisers.length}</p>
                    </div>
                </div>

                {/* Tabs */}
                <div className="flex gap-2 mb-6">
                    {tabs.map(({id, label, icon: Icon}) => (
                        <button
                            key={id}
                            onClick={() => setActiveTab(id)}
                            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                                activeTab === id
                                    ? "bg-[var(--brand-primary)] text-white"
                                    : "bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-50"
                            }`}
                        >
                            <Icon className="h-4 w-4"/>
                            {label}
                        </button>
                    ))}
                </div>

                {/* Tenants Tab */}
                {activeTab === "tenants" && (
                    <div
                        className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-hidden">
                        <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800">
                            <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">Alle Tenants
                                ({tenants.length})</h2>
                        </div>
                        <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                            {tenants.length === 0 ? (
                                <p className="px-6 py-8 text-sm text-zinc-500 text-center">Keine Tenants gefunden</p>
                            ) : tenants.map((t) => (
                                <div key={t.id} className="flex items-center justify-between px-6 py-4">
                                    <div>
                                        <p className="text-sm font-medium text-zinc-900 dark:text-white">{t.name}</p>
                                        <p className="text-xs text-zinc-500 font-mono">{t.id}</p>
                                    </div>
                                    <div className="flex items-center gap-3">
                                        <span className={`text-xs px-2 py-1 rounded-full ${
                                            t.tier === "enterprise"
                                                ? "bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300"
                                                : t.tier === "standard"
                                                    ? "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300"
                                                    : "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                                        }`}>
                                            {t.tier}
                                        </span>
                                        <button
                                            onClick={() => router.push(`/platform/tenant/${t.id}`)}
                                            className="text-xs text-[var(--brand-primary)] hover:underline"
                                        >
                                            Details →
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* Advertisers Tab */}
                {activeTab === "advertisers" && (
                    <div className="space-y-6">
                        {/* Neuen Advertiser anlegen */}
                        <div
                            className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-6">
                            <h2 className="text-sm font-semibold text-zinc-900 dark:text-white mb-4 flex items-center gap-2">
                                <Plus className="h-4 w-4"/>
                                Neuen Advertiser anlegen
                            </h2>
                            <form onSubmit={handleCreateAdvertiser} className="space-y-3">
                                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                                    <input
                                        type="text"
                                        placeholder="Name"
                                        value={newAdv.name}
                                        onChange={(e) => setNewAdv((p) => ({...p, name: e.target.value}))}
                                        required
                                        className="px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                    />
                                    <input
                                        type="email"
                                        placeholder="E-Mail"
                                        value={newAdv.email}
                                        onChange={(e) => setNewAdv((p) => ({...p, email: e.target.value}))}
                                        required
                                        className="px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                    />
                                    <input
                                        type="text"
                                        placeholder="Firebase UID (optional)"
                                        value={newAdv.firebaseUid}
                                        onChange={(e) => setNewAdv((p) => ({...p, firebaseUid: e.target.value}))}
                                        className="px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                    />
                                </div>
                                <button
                                    type="submit"
                                    disabled={creating}
                                    className="px-4 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                >
                                    {creating ? "Wird erstellt..." : "Advertiser anlegen"}
                                </button>
                            </form>
                        </div>

                        {/* Advertiser Liste */}
                        <div
                            className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-hidden">
                            <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800">
                                <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">Advertiser
                                    ({advertisers.length})</h2>
                            </div>
                            <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                                {advertisers.length === 0 ? (
                                    <p className="px-6 py-8 text-sm text-zinc-500 text-center">Keine Advertiser</p>
                                ) : advertisers.map((adv) => (
                                    <div key={adv.id} className="px-6 py-4 space-y-3">
                                        <div className="flex items-center justify-between">
                                            <div>
                                                <p className="text-sm font-medium text-zinc-900 dark:text-white">{adv.name}</p>
                                                <p className="text-xs text-zinc-500">{adv.email}</p>
                                            </div>
                                            <button
                                                onClick={() => setAssigningTo(assigningTo === adv.id ? null : adv.id)}
                                                className="flex items-center gap-1.5 text-xs px-3 py-1.5 border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800"
                                            >
                                                <Link className="h-3 w-3"/>
                                                Tenant zuweisen
                                            </button>
                                        </div>

                                        {/* Tenant zuweisen */}
                                        {assigningTo === adv.id && (
                                            <div className="flex gap-2">
                                                <select
                                                    value={assignTenantId}
                                                    onChange={(e) => setAssignTenantId(e.target.value)}
                                                    className="flex-1 px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                                >
                                                    <option value="">Tenant auswählen...</option>
                                                    {tenants
                                                        .filter((t) => !(adv.tenants ?? []).includes(t.id))
                                                        .map((t) => (
                                                            <option key={t.id} value={t.id}>{t.name}</option>
                                                        ))
                                                    }
                                                </select>
                                                <button
                                                    onClick={() => handleAssignTenant(adv.id)}
                                                    disabled={!assignTenantId}
                                                    className="px-3 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg disabled:opacity-50"
                                                >
                                                    Zuweisen
                                                </button>
                                            </div>
                                        )}

                                        {/* Zugewiesene Tenants */}
                                        {(adv.tenants ?? []).length > 0 && (
                                            <div className="flex flex-wrap gap-2">
                                                {adv.tenants.map((tid) => {
                                                    const tenant = tenants.find((t) => t.id === tid);
                                                    return (
                                                        <span key={tid}
                                                              className="flex items-center gap-1.5 text-xs px-2 py-1 bg-zinc-100 dark:bg-zinc-800 rounded-full">
                                                            {tenant?.name ?? tid}
                                                            <button onClick={() => handleRemoveTenant(adv.id, tid)}>
                                                                <Trash2
                                                                    className="h-3 w-3 text-red-400 hover:text-red-600"/>
                                                            </button>
                                                        </span>
                                                    );
                                                })}
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                )}

                {activeTab === "pricing" && (
                    <div
                        className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-6">
                        <h2 className="text-sm font-semibold text-zinc-900 dark:text-white mb-6">Preiskonfiguration</h2>
                        {!config ? (
                            <p className="text-sm text-zinc-500">Lade Konfiguration...</p>
                        ) : (
                            <form onSubmit={handleSaveConfig} className="space-y-8">
                                {(["free", "standard", "enterprise"] as const).map((tier) => (
                                    <div key={tier} className="space-y-4">
                                        <h3 className={`text-sm font-semibold capitalize ${
                                            tier === "enterprise" ? "text-purple-600" :
                                                tier === "standard" ? "text-[var(--brand-primary)]" : "text-zinc-500"
                                        }`}>
                                            {tier === "free" ? "Free" : tier === "standard" ? "Standard" : "Enterprise 🏆"}
                                        </h3>
                                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                            <div>
                                                <label className="block text-xs text-zinc-500 mb-1">Basispreis
                                                    (€/Monat)</label>
                                                <input
                                                    type="number"
                                                    min={0}
                                                    step={0.01}
                                                    value={config[tier].basePrice}
                                                    onChange={(e) => setConfig((prev) => prev ? {
                                                        ...prev,
                                                        [tier]: {
                                                            ...prev[tier],
                                                            basePrice: parseFloat(e.target.value) || 0
                                                        }
                                                    } : prev)}
                                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                                />
                                            </div>
                                            <div>
                                                <label className="block text-xs text-zinc-500 mb-1">Freie API
                                                    Calls</label>
                                                <input
                                                    type="number"
                                                    min={0}
                                                    value={config[tier].freeApiCalls}
                                                    onChange={(e) => setConfig((prev) => prev ? {
                                                        ...prev,
                                                        [tier]: {
                                                            ...prev[tier],
                                                            freeApiCalls: parseInt(e.target.value) || 0
                                                        }
                                                    } : prev)}
                                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                                />
                                            </div>
                                            <div>
                                                <label className="block text-xs text-zinc-500 mb-1">Preis pro Call
                                                    (€)</label>
                                                <input
                                                    type="number"
                                                    min={0}
                                                    step={0.0001}
                                                    value={config[tier].pricePerCall}
                                                    onChange={(e) => setConfig((prev) => prev ? {
                                                        ...prev,
                                                        [tier]: {
                                                            ...prev[tier],
                                                            pricePerCall: parseFloat(e.target.value) || 0
                                                        }
                                                    } : prev)}
                                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                                />
                                            </div>
                                        </div>
                                        {tier !== "enterprise" && (
                                            <div className="border-b border-zinc-100 dark:border-zinc-800"/>
                                        )}
                                    </div>
                                ))}
                                <div className="flex items-center gap-3 pt-4">
                                    <button
                                        type="submit"
                                        disabled={savingConfig}
                                        className="px-4 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                    >
                                        {savingConfig ? "Speichern..." : "Preise speichern"}
                                    </button>
                                    {configMsg && <span className="text-sm text-green-600">{configMsg}</span>}
                                </div>
                            </form>
                        )}
                    </div>
                )}
                {activeTab === "analytics" && (
                    <div className="space-y-6">
                        <div
                            className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-6">
                            <h2 className="text-sm font-semibold text-zinc-900 dark:text-white mb-4">Tenant
                                Analytics</h2>
                            <div className="flex gap-3 mb-6">
                                <select
                                    value={analyticstenantId}
                                    onChange={(e) => {
                                        setAnalyticsTenantId(e.target.value);
                                        loadAnalytics(e.target.value);
                                    }}
                                    className="flex-1 px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                >
                                    <option value="">Tenant auswählen...</option>
                                    {tenants.map((t) => (
                                        <option key={t.id} value={t.id}>{t.name}</option>
                                    ))}
                                </select>
                            </div>

                            {loadingAnalytics ? (
                                <p className="text-sm text-zinc-500 text-center py-8">Lade Daten...</p>
                            ) : analyticsData.length === 0 ? (
                                <p className="text-sm text-zinc-500 text-center py-8">
                                    {analyticstenantId ? "Keine Daten verfügbar" : "Bitte einen Tenant auswählen"}
                                </p>
                            ) : (
                                <ResponsiveContainer width="100%" height={300}>
                                    <RechartsBarChart data={analyticsData}>
                                        <CartesianGrid strokeDasharray="3 3" stroke="#e4e4e7"/>
                                        <XAxis
                                            dataKey="date"
                                            tick={{fontSize: 11, fill: "#71717a"}}
                                            tickFormatter={(v) => new Date(v).toLocaleDateString("de-DE", {
                                                day: "2-digit",
                                                month: "2-digit"
                                            })}
                                        />
                                        <YAxis tick={{fontSize: 11, fill: "#71717a"}}/>
                                        <Tooltip
                                            formatter={(value) => [Number(value).toLocaleString("de-DE"), "API Calls"]}
                                            labelFormatter={(label) => new Date(label).toLocaleDateString("de-DE")}
                                        />
                                        <Bar dataKey="calls" fill="#8b5cf6" radius={[4, 4, 0, 0]}/>
                                    </RechartsBarChart>
                                </ResponsiveContainer>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}