"use client"
import { useRouter } from "next/navigation";
import TripForm from "@/components/trips/TripForm";
import { createTrip } from "@/lib/api/trips";
import { components } from "@/generated/types";

type CreateTripRequest = components["schemas"]["CreateTripRequest"];

export default function NewTripPage() {
  const router = useRouter();

  async function handleCreateTrip(createTripRequest: CreateTripRequest) {
    try {
      await createTrip(createTripRequest);
      router.push("/");
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <TripForm onCreateTripAction={handleCreateTrip} />
  );
}