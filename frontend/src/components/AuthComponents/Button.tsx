import type { ReactNode } from "react";

interface ButtonProps {
  provider?: string;
  children?: ReactNode;
  type?: "button" | "submit";
  text?: string;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
  className?: string;
}

export default function Button(props: ButtonProps) {
  return (
    <button
      type={props.type}
      className={`w-64 flex items-center justify-center gap-1.5 px-5 py-2 border border-stone-700 rounded-xl transition duration-150 ease-in-out hover:bg-accent ${props.className}`}
      onClick={props.onClick}
    >
      {props.text}
      {props.children}
      {props.provider && <span>Continue with {props.provider}</span>}
    </button>
  );
}
