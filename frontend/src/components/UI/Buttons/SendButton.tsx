import { SendHorizonalIcon } from "lucide-react";

interface SendButtonProps {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function SendButton(props: SendButtonProps) {
  return (
    <button
      className="absolute right-2 top-1/2 -translate-y-1/2"
      onClick={props.onClick}
      type="button"
    >
      <SendHorizonalIcon size={25} />
    </button>
  );
}
