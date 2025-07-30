import { SendIcon } from "lucide-react";
interface ShareButtonProps {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function ShareButton(props: ShareButtonProps) {
  return (
    <button className="share" onClick={props.onClick}>
      <SendIcon
        className="transition duration-150 ease-in-out hover:text-accent-secondary"
        width={28}
        height={28}
      />
    </button>
  );
}
