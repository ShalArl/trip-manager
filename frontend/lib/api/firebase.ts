import { initializeApp, getApps } from "firebase/app";
import { getAuth, connectAuthEmulator } from "firebase/auth";

const firebaseConfig = {
    apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY || "fake-api-key",  // Emulator akzeptiert alles
    authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN || "localhost",
    appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID || "fake-app-id",
    projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID!,
};

const app = getApps().length === 0 ? initializeApp(firebaseConfig) : getApps()[0];
export const firebaseAuth = getAuth(app);


if (process.env.NEXT_PUBLIC_USE_FIREBASE_EMULATOR === "true") {
    connectAuthEmulator(firebaseAuth, "http://localhost:9099", { disableWarnings: true });
}