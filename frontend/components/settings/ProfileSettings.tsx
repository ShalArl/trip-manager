"use client";

import {useRef, useState} from "react";
import {UserResponse, UpdateUserRequest} from "@/types/user";
import {updateMe} from "@/lib/api/auth";
import {uploadAvatar} from "@/lib/api/uploads";
import {useUserContext} from "@/lib/context/UserContext";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Label} from "@/components/ui/label";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {AlertCircle, CheckCircle, ImagePlus, Mail, Trash2} from "lucide-react";

type ProfileSettingsProps = {
    user: UserResponse;
};

const ProfileSettings = ({user}: ProfileSettingsProps) => {
    const [name, setName] = useState(user.name);
    const [email, setEmail] = useState(user.email);
    const [bio, setBio] = useState(user.bio || "");
    const [avatarPreview, setAvatarPreview] = useState<string | null>(user.avatarUrl || null);
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const { updateUser: updateUserContext } = useUserContext();

    const handleAvatarSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            // Validate file type
            if (!file.type.startsWith("image/")) {
                setError("Bitte wähle eine Bilddatei aus");
                return;
            }

            // Validate file size (max 5MB)
            if (file.size > 5 * 1024 * 1024) {
                setError("Datei muss kleiner als 5MB sein");
                return;
            }

            setAvatarFile(file);
            const reader = new FileReader();
            reader.onload = (e) => {
                setAvatarPreview(e.target?.result as string);
            };
            reader.readAsDataURL(file);
            setError(null);
        }
    };

    const handleRemoveAvatar = () => {
        setAvatarFile(null);
        setAvatarPreview(null);
        if (fileInputRef.current) {
            fileInputRef.current.value = "";
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(false);
        setLoading(true);

        try {
            let avatarUrl: string | undefined = avatarPreview || undefined;

            // If avatar file was selected, upload it directly to S3/MinIO using presigned URL
            if (avatarFile) {
                console.log("[ProfileSettings] Starting presigned avatar upload...");
                const userId = user.id;
                avatarUrl = await uploadAvatar(avatarFile, userId);
                console.log("[ProfileSettings] Avatar uploaded successfully:", avatarUrl);
                setAvatarFile(null);
            }

            // Prepare data for profile update (no file, just metadata + avatar URL)
            const updateData: UpdateUserRequest = {
                name,
                email,
                bio,
                ...(avatarUrl && { avatarUrl }),
            };

            console.log("[ProfileSettings] Updating user profile...");
            const data = await updateMe(updateData);
            console.log("[ProfileSettings] Profile updated:", data);

            // Update user context - broadcasts to all components using UserContext
            updateUserContext(data);
            setSuccess(true);
            setAvatarPreview(data.avatarUrl || null);
            setTimeout(() => setSuccess(false), 3000);
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Fehler beim Aktualisieren des Profils"
            );
        } finally {
            setLoading(false);
        }
    };

    const hasChanges =
        name !== user.name ||
        email !== user.email ||
        bio !== (user.bio || "") ||
        avatarPreview !== (user.avatarUrl || null) ||
        avatarFile !== null;

    return (
        <div className="rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 shadow-lg p-8">
            <form onSubmit={handleSubmit} className="space-y-8">
                {/* Success Message */}
                {success && (
                    <div
                        className="flex items-center gap-3 rounded-lg bg-green-50 dark:bg-green-950/50 border border-green-200 dark:border-green-900 p-4 animate-in fade-in duration-300">
                        <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400 flex-shrink-0"/>
                        <p className="text-sm font-medium text-green-800 dark:text-green-200">
                            Profil erfolgreich aktualisiert
                        </p>
                    </div>
                )}

                {/* Error Message */}
                {error && (
                    <div
                        className="flex items-center gap-3 rounded-lg bg-red-50 dark:bg-red-950/50 border border-red-200 dark:border-red-900 p-4 animate-in fade-in duration-300">
                        <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0"/>
                        <p className="text-sm font-medium text-red-800 dark:text-red-200">
                            {error}
                        </p>
                    </div>
                )}

                {/* Avatar Section */}
                <div className="space-y-4">
                    <Label className="text-sm font-semibold text-zinc-900 dark:text-white flex items-center gap-2">
                        <ImagePlus className="h-4 w-4"/>
                        Profilbild
                    </Label>

                    <div className="flex items-center gap-6">
                        {/* Avatar Preview - Same as Navbar */}
                        <div className="flex-shrink-0">
                            <Avatar className="h-32 w-32 border-4 border-zinc-200 dark:border-zinc-800">
                                {avatarPreview && <AvatarImage src={avatarPreview} alt={name}/>}
                                <AvatarFallback className="bg-blue-500 text-white text-lg font-semibold">
                                    {name.charAt(0).toUpperCase()}
                                </AvatarFallback>
                            </Avatar>
                        </div>

                        {/* Upload Controls */}
                        <div className="space-y-3">
                            <input
                                ref={fileInputRef}
                                type="file"
                                accept="image/*"
                                onChange={handleAvatarSelect}
                                className="hidden"
                            />
                            <div className="flex gap-2">
                                <Button
                                    type="button"
                                    onClick={() => fileInputRef.current?.click()}
                                    className="bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-all duration-200"
                                >
                                    Bild hochladen
                                </Button>
                                {(avatarPreview || user.avatarUrl) && (
                                    <Button
                                        type="button"
                                        onClick={handleRemoveAvatar}
                                        variant="outline"
                                        className="border-red-300 text-red-600 hover:bg-red-50 dark:hover:bg-red-950/50"
                                    >
                                        <Trash2 className="h-4 w-4"/>
                                    </Button>
                                )}
                            </div>
                            <p className="text-xs text-zinc-500 dark:text-zinc-400">
                                JPG, PNG oder GIF bis 5MB
                            </p>
                        </div>
                    </div>
                </div>

                {/* Divider */}
                <div className="border-t border-zinc-200 dark:border-zinc-800"/>

                {/* Name Field */}
                <div className="space-y-3">
                    <Label htmlFor="name"
                           className="text-sm font-semibold text-zinc-900 dark:text-white flex items-center gap-2">
                        📝
                        Name
                    </Label>
                    <Input
                        id="name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        minLength={1}
                        placeholder="Dein Name"
                        className="bg-zinc-50 dark:bg-zinc-900 border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
                    />
                </div>

                {/* Email Field */}
                <div className="space-y-3">
                    <Label htmlFor="email"
                           className="text-sm font-semibold text-zinc-900 dark:text-white flex items-center gap-2">
                        <Mail className="h-4 w-4"/>
                        E-Mail
                    </Label>
                    <Input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        placeholder="deine@email.com"
                        className="bg-zinc-50 dark:bg-zinc-900 border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
                    />
                </div>

                {/* Bio Field */}
                <div className="space-y-3">
                    <Label htmlFor="bio" className="text-sm font-semibold text-zinc-900 dark:text-white">
                        Bio
                    </Label>
                    <textarea
                        id="bio"
                        value={bio}
                        onChange={(e) => setBio(e.target.value)}
                        placeholder="Erzähle etwas über dich selbst... (optional)"
                        maxLength={500}
                        rows={4}
                        className="w-full px-4 py-2 rounded-lg bg-zinc-50 dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 text-zinc-900 dark:text-white placeholder-zinc-500 dark:placeholder-zinc-400 resize-none"
                    />
                    <p className="text-xs text-zinc-500 dark:text-zinc-400">
                        {bio.length}/500 Zeichen
                    </p>
                </div>

                {/* Info Text */}
                <div
                    className="rounded-lg bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-900 p-4">
                    <p className="text-sm text-blue-800 dark:text-blue-200">
                        <span className="font-semibold">Tipp:</span> Deine Änderungen werden sofort gespeichert.
                    </p>
                </div>

                {/* Submit Button */}
                <div className="pt-4 flex gap-3">
                    <Button
                        type="submit"
                        disabled={loading || !hasChanges}
                        className="flex-1 sm:flex-none bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-all duration-200"
                    >
                        {loading ? (
                            <div className="flex items-center gap-2">
                                <div
                                    className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"/>
                                Wird gespeichert...
                            </div>
                        ) : (
                            "Änderungen speichern"
                        )}
                    </Button>
                </div>
            </form>
        </div>
    );
};

export default ProfileSettings;

