import { GoogleOAuthProvider } from "@react-oauth/google";
import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { BrowserRouter, Route, Routes } from "react-router";

// Layout
import HomeLayout from "./layouts/HomeLayout";
// Pages
import Brainrot from "./pages/BrainrotPage";
import ExplorePage from "./pages/ExplorePage";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/loginPage";
import MessagePage from "./pages/MessagesPage";
import NotFoundPage from "./pages/NotFoundPage";
import ProfilePage from "./pages/ProfilePage";
import RegisterPage from "./pages/registerPage";
import SearchPage from "./pages/SearchPage";

let root = document.getElementById("root");

if (!root) {
  root = document.createElement("div");
  document.body.appendChild(root);
}

ReactDOM.createRoot(root).render(
  <StrictMode>
    <GoogleOAuthProvider clientId={import.meta.env.VITE_GOOGLE_AUTH_CLIENT_ID}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

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
    </GoogleOAuthProvider>
  </StrictMode>,
);
