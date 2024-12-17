import { useEffect, useState } from "react";
import { Link, NavLink, useNavigate, useLocation } from "react-router-dom";

import { RiMenu2Fill } from "react-icons/ri";

import usePostStore from "../stores/post.store";
import useAuthStore from "../stores/auth.store";

import { checkAuth } from "../lib/utils";

const NavItems = [
    { name: "Posts", href: "/", protected: false },
    { name: "*Your* posts", href: "/you", protected: true },
    { name: "What is this?!", href: "/about", protected: false },
];

const Navbar = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const [showDropdown, setShowDropdown] = useState(false);

    const accessToken = useAuthStore((state) => state.accessToken);
    const removeAuth = useAuthStore((state) => state.removeAuth);

    const setFetchMode = usePostStore((state) => state.setFetchMode);

    const toggleDropdown = () => {
        setShowDropdown(!showDropdown);
        if (showDropdown === false) {
            document.getElementById("dropdown-close-helper")?.focus({ preventScroll: true });
        }
    };

    const handleLogout = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/logout`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
            });

            if (response.ok) {
                removeAuth();
                localStorage.setItem("isLoggedIn", "false");
                navigate("/login");
            } else {
                const data = await response.json();
                console.log(data.errors[0].message);
            }
        } catch (error) {
            console.log(error);
        }
    };

    useEffect(() => {
        // call this once when the page loads for the first time
        // because for some reason, the dropdown automatically closes on the first open
        toggleDropdown();
    }, []);

    useEffect(() => {
        // refreshing the page clears the store
        // so we need to restore the login state if isLoggedIn (set in localStorage) is true
        checkAuth();
        setFetchMode(location.pathname === "/you" ? "user" : "all");
    }, []);

    useEffect(() => {
        setFetchMode(location.pathname === "/you" ? "user" : "all");
    }, [location.pathname, setFetchMode]);

    return (
        <div className="navbar sticky top-4 z-50 bg-base-200 mt-4 rounded-box border border-base-content">
            <div className="navbar-start">
                <div className="dropdown">
                    <div
                        tabIndex={0}
                        role="button"
                        className="btn lg:hidden"
                        onClick={toggleDropdown}
                    >
                        <RiMenu2Fill size={20} />
                    </div>

                    <ul
                        tabIndex={0}
                        className="menu menu-sm dropdown-content bg-base-200 rounded-box z-[1] mt-3 w-52 p-2 border border-base-content"
                    >
                        {NavItems.map((item) =>
                            item.protected && accessToken === "" ? null : (
                                <li key={item.name} className="my-1">
                                    <NavLink to={item.href} onClick={toggleDropdown}>
                                        {item.name}
                                    </NavLink>
                                </li>
                            )
                        )}
                    </ul>
                </div>

                <Link className="btn btn-ghost text-xl" to="/">
                    Minimal React + Go SNS
                </Link>
            </div>

            <div className="navbar-center hidden lg:flex">
                <ul className="menu menu-horizontal px-1">
                    {NavItems.map((item) =>
                        item.protected && accessToken === "" ? null : (
                            <li key={item.name} className="mx-1">
                                <NavLink to={item.href}>{item.name}</NavLink>
                            </li>
                        )
                    )}
                </ul>
            </div>

            <div className="navbar-end">
                {accessToken !== "" ? (
                    <button className="btn" onClick={handleLogout}>
                        Log out
                    </button>
                ) : (
                    <Link className="btn btn-primary" to="/login">
                        Log in
                    </Link>
                )}
            </div>

            {/* invisible button to close the dropdown after selection */}
            <button type="button" id="dropdown-close-helper" className="h-0 w-0" />
        </div>
    );
};

export default Navbar;
