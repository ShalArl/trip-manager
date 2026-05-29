"use client";

import { useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { getDownloadUrl } from "@/lib/api/uploads";

type UserAvatarProps = {
    avatarKey?: string | null;
    name: string;
    className?: string;
};

export function UserAvatar({ avatarKey, name, className }: UserAvatarProps) {
    const [src, setSrc] = useState<string | null>(null);

    useEffect(() => {
        let cancelled = false;

        if (avatarKey) {
            getDownloadUrl(avatarKey).then((url) => {
                if (!cancelled) setSrc(url);
            });
        } else {
            setSrc(null);
        }

        return () => {
            cancelled = true;
        };
    }, [avatarKey]);

    return (
        <Avatar className={className}>
            {src && <AvatarImage src={src} alt={name} />}
            <AvatarFallback className="bg-blue-500 text-white">
                {name.charAt(0).toUpperCase()}
            </AvatarFallback>
        </Avatar>
    );
}