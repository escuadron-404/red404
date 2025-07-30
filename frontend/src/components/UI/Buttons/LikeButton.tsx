import { useState } from "react";
import { HeartIcon } from "lucide-react";

interface LikeButtonProps {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
  size: number;
}
export default function LikeButton(props: LikeButtonProps) {
  const [liked, setLiked] = useState(false);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setLiked((prev) => !prev);
    props.onClick(event);
  };

  return (
    <button className="like" onClick={handleClick}>
      <HeartIcon
        className="transition duration-150 ease-in-out hover:text-accent-secondary"
        size={props.size}
        strokeWidth={liked ? 0 : undefined}
        fill={liked ? "oklch(63.7% 0.237 25.331)" : "none"}
      />
    </button>
  );
}
