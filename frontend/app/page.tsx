"use client";

import { useState, useEffect } from "react";

import { User } from "@/types/user";
import { Trip } from "@/types/trip";
import { getTrips } from "@/lib/api/trips";
import { components } from "@/generated/types";

//import { mockTrips } from "@/lib/mock-trips";
import { register, login } from "@/lib/api/auth";

import AuthPage from "@/components/auth/AuthPage";
import Navbar from "@/components/global/Navbar";
import Hero from "@/components/home/Hero";
import FeatureGrid from "@/components/home/FeatureGrid";
import TripList from "@/components/trips/TripList";

type CreateUserRequest = components["schemas"]["CreateUserRequest"];
type LoginRequest = components["schemas"]["LoginRequest"];
type AuthResponse = components["schemas"]["AuthResponse"];
type TripResponse = components["schemas"]["TripResponse"];


export default function Home() {
  
  const [user, setUser] = useState<User | null>(null);
  useEffect(() => {
    const saved = localStorage.getItem("user");
    if (saved) setUser(JSON.parse(saved));
  }, []);

  const [trips, setTrips] = useState<TripResponse[]>([]);
  useEffect(() => {
    if (user) {
      getTrips().then(setTrips).catch(console.error);
    }
  }, [user]);

  const handleRegister = async (createUserRequest: CreateUserRequest) => {
    const response: AuthResponse = await register(createUserRequest)
    console.log(response)
    // Store token FIRST, then user - this ensures token is available when useEffect runs
    localStorage.setItem("token", response.token);
    localStorage.setItem("user", JSON.stringify({ 
      id: response.user.id,
      name: response.user.name, 
      email: response.user.email 
    }));
    setUser({ name: response.user.name, email: response.user.email });
  }

  const handleLogin = async (loginRequest: LoginRequest) => {
    const response = await login(loginRequest);
    console.log(response)
    // Store token FIRST, then user - this ensures token is available when useEffect runs
    localStorage.setItem("token", response.token);
    localStorage.setItem("user", JSON.stringify({ 
      id: response.user.id,
      name: response.user.name, 
      email: response.user.email 
    }));
    console.log("Token: " + response.token);
    setUser({ name: response.user.name, email: response.user.email });
  }

  const handleLogout = () => {
    localStorage.removeItem("user");
    setUser(null);
  };

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