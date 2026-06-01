import { getAuthHeaders } from "./auth";

export interface NewsletterTrip {
  tripId: string;
  title: string;
  description?: string;
  destination?: string;
  coverImageUrl?: string;
  creatorId: string;
  creatorName: string;
  likeCount?: number;
  commentCount?: number;
  relevanceReason: "creator_you_follow" | "liked_by_similar_users" | "trending_in_network";
  createdAt: string;
}

export interface NewsletterSection {
  title: string;
  description?: string;
  trips: NewsletterTrip[];
}

export interface NewsletterResponse {
  sections: NewsletterSection[];
  generatedAt: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getNewsletter(limit = 10): Promise<NewsletterResponse> {
  const response = await fetch(
    `${API_URL}/api/newsletter?limit=${limit}`,
    {
      method: "GET",
      headers: await getAuthHeaders(),
    }
  );
  
  if (!response.ok) {
    throw new Error(`Fehler beim Laden des Newsletters: ${response.status}`);
  }
  return response.json();
}