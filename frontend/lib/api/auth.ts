import {
    signInWithEmailAndPassword,
    createUserWithEmailAndPassword,
    signOut as firebaseSignOut,
    updatePassword as firebaseUpdatePassword,
    updateProfile,
    reauthenticateWithCredential,
    EmailAuthProvider,
} from "firebase/auth";
import { firebaseAuth } from "@/lib/api/firebase";
import {
    UserResponse,
    UpdateUserRequest,
    ChangePasswordRequest,
    ProvisionUserRequest,
    CreateUserRequest,
    LoginRequest,
} from "@/types/user";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

/**
 * Get middleware headers with a fresh Firebase ID token.
 * Always async: the token may need a refresh round-trip.
 */
export async function getAuthHeaders(): Promise<HeadersInit> {
    const user = firebaseAuth.currentUser;
    if (!user) {
        return { "Content-Type": "application/json" };
    }
    const token = await user.getIdToken();
    return {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
    };
}

/**
 * Register a new user with Firebase and provision the backend record.
 */
export async function register(req: CreateUserRequest, tenantId?: string): Promise<UserResponse> {
    const credential = await createUserWithEmailAndPassword(
        firebaseAuth,
        req.email,
        req.password,
    );

    if (req.name) {
        await updateProfile(credential.user, { displayName: req.name });
    }

    return provisionMe({ name: req.name }, tenantId);
}

/**
 * Log in with email and password.
 * Also ensures backend provisioning (idempotent) in case the user
 * was created out-of-band.
 */
export async function login(req: LoginRequest): Promise<UserResponse> {
    await signInWithEmailAndPassword(firebaseAuth, req.email, req.password);
    return provisionMe();
}

/**
 * Log out of Firebase. Clears SDK state.
 */
export async function logout(): Promise<void> {
    await firebaseSignOut(firebaseAuth);
}

/**
 * Create the backend user record after Firebase sign-up/sign-in.
 * Idempotent: returns existing record if already provisioned.
 */
async function provisionMe(data?: ProvisionUserRequest, tenantId?: string): Promise<UserResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/users/provision`, {
        method: "POST",
        headers,
        body: JSON.stringify({ ...data ?? {}, tenantId }),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`Provisioning failed (${response.status}):`, errorData);
        throw new Error(`Provisioning failed: ${response.status}`);
    }

    await firebaseAuth.currentUser?.getIdToken(true);
    return response.json();
}

export async function getMe(): Promise<UserResponse | null> {
    if (!firebaseAuth.currentUser) return null;

    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/users/me`, {
        method: "GET",
        headers,
    });

    if (!response.ok) {
        console.error(`getMe failed (${response.status})`);
        return null;
    }

    return response.json();
}

export async function updateMe(data: UpdateUserRequest): Promise<UserResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/users/me`, {
        method: "PUT",
        headers,
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.text();
        console.error(`updateMe failed (${response.status}):`, errorData);
        throw new Error(`Fehler beim Aktualisieren: ${response.status}`);
    }

    return response.json();
}

/**
 * Change password directly via Firebase. Requires recent authentication —
 * Firebase throws `middleware/requires-recent-login` if the session is too old,
 * in which case the caller must re-authenticate first.
 */
export async function changePassword(data: ChangePasswordRequest): Promise<void> {
    const user = firebaseAuth.currentUser;
    if (!user || !user.email) {
        throw new Error("Not authenticated");
    }

    const credential = EmailAuthProvider.credential(user.email, data.currentPassword);
    try {
        await reauthenticateWithCredential(user, credential);
    } catch (err: unknown) {
        // @ts-expect-error - Firebase Error has attribute code.
        if (err.code === "middleware/wrong-password" || err.code === "middleware/invalid-credential") {
            throw new Error("Aktuelles Passwort ist falsch");
        }
        throw err;
    }

    await firebaseUpdatePassword(user, data.newPassword);
}

// Tenant Related:

export interface RegisterTenantRequest {
    tenantName: string;
    tier: "free" | "standard";
}

export interface TenantResponse {
    tenantId: string;
    name: string;
    tier: string;
    slug: string;
}

export async function registerTenant(req: RegisterTenantRequest): Promise<TenantResponse> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/tenants/register`, {
        method: "POST",
        headers,
        body: JSON.stringify(req),
    });

    if (!response.ok) {
        await response.text();
        throw new Error(`Tenant registration failed: ${response.status}`);
    }

    // Token refreshen damit neue Claims aktiv werden
    await firebaseAuth.currentUser?.getIdToken(true);

    return response.json();
}

export async function getTenantBySlug(slug: string): Promise<TenantResponse | null> {
    const response = await fetch(`${API_URL}/api/tenants/by-slug/${slug}`);
    if (!response.ok) return null;
    return response.json();
}

export async function getTenant(): Promise<TenantResponse | null> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/tenants/me`, {
        method: "GET",
        headers,
    });

    if (!response.ok) return null;
    return response.json();
}

export interface TenantSettings {
    maxActiveTrips: number;
}

export async function getTenantSettings(): Promise<TenantSettings | null> {
    const headers = await getAuthHeaders();
    const response = await fetch(`${API_URL}/api/tenants/me/settings`, {
        method: "GET",
        headers,
    });
    if (!response.ok) return null;
    return response.json();
}

/**
 * Send a password reset email. Firebase handles the reset link + flow.
 */
export async function sendPasswordReset(email: string): Promise<void> {
    const { sendPasswordResetEmail } = await import("firebase/auth");
    await sendPasswordResetEmail(firebaseAuth, email);
}