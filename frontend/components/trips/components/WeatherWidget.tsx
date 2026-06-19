"use client";

import {useEffect, useState} from "react";

type ForecastDay = {
    date: string;
    tempMax: number;
    tempMin: number;
    precipitationMm: number;
    weatherCode: number;
    description: string;
};

type WeatherResponse = {
    latitude: number;
    longitude: number;
    forecast: ForecastDay[];
    updatedAt: string;
};

type Props = {
    lat: number;
    lng: number;
    locationName: string;
    startDate?: string;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

const SUPPORTED_PROVIDERS = [
    {name: "DWD Germany", url: "https://open-meteo.com/en/docs/dwd-api"},
    {name: "NOAA U.S.", url: "https://open-meteo.com/en/docs/gfs-api"},
    {name: "Météo-France", url: "https://open-meteo.com/en/docs/meteofrance-api"},
    {name: "ECMWF", url: "https://open-meteo.com/en/docs/ecmwf-api"},
    {name: "UK Met Office", url: "https://open-meteo.com/en/docs/ukmo-api"},
    {name: "KMA Korea", url: "https://open-meteo.com/en/docs/kma-api"},
    {name: "JMA Japan", url: "https://open-meteo.com/en/docs/jma-api"},
    {name: "MeteoSwiss", url: "https://open-meteo.com/en/docs/meteoswiss-api"},
    {name: "MET Norway", url: "https://open-meteo.com/en/docs/metno-api"},
    {name: "GEM Canada", url: "https://open-meteo.com/en/docs/gem-api"},
    {name: "BOM Australia", url: "https://open-meteo.com/en/docs/bom-api"},
    {name: "CMA China", url: "https://open-meteo.com/en/docs/cma-api"},
    {name: "KNMI Netherlands", url: "https://open-meteo.com/en/docs/knmi-api"},
    {name: "DMI Denmark", url: "https://open-meteo.com/en/docs/dmi-api"},
    {name: "ItaliaMeteo", url: "https://open-meteo.com/en/docs/italia-meteo-arpae-api"},
    {name: "GeoSphere Austria", url: "https://open-meteo.com/en/docs/geosphere-austria-api"},
];

function weatherEmoji(description: string): string {
    const d = (description ?? "").toLowerCase();
    if (d.includes("gewitter") || d.includes("thunder")) return "⛈️";
    if (d.includes("regen") || d.includes("rain") || d.includes("drizzle") || d.includes("niesel")) return "🌧️";
    if (d.includes("schnee") || d.includes("snow")) return "❄️";
    if (d.includes("nebel") || d.includes("fog") || d.includes("mist")) return "🌫️";
    if (d.includes("bewölkt") || d.includes("cloud") || d.includes("bedeckt")) return "☁️";
    if (d.includes("heiter") || d.includes("partly")) return "⛅";
    if (d.includes("klar") || d.includes("sonnig") || d.includes("clear") || d.includes("sunny")) return "☀️";
    return "🌤️";
}

function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString("de-DE", {weekday: "short", day: "2-digit", month: "2-digit"});
}

export default function WeatherWidget({lat, lng, locationName, startDate}: Props) {
    const [data, setData] = useState<WeatherResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [noData, setNoData] = useState(false);

    useEffect(() => {
        setLoading(true);
        setNoData(false);
        setData(null);

        const params = new URLSearchParams({lat: String(lat), lng: String(lng)});
        if (startDate) {
            params.set("date", startDate);
        }

        fetch(`${API_URL}/api/info/weather/?lat=${lat}&lng=${lng}${startDate ? `&date=${startDate}` : ""}`)
            .then((r) => {
                if (!r.ok) throw new Error("no data");
                return r.json();
            })
            .then((json: WeatherResponse) => {
                if (!json.forecast || json.forecast.length === 0) {
                    setNoData(true);
                } else {
                    setData(json);
                }
                setLoading(false);
            })
            .catch(() => {
                setNoData(true);
                setLoading(false);
            });
    }, [lat, lng, startDate]);

    const header = (
        <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-3">
            Wetter · {locationName}
        </p>
    );

    if (loading) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                {header}
                <div className="animate-pulse space-y-2">
                    {[...Array(3)].map((_, i) => (
                        <div key={i} className="h-14 bg-zinc-100 dark:bg-zinc-800 rounded-xl"/>
                    ))}
                </div>
            </div>
        );
    }

    if (noData || !data) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                {header}
                <p className="text-sm text-zinc-500 dark:text-zinc-400 mb-3">
                    Für diesen Ort sind keine Wetterdaten verfügbar.
                </p>
                <p className="text-xs text-zinc-400 dark:text-zinc-500 mb-2">Unterstützte Anbieter:</p>
                <ul className="space-y-1">
                    {SUPPORTED_PROVIDERS.map((p) => (
                        <li key={p.name}>
                            <a
                                href={p.url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-xs text-[var(--brand-primary)] dark:text-[var(--brand-primary-light)] hover:underline"
                            >
                                {p.name}
                            </a>
                        </li>
                    ))}
                </ul>
            </div>
        );
    }

    const days = data.forecast.slice(0, 3);

    return (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
            {header}
            <div className="space-y-2">
                {days.map((day, i) => (
                    <div
                        key={day.date}
                        className={`flex items-center gap-3 p-3 rounded-xl ${
                            i === 0
                                ? "bg-sky-50 dark:bg-sky-950/30 border border-sky-100 dark:border-sky-900"
                                : "bg-zinc-50 dark:bg-zinc-800/50"
                        }`}
                    >
                        <span className="text-2xl leading-none shrink-0">{weatherEmoji(day.description)}</span>
                        <div className="flex-1 min-w-0">
                            <p className="text-xs text-zinc-400 dark:text-zinc-500">{formatDate(day.date)}</p>
                            <p className="text-sm text-zinc-600 dark:text-zinc-400 truncate">{day.description}</p>
                        </div>
                        <div className="text-right shrink-0">
                            <p className="text-sm font-medium text-zinc-900 dark:text-white">
                                {Math.round(day.tempMax)}° <span
                                className="text-zinc-400 font-normal">/ {Math.round(day.tempMin)}°</span>
                            </p>
                            {day.precipitationMm > 0 && (
                                <p className="text-xs text-sky-500 dark:text-[var(--brand-primary-light)]">{day.precipitationMm} mm</p>
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}