import { components } from "@/generated/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

type CreateTripRequest = components["schemas"]["CreateTripRequest"];
type TripResponse = components["schemas"]["TripResponse"];

export async function createTrip(createTripRequest: CreateTripRequest) {
  const response = await fetch(`${API_URL}/api/trips`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${localStorage.getItem("token")}`,
    },
    body: JSON.stringify(createTripRequest),
  });

  if (!response.ok) {
    const errorData = await response.json()
    console.log("Fehler Details:", errorData);
    throw new Error("Fehler beim Erstellen der Reise");
  }

  const data = await response.json();
  return data as TripResponse;
}

export async function getTrips(): Promise<TripResponse[]> {
  const token = localStorage.getItem("token");
  console.log("Token vorhanden:", !!token);
  console.log("API_URL:", API_URL);

  const response = await fetch(`${API_URL}/api/trips`, {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    const errorData = await response.text();
    console.error(`Fehler beim Laden der Reisen (${response.status}):`, errorData);
    throw new Error(`Fehler beim Laden der Reisen: ${response.status} ${response.statusText}`);
  }

  const data = await response.json();
  console.log("Reisen werden aus der Datenbank geholt:", data)
  return data.data as TripResponse[];
}

export async function getTrip(tripId: string): Promise<TripResponse> {
  const token = localStorage.getItem("token");
  console.log("Token vorhanden:", !!token);
  console.log("API_URL:", API_URL);

  const response = await fetch(`${API_URL}/api/trips/${tripId}`, {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    const errorData = await response.text();
    console.error(`Fehler beim Laden der Reise (${response.status}):`, errorData);
    throw new Error(`Fehler beim Laden der Reise: ${response.status} ${response.statusText}`);
  }

  const data = await response.json();
  return data as TripResponse;
}
