export function formatMonth(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString("de-DE", { month: "short" }).replace(".", "");
}

export function formatDay(dateStr: string): string {
    return String(new Date(dateStr).getDate());
}

export function formatYear(dateStr: string): string {
    return String(new Date(dateStr).getFullYear());
}

export function getDuration(start: string, end: string): string {
    const days = Math.ceil((new Date(end).getTime() - new Date(start).getTime()) / (1000 * 60 * 60 * 24));
    if (days === 0) return "1 Tag";
    if (days === 1) return "2 Tage";
    return `${days + 1} Tage`;
}