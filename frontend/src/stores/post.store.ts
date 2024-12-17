import { create } from "zustand";
import useAuthStore from "./auth.store";
import { Post } from "../lib/types";

interface PostState {
    posts: Post[];
    fetchMode: "all" | "user";
    setFetchMode: (mode: "all" | "user") => void;
    fetchPosts: () => void;
}

const usePostStore = create<PostState>((set) => ({
    posts: [],
    fetchMode: "all",
    fetchPosts: async () => {
        let url = "";
        const headers: HeadersInit = {
            "Content-Type": "application/json",
        };

        if (usePostStore.getState().fetchMode === "user") {
            url = `${import.meta.env.VITE_API_URL}/posts/by-user/${useAuthStore.getState().id}`;
            const accessToken = useAuthStore.getState().accessToken;
            headers.Authorization = `Bearer ${accessToken}`;
        } else {
            url = `${import.meta.env.VITE_API_URL}/posts`;
        }

        try {
            const response = await fetch(url, {
                method: "GET",
                headers: headers,
                credentials: "include",
            });

            if (response.ok) {
                const data = await response.json();
                if (data.data != null) {
                    set({ posts: data.data });
                } else {
                    // user might have no posts
                    set({ posts: [] });
                }
            } else {
                const data = await response.json();
                console.log(data.errors[0].message);
            }
        } catch (error) {
            console.log(error);
        }
    },
    setFetchMode: (mode) => set({ fetchMode: mode }),
}));

export default usePostStore;
