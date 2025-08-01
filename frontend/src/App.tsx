import "./App.css";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import Auth from "./pages/Auth";
import Brainrot from "./pages/BrainrotPage";
import ExplorePage from "./pages/ExplorePage";
import HomePage from "./pages/HomePage";
import MessagePage from "./pages/MessagesPage";
import NotFoundPage from "./pages/NotFoundPage";

function App() {
	return (
		<Router>
			<Routes>
				<Route path="/" element={<Auth />} />
				<Route path="/home" element={<HomePage />} />
				<Route path="/explore" element={<ExplorePage />} />
				<Route path="/messages" element={<MessagePage />} />
				<Route path="/brainrot" element={<Brainrot />} />
				<Route path="*" element={<NotFoundPage />} />
			</Routes>
		</Router>
	);
}

export default App;
