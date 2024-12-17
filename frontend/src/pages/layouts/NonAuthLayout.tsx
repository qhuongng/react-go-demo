import { Outlet } from "react-router-dom";
import Navbar from "../../components/Navbar";

const NonAuthLayout = () => {
    return (
        <div className="h-full min-h-screen w-[calc(100%-34px)]">
            <Navbar />
            <div className="flex flex-col items-center py-8">
                <Outlet />
            </div>
        </div>
    );
};

export default NonAuthLayout;
