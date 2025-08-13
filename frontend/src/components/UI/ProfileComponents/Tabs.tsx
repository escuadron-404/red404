import { BookmarkIcon, BrainIcon, ImagesIcon } from "lucide-react";
import type { ReactNode } from "react";
// TODO: ?
// import { useState } from "react";
import { NavLink } from "react-router";

interface TabButtonProps {
  children: ReactNode;
  link: string;
}

function TabButton(props: TabButtonProps) {
  const isRoot = props.link === "/profile";
  return (
    <NavLink
      to={props.link}
      end={isRoot}
      className={({ isActive }) =>
        `flex items-center border-b-2 hover:border-b-accent ${isActive ? "border-b-accent-secondary" : "border-b-transparent"} border-base px-5 pb-2 transition duration-150 ease-in-out`
      }
    >
      {props.children}
    </NavLink>
  );
}

export default function Tabs() {
  return (
    <div className="flex w-full justify-evenly">
      <TabButton link="/profile">
        <ImagesIcon size={25} />
      </TabButton>
      <TabButton link="/profile/brainrot">
        <BrainIcon size={25} />
      </TabButton>
      <TabButton link="/profile/saved">
        <BookmarkIcon size={25} />
      </TabButton>
    </div>
  );
}
