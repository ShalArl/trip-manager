"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { Pencil, Plus, Trash2, ChevronDown, ChevronUp } from "lucide-react";

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
import { UserAvatar } from "@/components/global/UserAvatar";
import WeatherWidget from "@/components/trips/components/WeatherWidget";
import TravelWarningWidget from "@/components/trips/components/TravelWarningWidget";

type TripResponse = components["schemas"]["TripResponse"];

type Props = {
    trip: TripResponse;
    isEditable?: boolean;
    onTripUpdateAction: (trip: TripResponse) => void;
    currentUser?: UserResponse | null;
};

const STATUS_CONFIG: Record<string, { label: string; className: string }> = {
    planned:   { label: "Geplant",     className: "bg-sky-50 text-sky-700 dark:bg-sky-950/50 dark:text-sky-400 border border-sky-200 dark:border-sky-800" },
    ongoing:   { label: "Unterwegs",   className: "bg-green-50 text-green-700 dark:bg-green-950/50 dark:text-green-400 border border-green-200 dark:border-green-800" },
    completed: { label: "Abgeschlossen", className: "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700" },
    cancelled: { label: "Abgesagt",    className: "bg-red-50 text-red-700 dark:bg-red-950/50 dark:text-red-400 border border-red-200 dark:border-red-800" },
};

const TRANSPORT_EMOJI: Record<string, string> = { flight: "✈️", train: "🚂", car: "🚗", bus: "🚌" };

function SectionHeader({ icon, title, count, onAdd, isEditable }: {
    icon: string; title: string; count: number;
    onAdd?: () => void; isEditable?: boolean;
}) {
    return (
        <div className="flex items-center justify-between mb-5">
            <div className="flex items-center gap-2">
                <span className="text-lg">{icon}</span>
                <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">{title}</h2>
                <span className="text-xs text-zinc-400 dark:text-zinc-500 bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded-full">{count}</span>
            </div>
            {isEditable && onAdd && (
                <button
                    onClick={onAdd}
                    className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                >
                    <Plus className="w-3.5 h-3.5" /> Hinzufügen
                </button>
            )}
        </div>
    );
}

function EmptyState({ icon, text }: { icon: string; text: string }) {
    return (
        <div className="flex flex-col items-center justify-center py-10 text-zinc-300 dark:text-zinc-600">
            <span className="text-3xl mb-2 opacity-50">{icon}</span>
            <p className="text-sm">{text}</p>
        </div>
    );
}

