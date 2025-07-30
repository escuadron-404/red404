import { EyeOff } from "lucide-react";
interface HideButtonProps {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function HideButton(props: HideButtonProps) {
  return (
    <button
      className="flex items-center gap-2 py-2 px-8 w-full rounded-md hover:bg-accent"
      onClick={props.onClick}
    >
      <EyeOff size={20} />
      <span className="text-red-500/90">Hide</span>
    </button>
  );
}
