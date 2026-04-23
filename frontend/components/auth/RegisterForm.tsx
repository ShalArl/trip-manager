"use client";

import React, { useState } from "react";
import { CreateUserRequest } from "@/types/user";
import { validatePassword, getPasswordStrengthLabel, getPasswordStrengthColor, getPasswordStrengthBarColor } from "@/lib/validators/passwordValidator";
import { Eye, EyeOff, Check, X } from "lucide-react";


type Props = {
  onRegisterAction: (createUserRequest: CreateUserRequest) => void;
  onSwitchToLoginAction: () => void;
};

export default function RegisterForm({ onRegisterAction, onSwitchToLoginAction }: Props) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const passwordValidation = validatePassword(password);

  async function handleSubmit(e: React.SubmitEvent) {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    if (!name.trim() || !email.trim()) {
      setError("Bitte alle Felder ausfüllen.");
      setIsLoading(false);
      return;
    }

    if (!passwordValidation.isValid) {
      setError(passwordValidation.errors[0] || "Passwort erfüllt nicht alle Anforderungen.");
      setIsLoading(false);
      return;
    }

    try {
      await onRegisterAction({ name, email, password });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registrierung fehlgeschlagen. Bitte versuche es erneut.");
      setIsLoading(false);
    } finally {
      setIsLoading(false);
    }
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
          <div className="relative">
            <input
              type={showPassword ? "text" : "password"}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              className="w-full h-12 px-4 pr-10 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition text-sm"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            >
              {showPassword ? (
                <EyeOff className="h-5 w-5" />
              ) : (
                <Eye className="h-5 w-5" />
              )}
            </button>
          </div>

          {/* Password Strength Indicator */}
          {password && (
            <div className="mt-3 space-y-3">
              {/* Strength Bar */}
              <div className="space-y-1">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">
                    Passwort-Stärke
                  </span>
                  <span className={`text-xs font-semibold ${getPasswordStrengthColor(passwordValidation.score)}`}>
                    {getPasswordStrengthLabel(passwordValidation.score)}
                  </span>
                </div>
                <div className="h-2 bg-zinc-200 dark:bg-zinc-700 rounded-full overflow-hidden">
                  <div
                    className={`h-full transition-all duration-300 ${getPasswordStrengthBarColor(passwordValidation.score)}`}
                    style={{ width: `${(passwordValidation.score / 4) * 100}%` }}
                  />
                </div>
              </div>

              {/* Requirements */}
              <div className="space-y-2 p-3 rounded-lg bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800">
                <div className="space-y-1.5">
                  <div className="flex items-center gap-2">
                    {password.length >= 8 ? (
                      <Check className="h-4 w-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-red-600 dark:text-red-400 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${password.length >= 8 ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}>
                      Mindestens 8 Zeichen
                    </span>
                  </div>

                  <div className="flex items-center gap-2">
                    {/[A-Z]/.test(password) ? (
                      <Check className="h-4 w-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-red-600 dark:text-red-400 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${/[A-Z]/.test(password) ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}>
                      Mindestens ein Großbuchstabe (A-Z)
                    </span>
                  </div>

                  <div className="flex items-center gap-2">
                    {/\d/.test(password) ? (
                      <Check className="h-4 w-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-red-600 dark:text-red-400 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${/\d/.test(password) ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}>
                      Mindestens eine Zahl (0-9)
                    </span>
                  </div>

                  {/* Optional: Special characters bonus */}
                  <div className="flex items-center gap-2">
                    {/[!@#$%^&*()_+\-=\[\]{};:'",.<>?\/]/.test(password) ? (
                      <Check className="h-4 w-4 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-zinc-400 dark:text-zinc-600 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${/[!@#$%^&*()_+\-=\[\]{};:'",.<>?\/]/.test(password) ? "text-blue-700 dark:text-blue-300" : "text-zinc-500 dark:text-zinc-400"}`}>
                      Sonderzeichen (optional, aber empfohlen)
                    </span>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>

        {error && (
          <p className="text-sm text-red-500 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-800/50 rounded-xl px-4 py-3">
            {error}
          </p>
        )}

        <button
          type="submit"
          disabled={!name.trim() || !email.trim() || !passwordValidation.isValid || isLoading}
          className="w-full h-12 rounded-xl bg-sky-600 hover:bg-sky-700 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-sky-600 text-white font-semibold text-sm transition-all shadow-md shadow-sky-500/20 mt-2 flex items-center justify-center"
        >
          {isLoading ? (
            <div className="flex items-center gap-2">
              <div className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              Wird erstellt...
            </div>
          ) : (
            "Konto erstellen"
          )}
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