import FeatureCard from "./FeatureCard";

const FEATURES = [
  {
    emoji: "📅",
    title: "Intelligenter Planer",
    description: "Ziehe Sehenswürdigkeiten, Restaurants und Hotels einfach in deinen Tagesablauf.",
    bgColor: "bg-amber-50 dark:bg-amber-950/30",
    borderColor: "border-amber-100 dark:border-amber-800/30",
    iconBg: "bg-amber-100 dark:bg-amber-900/40",
    iconBorder: "border-amber-200 dark:border-amber-700/40",
  },
  {
    emoji: "🎒",
    title: "Nie wieder vergessen",
    description: "Dynamische Packlisten, die sich an dein Reiseziel und das Wetter anpassen.",
    bgColor: "bg-emerald-50 dark:bg-emerald-950/30",
    borderColor: "border-emerald-100 dark:border-emerald-800/30",
    iconBg: "bg-emerald-100 dark:bg-emerald-900/40",
    iconBorder: "border-emerald-200 dark:border-emerald-700/40",
  },
  {
    emoji: "🪙",
    title: "Budget im Blick",
    description: "Erfasse Ausgaben unterwegs und teile Kosten fair mit Mitreisenden.",
    bgColor: "bg-rose-50 dark:bg-rose-950/30",
    borderColor: "border-rose-100 dark:border-rose-800/30",
    iconBg: "bg-rose-100 dark:bg-rose-900/40",
    iconBorder: "border-rose-200 dark:border-rose-700/40",
  },
];

export default function FeatureGrid() {
  return (
    <main className="mx-auto max-w-7xl px-6 py-24">
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-3">
        {FEATURES.map((feature) => (
          <FeatureCard key={feature.title} {...feature} />
        ))}
      </div>
    </main>
  );
}