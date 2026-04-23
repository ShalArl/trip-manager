import { getAuthHeaders } from "./auth";
import { FileUploadRequest, PresignedURLRequest, PresignedURLResponse } from "@/types/upload";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

/**
 * Get a presigned URL for direct file upload to the storage backend.
 * The storage backend (GCS in prod, MinIO locally) is opaque to the client.
 */
export async function getPresignedUrl(req: PresignedURLRequest): Promise<PresignedURLResponse> {
    const response = await fetch(`${API_URL}/api/uploads/presigned`, {
        method: "POST",
        headers: await getAuthHeaders(),
        body: JSON.stringify(req),
    });

    if (!response.ok) {
        const errorData = await response.json();
        console.error(`Failed to get presigned URL (${response.status}):`, errorData);
        throw new Error(errorData.error || `Failed to get presigned URL: ${response.status}`);
    }

    return response.json();
}

/**
 * Upload a file directly to the storage backend using a presigned PUT URL.
 * The URL is used verbatim — no rewriting.
 */
export async function uploadToPresignedUrl(req: FileUploadRequest): Promise<void> {
    const response = await fetch(req.url, {
        method: "PUT",
        headers: {
            "Content-Type": req.file.type || "application/octet-stream",
        },
        body: req.file,
    });

    if (!response.ok) {
        const errorText = await response.text();
        console.error(`Upload failed (${response.status}):`, errorText);
        throw new Error(`Failed to upload file: ${response.status}`);
    }
}

/**
 * Upload an avatar and return the storage key.
 * The caller is responsible for associating the key with a user via
 * PUT /api/users/me with { avatarKey }.
 */
export async function uploadAvatar(file: File): Promise<string> {
    const ticket = await getPresignedUrl({
        fileName: file.name,
        mediaType: "avatar",
    });

    await uploadToPresignedUrl({
        url: ticket.presignedUrl,
        file,
    });

    return ticket.key;
}