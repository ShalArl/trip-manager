"use client";
import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import { firebaseAuth } from "@/lib/api/firebase";
import { onAuthStateChanged, User as FirebaseUser } from "firebase/auth";
import { getTenant } from "@/lib/api/auth";

type Branding = {
    logoUrl: string;
    primaryColor: string;
    companyName: string;
    customDomain: string;
} | null;

type TenantContextType = {
    tenantId: string;
    tenantName: string;
    role: string;
    isAdmin: boolean;
    isOwner: boolean;
    isPlatformAdmin: boolean;
    isAdvertiser: boolean;
    branding: Branding;
};

const DEFAULT_CLAIMS: TenantContextType = {
    tenantId: "default",
    tenantName: "",
    role: "tenant_member",
    isAdmin: false,
    isOwner: false,
    isPlatformAdmin: false,
    isAdvertiser: false,
    branding: null,
};

const TenantContext = createContext<TenantContextType>(DEFAULT_CLAIMS);

export function TenantProvider({ children }: { children: ReactNode }) {
    const [claims, setClaims] = useState<TenantContextType>(DEFAULT_CLAIMS);

    const loadClaims = async (user: FirebaseUser) => {
        const token = await user.getIdTokenResult(true);
        const tenantId = (token.claims.tenant_id as string) ?? "default";
        const role = (token.claims.role as string) ?? "tenant_member";

        console.log("[TenantContext] claims:", { tenantId, role }); // ← temporär


        let tenantName = "";
        let branding: Branding = null;

        if (tenantId !== "default") {
            try {
                const tenant = await getTenant();
                tenantName = tenant?.name ?? "";
            } catch {}

            try {
                const idToken = await firebaseAuth.currentUser?.getIdToken();
                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/tenants/me/branding`, {
                    headers: { Authorization: `Bearer ${idToken}` }
                });
                if (res.ok) {
                    branding = await res.json();
                }
            } catch {}

            if (branding?.primaryColor) {
                document.documentElement.style.setProperty("--brand-primary", branding!.primaryColor);
                document.documentElement.style.setProperty("--brand-primary-dark", branding!.primaryColor);
            }
        } else {
            // Reset CSS-Variablen wenn kein Tenant
            document.documentElement.style.removeProperty("--brand-primary");
            document.documentElement.style.removeProperty("--brand-primary-dark");
            document.documentElement.style.removeProperty("--brand-primary-light");
        }

        setClaims({
            tenantId,
            tenantName,
            role,
            branding,
            isAdmin: ["tenant_admin", "tenant_owner", "platform_admin"].includes(role),
            isOwner: ["tenant_owner", "platform_admin"].includes(role),
            isPlatformAdmin: role === "platform_admin",
            isAdvertiser: role === "advertiser",
        });
    };

    useEffect(() => {
        return onAuthStateChanged(firebaseAuth, async (user) => {
            if (!user) {
                document.documentElement.style.removeProperty("--brand-primary");
                document.documentElement.style.removeProperty("--brand-primary-dark");
                document.documentElement.style.removeProperty("--brand-primary-light");
                setClaims(DEFAULT_CLAIMS);
                return;
            }
            await loadClaims(user);
        });
    }, []);

    return (
        <TenantContext.Provider value={claims}>
            {children}
        </TenantContext.Provider>
    );
}

export function useTenantContext() {
    return useContext(TenantContext);
}