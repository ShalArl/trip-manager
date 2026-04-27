import Link from "next/link";
import { useState, useEffect } from "react";
import { updateTrip } from "@/lib/api/trips";
import { likeTrip, unlikeTrip, getTripComments, createTripComment, deleteTripComment, getTripLikes } from "@/lib/api/social";
import { components } from "@/generated/types";
import { TransportResponse, CreateTransportRequest, UpdateTransportRequest } from "@/types/transport";
import { createTransport, getTransports } from "@/lib/api/transports";
import { LocationResponse, CreateLocationRequest, UpdateLocationRequest } from "@/types/location";
import { getLocations, createLocation, updateLocation, deleteLocation } from "@/lib/api/locations";
import EditTransportModal from "./modals/EditTransportModal";
import { updateTransport, deleteTransport } from "@/lib/api/transports";
import AddLocationModal from "./modals/AddLocationModal";
import AddActivityModal from "./modals/AddActivityModal";
import EditTripModal from "./modals/EditTripModal";
import AddTransportModal from "./modals/AddTransportModal";
import LocationDetailModal from "./modals/LocationDetailModal";
import { UserResponse } from "@/types/user";
import { TripLikeResponse, TripCommentResponse, TripCommentListResponse, CreateTripCommentRequest } from "@/types/social";

type TripResponse = components["schemas"]["TripResponse"];

type Props = {
    trip: TripResponse;
    isEditable?: boolean;
    onTripUpdate: (trip: TripResponse) => void;
    currentUser?: UserResponse | null;
};

