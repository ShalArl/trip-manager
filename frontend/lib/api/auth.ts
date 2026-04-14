import {AuthResponse, CreateUserRequest, LoginRequest, UserResponse, UpdateUserRequest, ChangePasswordRequest} from "@/types/user";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

// Helper function to get auth token from localStorage
function getAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
}

// Helper function to get auth headers
export function getAuthHeaders(): HeadersInit {
  const token = getAuthToken();
  return {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
  };
}


export async function register(createUserRequest: CreateUserRequest) {
    const response = await fetch(`${API_URL}/api/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(createUserRequest),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Erstellen des Users (${response.status}):`, errorData);
        throw new Error(`Fehler beim Registrieren: ${response.status}`);
    }

    const data = await response.json();
    return data as AuthResponse;
}

export async function login(loginRequest: LoginRequest) {
    console.log("Login Request:", JSON.stringify(loginRequest));
    const response = await fetch(`${API_URL}/api/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(loginRequest),
    });
    
    console.log("Login Response Status:", response.status);
    
    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Einloggen (${response.status}):`, errorData);
        throw new Error(`Fehler beim Einloggen: ${response.status}`);
    }

    const data = await response.json();
    console.log("Login erfolgreich:", data);
    return data as AuthResponse;
}

export async function updateMe(data: UpdateUserRequest) {
    const response = await fetch(`${API_URL}/api/users/me`, {
        method: "PUT",
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Aktualisieren des Profils (${response.status}):`, errorData);
        throw new Error(`Fehler beim Aktualisieren: ${response.status}`);
    }

    const userData = await response.json();
    return userData as UserResponse;
}

export async function changePassword(data: ChangePasswordRequest) {
    const response = await fetch(`${API_URL}/api/users/me/password`, {
        method: "PUT",
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Fehler beim Ändern des Passworts (${response.status}):`, errorData);
        throw new Error(`Fehler beim Ändern des Passworts: ${response.status}`);
    }

    return true;
}

export async function getMe() {
    const req = {
        method: "GET",
        headers: getAuthHeaders(),
    }

    console.log("Req:", req);

    const response = await fetch(`${API_URL}/api/users/me`, req);

    if (!response.ok) {
        const errorData = await response.text();
        // console.error(`Fehler beim Abrufen des Profils (${response.status}):`, errorData);
        // throw new Error(`Fehler beim Abrufen des Profils: ${response.status}`);
        // Do nothing
    }

    const userData = await response.json();
    return userData as UserResponse;
}

