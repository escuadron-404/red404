import { useState, useRef } from "react";
import BookmarkButton from "../Buttons/BookmarkButton";
import CommentButton from "../Buttons/CommentButton";
import LikeButton from "../Buttons/LikeButton";
import ShareButton from "../Buttons/ShareButton";
import useClickOutside from "../../../hooks/useClickOutside";
import CommentsSection from "./CommentsSection";
import Comment from "./CommentsSectionComponents/Comment";
import ShareOptions from "./ShareOptions";
interface InteractionProps {
  linkToCopy: string;
  comments: {
    id: number;
    author: string;
    authorAvatar: string;
    text: string;
  }[];
}
export default function Interaction(props: InteractionProps) {
  const [showComments, setShowComments] = useState(false);
  const [showShareOptions, setShowShareOptions] = useState(false);
  const showCommentsRef = useRef(null);
  const shareRef = useRef(null);
  useClickOutside(showCommentsRef, () => {
    showComments && setShowComments(false);
  });
  useClickOutside(shareRef, () => {
    showShareOptions && setShowShareOptions(false);
  });
  return (
    <div className="interact flex gap-4">
      <LikeButton onClick={() => console.log("liked")} />
      <CommentButton onClick={() => setShowComments((prev) => !prev)} />
      {showComments && (
        <CommentsSection ref={showCommentsRef}>
          {props.comments.map((comment) => (
            <Comment {...comment} />
          ))}
        </CommentsSection>
      )}
      <div className="flex items-center relative" ref={shareRef}>
        <ShareButton onClick={() => setShowShareOptions((prev) => !prev)} />
        {showShareOptions && <ShareOptions linkToCopy={props.linkToCopy} />}
      </div>
      <BookmarkButton onClick={() => console.log("bookmarked")} />
    </div>
  );
}
