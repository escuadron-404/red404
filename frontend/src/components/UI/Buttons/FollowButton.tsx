import Button from "../ProfileComponents/Buttons/Button";

interface FollowButtonProps {
  onFollow: (event: React.MouseEvent<HTMLButtonElement>) => void;
  followStatus: boolean;
}

function FollowButton(props: FollowButtonProps) {
  return (
    <Button
      onClick={props.onFollow}
      text={props.followStatus ? "followed" : "follow"}
      className="w-24"
    />
  );
}

export default FollowButton;
