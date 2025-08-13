import HideButton from "../Buttons/HideButton";
import ReportButton from "../Buttons/ReportButton";
import UnfollowButton from "../Buttons/UnfollowButton";
export default function PostOptions() {
  return (
    <div className="absolute right-0  border border-accent p-2 rounded-md bg-base">
      <HideButton onClickHide={() => console.log("hidden")} />
      <ReportButton onClick={() => console.log("reported")} />
      <UnfollowButton onClick={() => console.log("unfollowed")} />
    </div>
  );
}
