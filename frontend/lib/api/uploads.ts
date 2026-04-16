import {getAuthHeaders} from "./auth";
import {FileUploadRequest, PresignedURLRequest, PresignedURLResponse} from "@/types/upload";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";
const S3_URL = process.env.NEXT_PUBLIC_S3_PUBLIC_URL || "http://localhost:9000";
/**
 * Get a presigned URL for direct file upload to S3/MinIO
 * @returns Presigned URL that can be used for PUT requests
 * @param req PresignedURLRequest
 */
export async function getPresignedUrl(req: PresignedURLRequest): Promise<PresignedURLResponse> {

    const response = await fetch(`${API_URL}/api/uploads/presigned`, {
        method: "POST",
        headers: getAuthHeaders(),
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
 * Upload a file directly to S3/MinIO using a presigned URL
 * @param fileUploadRequest of type FileUploadRequest
 */
export async function uploadToPresignedUrl(fileUploadRequest: FileUploadRequest): Promise<void> {
    console.log("[uploadToPresignedUrl] Starting upload to presigned URL...");

    const file = fileUploadRequest.file;
    const presignedUrl = fileUploadRequest.url;

    // Read file as binary
    const arrayBuffer = await file.arrayBuffer();
    // Workaround because caddy strips minio prefix leading to incorrect bucket resolution
    const uploadUrl = presignedUrl.replace(S3_URL, `${S3_URL}/minio/`);

    const response = await fetch(uploadUrl, {
        method: "PUT",
        headers: {
            "Content-Type": file.type || "application/octet-stream",
        },
        body: arrayBuffer,
    });

    if (!response.ok) {
        console.error(`Upload failed (${response.status}):`, await response.text());
        throw new Error(`Failed to upload file: ${response.status}`);
    }

    console.log("[uploadToPresignedUrl] File uploaded successfully");
}

/**
 * Upload an avatar directly to S3/MinIO
 * Combines getPresignedUrl and uploadToPresignedUrl
 * @param file - Avatar file to upload
 * @param _userId unused for now as userid is already part of the presignedUrl
 * @returns The public URL of the uploaded avatar
 */
export async function uploadAvatar(file: File, _userId: string): Promise<string> {
    // Get presigned URL
    console.log("[uploadAvatar] Getting presigned URL...");
    const req: PresignedURLRequest = {
        fileName: file.name,
        mediaType: "avatar"
    };
    const response = await getPresignedUrl(req);
    console.log("[uploadAvatar] Got presigned URL, expires in", response.expiresIn, "seconds");

    // Upload to S3/MinIO directly
    console.log("[uploadAvatar] Uploading to S3/MinIO...");
    const fileUploadRequest: FileUploadRequest = {
        url: response.presignedUrl,
        file: file,
    };
    await uploadToPresignedUrl(fileUploadRequest);

    // Extract the public URL from the presigned URL by removing query parameters
    // Presigned URL format: https://domain/minio/trip-manager/avatars/userId.ext?X-Amz-Algorithm=...
    // Public URL format:    https://domain/minio/trip-manager/avatars/userId.ext
    const baseImageUrl = response.presignedUrl.split('?')[0];
    const correctAvatarUrl = baseImageUrl.replace('travel-nugget.duckdns.org/', 'travel-nugget.duckdns.org/minio/');

    console.log("[uploadAvatar] Avatar uploaded successfully:", correctAvatarUrl);
    return correctAvatarUrl;
}

