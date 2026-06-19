"use client";
import {useTenantContext} from "@/lib/context/TenantContext";
import {getTenant} from "@/lib/api/auth";
import {useEffect, useState} from "react";
import {X} from "lucide-react";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function getAuthHeaders(): Promise<HeadersInit> {
    const {firebaseAuth} = await import("@/lib/api/firebase");
    const user = firebaseAuth.currentUser;
    if (!user) return {"Content-Type": "application/json"};
    const token = await user.getIdToken();
    return {"Content-Type": "application/json", Authorization: `Bearer ${token}`};
}

export default function TenantSettings() {
    const {tenantId, role, isOwner, isAdmin} = useTenantContext();
    const [tenant, setTenant] = useState<{ tenantId: string; name: string; tier: string } | null>(null);
    const [showUpgradeModal, setShowUpgradeModal] = useState(false);
    const [upgrading, setUpgrading] = useState(false);
    const [branding, setBranding] = useState({logoUrl: "", primaryColor: "#0284c7", companyName: "", customDomain: ""});
    const [savingBranding, setSavingBranding] = useState(false);
    const [brandingMsg, setBrandingMsg] = useState<string | null>(null);
    const [settings, setSettings] = useState<{ maxActiveTrips: number }>({maxActiveTrips: 0});
    const [savingSettings, setSavingSettings] = useState(false);
    const [settingsMsg, setSettingsMsg] = useState<string | null>(null);
    const [deleting, setDeleting] = useState(false);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [usage, setUsage] = useState<{
        apiCalls: number;
        breakdown: { service: string; calls: number }[];
        pricing: { basePrice: number; apiCallCost: number; totalCost: number; currency: string };
    } | null>(null);

    useEffect(() => {
        getTenant().then((t) => {
            setTenant(t);
            if (t) {
                getAuthHeaders().then((headers) => {
                    fetch(`${API_URL}/api/tenants/me/branding`, {headers})
                        .then((r) => r.json())
                        .then((b) => setBranding((prev) => ({...prev, ...b})))
                        .catch(() => {});

                    fetch(`${API_URL}/api/tenants/me/settings`, {headers})
                        .then((r) => r.json())
                        .then((s) => setSettings((prev) => ({...prev, ...s})))
                        .catch(() => {});
                });
            }
        });
    }, []);

    useEffect(() => {
        if (!tenant || tenant.tier === "free") return;
        getAuthHeaders().then((headers) =>
            fetch(`${API_URL}/api/tenants/me/usage`, {headers})
                .then((r) => r.json())
                .then(setUsage)
                .catch(() => {})
        );
    }, [tenant]);

    const handleUpgrade = async (tier: string) => {
        setUpgrading(true);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/tenants/me/tier`, {
                method: "PUT",
                headers,
                body: JSON.stringify({tier}),
            });
            if (res.ok) {
                const updated = await res.json();
                setTenant((prev) => prev ? {...prev, tier: updated.tier} : prev);
                setShowUpgradeModal(false);
            }
        } finally {
            setUpgrading(false);
        }
    };

    const handleSaveBranding = async (e: React.FormEvent) => {
        e.preventDefault();
        setSavingBranding(true);
        setBrandingMsg(null);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/tenants/me/branding`, {
                method: "PUT",
                headers,
                body: JSON.stringify(branding),
            });
            if (res.ok) setBrandingMsg("Gespeichert ✓");
            else setBrandingMsg("Fehler beim Speichern");
        } finally {
            setSavingBranding(false);
        }
    };

    const handleSaveSettings = async (e: React.FormEvent) => {
        e.preventDefault();
        setSavingSettings(true);
        setSettingsMsg(null);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/tenants/me/settings`, {
                method: "PUT",
                headers,
                body: JSON.stringify({maxActiveTrips: settings.maxActiveTrips}),
            });
            if (res.ok) setSettingsMsg("Gespeichert ✓");
            else setSettingsMsg("Fehler beim Speichern");
        } finally {
            setSavingSettings(false);
        }
    };

    const handleDelete = async () => {
        setDeleting(true);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/tenants/me`, {
                method: "DELETE",
                headers,
            });
            if (res.ok) {
                const {firebaseAuth} = await import("@/lib/api/firebase");
                await firebaseAuth.currentUser?.getIdToken(true);
                window.location.href = "/";
            }
        } finally {
            setDeleting(false);
            setShowDeleteConfirm(false);
        }
    };

    return (
        <div className="bg-white dark:bg-zinc-950 rounded-xl border border-zinc-200 dark:border-zinc-800 shadow-lg p-6">
            <h2 className="text-xl font-bold text-zinc-900 dark:text-white mb-4">
                {tenant?.name ?? "Reisebüro"}
            </h2>

            {tenant ? (
                <div className="space-y-4">
                    <div className="flex justify-between py-3 border-b border-zinc-100 dark:border-zinc-800">
                        <span className="text-sm text-zinc-500">Name</span>
                        <span className="text-sm font-medium text-zinc-900 dark:text-white">{tenant.name}</span>
                    </div>

                    <div className="flex justify-between items-center py-3 border-b border-zinc-100 dark:border-zinc-800">
                        <span className="text-sm text-zinc-500">Plan</span>
                        <div className="flex items-center gap-2">
                            <span className="text-sm font-medium capitalize text-[var(--brand-primary)]">{tenant.tier}</span>
                            {isOwner && tenant.tier === "free" && (
                                <button
                                    onClick={() => setShowUpgradeModal(true)}
                                    className="text-xs px-2 py-1 bg-[var(--brand-primary)] text-white rounded-md hover:bg-[var(--brand-primary-dark)]"
                                >
                                    Upgrade
                                </button>
                            )}
                            {isOwner && tenant.tier !== "free" && (
                                <button
                                    onClick={() => handleUpgrade("free")}
                                    disabled={upgrading}
                                    className="text-xs px-2 py-1 border border-zinc-300 dark:border-zinc-600 text-zinc-500 rounded-md hover:bg-zinc-50 dark:hover:bg-zinc-800 disabled:opacity-50"
                                >
                                    {upgrading ? "..." : "Downgrade auf Free"}
                                </button>
                            )}
                        </div>
                    </div>

                    <div className="flex justify-between py-3 border-b border-zinc-100 dark:border-zinc-800">
                        <span className="text-sm text-zinc-500">Deine Rolle</span>
                        <span className="text-sm font-medium text-zinc-900 dark:text-white">{role}</span>
                    </div>
                    <div className="flex justify-between py-3 border-b border-zinc-100 dark:border-zinc-800">
                        <span className="text-sm text-zinc-500">Mandant-ID</span>
                        <span className="text-xs font-mono text-zinc-500">{tenantId}</span>
                    </div>

                    {/* Mitarbeiter */}
                    {isOwner && (
                        <div className="pt-4">
                            <p className="text-sm text-zinc-500 mb-3">Mitarbeiter einladen (coming soon)</p>
                            <button disabled
                                    className="px-4 py-2 text-sm bg-sky-100 text-[var(--brand-primary)] rounded-lg opacity-50 cursor-not-allowed">
                                Einladung senden
                            </button>
                        </div>
                    )}

                    {/* Settings – nur Standard+, nur Owner/Platform-Admin */}
                    {isOwner && tenant.tier !== "free" && (
                        <form onSubmit={handleSaveSettings}
                              className="pt-4 border-t border-zinc-100 dark:border-zinc-800 space-y-4">
                            <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">Einstellungen</h3>
                            <div>
                                <label className="block text-xs text-zinc-500 mb-1">
                                    Maximale gleichzeitig aktive Reisen pro Mitarbeiter
                                </label>
                                <input
                                    type="number"
                                    min={0}
                                    value={settings.maxActiveTrips}
                                    onChange={(e) =>
                                        setSettings({maxActiveTrips: Math.max(0, parseInt(e.target.value, 10) || 0)})
                                    }
                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                                <p className="text-xs text-zinc-400 mt-1">
                                    0 bedeutet unbegrenzt.
                                </p>
                            </div>
                            <div className="flex items-center gap-3">
                                <button
                                    type="submit"
                                    disabled={savingSettings}
                                    className="px-4 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                >
                                    {savingSettings ? "Speichern..." : "Speichern"}
                                </button>
                                {settingsMsg && <span className="text-sm text-green-600">{settingsMsg}</span>}
                            </div>
                        </form>
                    )}

                    {isOwner && (
                        <div className="pt-4 border-t border-zinc-100 dark:border-zinc-800">
                            {!showDeleteConfirm ? (
                                <button
                                    onClick={() => setShowDeleteConfirm(true)}
                                    className="text-sm text-red-500 hover:underline"
                                >
                                    Reisebüro löschen
                                </button>
                            ) : (
                                <div className="bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded-lg p-4 space-y-3">
                                    <p className="text-sm text-red-700 dark:text-red-300 font-medium">
                                        Bist du sicher? Diese Aktion kann nicht rückgängig gemacht werden.
                                    </p>
                                    <div className="flex gap-2">
                                        <button
                                            onClick={handleDelete}
                                            disabled={deleting}
                                            className="px-3 py-1.5 text-sm bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
                                        >
                                            {deleting ? "Wird gelöscht..." : "Ja, löschen"}
                                        </button>
                                        <button
                                            onClick={() => setShowDeleteConfirm(false)}
                                            className="px-3 py-1.5 text-sm border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800"
                                        >
                                            Abbrechen
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Branding – nur Standard+ */}
                    {(isOwner || isAdmin) && tenant.tier !== "free" && (
                        <form onSubmit={handleSaveBranding}
                              className="pt-4 border-t border-zinc-100 dark:border-zinc-800 space-y-4">
                            <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">Branding</h3>
                            <div>
                                <label className="block text-xs text-zinc-500 mb-1">Logo URL</label>
                                <input
                                    type="url"
                                    value={branding.logoUrl}
                                    onChange={(e) => setBranding((p) => ({...p, logoUrl: e.target.value}))}
                                    placeholder="https://..."
                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                            </div>
                            <div>
                                <label className="block text-xs text-zinc-500 mb-1">Primärfarbe</label>
                                <div className="flex items-center gap-2">
                                    <input
                                        type="color"
                                        value={branding.primaryColor}
                                        onChange={(e) => setBranding((p) => ({...p, primaryColor: e.target.value}))}
                                        className="h-9 w-16 rounded border border-zinc-200 dark:border-zinc-700 cursor-pointer"
                                    />
                                    <span className="text-sm font-mono text-zinc-500">{branding.primaryColor}</span>
                                </div>
                            </div>
                            <div>
                                <label className="block text-xs text-zinc-500 mb-1">Firmenname</label>
                                <input
                                    type="text"
                                    value={branding.companyName}
                                    onChange={(e) => setBranding((p) => ({...p, companyName: e.target.value}))}
                                    placeholder="Muster Reisen GmbH"
                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                            </div>
                            <div>
                                <label className="block text-xs text-zinc-500 mb-1">Custom Domain (coming soon)</label>
                                <input
                                    type="text"
                                    value={branding.customDomain}
                                    disabled
                                    placeholder="reisen.muster.de"
                                    className="w-full px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-900 opacity-50 cursor-not-allowed"
                                />
                            </div>
                            <div className="flex items-center gap-3">
                                <button
                                    type="submit"
                                    disabled={savingBranding}
                                    className="px-4 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                >
                                    {savingBranding ? "Speichern..." : "Speichern"}
                                </button>
                                <button
                                    type="button"
                                    onClick={() => {
                                        setBranding({
                                            logoUrl: "",
                                            primaryColor: "#0284c7",
                                            companyName: "",
                                            customDomain: ""
                                        });
                                        setBrandingMsg("Zurückgesetzt – bitte speichern um zu bestätigen");
                                    }}
                                    className="px-4 py-2 text-sm border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800"
                                >
                                    Zurücksetzen
                                </button>
                                {brandingMsg && <span className="text-sm text-green-600">{brandingMsg}</span>}
                            </div>
                        </form>
                    )}

                    {(isOwner || isAdmin) && tenant.tier !== "free" && usage && (
                        <div className="pt-4 border-t border-zinc-100 dark:border-zinc-800 space-y-4">
                            <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">Nutzung & Kosten</h3>

                            <div className="bg-zinc-50 dark:bg-zinc-900 rounded-lg p-4 space-y-3">
                                <div className="flex justify-between text-sm">
                                    <span className="text-zinc-500">API Calls gesamt</span>
                                    <span className="font-medium text-zinc-900 dark:text-white">
                                        {usage.apiCalls.toLocaleString("de-DE")}
                                    </span>
                                </div>

                                {usage.breakdown?.map((b) => (
                                    <div key={b.service} className="flex justify-between text-xs">
                                        <span className="text-zinc-400">{b.service}</span>
                                        <span className="text-zinc-500">{b.calls.toLocaleString("de-DE")} calls</span>
                                    </div>
                                ))}
                            </div>

                            <div className="bg-sky-50 dark:bg-sky-950 rounded-lg p-4 space-y-2">
                                <h4 className="text-xs font-medium text-zinc-500 uppercase tracking-wide">Aktuelle Abrechnung</h4>
                                <div className="flex justify-between text-sm">
                                    <span className="text-zinc-500">Basispreis</span>
                                    <span className="text-zinc-900 dark:text-white">
                                        {usage.pricing.basePrice.toFixed(2)} {usage.pricing.currency}
                                    </span>
                                </div>
                                <div className="flex justify-between text-sm">
                                    <span className="text-zinc-500">API Calls (über 10.000)</span>
                                    <span className="text-zinc-900 dark:text-white">
                                        {usage.pricing.apiCallCost.toFixed(2)} {usage.pricing.currency}
                                    </span>
                                </div>
                                <div className="flex justify-between text-sm font-semibold border-t border-sky-100 dark:border-sky-900 pt-2 mt-2">
                                    <span className="text-zinc-900 dark:text-white">Gesamt</span>
                                    <span className="text-sky-600">
                                        {usage.pricing.totalCost.toFixed(2)} {usage.pricing.currency}
                                    </span>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            ) : (
                <p className="text-sm text-zinc-500">Lade Tenant-Daten...</p>
            )}

            {/* Upgrade Modal */}
            {showUpgradeModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-white dark:bg-zinc-900 rounded-xl p-6 w-full max-w-md shadow-xl">
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-bold text-zinc-900 dark:text-white">Plan upgraden</h3>
                            <button onClick={() => setShowUpgradeModal(false)}>
                                <X className="h-5 w-5 text-zinc-500"/>
                            </button>
                        </div>
                        <p className="text-sm text-zinc-500 mb-6">
                            Wähle einen Plan. Das Upgrade wird in einer echten Umgebung über ein Zahlungssystem
                            abgewickelt.
                        </p>
                        <div className="space-y-3">
                            <div className="border-2 border-[var(--brand-primary)] rounded-lg p-4">
                                <div className="flex justify-between items-center mb-2">
                                    <span className="font-semibold text-zinc-900 dark:text-white">Standard</span>
                                    <span className="text-[var(--brand-primary)] font-medium">€29 / Monat</span>
                                </div>
                                <ul className="text-xs text-zinc-500 space-y-1 mb-3">
                                    <li>✓ Eigenes Branding</li>
                                    <li>✓ Unbegrenzte Reisepläne</li>
                                    <li>✓ SLA 99.5% Uptime</li>
                                    <li>✓ Priority Support</li>
                                </ul>
                                <button
                                    onClick={() => handleUpgrade("standard")}
                                    disabled={upgrading}
                                    className="w-full py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                >
                                    {upgrading ? "Wird verarbeitet..." : "Auf Standard upgraden (Demo)"}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}