import { components } from "@/generated/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

type CreateUserRequest = components["schemas"]["CreateUserRequest"];
type AuthResponse = components["schemas"]["AuthResponse"];
type LoginRequest = components["schemas"]["LoginRequest"];

export async function createUser(createUserRequest: CreateUserRequest) {
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