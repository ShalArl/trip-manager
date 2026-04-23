"use client";

import { useState } from "react";
import LoginForm from "./LoginForm";
import RegisterForm from "./RegisterForm";
import {CreateUserRequest, LoginRequest} from "@/types/user";


type AuthMode = "login" | "register";


type Props = {
  onLoginAction: (loginRequest: LoginRequest) => void;
  onRegisterAction: (createUserRequest: CreateUserRequest) => void;
};

const BRAND_FEATURES = [
  { emoji: "📅", text: "Intelligente Reiseplanung" },
  { emoji: "🎒", text: "Dynamische Packlisten" },
  { emoji: "🪙", text: "Budget-Tracking" },
];

export default function AuthPage({ onLoginAction, onRegisterAction }: Props) {
  const [mode, setMode] = useState<AuthMode>("login");

  return (
    <div className="min-h-screen bg-stone-50 dark:bg-zinc-950 flex">

      {/* Left panel — branding (nur Desktop) */}
      <div className="hidden lg:flex lg:w-1/2 bg-sky-950 dark:bg-zinc-900 flex-col justify-between p-14 relative overflow-hidden">
        <div className="absolute inset-0 opacity-10">
          <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse">
                <path d="M 40 0 L 0 0 0 40" fill="none" stroke="white" strokeWidth="0.5" />
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#grid)" />
          </svg>
        </div>

        <div className="relative z-10">
          <div className="flex items-center gap-3 mb-16">
            <span className="text-2xl">🌍</span>
            <span className="text-xl font-bold tracking-tight text-white">TravelBuddy</span>
          </div>
          <h2 className="text-4xl font-bold text-white leading-snug mb-6">
            Deine nächste Reise<br />wartet auf dich.
          </h2>
          <p className="text-sky-200 text-lg leading-relaxed max-w-sm">
            Erstelle Reisepläne, Packlisten und verwalte dein Budget — alles an einem Ort.
          </p>
        </div>

        <div className="relative z-10 space-y-5">
          {BRAND_FEATURES.map((item) => (
            <div key={item.text} className="flex items-center gap-4">
              <div className="w-10 h-10 rounded-xl bg-white/10 flex items-center justify-center text-lg border border-white/20">
                {item.emoji}
              </div>
              <span className="text-sky-100 font-medium">{item.text}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Right panel — form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8">
        <div className="w-full max-w-md">

          {/* Mobile logo */}
          <div className="flex items-center gap-2 mb-10 lg:hidden">
            <span className="text-2xl">🌍</span>
            <span className="text-xl font-bold text-zinc-900 dark:text-white">TravelBuddy</span>
          </div>

          {mode === "login" ? (
            <LoginForm
              onLoginAction={onLoginAction}
              onSwitchToRegisterAction={() => setMode("register")}
            />
          ) : (
            <RegisterForm
              onRegisterAction={onRegisterAction}
              onSwitchToLoginAction={() => setMode("login")}
            />
          )}
        </div>
      </div>
    </div>
  );
}