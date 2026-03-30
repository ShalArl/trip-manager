"use client";

import { useState } from "react";
import { User } from "@/types/user";
import { components } from "@/generated/types";

type CreateUserRequest = components["schemas"]["CreateUserRequest"];

type Props = {
  onRegisterAction: (createUserRequest: CreateUserRequest) => void;
  onSwitchToLoginAction: () => void;
};

export default function RegisterForm({ onRegisterAction, onSwitchToLoginAction }: Props) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("")
  const [error, setError] = useState("");

  function handleSubmit(e: React.SyntheticEvent) {
    e.preventDefault();
    setError("");

    if (!name.trim() || !email.trim()) {
      setError("Bitte alle Felder ausfüllen.");
      return;
    }

    if (password.length < 8) {
      setError("Passwort muss mindestens 8 Zeichen lang sein.");
      return;
    }
    if (!/[A-Z]/.test(password)) {
      setError("Passwort muss mindestens einen Großbuchstaben enthalten.");
      return;
    }
    if (!/[0-9]/.test(password)) {
      setError("Passwort muss mindestens eine Zahl enthalten.");
      return;
    }

    // Kein Passwort, keine Überprüfung — direkt registrieren
    onRegisterAction({ name, email, password });
  }

  return (
    <div>
      <h1 className="text-3xl font-bold tracking-tight text-zinc-900 dark:text-white mb-2">
        Konto erstellen
      </h1>
      <p className="text-zinc-500 dark:text-zinc-400 mb-8 text-sm">
        Starte kostenlos und plane deine erste Reise.
      </p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1.5">
            Dein Name
          </label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Max Mustermann"
            className="w-full h-12 px-4 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
          />
        </div>

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
          className="w-full h-12 rounded-xl bg-sky-600 hover:bg-sky-700 active:scale-[0.98] text-white font-semibold text-sm transition-all shadow-md shadow-sky-500/20 mt-2"
        >
          Konto erstellen
        </button>
      </form>

      <p className="text-center text-sm text-zinc-500 dark:text-zinc-400 mt-6">
        Bereits registriert?{" "}
        <button
          onClick={onSwitchToLoginAction}
          className="text-sky-600 dark:text-sky-400 font-semibold hover:underline"
        >
          Anmelden
        </button>
      </p>
    </div>
  );
}