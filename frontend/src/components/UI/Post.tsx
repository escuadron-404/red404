import { useState, useRef } from "react";
import useClickOutside from "../../hooks/useClickOutside";
import { MapPinIcon, DotIcon, EllipsisIcon, MinusIcon } from "lucide-react";
import PostOptions from "./PostComponents/PostOptions";
import Interaction from "./PostComponents/Interaction";

const postData = {
  id: "123013",
  user: "pixshanghai",
  location: "lechería",
  avatar: "/testUserImage.jpg",
  postMedia: "/testPostImage.jpg",
  description: "going back to the phase where cows are in my head.",
  createdAt: "10h",
  url: "https://red404.app/post/123013",
};
const commentData = [
  {
    id: 1,
    author: "leonardo",
    authorAvatar: "/testUserImage.jpg",
    text: "how beautiful",
    isReply: false,
    replies: 1,
    likes: 0,
  },
  {
    id: 2,
    author: "mark",
    authorAvatar: "/testUserImage.jpg",
    text: "󰱱",
    isReply: true,
    replies: null,
    likes: 2,
  },
  {
    id: 3,
    author: "tobias",
    authorAvatar: "/testUserImage.jpg",
    text: "i need to buy some meat",
    isReply: false,
    replies: null,
    likes: 3,
  },
];
function PostHeader() {
  const [showPostOptions, setShowPostOptions] = useState(false);
  const dropdownRef = useRef(null);
  useClickOutside(dropdownRef, () => {
    showPostOptions && setShowPostOptions(false);
  });
  return (
    <div className="author justify-between flex gap-5">
      <div className="flex gap-2">
        <img
          className="w-10 h-10 rounded-full"
          src={postData.avatar}
          alt="user-profile-photo"
        />
        <div className="info flex flex-col gap-1">
          <div className="flex items-center gap-1">
            <span className="post-author font-bold">{postData.user}</span>
            <time className="flex items-center">
              <DotIcon />
              <span className="text-text-light">{postData.createdAt}</span>
            </time>
          </div>
          <address className="flex items-center gap-1.5">
            <MapPinIcon width={20} height={20} />
            <span className="text-sm">{postData.location}</span>
          </address>
        </div>
      </div>
      <div className="relative" ref={dropdownRef}>
        <button
          className="group"
          onClick={() => setShowPostOptions((prev) => !prev)}
        >
          <EllipsisIcon
            className="transition duration-150 ease-in-out hover:text-accent-secondary"
            width={20}
            height={20}
          />
        </button>
        {showPostOptions && <PostOptions />}
      </div>
    </div>
  );
}
function Post() {
  return (
    <div className="post flex flex-col gap-3 w-md" id={postData.id}>
      <PostHeader />
      <img
        className="w-full h-96 object-cover rounded-md"
        src={postData.postMedia}
        alt="post-media"
      />
      <Interaction linkToCopy={postData.url} comments={commentData} />
      {/* footer */}
      <p className="text-[0.96rem]">
        <span className="post-author font-bold">{postData.user}:</span>{" "}
        {postData.description}
      </p>
      <p className="pl-2 text-[0.96rem] flex gap-0.5">
        <MinusIcon />
        <span className="post-author font-bold">simonbolivar:</span> me fr dude
      </p>
    </div>
  );
}
export default Post;
