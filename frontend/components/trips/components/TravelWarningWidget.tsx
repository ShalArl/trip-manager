"use client";

import { useEffect, useState } from "react";

type WarningLevel = "none" | "low" | "medium" | "high" | "extreme";

type TravelWarningData = {
    level: WarningLevel;
    title: string;
    description: string;
    countryName: string;
};

type Props = {
    countryCode: string;
    countryName: string;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

const LEVEL_CONFIG: Record<WarningLevel, { label: string; color: string; bg: string; barColor: string; barWidth: string }> = {
    none:    { label: "Sicher",           color: "#3B6D11", bg: "#EAF3DE", barColor: "#639922", barWidth: "10%" },
    low:     { label: "Geringe Warnung",  color: "#3B6D11", bg: "#EAF3DE", barColor: "#639922", barWidth: "25%" },
    medium:  { label: "Mittlere Warnung", color: "#854F0B", bg: "#FAEEDA", barColor: "#EF9F27", barWidth: "50%" },
    high:    { label: "Hohe Warnung",     color: "#A32D2D", bg: "#FCEBEB", barColor: "#E24B4A", barWidth: "75%" },
    extreme: { label: "Reisewarnung",     color: "#A32D2D", bg: "#FCEBEB", barColor: "#E24B4A", barWidth: "100%" },
};

export default function TravelWarningWidget({ countryCode, countryName }: Props) {
    const [warning, setWarning] = useState<TravelWarningData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);

    useEffect(() => {
        setLoading(true);
        setError(false);
        fetch(`${API_URL}/api/warnings/${countryCode.toUpperCase()}`)
            .then((r) => r.json())
            .then((data) => {
                setWarning(data);
                setLoading(false);
            })
            .catch(() => {
                setError(true);
                setLoading(false);
            });
    }, [countryCode]);

    if (loading) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-3">Reisewarnung · {countryName}</p>
                <div className="animate-pulse space-y-3">
                    <div className="h-6 w-28 bg-zinc-100 dark:bg-zinc-800 rounded" />
                    <div className="h-2 bg-zinc-100 dark:bg-zinc-800 rounded-full" />
                    <div className="h-12 bg-zinc-100 dark:bg-zinc-800 rounded" />
                </div>
            </div>
        );
    }

    if (error || !warning) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                <p className="text-xs font-medium text-zinc-400 uppercase tracking-wider mb-2">Reisewarnung · {countryName}</p>
                <p className="text-sm text-zinc-400 dark:text-zinc-500">Keine Daten verfügbar</p>
            </div>
        );
    }

    const cfg = LEVEL_CONFIG[warning.level] ?? LEVEL_CONFIG.none;

    return (
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
            <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-3">
                Reisewarnung · {warning.countryName || countryName}
            </p>
            <div className="flex items-center gap-3 mb-3">
                <span
                    className="text-xs font-medium px-2.5 py-1 rounded-md"
                    style={{ background: cfg.bg, color: cfg.color }}
                >
                    {cfg.label}
                </span>
                <div className="flex-1 h-1.5 rounded-full bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
                    <div
                        className="h-full rounded-full transition-all duration-500"
                        style={{ width: cfg.barWidth, background: cfg.barColor }}
                    />
                </div>
            </div>
            {warning.title && (
                <p className="text-sm font-medium text-zinc-900 dark:text-white mb-1">{warning.title}</p>
            )}
            <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">{warning.description}</p>
        </div>
    );
}