export default function TripDetail({ trip, isEditable = false, onTripUpdate, currentUser }: Props) {
    const [isEditingTrip, setIsEditingTrip] = useState(false);
    const [currentTrip, setCurrentTrip] = useState<TripResponse>(trip);
    const [selectedLocationId, setSelectedLocationId] = useState<string | null>(null);

    const [showAddLocationModal, setShowAddLocationModal] = useState(false);
    const [showAddActivityModal, setShowAddActivityModal] = useState(false);
    const [showAddTransportModal, setShowAddTransportModal] = useState(false);
    const [showLocationDetailModal, setShowLocationDetailModal] = useState(false);

    const [transports, setTransports] = useState<TransportResponse[]>([]);
    const [detailTransport, setDetailTransport] = useState<TransportResponse | null>(null);
    const [showEditTransportModal, setShowEditTransportModal] = useState(false);
    const [locations, setLocations] = useState<LocationResponse[]>([]);
    const [activities, setActivities] = useState<any[]>([]);
    const [detailLocation, setDetailLocation] = useState<LocationResponse | null>(null);

    const activeLocation = locations.find((l) => l.id === selectedLocationId);
    const selectedLocationActivities = activities.filter(
        (a) => a.locationId === selectedLocationId
    );

    const [likeInfo, setLikeInfo] = useState<TripLikeResponse>({ likeCount: 0, hasLiked: false });
    const [comments, setComments] = useState<TripCommentResponse[]>([]);
    const [showComments, setShowComments] = useState(false);
    const [newComment, setNewComment] = useState("");
    const [isSubmittingComment, setIsSubmittingComment] = useState(false);

    useEffect(() => {
        getLocations(trip.id).then(setLocations).catch(console.error);
    }, [trip.id]);

    useEffect(() => {
        getTransports(trip.id).then(setTransports).catch(console.error);
    }, [trip.id]);
    // NEU
    useEffect(() => {
        getTripLikes(trip.id).then(setLikeInfo).catch(console.error);
    }, [trip.id, currentUser]);

    const handleAddLocation = async (newLocation: CreateLocationRequest) => {
        try {
            const created = await createLocation(trip.id, newLocation);
            setLocations([...locations, created]);
        } catch (error) {
            console.error("Fehler beim Erstellen der Location:", error);
        }
    };

    const handleLocationClick = (location: LocationResponse) => {
        setDetailLocation(location);
        setShowLocationDetailModal(true);
    };

    const handleUpdateLocation = async (req: UpdateLocationRequest) => {
        if (!detailLocation) return;
        try {
            const updated = await updateLocation(trip.id, detailLocation.id!, req);
            setLocations(locations.map((l) => l.id === updated.id ? updated : l));
            setShowLocationDetailModal(false);
        } catch (error) {
            console.error("Fehler beim Aktualisieren der Location:", error);
        }
    };

    const handleDeleteLocation = async () => {
        if (!detailLocation) return;
        try {
            await deleteLocation(trip.id, detailLocation.id!);
            setLocations(locations.filter((l) => l.id !== detailLocation.id));
            setShowLocationDetailModal(false);
        } catch (error) {
            console.error("Fehler beim Löschen der Location:", error);
        }
    };

    const handleAddActivity = (newActivity: any) => {
        const activity = {
            id: `act-${Date.now()}`,
            locationId: selectedLocationId!,
            ...newActivity,
        };
        setActivities([...activities, activity]);
    };

    const handleAddTransport = async (newTransport: CreateTransportRequest) => {
        try {
            const created = await createTransport(trip.id, newTransport);
            setTransports([...transports, created]);
        } catch (error) {
            console.error("Fehler beim Erstellen des Transports:", error);
        }
    };

    const handleUpdateTransport = async (req: UpdateTransportRequest) => {
        if (!detailTransport) return;
        try {
            const updated = await updateTransport(trip.id, detailTransport.id!, req);
            setTransports(transports.map((t) => t.id === updated.id ? updated : t));
            setShowEditTransportModal(false);
        } catch (error) {
            console.error("Fehler beim Aktualisieren des Transports:", error);
        }
    };

    const handleDeleteTransport = async () => {
        if (!detailTransport) return;
        try {
            await deleteTransport(trip.id, detailTransport.id!);
            setTransports(transports.filter((t) => t.id !== detailTransport.id));
            setShowEditTransportModal(false);
        } catch (error) {
            console.error("Fehler beim Löschen des Transports:", error);
        }
    };

    const handleEditTrip = async (updatedTrip: Partial<TripResponse>) => {
        try {
            const updated = await updateTrip(trip.id, updatedTrip);
            setCurrentTrip(updated);   // ← lokaler State
            onTripUpdate(updated);     // ← Parent informieren
            setIsEditingTrip(false);
        } catch (error) {
            console.error("Fehler beim Bearbeiten der Reise:", error);
        }
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
        } catch (error) {
            console.error("Fehler beim Liken:", error);
        }
    };

    const handleShowComments = async () => {
        if (!showComments && comments.length === 0) {
            try {
                const data = await getTripComments(trip.id);
                setComments(data.data ?? []);
            } catch (error) {
                console.error("Fehler beim Laden der Kommentare:", error);
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
            console.error("Fehler beim Erstellen des Kommentars:", error);
        } finally {
            setIsSubmittingComment(false);
        }
    };

    const handleDeleteComment = async (commentId: string) => {
        try {
            await deleteTripComment(trip.id, commentId);
            setComments(comments.filter((c) => c.id !== commentId));
        } catch (error) {
            console.error("Fehler beim Löschen des Kommentars:", error);
        }
    };

    return (
        <div className="max-w-5xl px-6 py-12">
            <Link
                href="/"
                className="inline-flex items-center gap-2 text-sm text-zinc-500 dark:text-zinc-400 hover:text-sky-600 dark:hover:text-sky-400 transition-colors mb-8"
            >
                ← Zurück zur Übersicht
            </Link>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 space-y-6">
                    {/* Trip Header */}
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

                    {/* Social: Likes & Comments */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6">
                        <div className="flex items-center gap-4">
                            {/* Like Button */}
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

                            {/* Comments Toggle */}
                            <button
                                onClick={handleShowComments}
                                className="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium bg-zinc-50 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700 transition-colors"
                            >
                                <span>💬</span>
                                <span>{showComments ? "Kommentare ausblenden" : "Kommentare anzeigen"}</span>
                            </button>
                        </div>

                        {/* Comments Section */}
                        {showComments && (
                            <div className="mt-6 space-y-4">
                                {comments.length === 0 ? (
                                    <p className="text-zinc-500 dark:text-zinc-400 text-sm text-center py-4">
                                        Noch keine Kommentare
                                    </p>
                                ) : (
                                    <div className="space-y-3">
                                        {comments.map((comment) => (
                                            <div
                                                key={comment.id}
                                                className="flex items-start justify-between gap-3 p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl"
                                            >
                                                <div>
                                                    <p className="text-sm font-medium text-zinc-900 dark:text-white">
                                                        {comment.user.name}
                                                    </p>
                                                    <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">
                                                        {comment.text}
                                                    </p>
                                                </div>
                                                {currentUser && comment.user.id === currentUser.id && (
                                                    <button
                                                        onClick={() => handleDeleteComment(comment.id)}
                                                        className="p-1.5 rounded-lg text-zinc-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-blue-950/30 transition-colors shrink-0"
                                                    >
                                                        🗑️
                                                    </button>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                )}

                                {/* New Comment Input */}
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

                    {/* Locations List */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-center justify-between mb-6">
                            <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                                Orte ({locations.length})
                            </h2>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddLocationModal(true)}
                                    className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                                >
                                    + Ort hinzufügen
                                </button>
                            )}
                        </div>
                        {locations.length === 0 ? (
                            <p className="text-zinc-500 dark:text-zinc-400 text-center py-8">
                                Keine Orte hinzugefügt
                            </p>
                        ) : (
                            <div className="space-y-2">
                                {locations.map((location) => (
                                    <div
                                        key={location.id}
                                        onClick={() => setSelectedLocationId(
                                            selectedLocationId === location.id ? null : location.id!
                                        )}
                                        className={`w-full text-left p-4 rounded-xl border-2 transition-colors cursor-pointer ${selectedLocationId === location.id
                                            ? "bg-sky-50 dark:bg-sky-950/30 border-sky-300 dark:border-sky-700"
                                            : "bg-zinc-50 dark:bg-zinc-800/50 border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700"
                                            }`}
                                    >
                                        <div className="flex items-center justify-between">
                                            <div>
                                                <p className="font-medium text-zinc-900 dark:text-white">{location.name}</p>
                                                <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                                    {location.city}, {location.country} · {location.dateFrom} – {location.dateTo}
                                                </p>
                                            </div>
                                            {isEditable && (
                                                <button
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        handleLocationClick(location);
                                                    }}
                                                    className="p-2 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                                >
                                                    ✏️
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Travel Plan */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-center justify-between mb-6">
                            <h2 className="text-lg font-bold text-zinc-900 dark:text-white">Travel Plan</h2>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddTransportModal(true)}
                                    className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                                >
                                    + Transport
                                </button>
                            )}
                        </div>
                        {transports.length === 0 ? (
                            <p className="text-zinc-500 dark:text-zinc-400 text-center py-8">
                                Kein Transport hinzugefügt
                            </p>
                        ) : (
                            <div className="space-y-2">
                                {transports.map((t) => {
                                    const fromLocation = locations.find((l) => l.id === t.fromLocationId);
                                    const toLocation = locations.find((l) => l.id === t.toLocationId);
                                    const typeEmoji = { flight: "✈️", train: "🚂", car: "🚗", bus: "🚌" }[t.type ?? "flight"] ?? "🚗";
                                    return (
                                        <div
                                            key={t.id}
                                            className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700"
                                        >
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-3">
                                                    <span className="text-2xl">{typeEmoji}</span>
                                                    <div>
                                                        <p className="font-medium text-zinc-900 dark:text-white">
                                                            {fromLocation?.name ?? t.fromLocationId} → {toLocation?.name ?? t.toLocationId}
                                                        </p>
                                                        <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                                            {t.departureTime ?? "Keine Abfahrtszeit"}
                                                        </p>
                                                    </div>
                                                </div>
                                                {isEditable && (
                                                    <button
                                                        onClick={() => { setDetailTransport(t); setShowEditTransportModal(true); }}
                                                        className="p-2 text-zinc-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                                    >
                                                        ✏️
                                                    </button>
                                                )}
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>
                </div>

                {/* Right: Activities */}
                {activeLocation && (
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 h-fit">
                        <div className="flex items-center justify-between mb-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1">
                                    Aktivitäten in
                                </p>
                                <h3 className="text-lg font-bold text-zinc-900 dark:text-white">
                                    {activeLocation.name}
                                </h3>
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
                            <p className="text-zinc-500 dark:text-zinc-400 text-sm text-center py-4">
                                Keine Aktivitäten
                            </p>
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
            </div>

            {/* Modals */}
            <AddLocationModal
                isOpen={showAddLocationModal}
                onCloseAction={() => setShowAddLocationModal(false)}
                onAddAction={handleAddLocation}
            />
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
                locations={locations}
                onCloseAction={() => setShowAddTransportModal(false)}
                onAddAction={handleAddTransport}
            />
            {detailLocation && (
                <LocationDetailModal
                    isOpen={showLocationDetailModal}
                    location={detailLocation}
                    onCloseAction={() => setShowLocationDetailModal(false)}
                    onSaveAction={handleUpdateLocation}
                    onDeleteAction={handleDeleteLocation}
                />
            )}
            {detailTransport && (
                <EditTransportModal
                    isOpen={showEditTransportModal}
                    transport={detailTransport}
                    locations={locations}
                    onCloseAction={() => setShowEditTransportModal(false)}
                    onSaveAction={handleUpdateTransport}
                    onDeleteAction={handleDeleteTransport}
                />
            )}
        </div>
    );
}