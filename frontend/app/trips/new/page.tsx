"use client"
import { useRouter } from "next/navigation";
import TripForm from "@/components/trips/TripForm";
import { createTrip } from "@/lib/api/trips";
import {CreateTripRequest, TripResponse} from "@/types/trip";


export default function NewTripPage() {
  const router = useRouter();

  async function handleCreateTrip(createTripRequest: CreateTripRequest) {
    try {
      const trip: TripResponse = await createTrip(createTripRequest);
      console.log("Trip response: ", trip);
      router.push(`/trips/${trip.id}`);
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <TripForm onCreateTripAction={handleCreateTrip} />
  );
}