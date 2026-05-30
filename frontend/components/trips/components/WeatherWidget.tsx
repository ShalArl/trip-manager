"use client";

import { useEffect, useState } from "react";

type WeatherData = {
    temperature: number;
    feelsLike: number;
    humidity: number;
    windSpeed: number;
    pressure: number;
    description: string;
    icon: string;
};

type Props = {
    lat: number;
    lng: number;
    locationName: string;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

function weatherEmoji(description: string): string {
    const d = (description ?? "").toLowerCase();
    if (d.includes("thunder")) return "⛈️";
    if (d.includes("rain") || d.includes("drizzle")) return "🌧️";
    if (d.includes("snow")) return "❄️";
    if (d.includes("fog") || d.includes("mist")) return "🌫️";
    if (d.includes("cloud")) return "⛅";
    if (d.includes("clear") || d.includes("sunny")) return "☀️";
    return "🌤️";
}

export default function WeatherWidget({ lat, lng, locationName }: Props) {
    const [weather, setWeather] = useState<WeatherData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);

    useEffect(() => {
        setLoading(true);
        setError(false);
        fetch(`${API_URL}/api/weather/current?lat=${lat}&lng=${lng}`)
            .then((r) => r.json())
            .then((data) => {
                setWeather(data);
                setLoading(false);
            })
            .catch(() => {
                setError(true);
                setLoading(false);
            });
    }, [lat, lng]);

    if (loading) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-3">Wetter · {locationName}</p>
                <div className="animate-pulse space-y-3">
                    <div className="h-10 w-24 bg-zinc-100 dark:bg-zinc-800 rounded" />
                    <div className="grid grid-cols-2 gap-2">
                        {[...Array(4)].map((_, i) => (
                            <div key={i} className="h-14 bg-zinc-100 dark:bg-zinc-800 rounded-lg" />
                        ))}
                    </div>
                </div>
            </div>
        );
    }

    if (error || !weather) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-2">Wetter · {locationName}</p>
                <p className="text-sm text-zinc-400 dark:text-zinc-500">Wetterdaten nicht verfügbar</p>
            </div>
        );
    }

    return (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
            <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-3">
                Wetter · {locationName}
            </p>
            <div className="flex items-end justify-between mb-4">
                <div>
                    <p className="text-4xl font-light text-zinc-900 dark:text-white leading-none">
                        {Math.round(weather.temperature)}°
                    </p>
                    <p className="text-sm text-zinc-500 dark:text-zinc-400 mt-1 capitalize">
                        {weather.description}
                    </p>
                </div>
                <span className="text-5xl leading-none">{weatherEmoji(weather.description)}</span>
            </div>
            <div className="grid grid-cols-2 gap-2">
                {[
                    { label: "Gefühlt", value: `${Math.round(weather.feelsLike)}°` },
                    { label: "Luftfeuchtigkeit", value: `${weather.humidity}%` },
                    { label: "Wind", value: `${Math.round(weather.windSpeed)} km/h` },
                    { label: "Druck", value: `${weather.pressure} hPa` },
                ].map(({ label, value }) => (
                    <div key={label} className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-3 text-center">
                        <p className="text-sm font-medium text-zinc-900 dark:text-white">{value}</p>
                        <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-0.5">{label}</p>
                    </div>
                ))}
            </div>
        </div>
    );
}