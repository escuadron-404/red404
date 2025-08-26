import Button from "./Button";

function MessageButton() {
  return <Button text="message" onClick={() => alert("message")} />;
}

export default MessageButton;
