"use client";

import { useEffect, useState } from "react";
import { getNewsletter, NewsletterResponse, NewsletterTrip } from "@/lib/api/newsletter";
import { useUserContext } from "@/lib/context/UserContext";
import Link from "next/link";
import { Heart, MessageCircle, ArrowRight, Compass, Users, TrendingUp, Mail } from "lucide-react";

const SECTION_META: Record<string, { icon: typeof Compass; accent: string; bg: string }> = {
  "From Travellers You Follow": {
    icon: Compass,
    accent: "text-amber-500",
    bg: "bg-amber-50 dark:bg-amber-950/20 border-amber-100 dark:border-amber-900/30",
  },
  "Popular in Your Network": {
    icon: Users,
    accent: "text-sky-500",
    bg: "bg-sky-50 dark:bg-sky-950/20 border-sky-100 dark:border-sky-900/30",
  },
  "Trending Among Your Peers": {
    icon: TrendingUp,
    accent: "text-emerald-500",
    bg: "bg-emerald-50 dark:bg-emerald-950/20 border-emerald-100 dark:border-emerald-900/30",
  },
};

function TripCard({ trip, accent }: { trip: NewsletterTrip; accent: string }) {
  return (
    <Link href={`/trips/${trip.tripId}`} className="group block">
      <div className="flex items-start gap-4 p-4 rounded-2xl hover:bg-zinc-50 dark:hover:bg-zinc-800/40 transition-all duration-200">
        {/* Color dot */}
        <div className={`mt-1 h-2 w-2 rounded-full flex-shrink-0 ${accent.replace("text-", "bg-")}`} />

        <div className="flex-1 min-w-0">
          <p className="font-semibold text-zinc-900 dark:text-zinc-100 group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors truncate">
            {trip.title || "Untitled Trip"}
          </p>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 mt-0.5">
            {trip.creatorName || "Unknown traveller"}
          </p>
          <div className="flex items-center gap-3 mt-2">
            <span className="flex items-center gap-1 text-xs text-zinc-400">
              <Heart className="h-3 w-3 text-rose-400" />
              {trip.likeCount ?? 0}
            </span>
            <span className="flex items-center gap-1 text-xs text-zinc-400">
              <MessageCircle className="h-3 w-3 text-sky-400" />
              {trip.commentCount ?? 0}
            </span>
          </div>
        </div>

        <ArrowRight className="h-4 w-4 text-zinc-300 dark:text-zinc-600 group-hover:text-sky-500 group-hover:translate-x-0.5 transition-all mt-1 flex-shrink-0" />
      </div>
    </Link>
  );
}

export default function NewsletterPage() {
  const { user, isLoading: userLoading } = useUserContext();
  const [newsletter, setNewsletter] = useState<NewsletterResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) return;
    const fetch = async () => {
      setIsLoading(true);
      setError(null);
      try {
        setNewsletter(await getNewsletter());
      } catch {
        setError("Newsletter konnte nicht geladen werden.");
      } finally {
        setIsLoading(false);
      }
    };
    fetch();
  }, [user]);

  if (userLoading || isLoading) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 space-y-8">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="space-y-3 animate-pulse">
            <div className="h-4 w-32 bg-zinc-200 dark:bg-zinc-700 rounded" />
            <div className="h-20 bg-zinc-100 dark:bg-zinc-800 rounded-2xl" />
          </div>
        ))}
      </div>
    );
  }

  if (!user) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-32 text-center">
        <div className="inline-flex items-center justify-center h-14 w-14 rounded-2xl bg-zinc-100 dark:bg-zinc-800 mb-6">
          <Mail className="h-6 w-6 text-zinc-400" />
        </div>
        <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-100 mb-2">
          Dein persönlicher Newsletter
        </h2>
        <p className="text-zinc-500 dark:text-zinc-400 mb-6">
          Melde dich an um Reiseempfehlungen basierend auf deinen Interessen zu erhalten.
        </p>
        <Link
          href="/auth"
          className="inline-flex items-center gap-2 px-5 py-2.5 bg-sky-600 hover:bg-sky-700 text-white text-sm font-medium rounded-xl transition-colors"
        >
          Jetzt anmelden
          <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <p className="text-red-500 text-sm">{error}</p>
      </div>
    );
  }

  const totalTrips = newsletter?.sections.reduce((acc, s) => acc + s.trips.length, 0) ?? 0;
  const generatedAt = newsletter?.generatedAt
    ? new Date(newsletter.generatedAt).toLocaleDateString("de-DE", {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
      })
    : null;

  return (
    <div className="max-w-2xl mx-auto px-4 py-12">

      {/* Header */}
      <div className="mb-10 pb-8 border-b border-zinc-100 dark:border-zinc-800">
        <div className="flex items-center gap-2 text-xs font-medium text-zinc-400 uppercase tracking-widest mb-4">
          <Mail className="h-3.5 w-3.5" />
          Dein wöchentlicher Newsletter
        </div>
        <h1 className="text-3xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">
          Hallo {user.name?.split(" ")[0]} 👋
        </h1>
        <p className="text-zinc-500 dark:text-zinc-400">
          {totalTrips > 0
            ? `Wir haben ${totalTrips} Reisen gefunden die dich interessieren könnten.`
            : "Interagiere mit Trips um personalisierte Empfehlungen zu erhalten."}
        </p>
        {generatedAt && (
          <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-3">
            {generatedAt}
          </p>
        )}
      </div>

      {/* Empty state */}
      {!isLoading && (!newsletter || newsletter.sections.length === 0) && (
        <div className="text-center py-16">
          <p className="text-zinc-400 dark:text-zinc-500 text-sm">
            Noch keine Empfehlungen – like und kommentiere Trips um personalisierte Inhalte zu sehen.
          </p>
        </div>
      )}

      {/* Sections */}
      <div className="space-y-10">
        {newsletter?.sections.map((section) => {
          const meta = SECTION_META[section.title];
          const Icon = meta?.icon ?? Compass;
          const accent = meta?.accent ?? "text-zinc-500";
          const bg = meta?.bg ?? "bg-zinc-50 border-zinc-100";

          return (
            <div key={section.title}>
              {/* Section header */}
              <div className={`flex items-start gap-3 p-4 rounded-2xl border mb-4 ${bg}`}>
                <div className={`mt-0.5 ${accent}`}>
                  <Icon className="h-5 w-5" />
                </div>
                <div>
                  <h2 className="font-semibold text-zinc-900 dark:text-zinc-100">
                    {section.title}
                  </h2>
                  {section.description && (
                    <p className="text-sm text-zinc-500 dark:text-zinc-400 mt-0.5">
                      {section.description}
                    </p>
                  )}
                </div>
              </div>

              {/* Trip list */}
              <div className="divide-y divide-zinc-100 dark:divide-zinc-800/60 border border-zinc-100 dark:border-zinc-800 rounded-2xl overflow-hidden">
                {section.trips.map((trip) => (
                  <TripCard key={trip.tripId} trip={trip} accent={accent} />
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {/* Footer */}
      {newsletter && newsletter.sections.length > 0 && (
        <div className="mt-12 pt-8 border-t border-zinc-100 dark:border-zinc-800 text-center">
          <p className="text-xs text-zinc-400 dark:text-zinc-500">
            Diese Empfehlungen basieren auf deinen Interaktionen auf Trip Manager.
          </p>
        </div>
      )}
    </div>
  );
}