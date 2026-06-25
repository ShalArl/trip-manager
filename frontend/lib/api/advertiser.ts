import { getAuthHeaders } from "./auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export type Advertiser = {
    id: string;
    email: string;
    name: string;
    firebaseUid: string;
    tenants: string[];
    createdAt: string;
};

export async function listAdvertisers(): Promise<Advertiser[]> {
    const headers = await getAuthHeaders();
    const res = await fetch(`${API_URL}/api/users/advertisers`, { headers });
    if (!res.ok) return [];
    return res.json();
}

export async function createAdvertiser(data: { firebaseUid: string; email: string; name: string }): Promise<Advertiser> {
    const headers = await getAuthHeaders();
    const res = await fetch(`${API_URL}/api/users/advertisers`, {
        method: "POST",
        headers,
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to create advertiser");
    return res.json();
}

export async function assignTenant(advertiserID: string, tenantId: string): Promise<void> {
    const headers = await getAuthHeaders();
    await fetch(`${API_URL}/api/users/advertisers/${advertiserID}/tenants`, {
        method: "POST",
        headers,
        body: JSON.stringify({ tenantId }),
    });
}

export async function removeTenant(advertiserID: string, tenantId: string): Promise<void> {
    const headers = await getAuthHeaders();
    await fetch(`${API_URL}/api/users/advertisers/${advertiserID}/tenants/${tenantId}`, {
        method: "DELETE",
        headers,
    });
}

export async function getAdvertiserMe(): Promise<Advertiser | null> {
    const headers = await getAuthHeaders();
    const res = await fetch(`${API_URL}/api/users/advertisers/me`, { headers });
    if (!res.ok) return null;
    return res.json();
}