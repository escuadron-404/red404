import type { MouseEvent } from "react";

interface ButtonProps {
  text: string;
  onClick?: (event: MouseEvent<HTMLButtonElement>) => void;
  className?: string;
}

function Button(props: ButtonProps) {
  return (
    <button
      onClick={props.onClick}
      className={`${props.className} py-1 font-bold px-2 items-center justify-center border border-accent-secondary rounded-md`}
      type="button"
    >
      <span>{props.text}</span>
    </button>
  );
}
export default Button;
