"use client";

import React, {useState} from "react";
import {LoginRequest} from "@/types/user";


type Props = {
    onLoginAction: (loginRequest: LoginRequest) => void;
    onSwitchToRegisterAction: () => void;
};

export default function LoginForm({onLoginAction, onSwitchToRegisterAction}: Props) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    async function handleSubmit(e: React.SubmitEvent) {
        e.preventDefault();
        setError("");
        setIsLoading(true);

        if (!email.trim() || !password.trim()) {
            setError("Bitte alle Felder ausfüllen.");
            setIsLoading(false);
            return;
        }

        try {
            await onLoginAction({email, password});
        } catch (err) {
            setError(err instanceof Error ? err.message : "Login fehlgeschlagen. Bitte versuche es erneut.");
            setIsLoading(false);
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div>
            <h1 className="text-3xl font-bold tracking-tight text-zinc-900 dark:text-white mb-2">
                Willkommen zurück
            </h1>
            <p className="text-zinc-500 dark:text-zinc-400 mb-8 text-sm">
                Gib deine E-Mail-Adresse ein, um fortzufahren.
            </p>

            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        E-Mail-Adresse
                    </label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="du@beispiel.de"
                        className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)] focus:border-transparent transition text-sm"
                    />
                </div>
                <div>
                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
                        Passwort
                    </label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="••••••••"
                        className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-[var(--brand-primary)] focus:border-transparent transition text-sm"
                    />
                </div>

                {error && (
                    <p className="text-sm text-red-500 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-800/50 rounded-xl px-4 py-3">
                        {error}
                    </p>
                )}

                <button
                    type="submit"
                    disabled={!email.trim() || !password.trim() || isLoading}
                    className="w-full h-12 rounded-xl bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-[var(--brand-primary)] text-white font-semibold text-sm transition-all shadow-md shadow-sky-500/20 mt-2 flex items-center justify-center"
                >
                    {isLoading ? (
                        <div className="flex items-center gap-2">
                            <div
                                className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"/>
                            Wird angemeldet...
                        </div>
                    ) : (
                        "Anmelden"
                    )}
                </button>
            </form>

            <p className="text-center text-sm text-zinc-500 dark:text-zinc-400 mt-6">
                Noch kein Konto?{" "}
                <button
                    onClick={onSwitchToRegisterAction}
                    className="text-[var(--brand-primary)] dark:text-[var(--brand-primary-light)] font-semibold hover:underline"
                >
                    Registrieren
                </button>
            </p>

            <div className="mt-4 pt-4 border-t border-zinc-100 dark:border-zinc-800">
                <p className="text-center text-xs text-zinc-400 dark:text-zinc-500">
                    Reisebüro oder Unternehmen?{" "}
                    <a
                        href="/business"
                        className="text-[var(--brand-primary)] dark:text-[var(--brand-primary-light)] hover:underline font-medium"
                    >
                        Eigenen Mandanten einrichten
                    </a>
                </p>
            </div>
        </div>
    );
}