import { Link2Icon, CopyCheckIcon, CheckCircle2Icon } from "lucide-react";
import { useState } from "react";
import testUserImage from "../../../assets/testUserImage.jpg";

interface ShareOptionsProps {
  linkToCopy: string;
}

function Friend() {
  return (
    <div className="flex flex-col items-center">
      <img
        className="w-10 h-10 rounded-full object-cover"
        src={testUserImage}
        alt="t"
      />
      <span className="mt-1 text-xs text-text-light group-hover:text-inherit">
        tobias
      </span>
    </div>
  );
}
function SentToFriend() {
  return (
    <div className="flex flex-col">
      <CheckCircle2Icon size={40} />{" "}
      <span className="mt-1 text-xs text-text-light group-hover:text-inherit">
        sent
      </span>
    </div>
  );
}

interface CopyButtonProps {
  linkToCopy: string;
}
function CopyButton(props: CopyButtonProps) {
  const [isCopied, setIsCopied] = useState(false);

  async function handleCopy() {
    try {
      setIsCopied(true);
      await navigator.clipboard.writeText(props.linkToCopy);
      setTimeout(() => setIsCopied(false), 3000);
    } catch (err) {
      console.log("error copying link", err);
    }
  }
  return (
    <div className="flex flex-col items-center gap-0.5 group">
      <button
        className="flex flex-col items-center justify-center w-10 h-10 p-2 rounded-full border border-accent"
        onClick={handleCopy}
      >
        {isCopied ? <CopyCheckIcon size={25} /> : <Link2Icon size={25} />}
      </button>
      <span className="text-xs text-text-light group-hover:text-inherit">
        {isCopied ? "copied" : "link"}
      </span>
    </div>
  );
}
export default function ShareOptions(props: ShareOptionsProps) {
  const users = Array(3).fill(null);
  const [submittedToFriends, setSubmittedToFriends] = useState(
    Array(users.length).fill(false)
  );
  {
    /* se debe añadir una funcion que muestre los 3 usuarios con el que el user interactua más*/
  }

  async function handleSubmitToFriend(index: number) {
    setSubmittedToFriends((prev) => {
      const updated = [...prev];
      updated[index] = true;
      return updated;
    });
  }
  return (
    <div className="flex items-center gap-6 absolute min-w-max bottom-10 left-0 bg-base border py-3 px-6 border-accent rounded-md">
      <CopyButton linkToCopy={props.linkToCopy} />
      {users.map((_, index) => (
        <button
          className="flex flex-col items-center group cursor-pointer"
          key={index}
          onClick={() => handleSubmitToFriend(index)}
        >
          {submittedToFriends[index] ? <SentToFriend /> : <Friend />}
        </button>
      ))}
    </div>
  );
}
