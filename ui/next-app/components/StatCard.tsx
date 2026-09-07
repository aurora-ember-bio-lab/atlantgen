interface StatCardProps {
  label: string;
  value: string | number;
  tone?: "default" | "amber" | "emerald" | "red";
}

const toneStyles: Record<NonNullable<StatCardProps["tone"]>, string> = {
  default: "text-neutral-100",
  amber: "text-amber-400",
  emerald: "text-emerald-400",
  red: "text-red-400",
};

export default function StatCard({ label, value, tone = "default" }: StatCardProps) {
  return (
    <div className="rounded-lg border border-neutral-800 bg-neutral-900/40 p-5">
      <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
        {label}
      </p>
      <p className={`mt-2 text-2xl font-semibold ${toneStyles[tone]}`}>{value}</p>
    </div>
  );
}
