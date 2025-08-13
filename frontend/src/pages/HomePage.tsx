import NavBar from "../components/NavBar";
import Post from "../components/UI/Post";

function HomePage() {
	return (
		<div className="w-full flex bg-base">
			<NavBar />
			<main className="w-full min-h-screen flex flex-col gap-10 pl-20 py-4">
				<Post />
				<Post />
				<Post />
			</main>
		</div>
	);
}
export default HomePage;
