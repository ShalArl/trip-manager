"use client";

import { useEffect, useState } from "react";

type WarningResponse = {
    countryCode: string;
    countryName: string;
    level: number;
    warning: boolean;
    partialWarning: boolean;
    updatedAt: string;
};

type Props = {
    countryCode: string;
    countryName: string;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

export default function TravelWarningWidget({ countryCode, countryName }: Props) {
    const [data, setData] = useState<WarningResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);

    useEffect(() => {
        if (!countryCode) return;
        setLoading(true);
        setError(false);

        fetch(`${API_URL}/api/warnings/${countryCode.toUpperCase()}`)
            .then((r) => { if (!r.ok) throw new Error(); return r.json(); })
            .then((json: WarningResponse) => { setData(json); setLoading(false); })
            .catch(() => { setError(true); setLoading(false); });
    }, [countryCode]);

    const header = (
        <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-3">
            Reisewarnung · {countryName}
        </p>
    );

    if (loading) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                {header}
                <div className="animate-pulse h-6 w-40 bg-zinc-100 dark:bg-zinc-800 rounded" />
            </div>
        );
    }

    if (error || !data) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5">
                {header}
                <p className="text-sm text-zinc-400 dark:text-zinc-500">Keine Daten verfügbar</p>
            </div>
        );
    }

    if (data.warning || data.level >= 3) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-red-200 dark:border-red-900 rounded-2xl p-5">
                {header}
                <div className="flex items-start gap-3">
                    <span className="text-2xl shrink-0">🔴</span>
                    <div>
                        <p className="text-sm font-medium text-red-700 dark:text-red-400 mb-1">Reisewarnung</p>
                        <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
                            Das Auswärtige Amt hat eine Reisewarnung für dieses Reiseziel ausgegeben.
                            Von nicht notwendigen Reisen wird abgeraten.
                        </p>
                        <a
                            href={`https://www.auswaertiges-amt.de/de/ReiseUndSicherheit/reise-und-sicherheitshinweise`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-block mt-2 text-xs text-red-600 dark:text-red-400 hover:underline"
                        >
                            Mehr Informationen →
                        </a>
                    </div>
                </div>
            </div>
        );
    }

    if (data.partialWarning || data.level === 2) {
        return (
            <div className="bg-white dark:bg-zinc-900 border border-amber-200 dark:border-amber-900 rounded-2xl p-5">
                {header}
                <div className="flex items-start gap-3">
                    <span className="text-2xl shrink-0">🟡</span>
                    <div>
                        <p className="text-sm font-medium text-amber-700 dark:text-amber-400 mb-1">Teilreisewarnung</p>
                        <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
                            Das Auswärtige Amt hat eine Teilreisewarnung für dieses Reiseziel ausgegeben.
                            Bitte informieren Sie sich vor Ihrer Reise.
                        </p>
                        <a
                            href={`https://www.auswaertiges-amt.de/de/ReiseUndSicherheit/reise-und-sicherheitshinweise`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-block mt-2 text-xs text-amber-600 dark:text-amber-400 hover:underline"
                        >
                            Mehr Informationen →
                        </a>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white dark:bg-zinc-900 border border-green-200 dark:border-green-900 rounded-2xl p-5">
            {header}
            <div className="flex items-center gap-3">
                <span className="text-2xl">🟢</span>
                <p className="text-sm text-zinc-600 dark:text-zinc-400">
                    Keine besonderen Sicherheitshinweise für dieses Reiseziel.
                </p>
            </div>
        </div>
    );
}