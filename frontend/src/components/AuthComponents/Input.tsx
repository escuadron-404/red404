interface InputProps {
  id: string;
  name: string;
  placeholder: string;
  type: "email" | "text";
}

export default function Input(props: InputProps) {
  return (
    <input
      className="p-2 outline outline-accent-secondary focus:outline-2 rounded-xl transition duration-150 ease-in-out hover:bg-accent"
      type={props.type}
      name={props.name}
      id={props.id}
      placeholder={props.placeholder}
    />
  );
}
