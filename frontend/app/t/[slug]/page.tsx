"use client";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getTenantBySlug } from "@/lib/api/auth";
import { useUserContext } from "@/lib/context/UserContext";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";

export default function TenantPage() {
    const { slug } = useParams<{ slug: string }>();
    const router = useRouter();
    const { user, isLoading } = useUserContext();
    const [tenant, setTenant] = useState<{ tenantId: string; name: string; tier: string; slug: string } | null>(null);
    const [notFound, setNotFound] = useState(false);

    useEffect(() => {
        getTenantBySlug(slug).then((t) => {
            if (!t) setNotFound(true);
            else setTenant(t);
        });
    }, [slug]);

    if (isLoading || (!tenant && !notFound)) return <LoadingSpinner />;

    if (notFound) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <p className="text-2xl font-bold text-zinc-900 dark:text-white mb-2">Reisebüro nicht gefunden</p>
                    <p className="text-zinc-500 mb-4">Der Link ist ungültig oder das Reisebüro existiert nicht mehr.</p>
                    <button onClick={() => router.push("/")} className="text-[var(--brand-primary)] hover:underline">
                        Zur Startseite
                    </button>
                </div>
            </div>
        );
    }

    // User ist eingeloggt → direkt zur App
    if (user) {
        router.push("/");
        return null;
    }

    // Nicht eingeloggt → zur Tenant-spezifischen Auth-Seite
    router.push(`/t/${slug}/auth`);
    return null;
}