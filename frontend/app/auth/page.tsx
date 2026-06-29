"use client";
import { Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import {register, login, getAuthHeaders} from "@/lib/api/auth";
import AuthPage from "@/components/auth/AuthPage";
import type { CreateUserRequest, LoginRequest } from "@/types/user";
import {createAdvertiser} from "@/lib/api/advertiser";

function AuthContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const redirect = searchParams.get("redirect") || "/";
    const { updateUser } = useUserContext();

    const type = searchParams.get("type");

    const handleRegister = async (createUserRequest: CreateUserRequest) => {
        try {
            const user = await register(createUserRequest);
            updateUser(user);

            if (type === "advertiser") {
                await createAdvertiser({email: user.email, name: user.name})
                router.push("/advertiser");
            } else {
                router.push(redirect);
            }
        } catch (error) {
            console.error("Registration failed:", error);
            throw error;
        }
    };

    const handleLogin = async (loginRequest: LoginRequest) => {
        try {
            const user = await login(loginRequest);
            updateUser(user);
            router.push(redirect);
        } catch (error) {
            console.error("Login failed:", error);
            throw error;
        }
    };

    return (
        <AuthPage
            onLoginAction={handleLogin}
            onRegisterAction={handleRegister}
        />
    );
}

export default function AuthRoute() {
    return (
        <Suspense>
            <AuthContent />
        </Suspense>
    );
}