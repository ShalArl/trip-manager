"use client";

import React, { createContext, useContext, useState, useCallback } from "react";

export interface Toast {
  id: string;
  message: string;
  type: "error" | "success" | "info" | "warning";
  duration?: number;
}

interface ErrorContextType {
  toasts: Toast[];
  addError: (message: string, duration?: number) => void;
  addSuccess: (message: string, duration?: number) => void;
  addInfo: (message: string, duration?: number) => void;
  addWarning: (message: string, duration?: number) => void;
  removeToast: (id: string) => void;
}

const ErrorContext = createContext<ErrorContextType | undefined>(undefined);

export function ErrorProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback((message: string, type: Toast["type"], duration = 4000) => {
    const id = Math.random().toString(36).substr(2, 9);
    const newToast: Toast = { id, message, type, duration };

    setToasts((prev) => [...prev, newToast]);

    if (duration) {
      setTimeout(() => {
        removeToast(id);
      }, duration);
    }

    return id;
  }, []);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addError = useCallback((message: string, duration?: number) => {
    addToast(message, "error", duration);
  }, [addToast]);

  const addSuccess = useCallback((message: string, duration?: number) => {
    addToast(message, "success", duration);
  }, [addToast]);

  const addInfo = useCallback((message: string, duration?: number) => {
    addToast(message, "info", duration);
  }, [addToast]);

  const addWarning = useCallback((message: string, duration?: number) => {
    addToast(message, "warning", duration);
  }, [addToast]);

  return (
    <ErrorContext.Provider
      value={{
        toasts,
        addError,
        addSuccess,
        addInfo,
        addWarning,
        removeToast,
      }}
    >
      {children}
    </ErrorContext.Provider>
  );
}

export function useError() {
  const context = useContext(ErrorContext);
  if (!context) {
    throw new Error("useError must be used within ErrorProvider");
  }
  return context;
}

