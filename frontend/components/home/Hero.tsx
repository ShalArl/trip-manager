import Link from "next/link";

export default function Hero() {
  return (
    <header className="bg-white dark:bg-black border-b border-zinc-100 dark:border-zinc-900 py-24">
      <div className="mx-auto max-w-4xl px-6 text-center flex flex-col items-center">
        <div className="inline-flex items-center rounded-full bg-sky-50 dark:bg-sky-950/50 px-4 py-1.5 text-sm font-medium text-sky-700 dark:text-sky-300 mb-6 border border-sky-100 dark:border-sky-900">
          🌴 Dein digitaler Reisebegleiter
        </div>
        <h1 className="mb-6 text-5xl sm:text-6xl font-bold tracking-tight text-zinc-950 dark:text-white leading-tight">
          Freizeitstress adé.<br />
          <span className="text-sky-600 dark:text-sky-400">Reiseplanung juche.</span>
        </h1>
        <p className="mb-10 max-w-2xl text-xl leading-relaxed text-zinc-500 dark:text-zinc-400">
          Organisiere deine Trips, erstelle Packlisten und behalte dein Budget im Griff.
        </p>
        <Link
          href="/trips/new"
          className="h-14 rounded-2xl bg-sky-600 px-12 font-semibold text-white hover:bg-sky-700 active:scale-[0.98] transition-all shadow-lg shadow-sky-500/25 text-base inline-flex items-center"
        >
          Reise planen
        </Link>
      </div>
    </header>
  );
}