import { components } from "@/generated/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

type CreateUserRequest = components["schemas"]["CreateUserRequest"];

export async function createUser(createUserRequest: CreateUserRequest) {
    console.log("createUser Request: " + JSON.stringify(createUserRequest));

    const response = await fetch(`${API_URL}/api/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(createUserRequest),
    });

    console.log("createUser Response: " + response);

    /*if (!response.ok) {
      throw new Error("Fehler beim Erstellen des Users");
    }*/

    return response.json();
}