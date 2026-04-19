"use client";

import { useError } from "@/lib/context/ErrorContext";
import { X } from "lucide-react";

export function ToastContainer() {
  const { toasts, removeToast } = useError();

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm">
      {toasts.map((toast) => (
        <Toast
          key={toast.id}
          id={toast.id}
          message={toast.message}
          type={toast.type}
          onClose={() => removeToast(toast.id)}
        />
      ))}
    </div>
  );
}

interface ToastProps {
  id: string;
  message: string;
  type: "error" | "success" | "info" | "warning";
  onClose: () => void;
}

function Toast({ id, message, type, onClose }: ToastProps) {
  const bgColor = {
    error: "bg-red-500",
    success: "bg-green-500",
    info: "bg-blue-500",
    warning: "bg-yellow-500",
  }[type];

  const icon = {
    error: "❌",
    success: "✅",
    info: "ℹ️",
    warning: "⚠️",
  }[type];

  return (
    <div
      className={`${bgColor} text-white rounded-lg shadow-lg p-4 flex items-center justify-between gap-3 animate-in slide-in-from-bottom-5 fade-in`}
      role="alert"
    >
      <div className="flex items-center gap-3 flex-1">
        <span className="text-lg">{icon}</span>
        <p className="text-sm font-medium">{message}</p>
      </div>
      <button
        onClick={onClose}
        className="flex-shrink-0 hover:opacity-80 transition-opacity"
        aria-label="Schließen"
      >
        <X size={18} />
      </button>
    </div>
  );
}

