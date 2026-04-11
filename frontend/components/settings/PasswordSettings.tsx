"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { changePassword } from "@/lib/api/auth";
import { validatePassword, getPasswordStrengthLabel, getPasswordStrengthColor, getPasswordStrengthBarColor } from "@/lib/validators/passwordValidator";
import { AlertCircle, CheckCircle, Eye, EyeOff, Lock, Shield, Check, X } from "lucide-react";

export default function PasswordSettings() {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPasswords, setShowPasswords] = useState(false);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const newPasswordValidation = validatePassword(newPassword);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);

    // Validation
    if (newPassword !== confirmPassword) {
      setError("Neue Passwörter stimmen nicht überein");
      return;
    }

    if (!newPasswordValidation.isValid) {
      setError(newPasswordValidation.errors[0] || "Passwort erfüllt nicht alle Anforderungen");
      return;
    }

    setLoading(true);

    try {
      await changePassword({
        currentPassword,
        newPassword,
      });
      setSuccess(true);
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Fehler beim Ändern des Passworts"
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 shadow-lg p-8">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Success Message */}
        {success && (
          <div className="flex items-center gap-3 rounded-lg bg-green-50 dark:bg-green-950/50 border border-green-200 dark:border-green-900 p-4 animate-in fade-in duration-300">
            <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400 flex-shrink-0" />
            <p className="text-sm font-medium text-green-800 dark:text-green-200">
              Passwort erfolgreich geändert
            </p>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="flex items-center gap-3 rounded-lg bg-red-50 dark:bg-red-950/50 border border-red-200 dark:border-red-900 p-4 animate-in fade-in duration-300">
            <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0" />
            <p className="text-sm font-medium text-red-800 dark:text-red-200">
              {error}
            </p>
          </div>
        )}

        {/* Current Password */}
        <div className="space-y-3">
          <Label htmlFor="current-password" className="text-sm font-semibold text-zinc-900 dark:text-white flex items-center gap-2">
            <Shield className="h-4 w-4" />
            Aktuelles Passwort
          </Label>
          <div className="relative">
            <Input
              id="current-password"
              type={showPasswords ? "text" : "password"}
              value={currentPassword}
              onChange={(e) => setCurrentPassword(e.target.value)}
              required
              placeholder="Gib dein aktuelles Passwort ein"
              className="pr-10 bg-zinc-50 dark:bg-zinc-900 border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
            />
            <button
              type="button"
              onClick={() => setShowPasswords(!showPasswords)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            >
              {showPasswords ? (
                <EyeOff className="h-5 w-5" />
              ) : (
                <Eye className="h-5 w-5" />
              )}
            </button>
          </div>
        </div>

        {/* Divider */}
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-zinc-200 dark:border-zinc-800" />
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-white dark:bg-zinc-950 text-zinc-500 dark:text-zinc-400">
              Neues Passwort
            </span>
          </div>
        </div>

        {/* New Password */}
        <div className="space-y-3">
          <Label htmlFor="new-password" className="text-sm font-semibold text-zinc-900 dark:text-white flex items-center gap-2">
            <Lock className="h-4 w-4" />
            Neues Passwort
          </Label>
          <Input
            id="new-password"
            type={showPasswords ? "text" : "password"}
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            minLength={8}
            placeholder="Mindestens 8 Zeichen mit Großbuchstabe und Zahl"
            className="bg-zinc-50 dark:bg-zinc-900 border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
          />

          {/* Password Strength Indicator */}
          {newPassword && (
            <div className="space-y-3">
              {/* Strength Bar */}
              <div className="space-y-1">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">
                    Passwort-Stärke
                  </span>
                  <span className={`text-xs font-semibold ${getPasswordStrengthColor(newPasswordValidation.score)}`}>
                    {getPasswordStrengthLabel(newPasswordValidation.score)}
                  </span>
                </div>
                <div className="h-2 bg-zinc-200 dark:bg-zinc-700 rounded-full overflow-hidden">
                  <div
                    className={`h-full transition-all duration-300 ${getPasswordStrengthBarColor(newPasswordValidation.score)}`}
                    style={{ width: `${(newPasswordValidation.score / 4) * 100}%` }}
                  />
                </div>
              </div>

              {/* Requirements */}
              <div className="space-y-2 p-3 rounded-lg bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    {newPassword.length >= 8 ? (
                      <Check className="h-4 w-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-red-600 dark:text-red-400 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${newPassword.length >= 8 ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}>
                      Mindestens 8 Zeichen
                    </span>
                  </div>

                  <div className="flex items-center gap-2">
                    {/[A-Z]/.test(newPassword) ? (
                      <Check className="h-4 w-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-red-600 dark:text-red-400 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${/[A-Z]/.test(newPassword) ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}>
                      Mindestens ein Großbuchstabe (A-Z)
                    </span>
                  </div>

                  <div className="flex items-center gap-2">
                    {/\d/.test(newPassword) ? (
                      <Check className="h-4 w-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-red-600 dark:text-red-400 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${/\d/.test(newPassword) ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}>
                      Mindestens eine Zahl (0-9)
                    </span>
                  </div>

                  {/* Optional: Special characters bonus */}
                  <div className="flex items-center gap-2">
                    {/[!@#$%^&*()_+\-=\[\]{};:'",.<>?\/]/.test(newPassword) ? (
                      <Check className="h-4 w-4 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                    ) : (
                      <X className="h-4 w-4 text-zinc-400 dark:text-zinc-600 flex-shrink-0" />
                    )}
                    <span className={`text-xs ${/[!@#$%^&*()_+\-=\[\]{};:'",.<>?\/]/.test(newPassword) ? "text-blue-700 dark:text-blue-300" : "text-zinc-500 dark:text-zinc-400"}`}>
                      Sonderzeichen (optional, aber empfohlen)
                    </span>
                  </div>
                </div>
              </div>

              {/* Suggestions */}
              {newPasswordValidation.suggestions.length > 0 && (
                <div className="text-xs text-green-700 dark:text-green-300 space-y-1">
                  {newPasswordValidation.suggestions.map((suggestion, idx) => (
                    <p key={idx}>💡 {suggestion}</p>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Confirm Password */}
        <div className="space-y-3">
          <Label htmlFor="confirm-password" className="text-sm font-semibold text-zinc-900 dark:text-white">
            Passwort bestätigen
          </Label>
          <Input
            id="confirm-password"
            type={showPasswords ? "text" : "password"}
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
            minLength={8}
            placeholder="Wiederhole das neue Passwort"
            className="bg-zinc-50 dark:bg-zinc-900 border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
          />
          {confirmPassword && newPassword === confirmPassword && (
            <p className="text-xs font-semibold text-green-600 dark:text-green-400">✓ Passwörter stimmen überein</p>
          )}
          {confirmPassword && newPassword !== confirmPassword && (
            <p className="text-xs font-semibold text-red-600 dark:text-red-400">✗ Passwörter stimmen nicht überein</p>
          )}
        </div>

        {/* Security Info */}
        <div className="rounded-lg bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-900 p-4">
          <p className="text-sm text-blue-800 dark:text-blue-200">
            <span className="font-semibold">Sicherheitshinweis:</span> Verwende ein starkes, eindeutiges Passwort mit mindestens 8 Zeichen.
          </p>
        </div>

        {/* Submit Button */}
        <div className="pt-4 flex gap-3">
          <Button
            type="submit"
            disabled={loading || !currentPassword || !newPassword || !confirmPassword || newPassword !== confirmPassword || !newPasswordValidation.isValid}
            className="flex-1 sm:flex-none bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-all duration-200"
          >
            {loading ? (
              <div className="flex items-center gap-2">
                <div className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Wird geändert...
              </div>
            ) : (
              "Passwort ändern"
            )}
          </Button>
        </div>
      </form>
    </div>
  );
}

