"use client";

import React, {useEffect, useState} from "react";
import {useUserContext} from "@/lib/context/UserContext";
import {logout} from "@/lib/api/auth";
import {getRecentPublicTrips, searchTrips} from "@/lib/api/trips";
import {components} from "@/generated/types";
import Link from "next/link";
import {Input} from "@/components/ui/input";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue,} from "@/components/ui/select";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {Button} from "@/components/ui/button";
import {Card} from "@/components/ui/card";
import {Skeleton} from "@/components/ui/skeleton";
import {ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, Search} from "lucide-react";
import {formatDay, formatMonth, formatYear, getDuration} from "@/utils/date"
import {router} from "next/client";
import {UserAvatar} from "@/components/global/UserAvatar";

type TripResponse = components["schemas"]["TripResponse"];

const PAGE_SIZE_OPTIONS = [10, 25, 50, 100];
const DEFAULT_PAGE_SIZE = 25;

export default function SearchPage() {
    const {user, updateUser} = useUserContext();
    const [query, setQuery] = useState("");
    const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
    const [trips, setTrips] = useState<TripResponse[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(0);
    const [isLoading, setIsLoading] = useState(false);

    const handleLogout = async () => {
        await logout();
        updateUser(null);
    };

    // Reset to page 0 when query or pageSize changes
    useEffect(() => {
        setPage(0);
    }, [query, pageSize]);

    useEffect(() => {
        const fetchTrips = async () => {
            const trimmedQuery = query.trim();
            const offset = page * pageSize;

            if (trimmedQuery.length === 0) {
                setIsLoading(true);
                try {
                    const result = await getRecentPublicTrips(pageSize, offset);
                    setTrips(result.data);
                    setTotal(result.total);
                } catch (error) {
                    console.error(error);
                } finally {
                    setIsLoading(false);
                }
                return;
            }

            if (trimmedQuery.length < 3) return;

            setIsLoading(true);
            try {
                const result = await searchTrips(trimmedQuery, pageSize, offset);
                setTrips(result.data);
                setTotal(result.total);
            } catch (error) {
                console.error(error);
            } finally {
                setIsLoading(false);
            }
        };

        const debounce = setTimeout(fetchTrips, 300);
        return () => clearTimeout(debounce);
    }, [query, page, pageSize]);

    const totalPages = Math.max(1, Math.ceil(total / pageSize));
    const currentPage = page + 1; // 1-indexed for display
    const firstVisibleItem = total === 0 ? 0 : page * pageSize + 1;
    const lastVisibleItem = Math.min((page + 1) * pageSize, total);

    return (
        <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
            <div className="mx-auto max-w-4xl px-6 py-12">
                <div className="mb-10">
                    <h1 className="text-3xl font-bold tracking-tight mb-2">
                        Reisen entdecken
                    </h1>
                    <p className="text-zinc-500 dark:text-zinc-400">
                        Entdecke Reisen von anderen Reisenden
                    </p>
                </div>

                {/* Search + Page Size Controls */}
                <div className="flex flex-col sm:flex-row gap-3 mb-8">
                    <div className="relative flex-1">
                        <Search
                            className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-400 pointer-events-none"
                        />
                        <Input
                            type="text"
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            placeholder="Nach Reisen suchen..."
                            className="pl-10 h-11"
                        />
                    </div>

                    <Select
                        value={String(pageSize)}
                        onValueChange={(v) => setPageSize(Number(v))}
                    >
                        <SelectTrigger className="w-full sm:w-[140px] h-11">
                            <SelectValue/>
                        </SelectTrigger>
                        <SelectContent>
                            {PAGE_SIZE_OPTIONS.map((size) => (
                                <SelectItem key={size} value={String(size)}>
                                    {size} pro Seite
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>

                {/* Pagination Header */}
                {total > 0 && (
                    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mb-8">
                        <p className="text-sm text-zinc-500 dark:text-zinc-400">
                            {firstVisibleItem}–{lastVisibleItem} von {total}
                        </p>

                        <div className="flex items-center gap-1">
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setPage(0)}
                                disabled={page === 0}
                                aria-label="Erste Seite"
                            >
                                <ChevronsLeft className="h-4 w-4"/>
                            </Button>
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setPage(page - 1)}
                                disabled={page === 0}
                                aria-label="Vorherige Seite"
                            >
                                <ChevronLeft className="h-4 w-4"/>
                            </Button>

                            <div className="px-4 text-sm font-medium min-w-[80px] text-center">
                                {currentPage} / {totalPages}
                            </div>

                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setPage(page + 1)}
                                disabled={page >= totalPages - 1}
                                aria-label="Nächste Seite"
                            >
                                <ChevronRight className="h-4 w-4"/>
                            </Button>
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setPage(totalPages - 1)}
                                disabled={page >= totalPages - 1}
                                aria-label="Letzte Seite"
                            >
                                <ChevronsRight className="h-4 w-4"/>
                            </Button>
                        </div>
                    </div>
                )}

                {/* Results */}
                {isLoading ? (
                    <div className="flex flex-col gap-3">
                        {Array.from({length: 5}).map((_, i) => (
                            <Skeleton key={i} className="h-20 rounded-2xl"/>
                        ))}
                    </div>
                ) : trips.length === 0 ? (
                    <Card className="p-12 text-center border-dashed">
                        <p className="text-zinc-500 dark:text-zinc-400">
                            {query ? "Keine Reisen gefunden" : "Gib etwas ein um zu suchen"}
                        </p>
                    </Card>
                ) : (
                    <>
                        <div className="flex flex-col gap-3">
                            {trips.map((trip) => (
                                <Link
                                    key={trip.id}
                                    href={`/trips/${encodeURIComponent(trip.id)}`}
                                    className="group block"
                                >
                                    <Card
                                        className="p-5 hover:border-sky-300 dark:hover:border-sky-700/50 hover:shadow-md transition-all cursor-pointer overflow-hidden">
                                        <div className="flex items-center gap-5">
                                            {/* Date badge - visual anchor left */}
                                            <div
                                                className="flex flex-col items-center justify-center flex-shrink-0 w-16 py-2 rounded-xl bg-sky-50 dark:bg-sky-950/40 border border-sky-100 dark:border-sky-900/50">
                                                <p className="text-[10px] uppercase tracking-wider text-sky-600 dark:text-sky-400 font-semibold">
                                                    {formatMonth(trip.startDate)}
                                                </p>
                                                <p className="text-xl font-bold text-sky-900 dark:text-sky-200 leading-none mt-0.5">
                                                    {formatDay(trip.startDate)}
                                                </p>
                                                <p className="text-[10px] text-sky-600/70 dark:text-sky-400/70 mt-0.5">
                                                    {formatYear(trip.startDate)}
                                                </p>
                                            </div>

                                            {/* Main content */}
                                            <div className="min-w-0 flex-1">
                                                <div className="flex items-center gap-2 mb-1">
                                                    <h3 className="font-semibold text-base truncate group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors">
                                                        {trip.title}
                                                    </h3>
                                                    <span
                                                        className="text-xs text-zinc-400 dark:text-zinc-600 flex-shrink-0"> · {getDuration(trip.startDate, trip.endDate)}</span>
                                                </div>
                                                {trip.shortDescription && (
                                                    <p className="text-sm text-zinc-500 dark:text-zinc-400 line-clamp-1">
                                                        {trip.shortDescription}
                                                    </p>
                                                )}

                                                {/* Author inline */}
                                                <div className="flex items-center gap-2 mt-2">
                                                    <UserAvatar name={trip.createdBy.name} avatarKey={trip.createdBy.avatarUrl} />
                                                    <button
                                                        type="button"
                                                        onClick={(e) => {
                                                            e.preventDefault();  // wichtig! verhindert dass outer Link triggert
                                                            e.stopPropagation();
                                                            router
                                                                .push(`/users/${trip.createdBy?.id || '#'}`)
                                                                .finally();
                                                        }}
                                                        className="text-xs text-zinc-500 dark:text-zinc-400 hover:text-sky-600 dark:hover:text-sky-400 transition-colors truncate"
                                                    >
                                                        {trip.createdBy?.name || 'Unbekannt'}
                                                    </button>
                                                </div>
                                            </div>

                                            {/* Chevron */}
                                            <ChevronRight
                                                className="h-5 w-5 text-zinc-300 dark:text-zinc-700 group-hover:text-sky-500 group-hover:translate-x-1 transition-all flex-shrink-0"/>
                                        </div>
                                    </Card>
                                </Link>
                            ))}
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}