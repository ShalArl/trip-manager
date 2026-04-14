"use client";

import { useState, useEffect } from "react";
import { getTrips } from "@/lib/api/trips";
import { components } from "@/generated/types";
import { register, login } from "@/lib/api/auth";
import { useUserContext } from "@/lib/context/UserContext";
import AuthPage from "@/components/auth/AuthPage";
import Navbar from "@/components/global/Navbar";
import Hero from "@/components/home/Hero";
import FeatureGrid from "@/components/home/FeatureGrid";
import TripList from "@/components/trips/TripList";
import { useQueryClient } from "@tanstack/react-query";

type CreateUserRequest = components["schemas"]["CreateUserRequest"];
type LoginRequest = components["schemas"]["LoginRequest"];
type AuthResponse = components["schemas"]["AuthResponse"];
type TripResponse = components["schemas"]["TripResponse"];


export default function Home() {
  
  const { user, isLoading, updateUser } = useUserContext();
  const [trips, setTrips] = useState<TripResponse[]>([]);
  const queryClient = useQueryClient();

  useEffect(() => {
    if (user) {
      getTrips().then(setTrips).catch(console.error);
    }
  }, [user]);

  const handleRegister = async (createUserRequest: CreateUserRequest) => {
    try {
      const response: AuthResponse = await register(createUserRequest)
      // Store token
      localStorage.setItem("token", response.token);
      localStorage.setItem("userId", response.user.id);
      // Update user context
      updateUser(response.user);
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  }

  const handleLogin = async (loginRequest: LoginRequest) => {
    try {
      const response = await login(loginRequest);
      // Store token
      localStorage.setItem("token", response.token);
      localStorage.setItem("userId", response.user.id);
      // Update user context
      updateUser(response.user);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  }

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("userId");
    localStorage.removeItem("user");
    updateUser(null);
  };


  if (isLoading) {
    return <div className="flex items-center justify-center min-h-screen">Loading...</div>;
  }

  if (!user) {
    return <AuthPage onLoginAction={handleLogin} onRegisterAction={handleRegister} />;
  }

  return (
    <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50">
      <Navbar user={user} onLogout={handleLogout} />
      <Hero />
      <TripList trips={trips} />
      <FeatureGrid />
    </div>
  );
}