"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { useTenantContext } from "@/lib/context/TenantContext";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";
import AuthPage from "@/components/auth/AuthPage";
import { register, login, getTenantBySlug } from "@/lib/api/auth";
import type { CreateUserRequest, LoginRequest } from "@/types/user";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function getAuthHeaders(): Promise<Record<string, string>> {
    const { firebaseAuth } = await import("@/lib/api/firebase");
    const user = firebaseAuth.currentUser;
    if (!user) return { "Content-Type": "application/json" };
    const token = await user.getIdToken();
    return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

export default function TenantJoinPage() {
    const router = useRouter();
    const { slug } = useParams<{ slug: string }>();
    const { user, isLoading, updateUser } = useUserContext();
    const { tenantId } = useTenantContext();
    const [tenant, setTenant] = useState<{ tenantId: string; name: string; tier: string } | null>(null);
    const [joining, setJoining] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (slug) {
            getTenantBySlug(slug).then(setTenant).catch(() => setError("Tenant nicht gefunden"));
        }
    }, [slug]);

    const joinTenant = async () => {
        if (!tenant) return;
        setJoining(true);
        try {
            const headers = await getAuthHeaders();
            const res = await fetch(`${API_URL}/api/users/tenants/join-by-slug`, {
                method: "POST",
                headers,
                body: JSON.stringify({ slug }),
            });
            if (res.ok) {
                // Token refreshen damit neue Claims geladen werden
                const { firebaseAuth } = await import("@/lib/api/firebase");
                await firebaseAuth.currentUser?.getIdToken(true);
                router.push("/trips");
            } else {
                setError("Beitritt fehlgeschlagen");
            }
        } finally {
            setJoining(false);
        }
    };

    const handleRegister = async (req: CreateUserRequest) => {
        const user = await register(req);
        updateUser(user);
        await joinTenant();
    };

    const handleLogin = async (req: LoginRequest) => {
        const user = await login(req);
        updateUser(user);
        await joinTenant();
    };

    if (isLoading) return <LoadingSpinner />;

    if (error) return (
        <div className="min-h-screen flex items-center justify-center">
            <p className="text-red-500">{error}</p>
        </div>
    );

    if (!tenant) return <LoadingSpinner />;

    // Bereits im richtigen Tenant
    if (user && tenantId === tenant.tenantId) {
        router.push("/trips");
        return null;
    }

    // Eingeloggt aber falscher Tenant
    if (user && tenantId !== tenant.tenantId) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="bg-white dark:bg-zinc-900 rounded-xl p-8 max-w-md w-full border border-zinc-200 dark:border-zinc-800 text-center">
                    <h2 className="text-xl font-bold mb-4">{tenant.name} beitreten</h2>
                    <p className="text-zinc-500 text-sm mb-6">
                        Möchtest du dem Reisebüro <strong>{tenant.name}</strong> beitreten?
                    </p>
                    {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
                    <button
                        onClick={joinTenant}
                        disabled={joining}
                        className="w-full py-3 bg-[var(--brand-primary)] text-white rounded-xl hover:bg-[var(--brand-primary-dark)] disabled:opacity-50"
                    >
                        {joining ? "Wird beigetreten..." : "Beitreten"}
                    </button>
                </div>
            </div>
        );
    }

    // Nicht eingeloggt – Auth anzeigen
    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <div className="text-center pt-10 pb-4">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-white">
                    {tenant.name} beitreten
                </h2>
                <p className="text-zinc-500 mt-1">
                    Registriere dich oder melde dich an um beizutreten
                </p>
            </div>
            <AuthPage
                onLoginAction={handleLogin}
                onRegisterAction={handleRegister}
            />
        </div>
    );
}