"use client";

import { useState, useRef, useEffect } from "react";
import { MapPin, Loader2 } from "lucide-react";

// ─── Types ────────────────────────────────────────────────────────────────────

export type PlaceValue = {
    name: string;
    city: string;
    country: string;
    lat?: number;
    lng?: number;
};

type NominatimResult = {
    place_id: number;
    display_name: string;
    name: string;
    address: {
        city?: string;
        town?: string;
        village?: string;
        municipality?: string;
        county?: string;
        country?: string;
        aerodrome?: string;
        road?: string;
    };
    lat: string;
    lon: string;
};

type Props = {
    label: string;
    value: PlaceValue | null;
    onChange: (place: PlaceValue) => void;
    placeholder?: string;
    required?: boolean;
};

// ─── Helpers ──────────────────────────────────────────────────────────────────

function extractCity(address: NominatimResult["address"]): string {
    return (
        address.city ??
        address.town ??
        address.village ??
        address.municipality ??
        address.county ??
        ""
    );
}

function extractName(result: NominatimResult): string {
    // Prefer the short name, fall back to first part of display_name
    if (result.name) return result.name;
    return result.display_name.split(",")[0].trim();
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function PlaceAutocomplete({ label, value, onChange, placeholder = "Ort suchen...", required }: Props) {
    const [query, setQuery] = useState(value ? `${value.name}, ${value.city}` : "");
    const [results, setResults] = useState<NominatimResult[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isOpen, setIsOpen] = useState(false);
    const [selected, setSelected] = useState<PlaceValue | null>(value);
    const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);

    // Close dropdown on outside click
    useEffect(() => {
        const handler = (e: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setIsOpen(false);
            }
        };
        document.addEventListener("mousedown", handler);
        return () => document.removeEventListener("mousedown", handler);
    }, []);

    const search = async (q: string) => {
        if (q.length < 2) {
            setResults([]);
            setIsOpen(false);
            return;
        }
        setIsLoading(true);
        try {
            const res = await fetch(
                `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(q)}&format=json&addressdetails=1&limit=5`,
                { headers: { "Accept-Language": "de" } }
            );
            const data: NominatimResult[] = await res.json();
            setResults(data);
            setIsOpen(data.length > 0);
        } catch {
            setResults([]);
        } finally {
            setIsLoading(false);
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const val = e.target.value;
        setQuery(val);
        setSelected(null); // reset selection when user types again

        if (debounceRef.current) clearTimeout(debounceRef.current);
        debounceRef.current = setTimeout(() => search(val), 400);
    };

    const handleSelect = (result: NominatimResult) => {
        const city = extractCity(result.address);
        const country = result.address.country ?? "";
        const name = extractName(result);

        const place: PlaceValue = {
            name,
            city,
            country,
            lat: parseFloat(result.lat),
            lng: parseFloat(result.lon),
        };

        setSelected(place);
        setQuery(`${name}${city ? ", " + city : ""}${country ? ", " + country : ""}`);
        setResults([]);
        setIsOpen(false);
        onChange(place);
    };

    return (
        <div ref={containerRef} className="relative">
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                {label} {required && <span className="text-red-500">*</span>}
            </label>

            <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input
                    type="text"
                    value={query}
                    onChange={handleInputChange}
                    onFocus={() => results.length > 0 && setIsOpen(true)}
                    placeholder={placeholder}
                    className={`w-full pl-9 pr-4 py-2 rounded-lg border bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white placeholder-zinc-400 focus:outline-none focus:ring-2 focus:ring-sky-500 transition-colors ${
                        selected
                            ? "border-sky-400 dark:border-sky-600"
                            : "border-zinc-200 dark:border-zinc-700"
                    }`}
                />
                {isLoading && (
                    <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400 animate-spin" />
                )}
            </div>

            {/* Selected badge */}
            {selected && (
                <p className="mt-1 text-xs text-sky-600 dark:text-sky-400">
                    {selected.city}, {selected.country}
                    {selected.lat && ` · ${selected.lat.toFixed(3)}, ${selected.lng?.toFixed(3)}`}
                </p>
            )}

            {/* Dropdown */}
            {isOpen && results.length > 0 && (
                <ul className="absolute z-50 mt-1 w-full bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl shadow-lg overflow-hidden">
                    {results.map((result) => {
                        const city = extractCity(result.address);
                        const country = result.address.country ?? "";
                        const name = extractName(result);
                        return (
                            <li key={result.place_id}>
                                <button
                                    type="button"
                                    onClick={() => handleSelect(result)}
                                    className="w-full text-left px-4 py-3 hover:bg-sky-50 dark:hover:bg-sky-950/30 transition-colors border-b border-zinc-100 dark:border-zinc-700 last:border-0"
                                >
                                    <p className="text-sm font-medium text-zinc-900 dark:text-white">{name}</p>
                                    <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">
                                        {[city, country].filter(Boolean).join(", ")}
                                    </p>
                                </button>
                            </li>
                        );
                    })}
                </ul>
            )}
        </div>
    );
}