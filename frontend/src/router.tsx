import { createBrowserRouter } from "react-router-dom";

import RootLayout from "./pages/layouts/RootLayout";
import NonAuthLayout from "./pages/layouts/NonAuthLayout";
import ProtectedLayout from "./pages/layouts/ProtectedLayout";

import AllPosts from "./pages/AllPosts";
import NotFound from "./pages/NotFound";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import About from "./pages/About";
import YourPosts from "./pages/YourPosts";

const router = createBrowserRouter([
    {
        element: <RootLayout />,
        errorElement: <NotFound />,
        children: [
            {
                element: <NonAuthLayout />,
                children: [
                    {
                        path: "/",
                        element: <AllPosts />,
                    },
                    {
                        path: "/about",
                        element: <About />,
                    },
                    {
                        element: <ProtectedLayout />,
                        children: [{ path: "/you", element: <YourPosts /> }],
                    },
                ],
            },
            {
                path: "/login",
                element: <Login />,
            },
            {
                path: "/signup",
                element: <Signup />,
            },
        ],
    },
]);

export default router;
