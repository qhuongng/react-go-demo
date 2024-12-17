import { Outlet } from "react-router-dom";

const RootLayout = () => {
    return (
        <div className="flex flex-col justify-center items-center min-h-screen">
            <Outlet />
        </div>
    );
};

export default RootLayout;
