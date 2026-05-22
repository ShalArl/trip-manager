"use client";

import Link from "next/link";
import { useState, useEffect } from "react";

import { updateTrip } from "@/lib/api/trips";
import { likeTrip, unlikeTrip, getTripComments, createTripComment, deleteTripComment, getTripLikes } from "@/lib/api/social";
import { createTransport, getTransports, updateTransport, deleteTransport } from "@/lib/api/transports";
import { getLocations } from "@/lib/api/locations";
import { getAccommodations, createAccommodation, updateAccommodation, deleteAccommodation } from "@/lib/api/accommodations";

import { components } from "@/generated/types";
import { TransportResponse, CreateTransportRequest, UpdateTransportRequest } from "@/types/transport";
import { LocationResponse } from "@/types/location";
import { AccommodationResponse, CreateAccommodationRequest, UpdateAccommodationRequest } from "@/types/accommodation";
import { UserResponse } from "@/types/user";
import { TripLikeResponse, TripCommentResponse } from "@/types/social";

import TripLocationsSection from "@/components/trips/sections/TripLocationsSection";
import AddActivityModal from "./modals/AddActivityModal";
import EditTripModal from "./modals/EditTripModal";
import AddTransportModal from "./modals/AddTransportModal";
import EditTransportModal from "./modals/EditTransportModal";
import AddAccommodationModal from "./modals/AddAccommodationModal";
import EditAccommodationModal from "./modals/EditAccommodationModal";

// ─── Types ────────────────────────────────────────────────────────────────────

type TripResponse = components["schemas"]["TripResponse"];

type Props = {
    trip: TripResponse;
    isEditable?: boolean;
    onTripUpdateAction: (trip: TripResponse) => void;
    currentUser?: UserResponse | null;
};

// ─── Component ────────────────────────────────────────────────────────────────

