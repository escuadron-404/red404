import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { BrowserRouter, Route, Routes } from "react-router";

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

const root = ReactDOM.createRoot(document.getElementById("root")!);

root.render(
	<StrictMode>
		<BrowserRouter>
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
		</BrowserRouter>
	</StrictMode>,
);
