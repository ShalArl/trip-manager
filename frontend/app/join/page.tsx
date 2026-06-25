"use client";

import {Suspense, useEffect, useState} from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { LoadingSpinner } from "@/components/global/LoadingSpinner";
import { Building2, CheckCircle, XCircle } from "lucide-react";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function getAuthHeaders(): Promise<HeadersInit> {
    const { firebaseAuth } = await import("@/lib/api/firebase");
    const user = firebaseAuth.currentUser;
    if (!user) return { "Content-Type": "application/json" };
    const token = await user.getIdToken();
    return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

function JoinContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const token = searchParams.get("token");
    const { user, isLoading } = useUserContext();
    const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
    const [errorMsg, setErrorMsg] = useState("");
    const [tenantName, setTenantName] = useState("");

    useEffect(() => {
        if (isLoading) return;
        if (!token) {
            setStatus("error");
            setErrorMsg("Kein Einladungstoken gefunden.");
            return;
        }
        if (!user) {
            // Nicht eingeloggt → zur Auth-Seite mit Redirect-Parameter
            router.push(`/auth?redirect=/join?token=${token}`);
            return;
        }

        // Eingeloggt → Token einlösen
        setStatus("loading");
        getAuthHeaders().then(async (headers) => {
            try {
                const res = await fetch(`${API_URL}/api/users/tenants/join?token=${token}`, {
                    method: "POST",
                    headers,
                });
                if (res.ok) {
                    const data = await res.json();
                    setTenantName(data.tenantId);
                    setStatus("success");
                    // Token refreshen damit neue Claims im JWT landen
                    const { firebaseAuth } = await import("@/lib/api/firebase");
                    await firebaseAuth.currentUser?.getIdToken(true);
                    setTimeout(() => router.push("/"), 2000);
                } else {
                    const err = await res.json();
                    setErrorMsg(err.error || "Ungültige oder abgelaufene Einladung.");
                    setStatus("error");
                }
            } catch {
                setErrorMsg("Verbindungsfehler. Bitte versuche es erneut.");
                setStatus("error");
            }
        });
    }, [user, isLoading, token, router]);

    if (isLoading || status === "loading") return <LoadingSpinner />;

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center px-4">
            <div className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 shadow-xl p-8 w-full max-w-md text-center">
                {status === "success" && (
                    <>
                        <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
                        <h1 className="text-xl font-bold text-zinc-900 dark:text-white mb-2">
                            Willkommen im Team!
                        </h1>
                        <p className="text-sm text-zinc-500">
                            Du wurdest erfolgreich zum Reisebüro hinzugefügt. Du wirst weitergeleitet...
                        </p>
                    </>
                )}
                {status === "error" && (
                    <>
                        <XCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
                        <h1 className="text-xl font-bold text-zinc-900 dark:text-white mb-2">
                            Einladung ungültig
                        </h1>
                        <p className="text-sm text-zinc-500 mb-6">{errorMsg}</p>
                        <button
                            onClick={() => router.push("/")}
                            className="px-4 py-2 text-sm bg-[var(--brand-primary)] text-white rounded-lg hover:bg-[var(--brand-primary-dark)]"
                        >
                            Zur Startseite
                        </button>
                    </>
                )}
                {status === "idle" && (
                    <>
                        <Building2 className="h-12 w-12 text-[var(--brand-primary)] mx-auto mb-4" />
                        <h1 className="text-xl font-bold text-zinc-900 dark:text-white mb-2">
                            Einladung wird verarbeitet...
                        </h1>
                    </>
                )}
            </div>
        </div>
    );
}

export default function JoinPage() {
    return (
        <Suspense>
            <JoinContent />
        </Suspense>
    );
}