export default function TripDetail({ trip, isEditable = false, onTripUpdateAction, currentUser }: Props) {
    // ── Trip ────────────────────────────────────────────────────────────────
    const [currentTrip, setCurrentTrip] = useState<TripResponse>(trip);
    const [isEditingTrip, setIsEditingTrip] = useState(false);

    // ── Locations ───────────────────────────────────────────────────────────
    const [locations, setLocations] = useState<LocationResponse[]>([]);
    const [selectedLocationId, setSelectedLocationId] = useState<string | null>(null);
    const activeLocation = locations.find((l) => l.id === selectedLocationId);

    // ── Activities ──────────────────────────────────────────────────────────
    const [activities, setActivities] = useState<any[]>([]);
    const [showAddActivityModal, setShowAddActivityModal] = useState(false);
    const selectedLocationActivities = activities.filter((a) => a.locationId === selectedLocationId);

    // ── Transports ──────────────────────────────────────────────────────────
    const [transports, setTransports] = useState<TransportResponse[]>([]);
    const [detailTransport, setDetailTransport] = useState<TransportResponse | null>(null);
    const [showAddTransportModal, setShowAddTransportModal] = useState(false);
    const [showEditTransportModal, setShowEditTransportModal] = useState(false);

    // ── Accommodations ──────────────────────────────────────────────────────
    const [accommodations, setAccommodations] = useState<AccommodationResponse[]>([]);
    const [detailAccommodation, setDetailAccommodation] = useState<AccommodationResponse | null>(null);
    const [showAddAccommodationModal, setShowAddAccommodationModal] = useState(false);
    const [showEditAccommodationModal, setShowEditAccommodationModal] = useState(false);

    // ── Social ──────────────────────────────────────────────────────────────
    const [likeInfo, setLikeInfo] = useState<TripLikeResponse>({ likeCount: 0, hasLiked: false });
    const [comments, setComments] = useState<TripCommentResponse[]>([]);
    const [showComments, setShowComments] = useState(false);
    const [newComment, setNewComment] = useState("");
    const [isSubmittingComment, setIsSubmittingComment] = useState(false);

    // ── Data fetching ────────────────────────────────────────────────────────

    useEffect(() => {
        getLocations(trip.id).then(setLocations).catch(console.error);
    }, [trip.id]);

    useEffect(() => {
        getTransports(trip.id).then(setTransports).catch(console.error);
    }, [trip.id]);

    useEffect(() => {
        getAccommodations(trip.id).then(setAccommodations).catch(console.error);
    }, [trip.id]);

    useEffect(() => {
        getTripLikes(trip.id).then(setLikeInfo).catch(console.error);
    }, [trip.id, currentUser]);

    // ── Trip handlers ────────────────────────────────────────────────────────

    const handleEditTrip = async (updatedTrip: Partial<TripResponse>) => {
        try {
            const updated = await updateTrip(trip.id, updatedTrip);
            setCurrentTrip(updated);
            onTripUpdateAction(updated);
            setIsEditingTrip(false);
        } catch (error) {
            console.error("[TripDetail] updateTrip:", error);
        }
    };

    // ── Activity handlers ────────────────────────────────────────────────────

    const handleAddActivity = (newActivity: any) => {
        setActivities([...activities, {
            id: `act-${Date.now()}`,
            locationId: selectedLocationId!,
            ...newActivity,
        }]);
    };

    // ── Transport handlers ───────────────────────────────────────────────────

    const handleAddTransport = async (req: CreateTransportRequest) => {
        try {
            const created = await createTransport(trip.id, req);
            setTransports([...transports, created]);
        } catch (error) {
            console.error("[TripDetail] createTransport:", error);
        }
    };

    const handleUpdateTransport = async (req: UpdateTransportRequest) => {
        if (!detailTransport) return;
        try {
            const updated = await updateTransport(trip.id, detailTransport.id!, req);
            setTransports(transports.map((t) => t.id === updated.id ? updated : t));
            setShowEditTransportModal(false);
        } catch (error) {
            console.error("[TripDetail] updateTransport:", error);
        }
    };

    const handleDeleteTransport = async () => {
        if (!detailTransport) return;
        try {
            await deleteTransport(trip.id, detailTransport.id!);
            setTransports(transports.filter((t) => t.id !== detailTransport.id));
            setShowEditTransportModal(false);
        } catch (error) {
            console.error("[TripDetail] deleteTransport:", error);
        }
    };

    // ── Accommodation handlers ───────────────────────────────────────────────

    const handleAddAccommodation = async (req: CreateAccommodationRequest) => {
        try {
            const created = await createAccommodation(trip.id, req);
            setAccommodations([...accommodations, created]);
        } catch (error) {
            console.error("[TripDetail] createAccommodation:", error);
        }
    };

    const handleUpdateAccommodation = async (req: UpdateAccommodationRequest) => {
        if (!detailAccommodation) return;
        try {
            const updated = await updateAccommodation(trip.id, detailAccommodation.id!, req);
            setAccommodations(accommodations.map((a) => a.id === updated.id ? updated : a));
            setShowEditAccommodationModal(false);
        } catch (error) {
            console.error("[TripDetail] updateAccommodation:", error);
        }
    };

    const handleDeleteAccommodation = async () => {
        if (!detailAccommodation) return;
        try {
            await deleteAccommodation(trip.id, detailAccommodation.id!);
            setAccommodations(accommodations.filter((a) => a.id !== detailAccommodation.id));
            setShowEditAccommodationModal(false);
        } catch (error) {
            console.error("[TripDetail] deleteAccommodation:", error);
        }
    };

    // ── Social handlers ──────────────────────────────────────────────────────

    const handleLike = async () => {
        if (!currentUser) return;
        try {
            if (likeInfo.hasLiked) {
                await unlikeTrip(trip.id);
                setLikeInfo({ likeCount: likeInfo.likeCount - 1, hasLiked: false });
            } else {
                await likeTrip(trip.id);
                setLikeInfo({ likeCount: likeInfo.likeCount + 1, hasLiked: true });
            }
        } catch (error) {
            console.error("[TripDetail] like:", error);
        }
    };

    const handleShowComments = async () => {
        if (!showComments && comments.length === 0) {
            try {
                const data = await getTripComments(trip.id);
                setComments(data.data ?? []);
            } catch (error) {
                console.error("[TripDetail] getTripComments:", error);
            }
        }
        setShowComments(!showComments);
    };

    const handleSubmitComment = async () => {
        if (!newComment.trim() || !currentUser) return;
        setIsSubmittingComment(true);
        try {
            const created = await createTripComment(trip.id, newComment.trim());
            setComments([...comments, created]);
            setNewComment("");
        } catch (error) {
            console.error("[TripDetail] createTripComment:", error);
        } finally {
            setIsSubmittingComment(false);
        }
    };

    const handleDeleteComment = async (commentId: string) => {
        try {
            await deleteTripComment(trip.id, commentId);
            setComments(comments.filter((c) => c.id !== commentId));
        } catch (error) {
            console.error("[TripDetail] deleteTripComment:", error);
        }
    };

    // ── Render ───────────────────────────────────────────────────────────────

    return (
        <div className="max-w-5xl px-6 py-12">
            <Link
                href="/"
                className="inline-flex items-center gap-2 text-sm text-zinc-500 dark:text-zinc-400 hover:text-sky-600 dark:hover:text-sky-400 transition-colors mb-8"
            >
                ← Zurück zur Übersicht
            </Link>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">

                {/* ── Left: Main Content ───────────────────────────────────── */}
                <div className="lg:col-span-2 space-y-6">

                    {/* ── Trip Header ─────────────────────────────────────── */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-start justify-between mb-6">
                            <div className="flex items-center gap-4">
                                <div className="w-14 h-14 rounded-2xl bg-sky-50 dark:bg-sky-950/50 border border-sky-100 dark:border-sky-900 flex items-center justify-center text-3xl">
                                    ✈️
                                </div>
                                <div>
                                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
                                        {currentTrip.title}
                                    </h1>
                                    <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                        {currentTrip.startDate} · {currentTrip.endDate}
                                    </p>
                                </div>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setIsEditingTrip(!isEditingTrip)}
                                    className="px-3 py-1.5 text-sm font-medium text-sky-600 dark:text-sky-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                >
                                    {isEditingTrip ? "Speichern" : "Bearbeiten"}
                                </button>
                            )}
                        </div>
                        <div className="border-t border-zinc-100 dark:border-zinc-800 pt-6 space-y-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-2">
                                    Kurzbeschreibung
                                </p>
                                <p className="text-zinc-700 dark:text-zinc-300">{currentTrip.shortDescription}</p>
                            </div>
                            {currentTrip.description && (
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-2">
                                        Details
                                    </p>
                                    <p className="text-zinc-700 dark:text-zinc-300 leading-relaxed">{currentTrip.description}</p>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* ── Social: Likes & Comments ─────────────────────────── */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6">
                        <div className="flex items-center gap-4">
                            <button
                                onClick={handleLike}
                                disabled={!currentUser}
                                className={`flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-colors ${likeInfo.hasLiked
                                    ? "bg-sky-100 dark:bg-sky-950/50 text-sky-600 dark:text-sky-400 border border-sky-200 dark:border-sky-800"
                                    : "bg-zinc-50 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700"
                                    } disabled:opacity-50 disabled:cursor-not-allowed`}
                            >
                                <span>{likeInfo.hasLiked ? "❤️" : "🤍"}</span>
                                <span>{likeInfo.likeCount} {likeInfo.likeCount === 1 ? "Like" : "Likes"}</span>
                            </button>
                            <button
                                onClick={handleShowComments}
                                className="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium bg-zinc-50 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700 transition-colors"
                            >
                                <span>💬</span>
                                <span>{showComments ? "Kommentare ausblenden" : "Kommentare anzeigen"}</span>
                            </button>
                        </div>

                        {showComments && (
                            <div className="mt-6 space-y-4">
                                {comments.length === 0 ? (
                                    <p className="text-zinc-500 dark:text-zinc-400 text-sm text-center py-4">
                                        Noch keine Kommentare
                                    </p>
                                ) : (
                                    <div className="space-y-3">
                                        {comments.map((comment) => (
                                            <div key={comment.id} className="flex items-start justify-between gap-3 p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl">
                                                <div className="flex items-start gap-3">
                                                    <div className="w-8 h-8 rounded-full bg-sky-500 flex items-center justify-center text-white text-sm font-semibold shrink-0 overflow-hidden">
                                                        {comment.user.avatarUrl ? (
                                                            <img src={comment.user.avatarUrl} alt={comment.user.name} className="w-full h-full object-cover" />
                                                        ) : (
                                                            comment.user.name.charAt(0).toUpperCase()
                                                        )}
                                                    </div>
                                                    <div>
                                                        <p className="text-sm font-medium text-zinc-900 dark:text-white">{comment.user.name}</p>
                                                        <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">{comment.text}</p>
                                                    </div>
                                                </div>
                                                {currentUser && comment.user.id === currentUser.id && (
                                                    <button
                                                        onClick={() => handleDeleteComment(comment.id)}
                                                        className="p-1.5 rounded-lg text-zinc-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/30 transition-colors shrink-0"
                                                    >
                                                        🗑️
                                                    </button>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                )}
                                {currentUser && (
                                    <div className="flex gap-2 pt-2 border-t border-zinc-100 dark:border-zinc-800">
                                        <input
                                            type="text"
                                            value={newComment}
                                            onChange={(e) => setNewComment(e.target.value)}
                                            onKeyDown={(e) => e.key === "Enter" && handleSubmitComment()}
                                            placeholder="Kommentar schreiben..."
                                            className="flex-1 px-4 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl focus:outline-none focus:border-sky-400 dark:focus:border-sky-600 text-zinc-900 dark:text-white placeholder-zinc-400"
                                        />
                                        <button
                                            onClick={handleSubmitComment}
                                            disabled={!newComment.trim() || isSubmittingComment}
                                            className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                        >
                                            {isSubmittingComment ? "..." : "Senden"}
                                        </button>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>

                    {/* ── Locations ────────────────────────────────────────── */}
                    <TripLocationsSection
                        tripId={trip.id}
                        isEditable={isEditable}
                        locations={locations}
                        selectedLocationId={selectedLocationId}
                        onLocationsChange={setLocations}
                        onLocationSelect={setSelectedLocationId}
                    />

                    {/* ── Transports ───────────────────────────────────────── */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-center justify-between mb-6">
                            <div className="flex items-center gap-2">
                                <span className="text-xl">🚀</span>
                                <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                                    Transporte
                                    <span className="ml-2 text-sm font-normal text-zinc-400 dark:text-zinc-500">
                                        ({transports.length})
                                    </span>
                                </h2>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddTransportModal(true)}
                                    className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                                >
                                    <span>+</span> Transport
                                </button>
                            )}
                        </div>

                        {transports.length === 0 ? (
                            <div className="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500">
                                <span className="text-4xl mb-3 opacity-30">🚌</span>
                                <p className="text-sm">Kein Transport hinzugefügt</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {transports.map((t) => {
                                    const typeEmoji = { flight: "✈️", train: "🚂", car: "🚗", bus: "🚌" }[t.type ?? "flight"] ?? "🚗";
                                    return (
                                        <div key={t.id} className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700 hover:border-sky-200 dark:hover:border-sky-800 transition-colors">
                                            <div className="flex items-center justify-between gap-4">
                                                <div className="flex items-center gap-3 min-w-0">
                                                    <div className="shrink-0 w-10 h-10 rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 flex items-center justify-center text-xl shadow-sm">
                                                        {typeEmoji}
                                                    </div>
                                                    <div className="min-w-0">
                                                        <div className="flex items-center gap-2 flex-wrap">
                                                            <span className="font-semibold text-zinc-900 dark:text-white truncate">
                                                                {t.from?.name || "–"}
                                                            </span>
                                                            <span className="text-zinc-400 dark:text-zinc-500 shrink-0">→</span>
                                                            <span className="font-semibold text-zinc-900 dark:text-white truncate">
                                                                {t.to?.name || "–"}
                                                            </span>
                                                        </div>
                                                        <div className="flex items-center gap-2 mt-0.5 flex-wrap">
                                                            <span className="text-xs text-zinc-400 dark:text-zinc-500">
                                                                {t.from?.city && t.from?.country ? `${t.from.city}, ${t.from.country}` : ""}
                                                            </span>
                                                            {t.from?.city && t.to?.city && (
                                                                <span className="text-xs text-zinc-300 dark:text-zinc-600">·</span>
                                                            )}
                                                            <span className="text-xs text-zinc-400 dark:text-zinc-500">
                                                                {t.to?.city && t.to?.country ? `${t.to.city}, ${t.to.country}` : ""}
                                                            </span>
                                                        </div>
                                                        {(t.departureTime || t.arrivalTime) && (
                                                            <div className="flex items-center gap-3 mt-1 flex-wrap">
                                                                {t.departureTime && (
                                                                    <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                                        🕐 {new Date(t.departureTime).toLocaleString("de-DE", {
                                                                            day: "2-digit", month: "2-digit", year: "numeric",
                                                                            hour: "2-digit", minute: "2-digit"
                                                                        })}
                                                                    </span>
                                                                )}
                                                                {t.departureTime && t.arrivalTime && (
                                                                    <span className="text-xs text-zinc-300 dark:text-zinc-600">→</span>
                                                                )}
                                                                {t.arrivalTime && (
                                                                    <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                                        {new Date(t.arrivalTime).toLocaleString("de-DE", {
                                                                            day: "2-digit", month: "2-digit", year: "numeric",
                                                                            hour: "2-digit", minute: "2-digit"
                                                                        })}
                                                                    </span>
                                                                )}
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                                {isEditable && (
                                                    <button
                                                        onClick={() => { setDetailTransport(t); setShowEditTransportModal(true); }}
                                                        className="shrink-0 p-2 text-zinc-400 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                                        aria-label="Transport bearbeiten"
                                                    >
                                                        <svg xmlns="http://www.w3.org/2000/svg" className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                                        </svg>
                                                    </button>
                                                )}
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>

                    {/* ── Accommodations ───────────────────────────────────── */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-center justify-between mb-6">
                            <div className="flex items-center gap-2">
                                <span className="text-xl">🏨</span>
                                <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                                    Unterkünfte
                                    <span className="ml-2 text-sm font-normal text-zinc-400 dark:text-zinc-500">
                                        ({accommodations.length})
                                    </span>
                                </h2>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddAccommodationModal(true)}
                                    className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                                >
                                    <span>+</span> Unterkunft
                                </button>
                            )}
                        </div>

                        {accommodations.length === 0 ? (
                            <div className="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500">
                                <span className="text-4xl mb-3 opacity-30">🏨</span>
                                <p className="text-sm">Keine Unterkunft hinzugefügt</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {accommodations.map((a) => (
                                    <div key={a.id} className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700 hover:border-sky-200 dark:hover:border-sky-800 transition-colors">
                                        <div className="flex items-center justify-between gap-4">
                                            <div className="flex items-center gap-3 min-w-0">
                                                <div className="shrink-0 w-10 h-10 rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 flex items-center justify-center text-xl shadow-sm">
                                                    🏨
                                                </div>
                                                <div className="min-w-0">
                                                    <p className="font-semibold text-zinc-900 dark:text-white truncate">
                                                        {a.name}
                                                    </p>
                                                    <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5 truncate">
                                                        📍 {a.location?.name
                                                            ? `${a.location.name}, ${a.location.city}, ${a.location.country}`
                                                            : "Kein Ort angegeben"}
                                                    </p>
                                                    {(a.checkIn || a.checkOut) && (
                                                        <div className="flex items-center gap-3 mt-1 flex-wrap">
                                                            {a.checkIn && (
                                                                <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                                    🛬 Check-in: {new Date(a.checkIn).toLocaleDateString("de-DE", {
                                                                        day: "2-digit", month: "2-digit", year: "numeric"
                                                                    })}
                                                                </span>
                                                            )}
                                                            {a.checkIn && a.checkOut && (
                                                                <span className="text-xs text-zinc-300 dark:text-zinc-600">·</span>
                                                            )}
                                                            {a.checkOut && (
                                                                <span className="text-xs text-zinc-500 dark:text-zinc-400">
                                                                    🛫 Check-out: {new Date(a.checkOut).toLocaleDateString("de-DE", {
                                                                        day: "2-digit", month: "2-digit", year: "numeric"
                                                                    })}
                                                                </span>
                                                            )}
                                                        </div>
                                                    )}
                                                    {a.pricePerNight && (
                                                        <p className="text-xs text-sky-600 dark:text-sky-400 mt-1 font-medium">
                                                            {a.pricePerNight} € / Nacht
                                                        </p>
                                                    )}
                                                </div>
                                            </div>
                                            {isEditable && (
                                                <button
                                                    onClick={() => { setDetailAccommodation(a); setShowEditAccommodationModal(true); }}
                                                    className="shrink-0 p-2 text-zinc-400 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                                    aria-label="Unterkunft bearbeiten"
                                                >
                                                    <svg xmlns="http://www.w3.org/2000/svg" className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                                    </svg>
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                </div>{/* end lg:col-span-2 */}

                {/* ── Right: Activities ────────────────────────────────────── */}
                {activeLocation && (
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 h-fit">
                        <div className="flex items-center justify-between mb-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1">
                                    Aktivitäten in
                                </p>
                                <h3 className="text-lg font-bold text-zinc-900 dark:text-white">{activeLocation.name}</h3>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddActivityModal(true)}
                                    className="p-2 text-sky-600 dark:text-sky-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                >
                                    +
                                </button>
                            )}
                        </div>
                        {selectedLocationActivities.length === 0 ? (
                            <p className="text-zinc-500 dark:text-zinc-400 text-sm text-center py-4">Keine Aktivitäten</p>
                        ) : (
                            <div className="space-y-3">
                                {selectedLocationActivities.map((activity) => (
                                    <div key={activity.id} className="p-3 bg-zinc-50 dark:bg-zinc-800 rounded-lg">
                                        <p className="font-medium text-sm text-zinc-900 dark:text-white">{activity.name}</p>
                                        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">{activity.category}</p>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

            </div>{/* end grid */}

            {/* ── Modals ───────────────────────────────────────────────────── */}
            <AddActivityModal
                isOpen={showAddActivityModal}
                locationId={selectedLocationId}
                locationName={activeLocation?.name || ""}
                tripStartDate={trip.startDate}
                onCloseAction={() => setShowAddActivityModal(false)}
                onAddAction={handleAddActivity}
            />
            <EditTripModal
                isOpen={isEditingTrip}
                trip={currentTrip}
                onCloseAction={() => setIsEditingTrip(false)}
                onSaveAction={handleEditTrip}
            />
            <AddTransportModal
                isOpen={showAddTransportModal}
                onCloseAction={() => setShowAddTransportModal(false)}
                onAddAction={handleAddTransport}
            />
            <AddAccommodationModal
                isOpen={showAddAccommodationModal}
                onCloseAction={() => setShowAddAccommodationModal(false)}
                onAddAction={handleAddAccommodation}
            />
            {detailTransport && (
                <EditTransportModal
                    isOpen={showEditTransportModal}
                    transport={detailTransport}
                    onCloseAction={() => setShowEditTransportModal(false)}
                    onSaveAction={handleUpdateTransport}
                    onDeleteAction={handleDeleteTransport}
                />
            )}
            {detailAccommodation && (
                <EditAccommodationModal
                    isOpen={showEditAccommodationModal}
                    accommodation={detailAccommodation}
                    onCloseAction={() => setShowEditAccommodationModal(false)}
                    onSaveAction={handleUpdateAccommodation}
                    onDeleteAction={handleDeleteAccommodation}
                />
            )}

        </div> // end max-w-5xl
    );
}