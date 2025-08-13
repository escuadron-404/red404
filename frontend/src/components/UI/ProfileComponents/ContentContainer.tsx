import type { ReactNode } from "react";

interface ContentContainerProps {
  children: ReactNode;
}

export default function ContentContainer(props: ContentContainerProps) {
  return <div className="grid grid-cols-3 gap-2">{props.children}</div>;
}
