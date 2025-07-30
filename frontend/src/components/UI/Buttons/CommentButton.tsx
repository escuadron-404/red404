import { MessageCircleIcon } from "lucide-react";
interface CommentButtonProps {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function CommentButton(props: CommentButtonProps) {
  return (
    <button className="comment" onClick={props.onClick}>
      <MessageCircleIcon
        className="transition duration-150 ease-in-out hover:text-accent-secondary"
        width={28}
        height={28}
      />
    </button>
  );
}
