import NavBar from "../components/NavBar";
import Post from "../components/UI/Post";

function HomePage() {
	return (
		<div className="w-full flex bg-base">
			<NavBar />
			<main className="w-full flex flex-col gap-10 pl-20 py-4">
				<Post PostID="t2sld21sl" />
				<Post PostID="t2sld21sl" />
				<Post PostID="t2sld21sl" />
			</main>
		</div>
	);
}
export default HomePage;
