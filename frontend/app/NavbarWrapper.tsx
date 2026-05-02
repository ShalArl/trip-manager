"use client";

import Navbar from "@/components/global/Navbar";
import { useUserContext } from "@/lib/context/UserContext";
import { logout } from "@/lib/api/auth";

export default function NavbarWrapper() {
    const { user, updateUser } = useUserContext();

    const handleLogout = async () => {
        await logout();
        updateUser(null);
    };

    return <Navbar user={user} onLogout={handleLogout} />;
}