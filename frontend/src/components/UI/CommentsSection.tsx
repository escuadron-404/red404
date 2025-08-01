import { type ReactNode } from "react";
import SendButton from "./Buttons/SendButton";

interface CommentsProps {
	children: ReactNode;
}
export default function CommentsSection(props: CommentsProps) {
	return (
		<div className="fixed z-10 bg-base w-lg flex flex-col justify-between top-0 right-0 h-screen border-l border-accent p-4">
			{!props.children ? (
				<p className="self-center my-auto">
					no comments found, be the first to comment!
				</p>
			) : (
				<div className="flex flex-col gap-7">{props.children}</div>
			)}

			<div className="bg-accent rounded-md relative">
				<input
					className=" outline-accent-secondary w-full rounded-md focus:outline p-2 placeholder:text-sm"
					type="text"
					name="comment"
					id="comment"
					placeholder="add a comment"
				/>
				<SendButton onClick={() => console.log("sent")} />
			</div>
		</div>
	);
}
