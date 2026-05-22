"use client";

import { useState, useEffect } from "react";
import { likeTrip, unlikeTrip, getTripComments, createTripComment, deleteTripComment } from "@/lib/api/social";
import { UserResponse } from "@/types/user";
import { TripLikeResponse, TripCommentResponse } from "@/types/social";

// ─── Types ────────────────────────────────────────────────────────────────────

interface Props {
    tripId: string;
    currentUser?: UserResponse | null;
    initialLikeInfo: TripLikeResponse;
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function TripSocialSection({ tripId, currentUser, initialLikeInfo }: Props) {
    const [likeInfo, setLikeInfo] = useState<TripLikeResponse>(initialLikeInfo);
    const [comments, setComments] = useState<TripCommentResponse[]>([]);
    const [showComments, setShowComments] = useState(false);
    const [newComment, setNewComment] = useState("");
    const [isSubmittingComment, setIsSubmittingComment] = useState(false);

    // ── Handlers ─────────────────────────────────────────────────────────────

    useEffect(() => {
        setLikeInfo(initialLikeInfo);
    }, [initialLikeInfo]);

    const handleLike = async () => {
        if (!currentUser) return;
        try {
            if (likeInfo.hasLiked) {
                await unlikeTrip(tripId);
                setLikeInfo({ likeCount: likeInfo.likeCount - 1, hasLiked: false });
            } else {
                await likeTrip(tripId);
                setLikeInfo({ likeCount: likeInfo.likeCount + 1, hasLiked: true });
            }
        } catch (err) {
            console.error("[TripSocialSection] like:", err);
        }
    };

    const handleShowComments = async () => {
        if (!showComments && comments.length === 0) {
            try {
                const data = await getTripComments(tripId);
                setComments(data.data ?? []);
            } catch (err) {
                console.error("[TripSocialSection] getTripComments:", err);
            }
        }
        setShowComments(!showComments);
    };

    const handleSubmitComment = async () => {
        if (!newComment.trim() || !currentUser) return;
        setIsSubmittingComment(true);
        try {
            const created = await createTripComment(tripId, newComment.trim());
            setComments([...comments, created]);
            setNewComment("");
        } catch (err) {
            console.error("[TripSocialSection] createTripComment:", err);
        } finally {
            setIsSubmittingComment(false);
        }
    };

    const handleDeleteComment = async (commentId: string) => {
        try {
            await deleteTripComment(tripId, commentId);
            setComments(comments.filter((c) => c.id !== commentId));
        } catch (err) {
            console.error("[TripSocialSection] deleteTripComment:", err);
        }
    };

    // ── Render ────────────────────────────────────────────────────────────────

    return (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6">
            {/* Like & Comments buttons */}
            <div className="flex items-center gap-4">
                <button
                    onClick={handleLike}
                    disabled={!currentUser}
                    className={`flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-colors ${likeInfo.hasLiked
                            ? "bg-sky-100 dark:bg-sky-950/50 text-sky-600 dark:text-sky-400 border border-sky-200 dark:border-sky-800"
                            : "bg-zinc-50 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700"
                        } disabled:opacity-50 disabled:cursor-not-allowed`}
                >
                    <span>{likeInfo.hasLiked ? "❤️" : "🤍"}</span>
                    <span>{likeInfo.likeCount} {likeInfo.likeCount === 1 ? "Like" : "Likes"}</span>
                </button>
                <button
                    onClick={handleShowComments}
                    className="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium bg-zinc-50 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 hover:border-sky-300 dark:hover:border-sky-700 transition-colors"
                >
                    <span>💬</span>
                    <span>{showComments ? "Kommentare ausblenden" : "Kommentare anzeigen"}</span>
                </button>
            </div>

            {/* Comments */}
            {showComments && (
                <div className="mt-6 space-y-4">
                    {comments.length === 0 ? (
                        <p className="text-zinc-500 dark:text-zinc-400 text-sm text-center py-4">
                            Noch keine Kommentare
                        </p>
                    ) : (
                        <div className="space-y-3">
                            {comments.map((comment) => (
                                <div key={comment.id} className="flex items-start justify-between gap-3 p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl">
                                    <div className="flex items-start gap-3">
                                        <div className="w-8 h-8 rounded-full bg-sky-500 flex items-center justify-center text-white text-sm font-semibold shrink-0 overflow-hidden">
                                            {comment.user.avatarUrl ? (
                                                <img src={comment.user.avatarUrl} alt={comment.user.name} className="w-full h-full object-cover" />
                                            ) : (
                                                comment.user.name.charAt(0).toUpperCase()
                                            )}
                                        </div>
                                        <div>
                                            <p className="text-sm font-medium text-zinc-900 dark:text-white">{comment.user.name}</p>
                                            <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">{comment.text}</p>
                                        </div>
                                    </div>
                                    {currentUser && comment.user.id === currentUser.id && (
                                        <button
                                            onClick={() => handleDeleteComment(comment.id)}
                                            className="p-1.5 rounded-lg text-zinc-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/30 transition-colors shrink-0"
                                        >
                                            🗑️
                                        </button>
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                    {currentUser && (
                        <div className="flex gap-2 pt-2 border-t border-zinc-100 dark:border-zinc-800">
                            <input
                                type="text"
                                value={newComment}
                                onChange={(e) => setNewComment(e.target.value)}
                                onKeyDown={(e) => e.key === "Enter" && handleSubmitComment()}
                                placeholder="Kommentar schreiben..."
                                className="flex-1 px-4 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl focus:outline-none focus:border-sky-400 dark:focus:border-sky-600 text-zinc-900 dark:text-white placeholder-zinc-400"
                            />
                            <button
                                onClick={handleSubmitComment}
                                disabled={!newComment.trim() || isSubmittingComment}
                                className="px-4 py-2 text-sm font-medium bg-sky-600 hover:bg-sky-700 text-white rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {isSubmittingComment ? "..." : "Senden"}
                            </button>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}