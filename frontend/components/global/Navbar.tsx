import React, {useEffect} from "react";
import {UserResponse as User} from "@/types/user";
import {Badge} from "@/components/ui/badge";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {Building2, LogOut, Mail, Settings, TrendingUp} from "lucide-react";
import {useRouter} from "next/navigation";
import {useUserContext} from "@/lib/context/UserContext";
import {UserAvatar} from "@/components/global/UserAvatar";
import {useTenantContext} from "@/lib/context/TenantContext";
import Image from "next/image";

type Props = {
    user?: User | null;
    onLogout?: () => void;
};

export default function Navbar({user: initialUser, onLogout}: Props) {
    const router = useRouter();
    const {user} = useUserContext();
    const {tenantId, tenantName, branding, role, isAdmin} = useTenantContext();

    const displayUser = user || initialUser;

    useEffect(() => {
        console.log("[Navbar] User context user:", user);
        if (displayUser) {
            console.log("[Navbar] DisplayUser.avatarUrl:", displayUser.avatarUrl);
        }
    }, [user, displayUser]);


    return (
        <nav className="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-black sticky top-0 z-50">
            <div className="mx-auto max-w-7xl px-6 py-4 flex items-center justify-between">
                {/* Links: Logo + Navigation */}
                <div className="flex items-center gap-6">
                    <div
                        className="flex items-center gap-2 cursor-pointer"
                        onClick={() => router.push("/")}
                    >
                        {branding?.logoUrl ? (
                            <Image
                                src={branding.logoUrl}
                                alt={branding.companyName || tenantName}
                                className="h-8 object-contain max-w-32"
                            />
                        ) : (
                            <>
                                <span className="text-xl">🌍</span>
                                <div className="flex flex-col items-start leading-tight">
                                    <span className="text-lg font-bold tracking-tight">
                                        {branding?.companyName || "Trip Manager"}
                                    </span>
                                    {tenantId !== "default" && tenantName && !branding?.companyName && (
                                        <span className="text-xs text-slate-400 dark:text-slate-500 font-normal mt-0.5">
                                            by {tenantName}
                                        </span>
                                    )}
                                </div>
                            </>
                        )}
                    </div>

                    <button
                        onClick={() => router.push("/search")}
                        className="text-sm text-zinc-500 dark:text-zinc-400 hover:text-[var(--brand-primary)] dark:hover:text-[var(--brand-primary-light)] transition-colors"
                    >
                        Reisen entdecken
                    </button>
                    <button
                        onClick={() => router.push("/feed")}
                        className="flex items-center gap-1.5 text-sm text-zinc-500 dark:text-zinc-400 hover:text-[var(--brand-primary)] dark:hover:text-[var(--brand-primary-light)] transition-colors"
                    >
                        <TrendingUp className="h-4 w-4"/>
                        Feed
                    </button>
                    {displayUser && (
                        <button
                            onClick={() => router.push("/newsletter")}
                            className="flex items-center gap-1.5 text-sm text-zinc-500 dark:text-zinc-400 hover:text-[var(--brand-primary)] dark:hover:text-[var(--brand-primary-light)] transition-colors"
                        >
                            <Mail className="h-4 w-4"/>
                            Newsletter
                        </button>
                    )}
                </div>

                {/* Rechts: User oder Anmelden */}
                <div className="flex items-center gap-4">
                    {displayUser ? (
                        <>
                            <span className="text-sm text-zinc-500 dark:text-zinc-400 hidden sm:block">
                                Hallo,{" "}
                                <span className="font-medium text-zinc-900 dark:text-white">
                                    {displayUser.name}
                                </span>
                            </span>
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <div className="relative cursor-pointer">
                                        <UserAvatar name={displayUser.name}
                                                    className={"bg-blue-500 text-white"}
                                                    avatarKey={displayUser.avatarUrl}/>
                                        <Badge
                                            className="absolute -bottom-1 -right-1 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs">
                                            ✓
                                        </Badge>
                                    </div>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end" className="w-56">
                                    <div className="px-2 py-1.5">
                                        <p className="text-sm font-medium text-zinc-900 dark:text-white">
                                            {displayUser.name}
                                        </p>
                                        <p className="text-xs text-zinc-500 dark:text-zinc-400">
                                            {displayUser.email}
                                        </p>
                                        {tenantId !== "default" && tenantName && (
                                            <div className="flex items-center gap-1 mt-1">
                                                <Building2 className="h-3 w-3 text-sky-500"/>
                                                <span className="text-xs text-[var(--brand-primary)] dark:text-[var(--brand-primary-light)] font-medium">
                                                    {tenantName}
                                                </span>
                                            </div>
                                        )}
                                    </div>
                                    <DropdownMenuSeparator/>
                                    <DropdownMenuItem
                                        onClick={() => router.push("/settings")}
                                        className="cursor-pointer"
                                    >
                                        <Settings className="mr-2 h-4 w-4"/>
                                        <span>Profileinstellungen</span>
                                    </DropdownMenuItem>
                                    <DropdownMenuSeparator/>
                                    <DropdownMenuItem
                                        onClick={onLogout}
                                        className="cursor-pointer text-red-600"
                                    >
                                        <LogOut className="mr-2 h-4 w-4"/>
                                        <span>Abmelden</span>
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </>
                    ) : (
                        <button
                            onClick={() => router.push("/auth")}
                            className="px-4 py-2 text-sm font-medium bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-dark)] text-white rounded-lg transition-colors"
                        >
                            Anmelden
                        </button>
                    )}
                </div>
            </div>
        </nav>
    );
}