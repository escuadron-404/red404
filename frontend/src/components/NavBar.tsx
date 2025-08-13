import {
	BrainIcon,
	CompassIcon,
	HomeIcon,
	MessageCircleDashedIcon,
	SearchIcon,
	UserCircle2Icon,
} from "lucide-react";
import NavBarLink from "./NavBarComponents/NavBarLink";

function NavBar() {
	return (
		<>
			<div className="flex flex-col items-center py-4 pl-2 pr-6 border-r border-accent-secondary">
				<header className="mb-9">
					<span className="font-bold text-xl text-primary">red404</span>
				</header>
				<nav className="flex flex-col gap-5">
					<NavBarLink link="/home">
						<HomeIcon size={30} />
					</NavBarLink>
					<NavBarLink link="/search">
						<SearchIcon size={30} />
					</NavBarLink>
					<NavBarLink link="/explore">
						<CompassIcon size={30} />
					</NavBarLink>
					<NavBarLink link="/messages">
						<MessageCircleDashedIcon size={30} />
					</NavBarLink>
					<NavBarLink link="/brainrot">
						<BrainIcon size={30} />
					</NavBarLink>
					<NavBarLink link="/profile">
						<UserCircle2Icon size={30} />
					</NavBarLink>
				</nav>
			</div>
		</>
	);
}

export default NavBar;
