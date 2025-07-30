interface CommentProps {
  text: string;
  author: string;
  authorAvatar: string;
}
export default function Comment(props: CommentProps) {
  return (
    <div className="w-full flex items-center gap-2">
      <img
        className="w-8 h-8 rounded-full"
        src={props.authorAvatar}
        alt="author-profile-photo"
      />
      <div className="flex">
        <span className="font-semibold">{props.author}:</span>
        <p className="pl-2">{props.text}</p>
      </div>
    </div>
  );
}
