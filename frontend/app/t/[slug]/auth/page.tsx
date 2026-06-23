"use client";
import { useParams, useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { getTenantBySlug, login, register } from "@/lib/api/auth";
import { useUserContext } from "@/lib/context/UserContext";
import type { LoginRequest, CreateUserRequest } from "@/types/user";

export default function TenantAuthPage() {
    const { slug } = useParams<{ slug: string }>();
    const router = useRouter();
    const { updateUser } = useUserContext();
    const [tenant, setTenant] = useState<{ tenantId: string; name: string; tier: string } | null>(null);
    const [mode, setMode] = useState<"login" | "register">("login");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        getTenantBySlug(slug).then((t) => {
            if (!t) router.push("/");
            else setTenant(t);
        });
    }, [slug, router]);

    const handleSubmit = async (e: React.SubmitEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);
        try {
            let user;
            if (mode === "login") {
                user = await login({ email, password } as LoginRequest);
            } else {
                user = await register({ email, password, name } as CreateUserRequest, tenant!.tenantId);
            }
            updateUser(user);
            router.push("/");
        } catch (err) {
            setError(err instanceof Error ? err.message : "Fehler beim Anmelden");
        } finally {
            setLoading(false);
        }
    };

    if (!tenant) return null;

    return (
        <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 flex items-center justify-center p-8">
            <div className="w-full max-w-md">
                <div className="text-center mb-8">
                    <span className="text-3xl">🏢</span>
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-white mt-2">{tenant.name}</h1>
                    <p className="text-sm text-zinc-500 mt-1">Melde dich an um fortzufahren</p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {mode === "register" && (
                        <div>
                            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">Name</label>
                            <input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                            />
                        </div>
                    )}
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">E-Mail</label>
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">Passwort</label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                        />
                    </div>

                    {error && <p className="text-sm text-red-500">{error}</p>}

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full py-2.5 rounded-lg bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white font-medium text-sm transition-colors disabled:opacity-50"
                    >
                        {loading ? "Wird verarbeitet..." : mode === "login" ? "Anmelden" : "Registrieren"}
                    </button>

                    <p className="text-center text-sm text-zinc-500">
                        {mode === "login" ? "Noch kein Konto? " : "Bereits ein Konto? "}
                        <button
                            type="button"
                            onClick={() => setMode(mode === "login" ? "register" : "login")}
                            className="text-[var(--brand-primary)] hover:underline"
                        >
                            {mode === "login" ? "Registrieren" : "Anmelden"}
                        </button>
                    </p>
                </form>
            </div>
        </div>
    );
}