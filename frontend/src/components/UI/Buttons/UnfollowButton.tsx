import { UserRoundXIcon } from "lucide-react";

interface UnfollowButtonProps {
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function UnfollowButton(props: UnfollowButtonProps) {
  return (
    <button
      className="flex items-center gap-2 py-2 px-8 w-full rounded-md hover:bg-accent"
      onClick={props.onClick}
      type="button"
    >
      <UserRoundXIcon width={20} />
      <span className="text-red-500/90">unfollow</span>
    </button>
  );
}
