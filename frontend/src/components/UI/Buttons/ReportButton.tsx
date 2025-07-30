import { LucideFlag } from "lucide-react";
interface ReportButtonProps {
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function ReportButton(props: ReportButtonProps) {
  return (
    <button
      className="flex items-center gap-2 py-2 px-8 w-full rounded-md hover:bg-accent"
      onClick={props.onClick}
    >
      <LucideFlag width={20} />
      <span className="text-red-500/90">report</span>
    </button>
  );
}
