"use client";

import { useRouter } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { register, login } from "@/lib/api/auth";
import AuthPage from "@/components/auth/AuthPage";
import type { CreateUserRequest, LoginRequest } from "@/types/user";

export default function AuthRoute() {
    const router = useRouter();
    const { updateUser } = useUserContext();

    const handleRegister = async (createUserRequest: CreateUserRequest) => {
        try {
            const user = await register(createUserRequest);
            updateUser(user);
            router.push("/");
        } catch (error) {
            console.error("Registration failed:", error);
            throw error;
        }
    };

    const handleLogin = async (loginRequest: LoginRequest) => {
        try {
            const user = await login(loginRequest);
            updateUser(user);
            router.push("/");
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