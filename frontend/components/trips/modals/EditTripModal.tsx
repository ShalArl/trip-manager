"use client";

import { useState } from "react";
import { components } from "@/generated/types";

type TripResponse = components["schemas"]["TripResponse"];

type Props = {
  isOpen: boolean;
  trip: TripResponse;
  onCloseAction: () => void;
  onSaveAction: (trip: Partial<TripResponse>) => void;
};

export default function EditTripModal({
  isOpen,
  trip,
  onCloseAction,
  onSaveAction,
}: Props) {
  const [formData, setFormData] = useState({
    title: trip.title,
    shortDescription: trip.shortDescription || "",
    description: trip.description || "",
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    
    if (!formData.title.trim()) {
      newErrors.title = "Titel ist erforderlich";
    }
    if (!formData.shortDescription.trim()) {
      newErrors.shortDescription = "Kurzbeschreibung ist erforderlich";
    }
    if (formData.shortDescription.length > 80) {
      newErrors.shortDescription = "Max. 80 Zeichen";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    onSaveAction({
      title: formData.title,
      shortDescription: formData.shortDescription,
      description: formData.description || undefined,
    });

    onCloseAction();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="bg-white dark:bg-zinc-900 rounded-2xl shadow-xl max-w-md w-full border border-zinc-100 dark:border-zinc-800 overflow-hidden">
        {/* Header mit Gradient */}
        <div className="bg-gradient-to-r from-sky-500 to-sky-600 px-8 py-6">
          <h2 className="text-2xl font-bold text-white flex items-center gap-2">
            ✏️ Reise bearbeiten
          </h2>
          <p className="text-sky-100 text-sm mt-1">
            Aktualisiere die Details deiner Reise
          </p>
        </div>

        <form onSubmit={handleSubmit} className="p-8 space-y-5">
          {/* Title */}
          <div>
            <label className="block text-sm font-semibold text-zinc-900 dark:text-white mb-2">
              Reise-Titel
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => {
                setFormData({ ...formData, title: e.target.value });
                if (errors.title) setErrors({ ...errors, title: "" });
              }}
              placeholder="z.B. Frankreich 2026"
              className={`w-full px-4 py-3 rounded-xl border-2 transition-all bg-zinc-50 dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 ${
                errors.title
                  ? "border-red-300 dark:border-red-700"
                  : "border-zinc-200 dark:border-zinc-700 focus:border-sky-500"
              }`}
            />
            {errors.title && (
              <p className="text-red-600 dark:text-red-400 text-sm mt-1">
                {errors.title}
              </p>
            )}
          </div>

          {/* Short Description */}
          <div>
            <label className="block text-sm font-semibold text-zinc-900 dark:text-white mb-2">
              Kurzbeschreibung
            </label>
            <input
              type="text"
              value={formData.shortDescription}
              onChange={(e) => {
                setFormData({
                  ...formData,
                  shortDescription: e.target.value,
                });
                if (errors.shortDescription)
                  setErrors({ ...errors, shortDescription: "" });
              }}
              placeholder="z.B. Wochenendtrip nach Paris"
              maxLength={80}
              className={`w-full px-4 py-3 rounded-xl border-2 transition-all bg-zinc-50 dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 ${
                errors.shortDescription
                  ? "border-red-300 dark:border-red-700"
                  : "border-zinc-200 dark:border-zinc-700 focus:border-sky-500"
              }`}
            />
            <div className="flex justify-between items-center mt-2">
              <p className={`text-xs font-medium ${
                formData.shortDescription.length > 70
                  ? "text-amber-600 dark:text-amber-400"
                  : "text-zinc-500 dark:text-zinc-400"
              }`}>
                {formData.shortDescription.length}/80 Zeichen
              </p>
              {formData.shortDescription.length > 70 && (
                <span className="text-xs text-amber-600 dark:text-amber-400">⚠️ Bald voll</span>
              )}
            </div>
            {errors.shortDescription && (
              <p className="text-red-600 dark:text-red-400 text-sm mt-1">
                {errors.shortDescription}
              </p>
            )}
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-semibold text-zinc-900 dark:text-white mb-2">
              Ausführliche Beschreibung
            </label>
            <textarea
              value={formData.description}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
              placeholder="Erzähle mehr über deine Reise..."
              rows={4}
              className="w-full px-4 py-3 rounded-xl border-2 border-zinc-200 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500 resize-none transition-all"
            />
          </div>

          {/* Buttons */}
          <div className="flex gap-3 pt-6 border-t border-zinc-100 dark:border-zinc-800">
            <button
              type="button"
              onClick={onCloseAction}
              className="flex-1 px-4 py-3 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-xl font-semibold transition-all duration-200 active:scale-95"
            >
              ← Abbrechen
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-3 text-white bg-gradient-to-r from-sky-600 to-sky-700 hover:from-sky-700 hover:to-sky-800 rounded-xl font-semibold transition-all duration-200 active:scale-95 shadow-lg shadow-sky-500/30"
            >
              ✓ Speichern
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

