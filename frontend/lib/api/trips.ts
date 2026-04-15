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
  return data.data as TripResponse;
}

export async function getPublicTrips(): Promise<TripResponse[]> {
    // TODO: Replace with real API call
    const mockTrips: TripResponse[] = [
        {
            id: "mock-1",
            title: "Abenteuer in Japan",
            shortDescription: "Tokio, Kyoto und mehr",
            startDate: "2026-05-01",
            endDate: "2026-05-15",
            status: "planned",
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            createdBy: {
                id: "user-1",
                name: "Anna Müller",
                email: "anna@example.com",
            },
        },
        {
            id: "mock-2",
            title: "Roadtrip durch Norwegen",
            shortDescription: "Fjorde und Nordlichter",
            startDate: "2026-07-10",
            endDate: "2026-07-20",
            status: "planned",
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            createdBy: {
                id: "user-2",
                name: "Max Schmidt",
                email: "max@example.com",
            },
        },
        {
            id: "mock-3",
            title: "Städtetrip nach Barcelona",
            shortDescription: "Kultur, Strand und Tapas",
            startDate: "2026-06-05",
            endDate: "2026-06-10",
            status: "planned",
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            createdBy: {
                id: "user-3",
                name: "Lisa Weber",
                email: "lisa@example.com",
            },
        },
    ];

    return mockTrips;
}

export async function searchTrips(query: string): Promise<TripResponse[]> {
    // TODO: Replace with real API call
    const allTrips = await getPublicTrips();
    
    return allTrips.filter((trip) =>
        trip.title.toLowerCase().includes(query.toLowerCase()) ||
        trip.shortDescription.toLowerCase().includes(query.toLowerCase())
    );
}
