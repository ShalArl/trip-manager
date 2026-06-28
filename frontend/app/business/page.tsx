"use client";
import {useEffect, useState} from "react";
import { useRouter } from "next/navigation";
import { registerTenant } from "@/lib/api/auth";
import { register, login } from "@/lib/api/auth";
import { useUserContext } from "@/lib/context/UserContext";
import type { CreateUserRequest, LoginRequest } from "@/types/user";
import {firebaseAuth} from "@/lib/api/firebase";

type Step = "account" | "agency" | "success";
type Tier = "free" | "standard";

const TIER_OPTIONS: {
    id: Tier;
    name: string;
    price: string;
    features: string[];
}[] = [
    {
        id: "free",
        name: "Free",
        price: "Kostenlos",
        features: [
            "Bis zu 3 Reisepläne",
            "Geteilte Infrastruktur",
            "Community-Support",
            "Trip Manager Branding",
        ],
    },
    {
        id: "standard",
        name: "Standard",
        price: "€29 / Monat",
        features: [
            "Unbegrenzte Reisepläne",
            "Eigenes Logo & Farben",
            "Priority-Support",
        ],
    },
];

export default function BusinessPage() {
    const router = useRouter();
    const { updateUser } = useUserContext();

    const [step, setStep] = useState<Step>("account");
    const [tier, setTier] = useState<Tier>("free");
    const [tenantName, setTenantName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [createdTenant, setCreatedTenant] = useState<{
        tenantId: string;
        name: string;
        tier: string;
    } | null>(null);

    useEffect(() => {
        const unsubscribe = firebaseAuth.onAuthStateChanged((user) => {
            if (user) {
                setStep("agency");
            }
        });
        return () => unsubscribe();
    }, []);

    const handleAccountSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);
        try {
            let user;
            try {
                user = await register({ email, password, name } as CreateUserRequest);
            } catch {
                user = await login({ email, password } as LoginRequest);
            }
            updateUser(user);
            setStep("agency");
        } catch (err) {
            setError("Anmeldung fehlgeschlagen. Bitte prüfe deine Zugangsdaten.");
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleAgencySubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);
        try {
            const tenant = await registerTenant({ tenantName, tier });
            setCreatedTenant(tenant);
            setStep("success");
        } catch (err) {
            setError("Tenant-Registrierung fehlgeschlagen.");
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 flex">
            {/* Left panel */}
            <div className="hidden lg:flex lg:w-1/2 bg-sky-950 dark:bg-zinc-900 flex-col justify-between p-14 relative overflow-hidden">
                <div className="absolute inset-0 opacity-10">
                    <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
                        <defs>
                            <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse">
                                <path d="M 40 0 L 0 0 0 40" fill="none" stroke="white" strokeWidth="0.5" />
                            </pattern>
                        </defs>
                        <rect width="100%" height="100%" fill="url(#grid)" />
                    </svg>
                </div>

                <div className="relative z-10">
                    <div className="flex items-center gap-3 mb-16">
                        <span className="text-2xl">🌍</span>
                        <span className="text-xl font-bold tracking-tight text-white">Trip Manager</span>
                    </div>
                    <h2 className="text-4xl font-bold text-white leading-snug mb-6">
                        Dein Reisebüro.<br />Deine Marke.
                    </h2>
                    <p className="text-sky-200 text-lg leading-relaxed max-w-sm">
                        Biete deinen Kunden professionelle Reiseplanung unter deinem eigenen Namen an — powered by Trip Manager.
                    </p>
                </div>

                <div className="relative z-10 space-y-5">
                    {[
                        { emoji: "🏢", text: "Eigener Mandant mit isolierten Daten" },
                        { emoji: "🎨", text: "Individuelles Branding (Standard+)" },
                        { emoji: "👥", text: "Mitarbeiter einladen & verwalten" },
                    ].map((item) => (
                        <div key={item.text} className="flex items-center gap-4">
                            <div className="w-10 h-10 rounded-xl bg-white/10 flex items-center justify-center text-lg border border-white/20">
                                {item.emoji}
                            </div>
                            <span className="text-sky-100 font-medium">{item.text}</span>
                        </div>
                    ))}
                </div>
            </div>

            {/* Right panel */}
            <div className="w-full lg:w-1/2 flex items-center justify-center p-8">
                <div className="w-full max-w-md">
                    {/* Mobile logo */}
                    <div className="flex items-center gap-2 mb-10 lg:hidden">
                        <span className="text-2xl">🌍</span>
                        <span className="text-xl font-bold text-zinc-900 dark:text-white">Trip Manager</span>
                    </div>

                    {/* Step: Account */}
                    {step === "account" && (
                        <form onSubmit={handleAccountSubmit} className="space-y-5">
                            <div>
                                <h1 className="text-2xl font-bold text-zinc-900 dark:text-white mb-1">
                                    Konto erstellen oder anmelden
                                </h1>
                                <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                    Du benötigst ein Trip Manager Konto um ein Reisebüro zu registrieren.
                                </p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                                    Name
                                </label>
                                <input
                                    type="text"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    placeholder="Max Mustermann"
                                    className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                                    E-Mail
                                </label>
                                <input
                                    type="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    placeholder="max@reisebuero.de"
                                    required
                                    className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                                    Passwort
                                </label>
                                <input
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    placeholder="Mindestens 8 Zeichen"
                                    required
                                    minLength={8}
                                    className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                            </div>

                            {error && (
                                <p className="text-sm text-red-500">{error}</p>
                            )}

                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full py-2.5 rounded-lg bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white font-medium text-sm transition-colors disabled:opacity-50"
                            >
                                {loading ? "Wird verarbeitet..." : "Weiter"}
                            </button>

                            <p className="text-sm text-center text-zinc-500 dark:text-zinc-400">
                                Bereits ein Konto?{" "}
                                <button
                                    type="button"
                                    onClick={() => router.push("/auth")}
                                    className="text-[var(--brand-primary)] hover:underline"
                                >
                                    Anmelden
                                </button>
                            </p>
                        </form>
                    )}

                    {/* Step: Agency */}
                    {step === "agency" && (
                        <form onSubmit={handleAgencySubmit} className="space-y-6">
                            <div>
                                <h1 className="text-2xl font-bold text-zinc-900 dark:text-white mb-1">
                                    Reisebüro einrichten
                                </h1>
                                <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                    Wähle einen Namen und einen Plan für dein Reisebüro.
                                </p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                                    Name des Reisebüros
                                </label>
                                <input
                                    type="text"
                                    value={tenantName}
                                    onChange={(e) => setTenantName(e.target.value)}
                                    placeholder="Muster Reisen GmbH"
                                    required
                                    className="w-full px-4 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)]"
                                />
                            </div>

                            {/* Tier selection */}
                            <div className="space-y-3">
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                                    Plan wählen
                                </label>
                                {TIER_OPTIONS.map((option) => (
                                    <button
                                        key={option.id}
                                        type="button"
                                        onClick={() => setTier(option.id)}
                                        className={`w-full text-left p-4 rounded-lg border-2 transition-colors ${
                                            tier === option.id
                                                ? "border-[var(--brand-primary)] bg-sky-50 dark:bg-sky-950"
                                                : "border-zinc-200 dark:border-zinc-700 hover:border-zinc-300 dark:hover:border-zinc-600"
                                        }`}
                                    >
                                        <div className="flex justify-between items-center mb-2">
                      <span className="font-semibold text-zinc-900 dark:text-white">
                        {option.name}
                      </span>
                                            <span className="text-sm font-medium text-[var(--brand-primary)] dark:text-[var(--brand-primary-light)]">
                        {option.price}
                      </span>
                                        </div>
                                        <ul className="space-y-1">
                                            {option.features.map((f) => (
                                                <li key={f} className="text-xs text-zinc-500 dark:text-zinc-400 flex items-center gap-1.5">
                                                    <span className="text-sky-500">✓</span> {f}
                                                </li>
                                            ))}
                                        </ul>
                                    </button>
                                ))}
                            </div>

                            {error && (
                                <p className="text-sm text-red-500">{error}</p>
                            )}

                            <button
                                type="submit"
                                disabled={loading || !tenantName}
                                className="w-full py-2.5 rounded-lg bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white font-medium text-sm transition-colors disabled:opacity-50"
                            >
                                {loading ? "Wird eingerichtet..." : "Reisebüro registrieren"}
                            </button>
                        </form>
                    )}

                    {/* Step: Success */}
                    {step === "success" && createdTenant && (
                        <div className="space-y-6 text-center">
                            <div className="w-16 h-16 rounded-full bg-green-100 dark:bg-green-900 flex items-center justify-center mx-auto text-3xl">
                                ✅
                            </div>
                            <div>
                                <h1 className="text-2xl font-bold text-zinc-900 dark:text-white mb-2">
                                    Reisebüro eingerichtet
                                </h1>
                                <p className="text-zinc-500 dark:text-zinc-400 text-sm">
                                    <strong className="text-zinc-900 dark:text-white">{createdTenant.name}</strong> ist jetzt aktiv.
                                    Du bist der Inhaber dieses Mandanten.
                                </p>
                            </div>

                            <div className="bg-zinc-100 dark:bg-zinc-800 rounded-lg p-4 text-left space-y-2">
                                <div className="flex justify-between text-sm">
                                    <span className="text-zinc-500">Mandant-ID</span>
                                    <span className="font-mono text-zinc-900 dark:text-white text-xs">{createdTenant.tenantId}</span>
                                </div>
                                <div className="flex justify-between text-sm">
                                    <span className="text-zinc-500">Plan</span>
                                    <span className="capitalize text-zinc-900 dark:text-white">{createdTenant.tier}</span>
                                </div>
                            </div>

                            <button
                                onClick={() => router.push("/")}
                                className="w-full py-2.5 rounded-lg bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white font-medium text-sm transition-colors"
                            >
                                Zur App
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}