export default function TripDetail({ trip, isEditable = false, onTripUpdateAction, currentUser }: Props) {
    const [currentTrip, setCurrentTrip] = useState<TripResponse>(trip);
    const [isEditingTrip, setIsEditingTrip] = useState(false);

    const [locations, setLocations] = useState<LocationResponse[]>([]);
    const [selectedLocationId, setSelectedLocationId] = useState<string | null>(null);
    const activeLocation = locations.find((l) => l.id === selectedLocationId);

    const [activities, setActivities] = useState<any[]>([]);
    const [showAddActivityModal, setShowAddActivityModal] = useState(false);
    const selectedLocationActivities = activities.filter((a) => a.locationId === selectedLocationId);

    const [transports, setTransports] = useState<TransportResponse[]>([]);
    const [detailTransport, setDetailTransport] = useState<TransportResponse | null>(null);
    const [showAddTransportModal, setShowAddTransportModal] = useState(false);
    const [showEditTransportModal, setShowEditTransportModal] = useState(false);

    const [accommodations, setAccommodations] = useState<AccommodationResponse[]>([]);
    const [detailAccommodation, setDetailAccommodation] = useState<AccommodationResponse | null>(null);
    const [showAddAccommodationModal, setShowAddAccommodationModal] = useState(false);
    const [showEditAccommodationModal, setShowEditAccommodationModal] = useState(false);

    const [likeInfo, setLikeInfo] = useState<TripLikeResponse>({ likeCount: 0, hasLiked: false });
    const [comments, setComments] = useState<TripCommentResponse[]>([]);
    const [showComments, setShowComments] = useState(false);
    const [newComment, setNewComment] = useState("");
    const [isSubmittingComment, setIsSubmittingComment] = useState(false);

    useEffect(() => { getLocations(trip.id).then(setLocations).catch(console.error); }, [trip.id]);
    useEffect(() => { getTransports(trip.id).then(setTransports).catch(console.error); }, [trip.id]);
    useEffect(() => { getAccommodations(trip.id).then(setAccommodations).catch(console.error); }, [trip.id]);
    useEffect(() => { getTripLikes(trip.id).then(setLikeInfo).catch(console.error); }, [trip.id, currentUser]);

    const handleEditTrip = async (updatedTrip: Partial<TripResponse>) => {
        try {
            const updated = await updateTrip(trip.id, updatedTrip);
            setCurrentTrip(updated);
            onTripUpdateAction(updated);
            setIsEditingTrip(false);
        } catch (error) { console.error("[TripDetail] updateTrip:", error); }
    };

    const handleAddActivity = (newActivity: any) => {
        setActivities([...activities, { id: `act-${Date.now()}`, locationId: selectedLocationId!, ...newActivity }]);
    };

    const handleAddTransport = async (req: CreateTransportRequest) => {
        try { const c = await createTransport(trip.id, req); setTransports([...transports, c]); }
        catch (error) { console.error("[TripDetail] createTransport:", error); }
    };

    const handleUpdateTransport = async (req: UpdateTransportRequest) => {
        if (!detailTransport) return;
        try {
            const u = await updateTransport(trip.id, detailTransport.id!, req);
            setTransports(transports.map((t) => t.id === u.id ? u : t));
            setShowEditTransportModal(false);
        } catch (error) { console.error("[TripDetail] updateTransport:", error); }
    };

    const handleDeleteTransport = async () => {
        if (!detailTransport) return;
        try {
            await deleteTransport(trip.id, detailTransport.id!);
            setTransports(transports.filter((t) => t.id !== detailTransport.id));
            setShowEditTransportModal(false);
        } catch (error) { console.error("[TripDetail] deleteTransport:", error); }
    };

    const handleAddAccommodation = async (req: CreateAccommodationRequest) => {
        try { const c = await createAccommodation(trip.id, req); setAccommodations([...accommodations, c]); }
        catch (error) { console.error("[TripDetail] createAccommodation:", error); }
    };

    const handleUpdateAccommodation = async (req: UpdateAccommodationRequest) => {
        if (!detailAccommodation) return;
        try {
            const u = await updateAccommodation(trip.id, detailAccommodation.id!, req);
            setAccommodations(accommodations.map((a) => a.id === u.id ? u : a));
            setShowEditAccommodationModal(false);
        } catch (error) { console.error("[TripDetail] updateAccommodation:", error); }
    };

    const handleDeleteAccommodation = async () => {
        if (!detailAccommodation) return;
        try {
            await deleteAccommodation(trip.id, detailAccommodation.id!);
            setAccommodations(accommodations.filter((a) => a.id !== detailAccommodation.id));
            setShowEditAccommodationModal(false);
        } catch (error) { console.error("[TripDetail] deleteAccommodation:", error); }
    };

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
        } catch (error) { console.error("[TripDetail] like:", error); }
    };

    const handleShowComments = async () => {
        if (!showComments && comments.length === 0) {
            try {
                const data = await getTripComments(trip.id);
                setComments(data.data ?? []);
            } catch (error) { console.error("[TripDetail] getTripComments:", error); }
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
        } catch (error) { console.error("[TripDetail] createTripComment:", error); }
        finally { setIsSubmittingComment(false); }
    };

    const handleDeleteComment = async (commentId: string) => {
        try {
            await deleteTripComment(trip.id, commentId);
            setComments(comments.filter((c) => c.id !== commentId));
        } catch (error) { console.error("[TripDetail] deleteTripComment:", error); }
    };

    const statusCfg = STATUS_CONFIG[currentTrip.status ?? "planned"];

    return (
        <div className="max-w-6xl mx-auto px-6 py-10">
            <Link
                href="/"
                className="inline-flex items-center gap-1.5 text-sm text-zinc-400 dark:text-zinc-500 hover:text-sky-600 dark:hover:text-sky-400 transition-colors mb-8"
            >
                ← Zurück
            </Link>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

                {/* ── Left column ── */}
                <div className="lg:col-span-2 space-y-5">

                    {/* Trip header */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6">
                        <div className="flex items-start justify-between mb-5">
                            <div className="flex items-center gap-4">
                                <div className="w-12 h-12 rounded-xl bg-sky-50 dark:bg-sky-950/50 border border-sky-100 dark:border-sky-900 flex items-center justify-center text-2xl shrink-0">
                                    ✈️
                                </div>
                                <div>
                                    <h1 className="text-xl font-semibold text-zinc-900 dark:text-white">{currentTrip.title}</h1>
                                    <p className="text-sm text-zinc-400 dark:text-zinc-500 mt-0.5">
                                        {currentTrip.startDate} · {currentTrip.endDate}
                                    </p>
                                </div>
                            </div>
                            <div className="flex items-center gap-2 shrink-0">
                                <span className={`text-xs font-medium px-2.5 py-1 rounded-md ${statusCfg.className}`}>
                                    {statusCfg.label}
                                </span>
                                {isEditable && (
                                    <button
                                        onClick={() => setIsEditingTrip(true)}
                                        className="p-2 text-zinc-400 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                    >
                                        <Pencil className="w-4 h-4" />
                                    </button>
                                )}
                            </div>
                        </div>

                        <div className="space-y-3 text-sm text-zinc-600 dark:text-zinc-400">
                            {currentTrip.shortDescription && <p>{currentTrip.shortDescription}</p>}
                            {currentTrip.description && (
                                <p className="leading-relaxed text-zinc-500 dark:text-zinc-500">{currentTrip.description}</p>
                            )}
                        </div>

                        <div className="flex items-center gap-3 mt-5 pt-5 border-t border-zinc-100 dark:border-zinc-800">
                            <button
                                onClick={handleLike}
                                disabled={!currentUser}
                                className={`flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-sm font-medium transition-colors border disabled:opacity-40 disabled:cursor-not-allowed ${
                                    likeInfo.hasLiked
                                        ? "bg-sky-50 dark:bg-sky-950/50 text-sky-600 dark:text-sky-400 border-sky-200 dark:border-sky-800"
                                        : "bg-white dark:bg-zinc-900 text-zinc-500 dark:text-zinc-400 border-zinc-200 dark:border-zinc-700 hover:border-sky-300"
                                }`}
                            >
                                {likeInfo.hasLiked ? "❤️" : "🤍"}
                                <span>{likeInfo.likeCount} {likeInfo.likeCount === 1 ? "Like" : "Likes"}</span>
                            </button>
                            <button
                                onClick={handleShowComments}
                                className="flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-sm font-medium bg-white dark:bg-zinc-900 text-zinc-500 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 hover:border-sky-300 transition-colors"
                            >
                                💬
                                <span>{showComments ? "Ausblenden" : "Kommentare"}</span>
                                {showComments ? <ChevronUp className="w-3.5 h-3.5" /> : <ChevronDown className="w-3.5 h-3.5" />}
                            </button>
                            <div className="flex items-center gap-2 ml-auto">
                                <UserAvatar
                                    avatarKey={currentTrip.createdBy.avatarUrl ?? null}
                                    name={currentTrip.createdBy.name}
                                    className="h-7 w-7"
                                />
                                <span className="text-xs text-zinc-400">{currentTrip.createdBy.name}</span>
                            </div>
                        </div>

                        {/* Comments */}
                        {showComments && (
                            <div className="mt-5 pt-5 border-t border-zinc-100 dark:border-zinc-800 space-y-3">
                                {comments.length === 0 ? (
                                    <p className="text-sm text-zinc-400 dark:text-zinc-500 text-center py-3">Noch keine Kommentare</p>
                                ) : (
                                    comments.map((comment) => (
                                        <div key={comment.id} className="flex items-start justify-between gap-3">
                                            <div className="flex items-start gap-2.5">
                                                <UserAvatar name={comment.user.name} avatarKey={comment.user.avatarUrl} className="h-7 w-7 shrink-0" />
                                                <div className="bg-zinc-50 dark:bg-zinc-800 rounded-xl px-3 py-2">
                                                    <p className="text-xs font-medium text-zinc-700 dark:text-zinc-300">{comment.user.name}</p>
                                                    <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-0.5">{comment.text}</p>
                                                </div>
                                            </div>
                                            {currentUser && comment.user.id === currentUser.id && (
                                                <button
                                                    onClick={() => handleDeleteComment(comment.id)}
                                                    className="p-1.5 text-zinc-300 hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/30 rounded-lg transition-colors shrink-0"
                                                >
                                                    <Trash2 className="w-3.5 h-3.5" />
                                                </button>
                                            )}
                                        </div>
                                    ))
                                )}
                                {currentUser && (
                                    <div className="flex gap-2 pt-2">
                                        <input
                                            type="text"
                                            value={newComment}
                                            onChange={(e) => setNewComment(e.target.value)}
                                            onKeyDown={(e) => e.key === "Enter" && handleSubmitComment()}
                                            placeholder="Kommentar schreiben..."
                                            className="flex-1 px-3.5 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl focus:outline-none focus:border-sky-400 dark:focus:border-sky-600 text-zinc-900 dark:text-white placeholder-zinc-400"
                                        />
                                        <button
                                            onClick={handleSubmitComment}
                                            disabled={!newComment.trim() || isSubmittingComment}
                                            className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-xl transition-colors disabled:opacity-40"
                                        >
                                            {isSubmittingComment ? "…" : "Senden"}
                                        </button>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>

                    {/* Locations */}
                    <TripLocationsSection
                        tripId={trip.id}
                        isEditable={isEditable}
                        locations={locations}
                        selectedLocationId={selectedLocationId}
                        onLocationsChange={setLocations}
                        onLocationSelect={setSelectedLocationId}
                    />

                    {/* Transports */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6">
                        <SectionHeader
                            icon="🚀" title="Transporte" count={transports.length}
                            isEditable={isEditable} onAdd={() => setShowAddTransportModal(true)}
                        />
                        {transports.length === 0 ? (
                            <EmptyState icon="🚌" text="Kein Transport hinzugefügt" />
                        ) : (
                            <div className="space-y-2">
                                {transports.map((t) => (
                                    <div key={t.id} className="flex items-center gap-3 p-3.5 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-100 dark:border-zinc-800">
                                        <div className="w-9 h-9 rounded-lg bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 flex items-center justify-center text-lg shrink-0">
                                            {TRANSPORT_EMOJI[t.type ?? "flight"] ?? "🚗"}
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center gap-2">
                                                <span className="text-sm font-medium text-zinc-900 dark:text-white truncate">{t.from?.name || "–"}</span>
                                                <span className="text-zinc-300 dark:text-zinc-600 text-xs shrink-0">→</span>
                                                <span className="text-sm font-medium text-zinc-900 dark:text-white truncate">{t.to?.name || "–"}</span>
                                            </div>
                                            <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-0.5">
                                                {t.from?.city && t.to?.city ? `${t.from.city} → ${t.to.city}` : ""}
                                                {t.departureTime ? ` · ${new Date(t.departureTime).toLocaleString("de-DE", { day: "2-digit", month: "2-digit", hour: "2-digit", minute: "2-digit" })}` : ""}
                                            </p>
                                        </div>
                                        {isEditable && (
                                            <button
                                                onClick={() => { setDetailTransport(t); setShowEditTransportModal(true); }}
                                                className="p-1.5 text-zinc-300 hover:text-sky-500 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors shrink-0"
                                            >
                                                <Pencil className="w-3.5 h-3.5" />
                                            </button>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Accommodations */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6">
                        <SectionHeader
                            icon="🏨" title="Unterkünfte" count={accommodations.length}
                            isEditable={isEditable} onAdd={() => setShowAddAccommodationModal(true)}
                        />
                        {accommodations.length === 0 ? (
                            <EmptyState icon="🏨" text="Keine Unterkunft hinzugefügt" />
                        ) : (
                            <div className="space-y-2">
                                {accommodations.map((a) => (
                                    <div key={a.id} className="flex items-center gap-3 p-3.5 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-100 dark:border-zinc-800">
                                        <div className="w-9 h-9 rounded-lg bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 flex items-center justify-center text-lg shrink-0">
                                            🏨
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <p className="text-sm font-medium text-zinc-900 dark:text-white truncate">{a.name}</p>
                                            <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-0.5 truncate">
                                                {a.location?.city ? `${a.location.city}, ${a.location.country}` : "Kein Ort"}
                                                {a.checkIn ? ` · ${new Date(a.checkIn).toLocaleDateString("de-DE", { day: "2-digit", month: "2-digit" })}` : ""}
                                                {a.pricePerNight ? ` · ${a.pricePerNight} €/Nacht` : ""}
                                            </p>
                                        </div>
                                        {isEditable && (
                                            <button
                                                onClick={() => { setDetailAccommodation(a); setShowEditAccommodationModal(true); }}
                                                className="p-1.5 text-zinc-300 hover:text-sky-500 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors shrink-0"
                                            >
                                                <Pencil className="w-3.5 h-3.5" />
                                            </button>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>

                {/* ── Right sidebar ── */}
                <div className="space-y-5">
                    {/* Weather & Warning for selected location */}
                    {activeLocation?.latitude != null && activeLocation?.longitude != null && (
                        <WeatherWidget
                            lat={activeLocation.latitude}
                            lng={activeLocation.longitude}
                            locationName={activeLocation.name}
                        />
                    )}
                    {activeLocation?.countryCode && (
                        <TravelWarningWidget
                            countryCode={activeLocation.countryCode}
                            countryName={activeLocation.country}
                        />
                    )}

                    {/* Activities for selected location */}
                    {activeLocation && (
                        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                            <div className="flex items-center justify-between mb-4">
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-0.5">Aktivitäten</p>
                                    <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">{activeLocation.name}</h3>
                                </div>
                                {isEditable && (
                                    <button
                                        onClick={() => setShowAddActivityModal(true)}
                                        className="p-1.5 text-sky-600 dark:text-sky-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                    >
                                        <Plus className="w-4 h-4" />
                                    </button>
                                )}
                            </div>
                            {selectedLocationActivities.length === 0 ? (
                                <p className="text-sm text-zinc-400 dark:text-zinc-500 text-center py-4">Keine Aktivitäten</p>
                            ) : (
                                <div className="space-y-2">
                                    {selectedLocationActivities.map((activity) => (
                                        <div key={activity.id} className="p-3 bg-zinc-50 dark:bg-zinc-800 rounded-lg">
                                            <p className="text-sm font-medium text-zinc-900 dark:text-white">{activity.name}</p>
                                            <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-0.5">{activity.category}</p>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}

                    {/* Placeholder when no location selected */}
                    {!activeLocation && (
                        <div className="bg-white dark:bg-zinc-900 border border-dashed border-zinc-200 dark:border-zinc-700 rounded-2xl p-6 text-center">
                            <p className="text-sm text-zinc-400 dark:text-zinc-500">
                                Wähle einen Ort aus um Wetter, Reisewarnung und Aktivitäten zu sehen
                            </p>
                        </div>
                    )}
                </div>
            </div>

            {/* Modals */}
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
        </div>
    );
}