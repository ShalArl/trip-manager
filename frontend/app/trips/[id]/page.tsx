import { mockTrips } from "@/lib/mock-trips";
import TripDetail from "@/components/trips/TripDetail";

export default async function TripDetailPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;

    const trip = mockTrips.find((t) => t.title === decodeURIComponent(id));

    if (!trip) {
        return <div>Reise nicht gefunden</div>
    }

    return <TripDetail trip={trip} />;
}