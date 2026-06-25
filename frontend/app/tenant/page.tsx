"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { useTenantContext } from "@/lib/context/TenantContext";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";
import { Building2, Users, Mail, Palette, BarChart3 } from "lucide-react";
import TenantSettings from "@/components/settings/TenantSettings";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function getAuthHeaders(): Promise<HeadersInit> {
    const { firebaseAuth } = await import("@/lib/api/firebase");
    const user = firebaseAuth.currentUser;
    if (!user) return { "Content-Type": "application/json" };
    const token = await user.getIdToken();
    return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

type Member = { id: string; email: string; name: string; role: string };
type Invitation = { id: string; email: string; role: string; inviteLink: string; expiresAt: string };

export default function TenantDashboardPage() {
    const router = useRouter();
    const { user, isLoading } = useUserContext();
    const { tenantId, isAdmin, isOwner } = useTenantContext();
    const [activeTab, setActiveTab] = useState("overview");
    const [members, setMembers] = useState<Member[]>([]);
    const [invitations, setInvitations] = useState<Invitation[]>([]);
    const [inviteEmail, setInviteEmail] = useState("");
    const [inviteRole, setInviteRole] = useState("tenant_member");
    const [inviting, setInviting] = useState(false);
    const [inviteLink, setInviteLink] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        if (!isLoading && (!user || (!isAdmin && !isOwner))) {
            router.push("/");
        }
    }, [user, isLoading, isAdmin, isOwner, router]);

    useEffect(() => {
        if (!tenantId || tenantId === "default") return;
        getAuthHeaders().then((headers) => {
            fetch(`${API_URL}/api/tenants/me/members`, { headers })
                .then((r) => r.json())
                .then(setMembers)
                .catch(() => {});
            fetch(`${API_URL}/api/tenants/me/invitations`, { headers })
                .then((r) => r.json())
                .then(setInvitations)
                .catch(() => {});
        });
    }, [tenantId]);

    const handleInvite = async (e: React.FormEvent) => {
        e.preventDefault();
        setInviting(true);
        setInviteLink(null);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/tenants/me/invitations`, {
                method: "POST",
                headers,
                body: JSON.stringify({ email: inviteEmail, role: inviteRole }),
            });
            if (res.ok) {
                const data = await res.json();
                setInviteLink(data.inviteLink);
                setInvitations((prev) => [data, ...prev]);
                setInviteEmail("");
            }
        } finally {
            setInviting(false);
        }
    };

    const handleRemoveMember = async (userId: string) => {
        const headers = await getAuthHeaders();
        const res = await fetch(`${API_URL}/api/tenants/me/members/${userId}`, {
            method: "DELETE",
            headers,
        });
        if (res.ok) {
            setMembers((prev) => prev.filter((m) => m.id !== userId));
        }
    };

    const handleDeleteInvitation = async (invId: string) => {
        const headers = await getAuthHeaders();
        const res = await fetch(`${API_URL}/api/tenants/me/invitations/${invId}`, {
            method: "DELETE",
            headers,
        });
        if (res.ok) {
            setInvitations((prev) => prev.filter((i) => i.id !== invId));
        }
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    if (isLoading) return <LoadingSpinner />;
    if (!user || tenantId === "default") return null;

    const tabs = [
        { id: "overview", label: "Übersicht", icon: Building2 },
        { id: "members", label: "Mitglieder", icon: Users },
        { id: "invitations", label: "Einladungen", icon: Mail },
        { id: "settings", label: "Branding & Einstellungen", icon: Palette },
        { id: "usage", label: "Nutzung", icon: BarChart3 },
    ];

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-zinc-900 dark:text-white flex items-center gap-3">
                        <Building2 className="h-8 w-8 text-[var(--brand-primary)]" />
                        Reisebüro-Dashboard
                    </h1>
                    <p className="text-zinc-500 dark:text-zinc-400 mt-1">
                        Verwalte dein Reisebüro, Mitglieder und Einstellungen
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-5 gap-6">
                    {/* Sidebar */}
                    <div className="md:col-span-1">
                        <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-2 sticky top-20">
                            {tabs.map(({ id, label, icon: Icon }) => (
                                <button
                                    key={id}
                                    onClick={() => setActiveTab(id)}
                                    className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left text-sm font-medium transition-all ${
                                        activeTab === id
                                            ? "bg-[var(--brand-primary)]/10 text-[var(--brand-primary)]"
                                            : "text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
                                    }`}
                                >
                                    <Icon className="h-4 w-4 shrink-0" />
                                    <span className="hidden md:block">{label}</span>
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Content */}
                    <div className="md:col-span-4 space-y-6">

                        {/* Übersicht */}
                        {activeTab === "overview" && (
                            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                    <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Mitglieder</p>
                                    <p className="text-3xl font-bold text-zinc-900 dark:text-white">{members.length}</p>
                                </div>
                                <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                    <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Admins</p>
                                    <p className="text-3xl font-bold text-zinc-900 dark:text-white">
                                        {members.filter((m) => m.role === "tenant_admin" || m.role === "tenant_owner").length}
                                    </p>
                                </div>
                                <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                    <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Offene Einladungen</p>
                                    <p className="text-3xl font-bold text-zinc-900 dark:text-white">{invitations.length}</p>
                                </div>
                            </div>
                        )}

                        {/* Mitglieder */}
                        {activeTab === "members" && (
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-hidden">
                                <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800">
                                    <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">Mitglieder ({members.length})</h2>
                                </div>
                                <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                                    {members.length === 0 ? (
                                        <p className="px-6 py-8 text-sm text-zinc-500 text-center">Keine Mitglieder gefunden</p>
                                    ) : members.map((m) => (
                                        <div key={m.id} className="flex items-center justify-between px-6 py-4">
                                            <div>
                                                <p className="text-sm font-medium text-zinc-900 dark:text-white">{m.name}</p>
                                                <p className="text-xs text-zinc-500">{m.email}</p>
                                            </div>
                                            <div className="flex items-center gap-3">
                                                <span className="text-xs px-2 py-1 rounded-full bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">
                                                    {m.role}
                                                </span>
                                                {isOwner && m.role !== "tenant_owner" && (
                                                    <button
                                                        onClick={() => handleRemoveMember(m.id)}
                                                        className="text-xs text-red-500 hover:underline"
                                                    >
                                                        Entfernen
                                                    </button>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Einladungen */}
                        {activeTab === "invitations" && (
                            <div className="space-y-4">
                                {/* Einladung erstellen */}
                                {(isOwner || isAdmin) && (
                                    <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-6">
                                        <h2 className="text-sm font-semibold text-zinc-900 dark:text-white mb-4">Mitarbeiter einladen</h2>
                                        <form onSubmit={handleInvite} className="space-y-3">
                                            <div className="flex gap-3">
                                                <input
                                                    type="email"
                                                    value={inviteEmail}
                                                    onChange={(e) => setInviteEmail(e.target.value)}
                                                    placeholder="email@beispiel.de"
                                                    required
                                                    className="flex-1 px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                                />
                                                <select
                                                    value={inviteRole}
                                                    onChange={(e) => setInviteRole(e.target.value)}
                                                    className="px-3 py-2 text-sm rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900"
                                                >
                                                    <option value="tenant_member">Mitarbeiter</option>
                                                    {isOwner && <option value="tenant_admin">Admin</option>}
                                                </select>
                                                <button
                                                    type="submit"
                                                    disabled={inviting}
                                                    className="px-4 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                                                >
                                                    {inviting ? "..." : "Einladen"}
                                                </button>
                                            </div>
                                        </form>

                                        {inviteLink && (
                                            <div className="mt-4 p-3 bg-green-50 dark:bg-green-950 rounded-lg border border-green-200 dark:border-green-800">
                                                <p className="text-xs text-green-700 dark:text-green-300 mb-2">Einladungslink erstellt:</p>
                                                <div className="flex items-center gap-2">
                                                    <code className="flex-1 text-xs text-zinc-600 dark:text-zinc-400 truncate">{inviteLink}</code>
                                                    <button
                                                        onClick={() => copyToClipboard(inviteLink)}
                                                        className="text-xs px-2 py-1 bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded"
                                                    >
                                                        {copied ? "✓" : "Kopieren"}
                                                    </button>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                )}

                                {/* Offene Einladungen */}
                                <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-hidden">
                                    <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800">
                                        <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">Offene Einladungen ({invitations.length})</h2>
                                    </div>
                                    <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                                        {invitations.length === 0 ? (
                                            <p className="px-6 py-8 text-sm text-zinc-500 text-center">Keine offenen Einladungen</p>
                                        ) : invitations.map((inv) => (
                                            <div key={inv.id} className="flex items-center justify-between px-6 py-4">
                                                <div>
                                                    <p className="text-sm text-zinc-900 dark:text-white">{inv.email}</p>
                                                    <p className="text-xs text-zinc-500">
                                                        {inv.role} · Läuft ab: {new Date(inv.expiresAt).toLocaleDateString("de-DE")}
                                                    </p>
                                                </div>
                                                <div className="flex items-center gap-2">
                                                    <button
                                                        onClick={() => copyToClipboard(inv.inviteLink)}
                                                        className="text-xs px-2 py-1 border border-zinc-200 dark:border-zinc-700 rounded hover:bg-zinc-50 dark:hover:bg-zinc-800"
                                                    >
                                                        Link kopieren
                                                    </button>
                                                    {isOwner && (
                                                        <button
                                                            onClick={() => handleDeleteInvitation(inv.id)}
                                                            className="text-xs text-red-500 hover:underline"
                                                        >
                                                            Löschen
                                                        </button>
                                                    )}
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Branding & Settings */}
                        {activeTab === "settings" && <TenantSettings />}

                        {/* Nutzung */}
                        {activeTab === "usage" && (
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-6">
                                <p className="text-sm text-zinc-500">Nutzungsdaten werden im Branding & Einstellungen Tab angezeigt.</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}