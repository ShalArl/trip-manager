import React from "react";
import {UserResponse as User} from "@/types/user";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {Badge} from "@/components/ui/badge";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {LogOut, Settings} from "lucide-react";
import {useRouter} from "next/navigation";
import {useUserContext} from "@/lib/context/UserContext";

type Props = {
    user: User;
    onLogout: () => void;
};

export default function Navbar({user: initialUser, onLogout}: Props) {
    const router = useRouter();
    const { user } = useUserContext();

    // Use context user if available, otherwise use initial user
    const displayUser = user || initialUser;

    // Debug logging
    React.useEffect(() => {
        console.log("[Navbar] User context user:", user);
        if (displayUser) {
            console.log("[Navbar] DisplayUser.avatarUrl:", displayUser.avatarUrl);
        }
    }, [user, displayUser]);

    // If no user, don't render the navbar (should redirect to login)
    if (!displayUser) {
        return null;
    }

    return (
        <nav className="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-black sticky top-0 z-50">
            <div className="mx-auto max-w-7xl px-6 py-4 flex items-center justify-between">
                <div className="flex items-center gap-2">
                    <span className="text-xl">🌍</span>
                    <span className="text-lg font-bold tracking-tight">TravelBuddy</span>
                </div>
                <div className="flex items-center gap-4">
          <span className="text-sm text-zinc-500 dark:text-zinc-400 hidden sm:block">
            Hallo, <span className="font-medium text-zinc-900 dark:text-white">{displayUser.name}</span>
          </span>

                    {/* Avatar mit Dropdown */}
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <div className="relative cursor-pointer">
                                <Avatar className="h-10 w-10">
                                    {displayUser.avatarUrl && <AvatarImage src={displayUser.avatarUrl} alt={displayUser.name}/>}
                                    <AvatarFallback className="bg-blue-500 text-white">
                                        {displayUser.name.charAt(0).toUpperCase()}
                                    </AvatarFallback>
                                </Avatar>
                                <Badge
                                    className="absolute -bottom-1 -right-1 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs">
                                    ✓
                                </Badge>
                            </div>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-56">
                            {/* User Info */}
                            <div className="px-2 py-1.5">
                                <p className="text-sm font-medium text-zinc-900 dark:text-white">{displayUser.name}</p>
                                <p className="text-xs text-zinc-500 dark:text-zinc-400">{displayUser.email}</p>
                            </div>

                            <DropdownMenuSeparator/>

                            {/* Settings */}
                            <DropdownMenuItem
                                onClick={() => router.push("/settings")}
                                className="cursor-pointer"
                            >
                                <Settings className="mr-2 h-4 w-4"/>
                                <span>Profileinstellungen</span>
                            </DropdownMenuItem>

                            <DropdownMenuSeparator/>

                            {/* Logout */}
                            <DropdownMenuItem onClick={onLogout} className="cursor-pointer text-red-600">
                                <LogOut className="mr-2 h-4 w-4"/>
                                <span>Abmelden</span>
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>
        </nav>
    );
}