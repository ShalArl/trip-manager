"use client";

import React, { useRef, useState, useEffect } from "react";
import { ImagePlus, Trash2, X } from "lucide-react";

import { LocationResponse, UpdateLocationRequest } from "@/types/location";
import { addLocationImage, deleteLocationImage } from "@/lib/api/locations";
import { components } from "@/generated/types";

type LocationImageResponse = components["schemas"]["LocationImageResponse"];

// ─── Types ────────────────────────────────────────────────────────────────────

type Props = {
    isOpen: boolean;
    location: LocationResponse;
    tripId: string;
    isEditable: boolean;
    initialEditing?: boolean;
    onCloseAction: () => void;
    onSaveAction: (req: UpdateLocationRequest) => void;
    onDeleteAction: () => void;
    onLocationUpdateAction?: (location: LocationResponse) => void;
};

// ─── Component ────────────────────────────────────────────────────────────────

export default function LocationDetailModal({
    isOpen,
    location,
    tripId,
    isEditable,
    initialEditing,
    onCloseAction,
    onSaveAction,
    onDeleteAction,
    onLocationUpdateAction,
}: Props) {
    const [isEditing, setIsEditing] = useState(initialEditing ?? false);
    const [formData, setFormData] = useState({
        name: location.name,
        city: location.city,
        country: location.country,
        shortDescription: location.shortDescription,
        dateFrom: location.dateFrom,
        dateTo: location.dateTo,
        notes: location.notes ?? "",
    });

    useEffect(() => {
        if (isOpen) {
            setIsEditing(initialEditing ?? false);
        }
    }, [isOpen, initialEditing]);

    // ── Image state ──────────────────────────────────────────────────────────
    const [images, setImages] = useState<LocationImageResponse[]>(
        location.images ?? []
    );
    const [isUploadingImage, setIsUploadingImage] = useState(false);
    const [imageError, setImageError] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    // ── Handlers ─────────────────────────────────────────────────────────────

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSaveAction({
            name: formData.name,
            city: formData.city,
            country: formData.country,
            shortDescription: formData.shortDescription,
            dateFrom: formData.dateFrom,
            dateTo: formData.dateTo,
            notes: formData.notes || undefined,
        });
        setIsEditing(false);
    };

    const handleImageSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        if (!file.type.startsWith("image/")) {
            setImageError("Bitte wähle eine Bilddatei aus");
            return;
        }
        if (file.size > 10 * 1024 * 1024) {
            setImageError("Datei muss kleiner als 10MB sein");
            return;
        }

        setImageError(null);
        setIsUploadingImage(true);
        try {
            const created = await addLocationImage(tripId, location.id!, file, images.length);
            const updatedImages = [...images, created];
            setImages(updatedImages);

            if (onLocationUpdateAction) {
                onLocationUpdateAction({ ...location, images: updatedImages });
            }
        } catch (err) {
            setImageError("Fehler beim Hochladen des Bildes");
            console.error("[LocationDetailModal] addLocationImage:", err);
        } finally {
            setIsUploadingImage(false);
            if (fileInputRef.current) fileInputRef.current.value = "";
        }
    };

    const handleDeleteImage = async (image: LocationImageResponse) => {
        try {
            await deleteLocationImage(tripId, location.id!, image.id.toString());
            const updatedImages = images.filter((i) => i.id !== image.id);
            setImages(updatedImages);

            if (onLocationUpdateAction) {
                onLocationUpdateAction({ ...location, images: updatedImages });
            }
        } catch (err) {
            setImageError("Fehler beim Löschen des Bildes");
            console.error("[LocationDetailModal] deleteLocationImage:", err);
        }
    };

    if (!isOpen) return null;

    // ── Render ───────────────────────────────────────────────────────────────

    return (
        <div className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
            <div className="bg-white dark:bg-zinc-900 rounded-2xl shadow-xl max-w-lg w-full border border-zinc-100 dark:border-zinc-800 overflow-hidden max-h-[90vh] flex flex-col">

                {/* Header */}
                <div className="bg-gradient-to-r from-sky-500 to-sky-600 px-8 py-6 shrink-0">
                    <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                        📍 {location.name}
                    </h2>
                    <p className="text-sky-100 text-sm mt-1">
                        {location.city}, {location.country}
                    </p>
                </div>

                <div className="p-8 overflow-y-auto flex-1">
                    {!isEditing ? (
                        <div className="space-y-6">
                            {/* Details */}
                            <div className="space-y-4">
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Kurzbeschreibung</p>
                                    <p className="text-zinc-700 dark:text-zinc-300">{location.shortDescription}</p>
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Von</p>
                                        <p className="text-zinc-700 dark:text-zinc-300">{location.dateFrom}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Bis</p>
                                        <p className="text-zinc-700 dark:text-zinc-300">{location.dateTo}</p>
                                    </div>
                                </div>
                                {location.notes && (
                                    <div>
                                        <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-1">Notizen</p>
                                        <p className="text-zinc-700 dark:text-zinc-300">{location.notes}</p>
                                    </div>
                                )}
                            </div>

                            {/* Images */}
                            <div className="border-t border-zinc-100 dark:border-zinc-800 pt-5">
                                <div className="flex items-center justify-between mb-3">
                                    <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider">
                                        Bilder ({images.length})
                                    </p>
                                    {isEditable && (
                                        <>
                                            <input
                                                ref={fileInputRef}
                                                type="file"
                                                accept="image/*"
                                                onChange={handleImageSelect}
                                                className="hidden"
                                            />
                                            <button
                                                onClick={() => fileInputRef.current?.click()}
                                                disabled={isUploadingImage}
                                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors disabled:opacity-50"
                                            >
                                                <ImagePlus className="w-3.5 h-3.5" />
                                                {isUploadingImage ? "Lädt hoch..." : "Bild hinzufügen"}
                                            </button>
                                        </>
                                    )}
                                </div>

                                {imageError && (
                                    <p className="text-sm text-red-500 mb-3">{imageError}</p>
                                )}

                                {images.length === 0 ? (
                                    <div className="flex flex-col items-center justify-center py-8 text-zinc-400 dark:text-zinc-500 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl">
                                        <ImagePlus className="w-8 h-8 mb-2 opacity-30" />
                                        <p className="text-sm">Keine Bilder vorhanden</p>
                                    </div>
                                ) : (
                                    <div className="grid grid-cols-2 gap-3">
                                        {images.map((image) => (
                                            <div key={image.id.toString()} className="relative group rounded-xl overflow-hidden aspect-video bg-zinc-100 dark:bg-zinc-800">
                                                <img
                                                    src={image.imageUrl}
                                                    alt="Location"
                                                    className="w-full h-full object-cover"
                                                />
                                                {isEditable && (
                                                    <button
                                                        onClick={() => handleDeleteImage(image)}
                                                        className="absolute top-2 right-2 p-1.5 bg-red-500 hover:bg-red-600 text-white rounded-lg opacity-0 group-hover:opacity-100 transition-opacity"
                                                    >
                                                        <X className="w-3.5 h-3.5" />
                                                    </button>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>

                            {/* Actions */}
                            <div className="flex gap-3 pt-2 border-t border-zinc-100 dark:border-zinc-800">
                                <button
                                    onClick={onCloseAction}
                                    className="flex-1 px-4 py-3 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-xl font-semibold transition-all"
                                >
                                    Schließen
                                </button>
                                {isEditable && (
                                    <>
                                        <button
                                            onClick={() => setIsEditing(true)}
                                            className="flex-1 px-4 py-3 text-white bg-gradient-to-r from-sky-600 to-sky-700 hover:from-sky-700 hover:to-sky-800 rounded-xl font-semibold transition-all"
                                        >
                                            ✏️ Bearbeiten
                                        </button>
                                        <button
                                            onClick={onDeleteAction}
                                            className="px-4 py-3 text-white bg-red-500 hover:bg-red-600 rounded-xl font-semibold transition-all"
                                        >
                                            <Trash2 className="w-4 h-4" />
                                        </button>
                                    </>
                                )}
                            </div>
                        </div>
                    ) : (
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Name *</label>
                                    <input type="text" value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Stadt *</label>
                                    <input type="text" value={formData.city} onChange={(e) => setFormData({ ...formData, city: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Land *</label>
                                <input type="text" value={formData.country} onChange={(e) => setFormData({ ...formData, country: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Kurzbeschreibung *</label>
                                <input type="text" value={formData.shortDescription} onChange={(e) => setFormData({ ...formData, shortDescription: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Von *</label>
                                    <input type="date" value={formData.dateFrom} onChange={(e) => setFormData({ ...formData, dateFrom: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Bis *</label>
                                    <input type="date" value={formData.dateTo} onChange={(e) => setFormData({ ...formData, dateTo: e.target.value })} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500" />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Notizen</label>
                                <textarea value={formData.notes} onChange={(e) => setFormData({ ...formData, notes: e.target.value })} rows={3} className="w-full px-4 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none" />
                            </div>
                            <div className="flex gap-3 pt-4 border-t border-zinc-100 dark:border-zinc-800">
                                <button type="button" onClick={() => setIsEditing(false)} className="flex-1 px-4 py-3 text-zinc-700 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded-xl font-semibold transition-all">
                                    Abbrechen
                                </button>
                                <button type="submit" className="flex-1 px-4 py-3 text-white bg-gradient-to-r from-sky-600 to-sky-700 hover:from-sky-700 hover:to-sky-800 rounded-xl font-semibold transition-all">
                                    ✓ Speichern
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}