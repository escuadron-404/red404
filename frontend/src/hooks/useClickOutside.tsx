import { useEffect } from "react";

function useClickOutside(
	ref: React.RefObject<HTMLDivElement | null>,
	callback: VoidFunction,
) {
	useEffect(() => {
		function handleClickOutside(event: MouseEvent) {
			const target = event.target instanceof Node ? event.target : null;
			if (ref.current && !ref.current.contains(target)) {
				callback();
			}
		}

		document.addEventListener("mousedown", handleClickOutside);
		return () => {
			document.removeEventListener("mousedown", handleClickOutside);
		};
	}, [ref, callback]);
}

export default useClickOutside;
