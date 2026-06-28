"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { useTenantContext } from "@/lib/context/TenantContext";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";
import { Building2, Users, Mail, ArrowLeft } from "lucide-react";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function getAuthHeaders(): Promise<Record<string, string>> {
    const { firebaseAuth } = await import("@/lib/api/firebase");
    const user = firebaseAuth.currentUser;
    if (!user) return { "Content-Type": "application/json" };
    const token = await user.getIdToken();
    return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

type Member = { id: string; email: string; name: string; role: string };
type TenantDetail = { id: string; name: string; tier: string; status: string; slug: string };

export default function PlatformTenantDetailPage() {
    const router = useRouter();
    const { id } = useParams<{ id: string }>();
    const { user, isLoading } = useUserContext();
    const { isPlatformAdmin } = useTenantContext();
    const [tenant, setTenant] = useState<TenantDetail | null>(null);
    const [members, setMembers] = useState<Member[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!isLoading && (!user || !isPlatformAdmin)) {
            router.push("/");
        }
    }, [user, isLoading, isPlatformAdmin, router]);

    useEffect(() => {
        if (!id || !isPlatformAdmin) return;

        getAuthHeaders().then(async (headers) => {
            // Tenant-Details laden
            const tenantRes = await fetch(`${API_URL}/api/tenants/all`, { headers });
            if (tenantRes.ok) {
                const tenants = await tenantRes.json();
                const found = tenants.find((t: TenantDetail) => t.id === id);
                setTenant(found ?? null);
            }

            // Members laden – als platform_admin mit tenant context
            const membersRes = await fetch(
                `${API_URL}/api/tenants/me/members?tenantId=${id}`,
                { headers }
            );
            
            if (membersRes.ok) {
                setMembers(await membersRes.json());
            }

            setLoading(false);
        });
    }, [id, isPlatformAdmin]);

    if (isLoading || loading) return <LoadingSpinner />;
    if (!user || !isPlatformAdmin) return null;

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
                <button
                    onClick={() => router.push("/platform")}
                    className="flex items-center gap-2 text-sm text-zinc-500 hover:text-zinc-900 dark:hover:text-white mb-6"
                >
                    <ArrowLeft className="h-4 w-4" />
                    Zurück zur Übersicht
                </button>

                {tenant ? (
                    <>
                        <div className="mb-8 flex items-center gap-4">
                            <Building2 className="h-8 w-8 text-[var(--brand-primary)]" />
                            <div>
                                <h1 className="text-3xl font-bold text-zinc-900 dark:text-white">{tenant.name}</h1>
                                <p className="text-zinc-500 text-sm mt-1">
                                    {tenant.slug} ·
                                    <span className={`ml-1 font-medium ${
                                        tenant.tier === "enterprise" ? "text-purple-600" :
                                            tenant.tier === "standard" ? "text-blue-600" : "text-zinc-500"
                                    }`}>{tenant.tier}</span>
                                    {" · "}{tenant.status}
                                </p>
                            </div>
                        </div>

                        {/* Stats */}
                        <div className="grid grid-cols-2 gap-4 mb-8">
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Mitglieder</p>
                                <p className="text-3xl font-bold text-zinc-900 dark:text-white">{members.length}</p>
                            </div>
                            <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 p-5">
                                <p className="text-xs text-zinc-500 uppercase tracking-wide mb-1">Admins</p>
                                <p className="text-3xl font-bold text-zinc-900 dark:text-white">
                                    {members.filter(m => m.role === "tenant_admin" || m.role === "tenant_owner").length}
                                </p>
                            </div>
                        </div>

                        {/* Members */}
                        <div className="bg-white dark:bg-zinc-900 rounded-xl border border-zinc-200 dark:border-zinc-800 overflow-hidden">
                            <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800 flex items-center gap-2">
                                <Users className="h-4 w-4 text-zinc-500" />
                                <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">Mitglieder ({members.length})</h2>
                            </div>
                            <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                                {members.length === 0 ? (
                                    <p className="px-6 py-8 text-sm text-zinc-500 text-center">Keine Mitglieder</p>
                                ) : members.map((m) => (
                                    <div key={m.id} className="flex items-center justify-between px-6 py-4">
                                        <div>
                                            <p className="text-sm font-medium text-zinc-900 dark:text-white">{m.name}</p>
                                            <p className="text-xs text-zinc-500">{m.email}</p>
                                        </div>
                                        <span className="text-xs px-2 py-1 rounded-full bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">
                                            {m.role}
                                        </span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </>
                ) : (
                    <p className="text-zinc-500">Tenant nicht gefunden.</p>
                )}
            </div>
        </div>
    );
}