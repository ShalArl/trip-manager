"use client";

import { useState, useEffect } from "react";
import { useUserContext } from "@/lib/context/UserContext";
import { logout } from "@/lib/api/auth";
import { searchTrips, getRecentPublicTrips } from "@/lib/api/trips";
import { components } from "@/generated/types";
import Navbar from "@/components/global/Navbar";
import Link from "next/link";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Search, ChevronsLeft, ChevronLeft, ChevronRight, ChevronsRight, Plane } from "lucide-react";

type TripResponse = components["schemas"]["TripResponse"];

const PAGE_SIZE_OPTIONS = [10, 25, 50, 100];
const DEFAULT_PAGE_SIZE = 25;

export default function SearchPage() {
    const { user, updateUser } = useUserContext();
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
            <Navbar user={user} onLogout={handleLogout} />

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
                            <SelectValue />
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

                {/* Results */}
                {isLoading ? (
                    <div className="flex flex-col gap-3">
                        {Array.from({ length: 5 }).map((_, i) => (
                            <Skeleton key={i} className="h-20 rounded-2xl" />
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
                                    className="group"
                                >
                                    <Card className="px-6 py-5 flex items-center justify-between hover:border-sky-400 dark:hover:border-sky-600 hover:shadow-sm transition-all cursor-pointer">
                                        <div className="flex items-center gap-4 min-w-0">
                                            <div className="w-10 h-10 rounded-xl bg-sky-50 dark:bg-sky-950/50 border border-sky-100 dark:border-sky-900 flex items-center justify-center flex-shrink-0">
                                                <Plane className="h-5 w-5 text-sky-600 dark:text-sky-400" />
                                            </div>
                                            <div className="min-w-0">
                                                <p className="font-semibold truncate group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors">
                                                    {trip.title}
                                                </p>
                                                <p className="text-sm text-zinc-500 dark:text-zinc-400">
                                                    {trip.startDate} · {trip.endDate}
                                                </p>
                                            </div>
                                        </div>
                                        <ChevronRight className="h-5 w-5 text-zinc-400 group-hover:text-sky-500 group-hover:translate-x-0.5 transition-all flex-shrink-0 ml-2" />
                                    </Card>
                                </Link>
                            ))}
                        </div>

                        {/* Pagination Footer */}
                        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mt-8">
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
                                    <ChevronsLeft className="h-4 w-4" />
                                </Button>
                                <Button
                                    variant="outline"
                                    size="icon"
                                    onClick={() => setPage(page - 1)}
                                    disabled={page === 0}
                                    aria-label="Vorherige Seite"
                                >
                                    <ChevronLeft className="h-4 w-4" />
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
                                    <ChevronRight className="h-4 w-4" />
                                </Button>
                                <Button
                                    variant="outline"
                                    size="icon"
                                    onClick={() => setPage(totalPages - 1)}
                                    disabled={page >= totalPages - 1}
                                    aria-label="Letzte Seite"
                                >
                                    <ChevronsRight className="h-4 w-4" />
                                </Button>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}