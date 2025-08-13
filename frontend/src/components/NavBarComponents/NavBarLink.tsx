import type { ReactNode } from "react";
import { Link } from "react-router";

interface NavBarLinkProps {
  link: string;
  children: ReactNode;
}

export default function NavBarLink(props: NavBarLinkProps) {
  const linkName = props.link.slice(1);
  return (
    <Link
      className={`${`${linkName}-link`}  hover:bg-accent w-full rounded-md py-3 pl-3 pr-16 flex gap-3 items-center transition duration-100 ease-in-out`}
      to={props.link}
    >
      {props.children}
      <span className="font-semibold text-lg">{linkName}</span>
    </Link>
  );
}
