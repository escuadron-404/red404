import type { ReactNode } from "react";
import { useState } from "react";
import { NavLink, Route, Routes } from "react-router";
import { UseAuth } from "@/auth/context/auth-context";
import ContentContainer from "@/components/UI/ProfileComponents/ContentContainer";
import Tabs from "@/components/UI/ProfileComponents/Tabs";
import profileData from "@/templates/ProfileTemplate";

interface ButtonProps {
  text: string;
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
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

interface FollowButtonProps {
  onFollow: (event: React.MouseEvent<HTMLButtonElement>) => void;
  followStatus: boolean;
}

function FollowButton(props: FollowButtonProps) {
  return (
    <Button
      onClick={props.onFollow}
      text={props.followStatus ? "followed" : "follow"}
      className="w-24"
    />
  );
}

function MessageButton() {
  return <Button text="message" onClick={() => alert("message")} />;
}

function MoreOptions() {
  return <Button text="..." />;
}

function LogoutButton() {
  const { logout } = UseAuth();
  return <Button text="logout" onClick={() => logout()} />;
}

function ProfileActions() {
  const [isFollowed, setIsFollowed] = useState(false);
  function handleOnFollow() {
    setIsFollowed((prev) => !prev);
    console.log(isFollowed);
  }
  return (
    <div className="flex gap-3 items-center">
      <FollowButton onFollow={handleOnFollow} followStatus={isFollowed} />
      <MessageButton />
      <MoreOptions />
      <LogoutButton />
    </div>
  );
}

interface ButtonDataProps {
  text: string;
  dataCount: number;
}

function ButtonData(props: ButtonDataProps) {
  // TODO: is this really just a button?
  return (
    <button className="flex gap-1 items-center" type="button">
      <span className="font-bold">{props.dataCount}</span> {props.text}
    </button>
  );
}

interface UserDataProps {
  userFullName: string;
  username: string;
  biography: string;
}
function UserData(props: UserDataProps) {
  return (
    <div className="flex flex-col items-center">
      <div className="flex gap-6 mb-4 items-center">
        <span className="font-bold text-xl">{props.userFullName}</span>
        <ProfileActions />
      </div>
      <span className="font-bold mb-3 self-start">@{props.username}</span>
      <div className="flex gap-5 items-center">
        <ButtonData text="following" dataCount={40} />
        <ButtonData text="followers" dataCount={21} />
        <ButtonData text="posts" dataCount={10} />
      </div>
      <p className="biography self-start mt-3">{props.biography}</p>
    </div>
  );
}

function Bio() {
  return (
    <div className="flex gap-24 w-full p-5 ">
      <img
        className="w-44 h-44 rounded-full object-fill"
        src="/testUserImage.jpg"
        alt="user profile avatar"
      />
      <UserData
        userFullName="tobias"
        username={profileData.userData.userName}
        biography="literally i dont even exist"
      />
    </div>
  );
}

interface HighlightProps {
  image?: string;
  title: string;
  link: string;
}

function Highlight(props: HighlightProps) {
  return (
    <a href={props.link} className="flex flex-col items-center gap-2.5">
      <img
        className="w-36 h-36 rounded-md"
        src="/testPostImage.jpg"
        alt="Hightlight"
      />
      <h3 className="font-semibold text-text-light">{props.title}</h3>
    </a>
  );
}
interface HighlightsProps {
  children: ReactNode;
}
function Highlights(props: HighlightsProps) {
  return (
    <div className="flex flex-col items-center">
      <NavLink to="highlights" className="font-bold text-lg mb-5">
        highlights
      </NavLink>
      {props.children}
    </div>
  );
}
export default function ProfilePage() {
  return (
    <main className="w-full min-h-screen flex flex-col gap-10 pl-20 py-4">
      <Bio />
      <div className="flex w-full items-start gap-10">
        <Highlights>
          <ul className="flex flex-col gap-5">
            {profileData.highlightsData.map((highlight) => (
              <li key={highlight.id}>
                <Highlight title={highlight.title} link={highlight.link} />
              </li>
            ))}
          </ul>
        </Highlights>
        <div className="flex w-full flex-col gap-4 pr-10">
          <Tabs />
          <Routes>
            <Route
              index
              element={
                <ContentContainer>
                  {profileData.posts.map((post) => (
                    <NavLink
                      key={post.id}
                      className="group h-96 flex"
                      to={`/${profileData.userData.userName}/posts/${post.id}`}
                    >
                      <img
                        className="rounded-sm w-full object-cover group-hover:opacity-40 transition duration-150 ease-in-out"
                        src={post.postMedia}
                        alt="post media"
                      />
                    </NavLink>
                  ))}
                </ContentContainer>
              }
            />
            <Route path="brainrot" element={<div>brainrot</div>} />
            <Route path="saved" element={<div>saved</div>} />
            <Route path="highlights" element={<div>highlight</div>} />
          </Routes>
        </div>
      </div>
    </main>
  );
}
