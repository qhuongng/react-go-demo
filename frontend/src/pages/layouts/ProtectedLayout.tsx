import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";

const ProtectedLayout = () => {
    const navigate = useNavigate();

    useEffect(() => {
        const isLoggedIn = localStorage.getItem("isLoggedIn");

        if (isLoggedIn == "false") {
            navigate("/login");
        }
    }, [navigate]);

    return (
        <div className="flex flex-col items-center py-8 h-full min-h-screen w-[calc(100%-34px)]">
            <Outlet />
        </div>
    );
};

export default ProtectedLayout;
