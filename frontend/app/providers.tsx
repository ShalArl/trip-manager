"use client";

import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {ReactNode, useMemo} from "react";
import {ErrorProvider} from "@/lib/context/ErrorContext";
import {UserProvider} from "@/lib/context/UserContext";
import {ToastContainer} from "@/components/global/ToastContainer";
import {TenantProvider} from "@/lib/context/TenantContext";

export function Providers({children}: { children: ReactNode }) {
    // Create a client once and reuse it - this ensures cache is shared across all routes
    const queryClient = useMemo(
        () =>
            new QueryClient({
                defaultOptions: {
                    queries: {
                        staleTime: 5 * 60 * 1000, // 5 minutes
                        gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
                    },
                },
            }),
        []
    );
    return (
        <ErrorProvider>
            <UserProvider>
                <TenantProvider>
                    <QueryClientProvider client={queryClient}>
                        {children}
                        <ToastContainer/>
                    </QueryClientProvider>
                </TenantProvider>
            </UserProvider>
        </ErrorProvider>
    );
}

