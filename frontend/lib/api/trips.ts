import { components } from "@/generated/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

type CreateTripRequest = components["schemas"]["CreateTripRequest"];
type TripResponse = components["schemas"]["TripResponse"];

export async function createTrip(createTripRequest: CreateTripRequest): Promise<TripResponse> {
  const response = await fetch(`${API_URL}/api/trips`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${localStorage.getItem("token")}`,
    },
    body: JSON.stringify(createTripRequest),
  });

  if (!response.ok) {
    throw new Error("Fehler beim Erstellen der Reise");
  }

  const data = await response.json();
  return data as TripResponse;
}