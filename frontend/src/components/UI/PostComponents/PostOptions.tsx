import HideButton from "@/components/UI/Buttons/HideButton";
import ReportButton from "@/components/UI/Buttons/ReportButton";
import UnfollowButton from "@/components/UI/Buttons/UnfollowButton";
export default function PostOptions() {
  return (
    <div className="absolute right-0  border border-accent p-2 rounded-md bg-base">
      <HideButton onClickHide={() => console.log("hidden")} />
      <ReportButton onClick={() => console.log("reported")} />
      <UnfollowButton onClick={() => console.log("unfollowed")} />
    </div>
  );
}
