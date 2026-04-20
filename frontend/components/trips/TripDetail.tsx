import Link from "next/link";
import { useState } from "react";
import { components } from "@/generated/types";
import { TransportResponse, CreateTransportRequest } from "@/types/transport";
import { createTransport } from "@/lib/api/transport";
import AddLocationModal from "./modals/AddLocationModal";
import AddActivityModal from "./modals/AddActivityModal";
import EditTripModal from "./modals/EditTripModal";
import AddTransportModal from "./modals/AddTransportModal";


type TripResponse = components["schemas"]["TripResponse"];

type Props = {
    trip: TripResponse;
    isEditable?: boolean;
};

export default function TripDetail({ trip, isEditable = false }: Props) {
    const [isEditingTrip, setIsEditingTrip] = useState(false);
    const [selectedLocationId, setSelectedLocationId] = useState<string | null>(null);

    const [showAddLocationModal, setShowAddLocationModal] = useState(false);
    const [showAddActivityModal, setShowAddActivityModal] = useState(false);

    const [showAddTransportModal, setShowAddTransportModal] = useState(false);
    const [showAddAccommodationModal, setShowAddAccommodationModal] = useState(false);

    const [transports, setTransports] = useState<TransportResponse[]>([]);
    const [accommodations, setAccommodations] = useState([]);

    // TODO: Mock locations and activities - replace with API calls
    const [locations, setLocations] = useState([
        {
            id: "loc-1",
            name: "Paris",
            city: "Paris",
            country: "France",
            sequence: 1,
        },
        {
            id: "loc-2",
            name: "Lyon",
            city: "Lyon",
            country: "France",
            sequence: 2,
        },
    ]);

    const [activities, setActivities] = useState([
        {
            id: "act-1",
            name: "Eiffelturm",
            locationId: "loc-1",
            date: trip.startDate,
            category: "sightseeing",
        },
        {
            id: "act-2",
            name: "Restaurant",
            locationId: "loc-1",
            date: trip.startDate,
            category: "dining",
        },
    ]);

    const selectedLocation = locations.find((l) => l.id === selectedLocationId);
    const selectedLocationActivities = activities.filter(
        (a) => a.locationId === selectedLocationId
    );

    // Modal handlers
    const handleAddLocation = (newLocation: any) => {
        const location = {
            id: `loc-${Date.now()}`,
            ...newLocation,
            sequence: locations.length + 1,
        };
        setLocations([...locations, location]);
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

    const handleEditTrip = (updatedTrip: any) => {
        // TODO: Call API to update trip
        console.log("Updated trip:", updatedTrip);
        setIsEditingTrip(false);
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
                {/* Left: Trip Info & Locations */}
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
                                        {trip.title}
                                    </h1>
                                    <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                        {trip.startDate} · {trip.endDate}
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

                        {/* Trip Description */}
                        <div className="border-t border-zinc-100 dark:border-zinc-800 pt-6 space-y-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-2">
                                    Kurzbeschreibung
                                </p>
                                <p className="text-zinc-700 dark:text-zinc-300">
                                    {trip.shortDescription}
                                </p>
                            </div>
                            {trip.description && (
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-2">
                                        Details
                                    </p>
                                    <p className="text-zinc-700 dark:text-zinc-300 leading-relaxed">
                                        {trip.description}
                                    </p>
                                </div>
                            )}
                        </div>
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
                                    <button
                                        key={location.id}
                                        onClick={() =>
                                            setSelectedLocationId(
                                                selectedLocationId === location.id ? null : location.id
                                            )
                                        }
                                        className={`w-full text-left p-4 rounded-xl border-2 transition-colors ${selectedLocationId === location.id
                                            ? "bg-sky-50 dark:bg-sky-950/30 border-sky-300 dark:border-sky-700"
                                            : "bg-zinc-50 dark:bg-zinc-800/50 border-zinc-200 dark:border-zinc-700 hover:border-zinc-300 dark:hover:border-zinc-600"
                                            }`}
                                    >
                                        <p className="font-medium text-zinc-900 dark:text-white">
                                            {location.name}
                                        </p>
                                        <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                            {location.city}, {location.country}
                                        </p>
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                    {/* Travel Plan */}
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-8">
                        <div className="flex items-center justify-between mb-6">
                            <h2 className="text-lg font-bold text-zinc-900 dark:text-white">
                                Travel Plan
                            </h2>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddTransportModal(true)}
                                    className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-lg transition-colors"
                                >
                                    + Transport
                                </button>
                            )}
                        </div>
                        {/* Liste der Einträge */}
                        {transports.length === 0 ? (
                            <p className="text-zinc-500 dark:text-zinc-400 text-center py-8">
                                Kein Transport hinzugefügt
                            </p>
                        ) : (
                            <div className="space-y-2">
                                {transports.map((t: any) => (
                                    <div key={t.id} className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-200 dark:border-zinc-700">
                                        <p className="font-medium text-zinc-900 dark:text-white">{t.from} → {t.to}</p>
                                        <p className="text-sm text-zinc-500 dark:text-zinc-400">{t.type} · {t.date}</p>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>

                {/* Right: Activities (for selected location) */}
                {selectedLocation && (
                    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 h-fit">
                        <div className="flex items-center justify-between mb-4">
                            <div>
                                <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1">
                                    Aktivitäten in
                                </p>
                                <h3 className="text-lg font-bold text-zinc-900 dark:text-white">
                                    {selectedLocation.name}
                                </h3>
                            </div>
                            {isEditable && (
                                <button
                                    onClick={() => setShowAddActivityModal(true)}
                                    className="p-2 text-sky-600 dark:text-sky-400 hover:bg-sky-50 dark:hover:bg-sky-950/30 rounded-lg transition-colors"
                                    title="Aktivität hinzufügen"
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
                                    <div
                                        key={activity.id}
                                        className="p-3 bg-zinc-50 dark:bg-zinc-800 rounded-lg"
                                    >
                                        <p className="font-medium text-sm text-zinc-900 dark:text-white">
                                            {activity.name}
                                        </p>
                                        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                                            {activity.category}
                                        </p>
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
                onClose={() => setShowAddLocationModal(false)}
                onAdd={handleAddLocation}
            />
            <AddActivityModal
                isOpen={showAddActivityModal}
                locationId={selectedLocationId}
                locationName={selectedLocation?.name || ""}
                tripStartDate={trip.startDate}
                onClose={() => setShowAddActivityModal(false)}
                onAdd={handleAddActivity}
            />
            <EditTripModal
                isOpen={isEditingTrip}
                trip={trip}
                onClose={() => setIsEditingTrip(false)}
                onSave={handleEditTrip}
            />
            <AddTransportModal
                isOpen={showAddTransportModal}
                locations={locations}
                onCloseAction={() => setShowAddTransportModal(false)}
                onAddAction={handleAddTransport}
            />
        </div>
    );
}
