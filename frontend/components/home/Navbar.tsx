import { User } from "@/types/user";

type Props = {
  user: User;
  onLogout: () => void;
};

export default function Navbar({ user, onLogout }: Props) {
  return (
    <nav className="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-black sticky top-0 z-50">
      <div className="mx-auto max-w-7xl px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="text-xl">🌍</span>
          <span className="text-lg font-bold tracking-tight">TravelBuddy</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-sm text-zinc-500 dark:text-zinc-400 hidden sm:block">
            Hallo, <span className="font-medium text-zinc-900 dark:text-white">{user.name}</span>
          </span>
          <button
            onClick={onLogout}
            className="rounded-full bg-zinc-100 dark:bg-zinc-800 px-5 py-2 text-sm font-medium text-zinc-700 dark:text-zinc-300 hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors"
          >
            Abmelden
          </button>
        </div>
      </div>
    </nav>
  );
}