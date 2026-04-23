"use client";

import {createContext, ReactNode, useContext, useEffect, useState} from "react";
import {onAuthStateChanged} from "firebase/auth";
import {UserResponse} from "@/types/user";
import {getMe} from "@/lib/api/auth";
import {firebaseAuth} from "@/lib/api/firebase";

type UserContextType = {
    user: UserResponse | null;
    isLoading: boolean;
    error: Error | null;
    updateUser: (user: UserResponse | null) => void;
    refetchUser: () => Promise<void>;
};

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<UserResponse | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        return onAuthStateChanged(firebaseAuth, async (firebaseUser) => {
            if (firebaseUser) {
                try {
                    const userData = await getMe();
                    setUser(userData);
                    setError(null);
                } catch (err) {
                    console.error("[UserContext] getMe failed:", err);
                    setError(err instanceof Error ? err : new Error("Failed to load user"));
                    setUser(null);
                }
            } else {
                setUser(null);
                setError(null);
            }
            setIsLoading(false);
        });
    }, []);

    const updateUser = (userData: UserResponse | null) => {
        setUser(userData);
    };

    const refetchUser = async () => {
        if (!firebaseAuth.currentUser) {
            setUser(null);
            return;
        }
        try {
            const userData = await getMe();
            setUser(userData);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error("Failed to fetch user"));
        }
    };

    return (
        <UserContext.Provider value={{ user, isLoading, error, updateUser, refetchUser }}>
            {children}
        </UserContext.Provider>
    );
}

export function useUserContext() {
    const context = useContext(UserContext);
    if (context === undefined) {
        throw new Error("useUserContext must be used within a UserProvider");
    }
    return context;
}