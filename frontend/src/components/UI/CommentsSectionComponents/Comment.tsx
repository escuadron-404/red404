import LikeButton from "../Buttons/LikeButton";

export interface CommentProps {
	id: number;
	text: string;
	author: string;
	authorAvatar: string;
	isReply: boolean;
	replies: CommentProps | null;
	likes: number;
}
export default function Comment(props: CommentProps) {
	return (
		<div className={`w-full flex justify-between" ${props.isReply && "pl-10"}`}>
			<img
				className="w-8 h-8 rounded-full mr-2"
				src={props.authorAvatar}
				alt="author-profile-photo"
			/>
			<div className="w-full flex justify-between">
				<div>
					<div className="flex">
						<span className="font-semibold">{props.author}:</span>
						<p className="pl-2">{props.text}</p>
					</div>
					<div className="flex items-center gap-3">
						<span className="text-text-light text-sm">3m</span>
						{props.likes >= 1 && (
							<span className="font-bold text-text-light text-sm">
								{props.likes} likes
							</span>
						)}
						<button className="font-bold text-text-light text-sm">reply</button>
						{props.replies && (
							<button className="font-bold text-text-light text-sm">
								view replies
							</button>
						)}
					</div>
				</div>
				<LikeButton size={15} onClick={() => console.log("liked")} />
			</div>
		</div>
	);
}
