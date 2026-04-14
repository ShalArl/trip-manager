"use client";

import React, { useState } from "react";
import { components } from "@/generated/types";

type LoginRequest = components["schemas"]["LoginRequest"];

type Props = {
  onLoginAction: (loginRequest: LoginRequest) => void;
  onSwitchToRegisterAction: () => void;
};

export default function LoginForm({ onLoginAction, onSwitchToRegisterAction }: Props) {
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
      // Kein Passwort, keine Überprüfung — einfach einloggen
      onLoginAction({ email, password });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login fehlgeschlagen. Bitte versuche es erneut.");
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
            className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
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
            className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
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
          className="w-full h-12 rounded-xl bg-sky-600 hover:bg-sky-700 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-sky-600 text-white font-semibold text-sm transition-all shadow-md shadow-sky-500/20 mt-2 flex items-center justify-center"
        >
          {isLoading ? (
            <div className="flex items-center gap-2">
              <div className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
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
          className="text-sky-600 dark:text-sky-400 font-semibold hover:underline"
        >
          Registrieren
        </button>
      </p>
    </div>
  );
}