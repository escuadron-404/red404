import ReportButton from "../Buttons/ReportButton";
import UnfollowButton from "../Buttons/UnfollowButton";
export default function PostOptions() {
  return (
    <div className="absolute right-0  border border-accent p-2 rounded-md bg-base">
      <UnfollowButton onClick={() => console.log("unfollowed")} />
      <ReportButton onClick={() => console.log("reported")} />
    </div>
  );
}
