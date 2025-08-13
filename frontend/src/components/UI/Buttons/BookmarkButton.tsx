import { BookmarkIcon } from "lucide-react";
import { useState } from "react";

interface BookmarkButtonProps {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
}
export default function BookmarkButton(props: BookmarkButtonProps) {
  const [bookmarked, setBookMarked] = useState(false);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setBookMarked((prev) => !prev);
    props.onClick(event);
  };

  return (
    <button className="bookmark" onClick={handleClick} type="button">
      <BookmarkIcon
        className="transition duration-150 ease-in-out hover:text-accent-secondary"
        width={28}
        height={28}
        strokeWidth={bookmarked ? 0 : undefined}
        fill={bookmarked ? "oklch(76.9% 0.188 70.08)" : "none"}
      />
    </button>
  );
}
