import { useRef, useState } from "react";
import useClickOutside from "../../../hooks/useClickOutside";
import BookmarkButton from "../Buttons/BookmarkButton";
import CommentButton from "../Buttons/CommentButton";
import LikeButton from "../Buttons/LikeButton";
import ShareButton from "../Buttons/ShareButton";
import CommentsSection from "../CommentsSection";
import type { CommentProps } from "../CommentsSectionComponents/Comment";
import Comment from "../CommentsSectionComponents/Comment";
import ShareOptions from "./ShareOptions";

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
