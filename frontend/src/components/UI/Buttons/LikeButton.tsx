import { HeartIcon } from "lucide-react";
import { useState } from "react";

interface LikeButtonProps {
  onLike: (event: React.MouseEvent<HTMLButtonElement>) => void;
  size: number;
}
export default function LikeButton(props: LikeButtonProps) {
  const [liked, setLiked] = useState(false);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setLiked((prev) => !prev);
    props.onLike(event);
  };

  return (
    <button className="like" onClick={handleClick} type="button">
      <HeartIcon
        className="transition duration-150 ease-in-out hover:text-accent-secondary"
        size={props.size}
        strokeWidth={liked ? 0 : undefined}
        fill={liked ? "oklch(63.7% 0.237 25.331)" : "none"}
      />
    </button>
  );
}
