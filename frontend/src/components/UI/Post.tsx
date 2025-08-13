import { DotIcon, EllipsisIcon, MapPinIcon, MinusIcon } from "lucide-react";
import { useRef, useState } from "react";
import useClickOutside from "../../hooks/useClickOutside";
import postData from "../../templates/PostTemplate.tsx";
import HideButton from "./Buttons/HideButton.tsx";
import ReportButton from "./Buttons/ReportButton";
import UnfollowButton from "./Buttons/UnfollowButton";
import Interaction from "./PostComponents/Interaction";

interface PostActions {
  onClickHide: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
function PostOptions(props: PostActions) {
  return (
    <div className="absolute right-0  border border-accent p-2 rounded-md bg-base">
      <HideButton {...props} />
      <ReportButton onClick={() => console.log("reported")} />
      <UnfollowButton onClick={() => console.log("unfollowed")} />
    </div>
  );
}

function PostHeader(props: PostActions) {
  const [showPostOptions, setShowPostOptions] = useState(false);
  const dropdownRef = useRef(null);
  useClickOutside(dropdownRef, () => {
    showPostOptions && setShowPostOptions(false);
  });

  return (
    <div className="author justify-between flex gap-5">
      <div className="flex gap-2">
        <a href="/profile">
          <img
            className="w-10 h-10 rounded-full"
            src={postData.avatar}
            alt="user-profile-photo"
          />
        </a>
        <div className="info flex flex-col gap-1">
          <div className="flex items-center gap-1">
            <a href="/profile" className="post-author font-bold">
              {postData.user}
            </a>
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
          type="button"
        >
          <EllipsisIcon
            className="transition duration-150 ease-in-out hover:text-accent-secondary"
            width={20}
            height={20}
          />
        </button>
        {showPostOptions && <PostOptions {...props} />}
      </div>
    </div>
  );
}
function Post() {
  const [hidePost, setHidePost] = useState(false);
  function handleHidePost() {
    setHidePost(!hidePost);
  }
  if (hidePost)
    return (
      <div className="flex bg-accent items-center justify-center w-md p-5 rounded-md">
        <span>Post has been hidden</span>
      </div>
    );
  return (
    <div className="flex flex-col gap-3 w-md" id={postData.id}>
      <PostHeader onClickHide={handleHidePost} />
      <img
        className="w-full h-96 object-cover rounded-md"
        src={postData.postMedia}
        alt="post-media"
      />
      <Interaction
        linkToCopy={postData.url}
        comments={postData.commentsData.comments}
      />
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
