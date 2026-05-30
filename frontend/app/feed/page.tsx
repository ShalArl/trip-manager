"use client";

import React, { useEffect, useState } from "react";
import { getFeed, getPersonalFeed } from "@/lib/api/feed";
import { components } from "@/generated/types";
import Link from "next/link";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { useUserContext } from "@/lib/context/UserContext";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  Heart,
  MessageCircle,
  TrendingUp,
  User,
} from "lucide-react";

type FeedTrip = components["schemas"]["FeedTrip"];

type FeedMode = "global" | "personal";

const PAGE_SIZE = 20;

export default function FeedPage() {
  const { user } = useUserContext();
  const [mode, setMode] = useState<FeedMode>("global");
  const [trips, setTrips] = useState<FeedTrip[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchFeed = async () => {
      setIsLoading(true);
      try {
        const result =
          mode === "personal" && user
            ? await getPersonalFeed(PAGE_SIZE, page * PAGE_SIZE)
            : await getFeed(PAGE_SIZE, page * PAGE_SIZE);
        setTrips(result.data);
        setTotal(result.total);
      } catch (error) {
        console.error(error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchFeed();
  }, [page, mode, user]);

  const handleModeChange = (newMode: FeedMode) => {
    setMode(newMode);
    setPage(0);
  };

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));
  const firstVisibleItem = total === 0 ? 0 : page * PAGE_SIZE + 1;
  const lastVisibleItem = Math.min((page + 1) * PAGE_SIZE, total);

  return (
    <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
      <div className="mx-auto max-w-4xl px-6 py-12">

        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <TrendingUp className="h-7 w-7 text-sky-500" />
            <h1 className="text-3xl font-bold tracking-tight">Feed</h1>
          </div>
          <p className="text-zinc-500 dark:text-zinc-400">
            {mode === "personal"
              ? "Reisen basierend auf deinen Interaktionen"
              : "Die beliebtesten Reisen der Community"}
          </p>
        </div>

        {/* Mode Toggle – nur für eingeloggte User */}
        {user && (
          <div className="flex items-center gap-2 mb-8 p-1 bg-zinc-100 dark:bg-zinc-900 rounded-xl w-fit">
            <button
              onClick={() => handleModeChange("global")}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                mode === "global"
                  ? "bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 shadow-sm"
                  : "text-zinc-500 dark:text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-300"
              }`}
            >
              <TrendingUp className="h-4 w-4" />
              Global
            </button>
            <button
              onClick={() => handleModeChange("personal")}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                mode === "personal"
                  ? "bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 shadow-sm"
                  : "text-zinc-500 dark:text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-300"
              }`}
            >
              <User className="h-4 w-4" />
              Für dich
            </button>
          </div>
        )}

        {/* Pagination Header */}
        {total > 0 && (
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mb-8">
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              {firstVisibleItem}–{lastVisibleItem} von {total}
            </p>
            <div className="flex items-center gap-1">
              <Button variant="outline" size="icon" onClick={() => setPage(0)} disabled={page === 0} aria-label="Erste Seite">
                <ChevronsLeft className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" onClick={() => setPage(page - 1)} disabled={page === 0} aria-label="Vorherige Seite">
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <div className="px-4 text-sm font-medium min-w-[80px] text-center">
                {page + 1} / {totalPages}
              </div>
              <Button variant="outline" size="icon" onClick={() => setPage(page + 1)} disabled={page >= totalPages - 1} aria-label="Nächste Seite">
                <ChevronRight className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" onClick={() => setPage(totalPages - 1)} disabled={page >= totalPages - 1} aria-label="Letzte Seite">
                <ChevronsRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        )}

        {/* Feed */}
        {isLoading ? (
          <div className="flex flex-col gap-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-24 rounded-2xl" />
            ))}
          </div>
        ) : trips.length === 0 ? (
          <Card className="p-12 text-center border-dashed">
            {mode === "personal" ? (
              <>
                <User className="h-10 w-10 text-zinc-300 dark:text-zinc-700 mx-auto mb-3" />
                <p className="text-zinc-500 dark:text-zinc-400">
                  Noch keine personalisierten Empfehlungen
                </p>
                <p className="text-sm text-zinc-400 dark:text-zinc-600 mt-1">
                  Like oder kommentiere Reisen damit wir dir passende empfehlen können
                </p>
              </>
            ) : (
              <>
                <TrendingUp className="h-10 w-10 text-zinc-300 dark:text-zinc-700 mx-auto mb-3" />
                <p className="text-zinc-500 dark:text-zinc-400">
                  Noch keine Reisen im Feed
                </p>
                <p className="text-sm text-zinc-400 dark:text-zinc-600 mt-1">
                  Like oder kommentiere Reisen damit sie hier erscheinen
                </p>
              </>
            )}
          </Card>
        ) : (
          <div className="flex flex-col gap-3">
            {trips.map((trip, index) => (
              <Link
                key={trip.tripId}
                href={`/trips/${encodeURIComponent(trip.tripId)}`}
                className="group block"
              >
                <Card className="p-5 hover:border-sky-300 dark:hover:border-sky-700/50 hover:shadow-md transition-all cursor-pointer overflow-hidden">
                  <div className="flex items-center gap-5">

                    {/* Rank Badge */}
                    <div className="flex flex-col items-center justify-center flex-shrink-0 w-12 py-2 rounded-xl bg-sky-50 dark:bg-sky-950/40 border border-sky-100 dark:border-sky-900/50">
                      <p className="text-[10px] uppercase tracking-wider text-sky-600 dark:text-sky-400 font-semibold">
                        #
                      </p>
                      <p className="text-xl font-bold text-sky-900 dark:text-sky-200 leading-none mt-0.5">
                        {page * PAGE_SIZE + index + 1}
                      </p>
                    </div>

                    {/* Main Content */}
                    <div className="min-w-0 flex-1">
                      <h3 className="font-semibold text-base truncate group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors mb-2">
                        {trip.title}
                      </h3>

                      {/* Stats */}
                      <div className="flex items-center gap-4">
                        <div className="flex items-center gap-1.5 text-sm text-zinc-500 dark:text-zinc-400">
                          <Heart className="h-3.5 w-3.5 text-rose-400" />
                          <span>{trip.likes}</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-sm text-zinc-500 dark:text-zinc-400">
                          <MessageCircle className="h-3.5 w-3.5 text-sky-400" />
                          <span>{trip.comments}</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-sm text-zinc-500 dark:text-zinc-400">
                          <TrendingUp className="h-3.5 w-3.5 text-emerald-400" />
                          <span>{trip.score}</span>
                        </div>
                      </div>
                    </div>

                    {/* Chevron */}
                    <ChevronRight className="h-5 w-5 text-zinc-300 dark:text-zinc-700 group-hover:text-sky-500 group-hover:translate-x-1 transition-all flex-shrink-0" />
                  </div>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}