"use client";

import { useRouter } from "next/navigation";
import { useUserContext } from "@/lib/context/UserContext";
import { register, login } from "@/lib/api/auth";
import { components } from "@/generated/types";
import AuthPage from "@/components/auth/AuthPage";

type CreateUserRequest = components["schemas"]["CreateUserRequest"];
type LoginRequest = components["schemas"]["LoginRequest"];
type AuthResponse = components["schemas"]["AuthResponse"];

export default function AuthRoute() {
    const router = useRouter();
    const { updateUser } = useUserContext();

    const handleRegister = async (createUserRequest: CreateUserRequest) => {
        try {
            const response: AuthResponse = await register(createUserRequest);
            localStorage.setItem("token", response.token);
            localStorage.setItem("userId", response.user.id);
            updateUser(response.user);
            router.push("/");
        } catch (error) {
            console.error("Registration failed:", error);
            throw error;
        }
    };

    const handleLogin = async (loginRequest: LoginRequest) => {
        try {
            const response = await login(loginRequest);
            localStorage.setItem("token", response.token);
            localStorage.setItem("userId", response.user.id);
            updateUser(response.user);
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