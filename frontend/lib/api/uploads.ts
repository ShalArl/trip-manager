import {getAuthHeaders} from "./auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";
const S3_PUBLIC_URL = process.env.NEXT_PUBLIC_S3_PUBLIC_URL || "http://localhost/minio";

/**
 * Get a presigned URL for direct file upload to S3/MinIO
 * @param fileName - Name of the file to upload
 * @param mediaType - Type of media (avatar, trip, location, activity)
 * @returns Presigned URL that can be used for PUT requests
 */
export async function getPresignedUrl(fileName: string, mediaType: string): Promise<{
  presignedUrl: string;
  expiresIn: number;
}> {

  const response = await fetch(`${API_URL}/api/uploads/presigned`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify({
      fileName,
      mediaType,
    }),
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
 * @param presignedUrl - The presigned URL from getPresignedUrl()
 * @param file - The file to upload
 */
export async function uploadToPresignedUrl(presignedUrl: string, file: File): Promise<void> {
  console.log("[uploadToPresignedUrl] Starting upload to presigned URL...");

  // Read file as binary
  const arrayBuffer = await file.arrayBuffer();

  const response = await fetch(presignedUrl, {
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
 * @param userId - User ID (used for the filename)
 * @returns The public URL of the uploaded avatar
 */
export async function uploadAvatar(file: File, userId: string): Promise<string> {
  // Get presigned URL
  console.log("[uploadAvatar] Getting presigned URL...");
  const { presignedUrl, expiresIn } = await getPresignedUrl(file.name, "avatar");
  console.log("[uploadAvatar] Got presigned URL, expires in", expiresIn, "seconds");

  // Upload to S3/MinIO directly
  console.log("[uploadAvatar] Uploading to S3/MinIO...");
  await uploadToPresignedUrl(presignedUrl, file);

   // Generate the public URL (same path as presigned URL but without query params)
   // The backend generates the path as: avatars/{userId}.{extension}
   const ext = file.name.substring(file.name.lastIndexOf("."));
   const publicUrl = `${S3_PUBLIC_URL}/avatars/${userId}${ext}`;

  console.log("[uploadAvatar] Avatar uploaded successfully:", publicUrl);
  return publicUrl;
}

