import Post from "@/components/UI/Post";

function HomePage() {
  return (
    <main className="w-full min-h-screen flex flex-col gap-10 pl-20 py-4">
      <Post />
      <Post />
      <Post />
    </main>
  );
}
export default HomePage;
