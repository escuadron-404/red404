import { Outlet } from "react-router";
import NavBar from "../components/NavBar";

const HomeLayout = () => {
  return (
    <div className="w-full h-screen flex bg-base">
      <NavBar />
      <Outlet />
    </div>
  );
};

export default HomeLayout;
