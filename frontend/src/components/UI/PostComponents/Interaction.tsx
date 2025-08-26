import { useRef, useState } from "react";
import BookmarkButton from "@/components/UI/Buttons/BookmarkButton";
import CommentButton from "@/components/UI/Buttons/CommentButton";
import LikeButton from "@/components/UI/Buttons/LikeButton";
import ShareButton from "@/components/UI/Buttons/ShareButton";
import CommentsSection from "@/components/UI/CommentsSection";
import type { CommentProps } from "@/components/UI/CommentsSectionComponents/Comment";
import Comment from "@/components/UI/CommentsSectionComponents/Comment";
import ShareOptions from "@/components/UI/PostComponents/ShareOptions";
import useClickOutside from "@/hooks/useClickOutside";

interface InteractionProps {
  linkToCopy: string;
  comments: CommentProps[];
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
    <div className="interact flex gap-4 ">
      <LikeButton size={25} onLike={() => console.log("liked")} />
      <div className="flex items-center relative" ref={showCommentsRef}>
        <CommentButton onClick={() => setShowComments((prev) => !prev)} />
        {showComments && (
          <CommentsSection>
            {props.comments.map((comment) => (
              <Comment key={comment.id} {...comment} />
            ))}
          </CommentsSection>
        )}
      </div>

      <div className="flex items-center relative" ref={shareRef}>
        <ShareButton onClick={() => setShowShareOptions((prev) => !prev)} />
        {showShareOptions && <ShareOptions linkToCopy={props.linkToCopy} />}
      </div>
      <BookmarkButton onClick={() => console.log("bookmarked")} />
    </div>
  );
}
