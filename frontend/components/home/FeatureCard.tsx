type Props = {
  emoji: string;
  title: string;
  description: string;
  bgColor: string;
  borderColor: string;
  iconBg: string;
  iconBorder: string;
};

export default function FeatureCard({
  emoji,
  title,
  description,
  bgColor,
  borderColor,
  iconBg,
  iconBorder,
}: Props) {
  return (
    <div className={`${bgColor} p-8 rounded-3xl border ${borderColor} hover:scale-[1.02] transition-transform cursor-default`}>
      <div className={`w-14 h-14 rounded-2xl ${iconBg} flex items-center justify-center text-3xl mb-6 border ${iconBorder}`}>
        {emoji}
      </div>
      <h3 className="mb-3 text-xl font-bold tracking-tight">{title}</h3>
      <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed text-sm">{description}</p>
    </div>
  );
}