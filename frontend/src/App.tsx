import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import "./App.css";

// Layouts
import HomeLayout from "./layouts/HomeLayout";
import AuthPage from "./pages/Auth";
import Brainrot from "./pages/BrainrotPage";
import ExplorePage from "./pages/ExplorePage";
import HomePage from "./pages/HomePage";
import MessagePage from "./pages/MessagesPage";
import NotFoundPage from "./pages/NotFoundPage";
import ProfilePage from "./pages/ProfilePage";
import SearchPage from "./pages/SearchPage";

function App() {
	return (
		<Router>
			<Routes>
				<Route path="/" element={<AuthPage />} />
				<Route path="/auth" />

				<Route element={<HomeLayout />}>
					<Route path="/home" element={<HomePage />} />
					<Route path="/search" element={<SearchPage />} />
					<Route path="/explore" element={<ExplorePage />} />
					<Route path="/messages" element={<MessagePage />} />
					<Route path="/brainrot" element={<Brainrot />} />
					<Route path="/profile/*" element={<ProfilePage />} />
				</Route>

				<Route path="*" element={<NotFoundPage />} />
			</Routes>
		</Router>
	);
}

export default App;
